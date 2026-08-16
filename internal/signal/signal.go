// Package signal provides structural-vibration time-series parsing and
// statistics: RMS, peak, discrete Fourier transform, and dominant frequency.
//
// A Series is a single labelled sensor channel. Long-format CSV rows of the form
// label,samplerate,value are grouped by label into one Series per label.
package signal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"strings"
)

// Series is one labelled sensor channel.
type Series struct {
	Label      string
	SampleRate float64
	Data       []float64
}

// ParseSeries reads a long-format CSV (header label,samplerate,value) and groups
// rows by label into one Series per label. It returns an error if the file cannot
// be opened, the header/rows are malformed, or the samplerate/value columns are
// not numeric.
func ParseSeries(path string) ([]Series, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open series file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = 3

	rec, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if len(rec) != 3 {
		return nil, fmt.Errorf("malformed header: expected 3 columns, got %d", len(rec))
	}
	// A header row has non-numeric samplerate/value columns; a data row has numeric ones.
	startIsData := isNumeric(rec[1]) && isNumeric(rec[2])

	byLabel := make(map[string]*Series)
	var order []string

	add := func(fields []string) error {
		if len(fields) != 3 {
			return fmt.Errorf("malformed row %v: expected 3 fields", fields)
		}
		label := strings.TrimSpace(fields[0])
		sr, e1 := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if e1 != nil {
			return fmt.Errorf("malformed samplerate %q: %v", fields[1], e1)
		}
		val, e2 := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if e2 != nil {
			return fmt.Errorf("malformed value %q: %v", fields[2], e2)
		}
		ss, ok := byLabel[label]
		if !ok {
			ss = &Series{Label: label, SampleRate: sr}
			byLabel[label] = ss
			order = append(order, label)
		}
		ss.Data = append(ss.Data, val)
		return nil
	}

	if startIsData {
		if err := add(rec); err != nil {
			return nil, err
		}
	}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record: %w", err)
		}
		if err := add(rec); err != nil {
			return nil, err
		}
	}

	result := make([]Series, 0, len(order))
	for _, l := range order {
		result = append(result, *byLabel[l])
	}
	if len(result) == 0 {
		return nil, errors.New("no data rows found")
	}
	return result, nil
}

// RMS returns sqrt(mean(x^2)) over Data. Empty Data returns 0.
func RMS(s Series) float64 {
	if len(s.Data) == 0 {
		return 0
	}
	var sum float64
	for _, v := range s.Data {
		sum += v * v
	}
	return math.Sqrt(sum / float64(len(s.Data)))
}

// Peak returns the maximum absolute value in Data. Empty Data returns 0.
func Peak(s Series) float64 {
	if len(s.Data) == 0 {
		return 0
	}
	m := 0.0
	for _, v := range s.Data {
		if v > m {
			m = v
		}
	}
	return m
}

// DFT computes the O(n^2) discrete Fourier transform of the real Data.
// The result length equals len(Data), with X[k] = sum_n x[n]*exp(-2πi kn/N).
func DFT(s Series) []complex128 {
	N := len(s.Data)
	X := make([]complex128, N)
	for k := 0; k < N; k++ {
		var re, im float64
		for n := 0; n < N; n++ {
			phi := -2 * math.Pi * float64(k) * float64(n) / float64(N)
			re += s.Data[n] * math.Cos(phi)
			im += s.Data[n] * math.Sin(phi)
		}
		X[k] = complex(re, im)
	}
	return X
}

// DominantFreq returns the frequency (k * SampleRate / N) of the maximum |X[k]|
// for k in 1..N/2 (DC component k=0 excluded). Empty Data returns 0.
func DominantFreq(s Series) float64 {
	N := len(s.Data)
	if N == 0 {
		return 0
	}
	X := DFT(s)
	limit := N / 2
	if limit < 1 {
		limit = 1
	}
	maxMag := 0.0
	maxK := 1
	for k := 1; k <= limit; k++ {
		mag := cmplx.Abs(X[k])
		if mag > maxMag {
			maxMag = mag
			maxK = k
		}
	}
	return float64(maxK) * s.SampleRate / float64(N)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil
}

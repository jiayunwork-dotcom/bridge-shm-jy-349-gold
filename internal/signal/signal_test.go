package signal

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

// TestParseSeries covers the error path (missing file, malformed value) and the
// slice path (grouping rows into one Series per label, preserving order).
func TestParseSeries(t *testing.T) {
	// error: missing file
	if _, err := ParseSeries("does_not_exist.csv"); err == nil {
		t.Fatal("expected error for missing file")
	}

	dir := t.TempDir()
	p := filepath.Join(dir, "s.csv")
	content := "label,samplerate,value\n" +
		"accel,32,1\naccel,32,2\naccel,32,3\n" +
		"strain,16,4\nstrain,16,5\n"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	series, err := ParseSeries(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(series) != 2 {
		t.Fatalf("expected 2 series, got %d", len(series))
	}
	if series[0].Label != "accel" || len(series[0].Data) != 3 {
		t.Fatalf("accel series wrong: %+v", series[0])
	}
	if series[1].Label != "strain" || len(series[1].Data) != 2 {
		t.Fatalf("strain series wrong: %+v", series[1])
	}

	// error: malformed (non-numeric) value
	bad := filepath.Join(dir, "bad.csv")
	os.WriteFile(bad, []byte("label,samplerate,value\naccel,32,notnum\n"), 0644)
	if _, err := ParseSeries(bad); err == nil {
		t.Fatal("expected error for malformed value")
	}
}

// TestParseSeries_HeaderOnlyErrors checks that a file with only a header and
// no data rows returns an error.
func TestParseSeries_HeaderOnlyErrors(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.csv")
	if err := os.WriteFile(p, []byte("label,samplerate,value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseSeries(p); err == nil {
		t.Fatal("expected error for header-only file")
	}
}

// TestRMS checks RMS of [3,4] = sqrt((9+16)/2) = sqrt(12.5) and that an empty
// series is 0.
func TestRMS(t *testing.T) {
	s := Series{Label: "x", SampleRate: 1, Data: []float64{3, 4}}
	got := RMS(s)
	want := math.Sqrt(12.5)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("RMS = %v, want %v", got, want)
	}
	if RMS(Series{}) != 0 {
		t.Fatal("RMS of empty series must be 0")
	}
}

// TestDFT checks the result length equals len(Data) and that the DC bin X[0].Re
// equals the sum of the data.
func TestDFT(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	X := DFT(Series{Data: data})
	if len(X) != len(data) {
		t.Fatalf("DFT len = %d, want %d", len(X), len(data))
	}
	wantRe := 0.0
	for _, v := range data {
		wantRe += v
	}
	if math.Abs(real(X[0])-wantRe) > 1e-9 {
		t.Fatalf("DC bin X[0].Re = %v, want %v", real(X[0]), wantRe)
	}
}

// TestDFT_EmptyNonNil checks that an empty series yields a non-nil empty
// coefficient slice.
func TestDFT_EmptyNonNil(t *testing.T) {
	X := DFT(Series{})
	if X == nil {
		t.Fatal("DFT of empty series must return empty non-nil slice")
	}
	if len(X) != 0 {
		t.Fatalf("DFT of empty series len = %d, want 0", len(X))
	}
}

// TestPeak_Negative checks Peak of an all-negative series is the max absolute
// value, not 0.
func TestPeak_Negative(t *testing.T) {
	s := Series{Label: "x", SampleRate: 1, Data: []float64{-3, -1, -2}}
	if got := Peak(s); math.Abs(got-3) > 1e-9 {
		t.Fatalf("Peak = %v, want 3", got)
	}
}

// TestDominantFreq checks a pure 2 Hz sine at fs=32 maps to dominant frequency 2 Hz,
// and that an empty series yields 0.
func TestDominantFreq(t *testing.T) {
	N := 64
	fs := 32.0
	data := make([]float64, N)
	for n := 0; n < N; n++ {
		t := float64(n) / fs
		data[n] = math.Sin(2 * math.Pi * 2 * t)
	}
	s := Series{Label: "accel", SampleRate: fs, Data: data}
	if got := DominantFreq(s); math.Abs(got-2.0) > 1e-6 {
		t.Fatalf("DominantFreq = %v, want 2.0", got)
	}
	if DominantFreq(Series{}) != 0 {
		t.Fatal("DominantFreq of empty series must be 0")
	}
}

package modal

import (
	"math"
	"testing"

	"bridge-shm/internal/signal"
)

// TestNaturalFreq checks NaturalFreq equals DominantFreq for a pure 2 Hz sine.
func TestNaturalFreq(t *testing.T) {
	s := signal.Series{Label: "a", SampleRate: 32, Data: make([]float64, 64)}
	for n := range s.Data {
		tv := float64(n) / 32
		s.Data[n] = math.Sin(2 * math.Pi * 2 * tv)
	}
	if got := NaturalFreq(s); math.Abs(got-2.0) > 1e-6 {
		t.Fatalf("NaturalFreq = %v, want 2.0", got)
	}
}

// TestDampingRatio checks the logarithmic-decrement formula on a clean geometric
// decay (ratio 0.9) and that a series with <2 samples is 0.
func TestDampingRatio(t *testing.T) {
	ratio := 0.9
	var data []float64
	v := 10.0
	for i := 0; i < 20; i++ {
		data = append(data, v)
		v *= ratio
	}
	s := signal.Series{Label: "decay", Data: data}
	got := DampingRatio(s)
	want := (math.Log(1.0 / ratio)) / (2 * math.Pi)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("DampingRatio = %v, want %v", got, want)
	}
	if DampingRatio(signal.Series{}) != 0 {
		t.Fatal("DampingRatio of len<2 must be 0")
	}
}

// TestDampingRatio_SkipsRisingPairs checks that a rising pair is skipped and
// only the decaying pair contributes to the mean log decrement.
func TestDampingRatio_SkipsRisingPairs(t *testing.T) {
	s := signal.Series{Label: "decay", Data: []float64{1, 2, 1.8}}
	got := DampingRatio(s)
	want := math.Log(2/1.8) / (2 * math.Pi)
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("DampingRatio = %v, want %v", got, want)
	}
}

// TestFreqShift checks a frequency drop yields a negative shift and that a zero
// baseline yields 0.
func TestFreqShift(t *testing.T) {
	if got := FreqShift(10, 9); math.Abs(got+0.1) > 1e-9 {
		t.Fatalf("FreqShift(10,9) = %v, want -0.1", got)
	}
	if FreqShift(10, 11) <= 0 {
		t.Fatal("FreqShift increase should be positive")
	}
	if FreqShift(0, 5) != 0 {
		t.Fatal("FreqShift f0=0 should be 0")
	}
}

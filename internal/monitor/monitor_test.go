package monitor

import (
	"math"
	"testing"
)

// TestControlLimits checks that zero spread yields lo=hi=mean and that a spread
// yields mean ± k*std.
func TestControlLimits(t *testing.T) {
	lo, hi := ControlLimits([]Sample{{Freq: 5}, {Freq: 5}, {Freq: 5}}, 2)
	if lo != 5 || hi != 5 {
		t.Fatalf("lo=%v hi=%v, want 5,5", lo, hi)
	}
	lo2, hi2 := ControlLimits([]Sample{{Freq: 4}, {Freq: 6}}, 1)
	mean := 5.0
	std := math.Sqrt(1.0)
	if math.Abs(lo2-(mean-std)) > 1e-9 || math.Abs(hi2-(mean+std)) > 1e-9 {
		t.Fatalf("lo2=%v hi2=%v, want %.3f,%.3f", lo2, hi2, mean-std, mean+std)
	}
}

// TestAnomalies checks the nil path (empty input) and the slice path (indices
// outside [lo,hi]).
func TestAnomalies(t *testing.T) {
	if got := Anomalies(nil, 0, 0); got != nil {
		t.Fatalf("empty samples should return nil, got %v", got)
	}
	samples := []Sample{{Freq: 5}, {Freq: 1}, {Freq: 5}, {Freq: 9}}
	idx := Anomalies(samples, 4, 6)
	if len(idx) != 2 || idx[0] != 1 || idx[1] != 3 {
		t.Fatalf("Anomalies = %v, want [1 3]", idx)
	}
}

// TestAnomalies_ResultsIndependent checks that a later Anomalies call does not
// overwrite the slice returned by an earlier call.
func TestAnomalies_ResultsIndependent(t *testing.T) {
	a := []Sample{{Freq: 1}, {Freq: 5}, {Freq: 9}}
	first := Anomalies(a, 4, 6)
	if len(first) != 2 || first[0] != 0 || first[1] != 2 {
		t.Fatalf("first = %v, want [0 2]", first)
	}
	b := []Sample{{Freq: 0}, {Freq: 0}, {Freq: 5}, {Freq: 0}}
	second := Anomalies(b, 4, 6)
	if len(second) != 3 || second[0] != 0 || second[1] != 1 || second[2] != 3 {
		t.Fatalf("second = %v, want [0 1 3]", second)
	}
	if len(first) != 2 || first[0] != 0 || first[1] != 2 {
		t.Fatalf("first mutated after second call: %v", first)
	}
}

// TestWarn checks the warning fires only when the frequency drop exceeds the
// threshold.
func TestWarn(t *testing.T) {
	if !Warn(8.9, 10, 10) {
		t.Fatal("should warn when frequency drops >10%")
	}
	if Warn(9.5, 10, 10) {
		t.Fatal("should not warn when drop <= 10%")
	}
}

// Package monitor supports long-term structural-health monitoring: control
// limits around the mean frequency, anomaly detection against those limits, and a
// baseline frequency-drop warning.
package monitor

import "math"

// Sample is one monitoring observation of a structure's dominant frequency.
type Sample struct {
	Time float64
	Freq float64
	RMS  float64
}

// ControlLimits returns mean(Freq) ± k*std(Freq). If the standard deviation is 0
// (all frequencies equal), both limits equal the mean. Empty input returns 0,0.
func ControlLimits(samples []Sample, k float64) (lo, hi float64) {
	if len(samples) == 0 {
		return 0, 0
	}
	var sum float64
	for _, s := range samples {
		sum += s.Freq
	}
	mean := sum / float64(len(samples))
	var varSum float64
	for _, s := range samples {
		d := s.Freq - mean
		varSum += d * d
	}
	std := math.Sqrt(varSum / float64(len(samples)))
	if std == 0 {
		return mean, mean
	}
	return mean - k*std, mean + k*std
}

// Anomalies returns the indices i where samples[i].Freq lies outside [lo,hi].
// Empty input returns nil (never panics).
func Anomalies(samples []Sample, lo, hi float64) []int {
	if len(samples) == 0 {
		return nil
	}
	var idx []int
	for i, s := range samples {
		if s.Freq < lo || s.Freq > hi {
			idx = append(idx, i)
		}
	}
	return idx
}

// Warn reports true when freq is below baseline*(1 - dropPct/100), i.e. the
// frequency has dropped beyond the given percentage threshold.
func Warn(freq, baseline, dropPct float64) bool {
	return freq < baseline*(1-dropPct/100)
}

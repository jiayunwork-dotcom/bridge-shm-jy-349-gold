// Package modal computes modal parameters from vibration Series: natural
// frequency, damping ratio (logarithmic decrement on a decaying segment), and a
// frequency-shift damage indicator.
package modal

import (
	"math"

	"bridge-shm/internal/signal"
)

// NaturalFreq returns the natural frequency of a Series, which is its dominant
// frequency.
func NaturalFreq(s signal.Series) float64 {
	return signal.DominantFreq(s)
}

// DampingRatio estimates the damping ratio from a decaying oscillation segment
// using the logarithmic decrement: zeta = (1/(2π)) * mean over i of ln(x[i]/x[i+1])
// for the monotonically positive-decaying pairs. Pairs that are non-positive or not
// decaying are skipped. The result is clamped to be >= 0. A series with fewer than
// 2 samples returns 0.
func DampingRatio(decay signal.Series) float64 {
	d := decay.Data
	if len(d) < 2 {
		return 0
	}
	var sum float64
	count := 0
	for i := 0; i+1 < len(d); i++ {
		x := d[i]
		y := d[i+1]
		if x <= 0 || y <= 0 {
			continue
		}
		if x <= y {
			continue
		}
		sum += math.Log(x / y)
		count++
	}
	if count == 0 {
		return 0
	}
	zeta := (1.0 / (2 * math.Pi)) * (sum / float64(count))
	if zeta < 0 {
		zeta = 0
	}
	return zeta
}

// FreqShift returns the relative frequency shift (f1-f0)/f0. It is negative when
// the frequency drops (a damage indicator). A zero baseline returns 0.
func FreqShift(f0, f1 float64) float64 {
	if f0 == 0 {
		return 0
	}
	return (f1 - f0) / f0
}

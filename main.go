// Command bridge-shm is a CLI for bridge structural-health monitoring vibration
// analysis. It reads long-format sensor CSV, prints per-series signal statistics
// (RMS/peak/dominant frequency), the damping ratio for any decay series, and an
// optional frequency-shift damage warning against a baseline.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bridge-shm/internal/modal"
	"bridge-shm/internal/monitor"
	"bridge-shm/internal/signal"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("bridge-shm", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	seriesPath := fs.String("series", "", "path to sensor series CSV (label,samplerate,value)")
	baseline := fs.Float64("baseline", 0, "baseline natural frequency (Hz) for damage warning")
	drop := fs.Float64("drop", 0, "drop percentage threshold for the warning")
	if err := fs.Parse(args); err != nil {
		usage(fs)
		return 2
	}
	if *seriesPath == "" {
		usage(fs)
		return 2
	}

	series, err := signal.ParseSeries(*seriesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	for _, s := range series {
		fmt.Printf("series=%s samplerate=%.2f rms=%.4f peak=%.4f dominant_freq=%.4f\n",
			s.Label, s.SampleRate, signal.RMS(s), signal.Peak(s), signal.DominantFreq(s))
	}

	for _, s := range series {
		if strings.Contains(strings.ToLower(s.Label), "decay") {
			fmt.Printf("decay_label=%s damping_ratio=%.4f\n", s.Label, modal.DampingRatio(s))
		}
	}

	if *baseline != 0 {
		f1 := modal.NaturalFreq(series[0])
		shift := modal.FreqShift(*baseline, f1)
		fmt.Printf("baseline=%.4f measured=%.4f freq_shift=%.4f\n", *baseline, f1, shift)
		if monitor.Warn(f1, *baseline, *drop) {
			fmt.Printf("WARNING: dominant frequency dropped beyond %.2f%%\n", *drop)
		}
	}

	return 0
}

func usage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "usage: bridge-shm -series <path> [-baseline <float> -drop <float>]")
	fs.PrintDefaults()
}

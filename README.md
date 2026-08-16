# bridge-shm

`bridge-shm` is a gold-standard (correct, bug-free) Go module for **bridge
structural-health monitoring (SHM) vibration analysis**. It imports long-format
sensor time-series (CSV grouped by label) and computes:

- **Signal statistics** — RMS, peak amplitude, discrete Fourier transform (DFT),
  and dominant frequency.
- **Modal parameters** — natural frequency (dominant frequency), damping ratio
  (logarithmic decrement on a decaying-oscillation segment), and a frequency-shift
  damage indicator.
- **Long-term monitoring** — control limits (mean ± k·σ), anomaly detection, and a
  frequency-drop warning against a baseline.

It is pure Go standard library only (no third-party imports, no network), so it is
easy to vendor, audit, and use as training-data reference material.

## CLI usage

Build and run the command-line tool:

```sh
go build -o bridge-shm .
./bridge-shm -series example/vibration.csv
```

Flags:

- `-series <path>` — (required) path to a long-format CSV with header
  `label,samplerate,value`. Rows are grouped by `label`.
- `-baseline <float>` — baseline natural frequency (Hz) used for the damage warning.
- `-drop <float>` — drop-percentage threshold; a warning is printed when the
  measured dominant frequency falls below `baseline * (1 - drop/100)`.

Exit codes:

- `2` — missing `-series` (usage printed to stderr).
- `1` — malformed / unreadable input (clear message to stderr).
- `0` — success; per-series RMS/peak/dominant frequency are printed, the damping
  ratio for any `decay` series, and (with `-baseline`) the frequency-shift warning.

## Example files

The `example/vibration.csv` file contains three deterministic series:

- `accel` — pure 2 Hz sine at 32 Hz sample rate (dominant frequency ≈ 2 Hz).
- `strain` — pure 5 Hz sine at 32 Hz sample rate.
- `decay` — exponentially decaying 3 Hz sine, used to demonstrate damping ratio.

## Library usage

```go
import "bridge-shm/internal/signal"

series, err := signal.ParseSeries("example/vibration.csv")
// series[0].Label, signal.RMS(series[0]), signal.DominantFreq(series[0]) ...
```

See the `internal/signal`, `internal/modal`, and `internal/monitor` packages for
the full exported API.

package report

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
)

type Comparison struct {
	BaselinePath  string
	CandidatePath string
	Baseline      CatbenchResult
	Candidate     CatbenchResult
}

type compareJSON struct {
	Baseline             string             `json:"baseline"`
	Candidate            string             `json:"candidate"`
	RPSChangePercent     float64            `json:"rps_change_percent"`
	LatencyChangePercent compareLatencyJSON `json:"latency_change_percent"`
	Errors               compareCountJSON   `json:"errors"`
}

type compareLatencyJSON struct {
	P50 float64 `json:"p50"`
	P95 float64 `json:"p95"`
	P99 float64 `json:"p99"`
	Max float64 `json:"max"`
}

type compareCountJSON struct {
	Baseline  int `json:"baseline"`
	Candidate int `json:"candidate"`
}

func CompareFiles(baselinePath, candidatePath string) (Comparison, error) {
	baseline, err := LoadResultFile(baselinePath)
	if err != nil {
		return Comparison{}, fmt.Errorf("load baseline %q: %w", baselinePath, err)
	}

	candidate, err := LoadResultFile(candidatePath)
	if err != nil {
		return Comparison{}, fmt.Errorf("load candidate %q: %w", candidatePath, err)
	}

	return Comparison{
		BaselinePath:  baselinePath,
		CandidatePath: candidatePath,
		Baseline:      baseline,
		Candidate:     candidate,
	}, nil
}

func PrintCompare(w io.Writer, comparison Comparison) {
	baseline := comparison.Baseline
	candidate := comparison.Candidate

	fmt.Fprintln(w, "Catbench Compare")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Target:")
	fmt.Fprintf(w, "baseline:  %s\n", baseline.Target)
	fmt.Fprintf(w, "candidate: %s\n\n", candidate.Target)

	fmt.Fprintln(w, "RPS:")
	fmt.Fprintf(w, "baseline:  %.2f\n", baseline.RPS)
	fmt.Fprintf(w, "candidate: %.2f\n", candidate.RPS)
	fmt.Fprintf(w, "change:    %s\n\n", formatPercentChange(baseline.RPS, candidate.RPS))

	fmt.Fprintln(w, "Latency:")
	fmt.Fprintf(w, "p50: baseline %s -> candidate %s (%s)\n", formatMilliseconds(baseline.Latency.P50MS), formatMilliseconds(candidate.Latency.P50MS), formatPercentChange(baseline.Latency.P50MS, candidate.Latency.P50MS))
	fmt.Fprintf(w, "p95: baseline %s -> candidate %s (%s)\n", formatMilliseconds(baseline.Latency.P95MS), formatMilliseconds(candidate.Latency.P95MS), formatPercentChange(baseline.Latency.P95MS, candidate.Latency.P95MS))
	fmt.Fprintf(w, "p99: baseline %s -> candidate %s (%s)\n", formatMilliseconds(baseline.Latency.P99MS), formatMilliseconds(candidate.Latency.P99MS), formatPercentChange(baseline.Latency.P99MS, candidate.Latency.P99MS))
	fmt.Fprintf(w, "max: baseline %s -> candidate %s (%s)\n\n", formatMilliseconds(baseline.Latency.MaxMS), formatMilliseconds(candidate.Latency.MaxMS), formatPercentChange(baseline.Latency.MaxMS, candidate.Latency.MaxMS))

	fmt.Fprintln(w, "Success:")
	fmt.Fprintf(w, "baseline:  %d\n", baseline.Success)
	fmt.Fprintf(w, "candidate: %d\n\n", candidate.Success)

	fmt.Fprintln(w, "Errors:")
	fmt.Fprintf(w, "baseline:  %d\n", baseline.Errors)
	fmt.Fprintf(w, "candidate: %d\n", candidate.Errors)
	fmt.Fprintf(w, "rate diff: %+0.2f percentage points\n", errorRate(candidate)-errorRate(baseline))
}

func WriteCompareJSON(w io.Writer, comparison Comparison) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(compareJSON{
		Baseline:         comparison.BaselinePath,
		Candidate:        comparison.CandidatePath,
		RPSChangePercent: roundPercent(finitePercentChange(comparison.Baseline.RPS, comparison.Candidate.RPS)),
		LatencyChangePercent: compareLatencyJSON{
			P50: roundPercent(finitePercentChange(comparison.Baseline.Latency.P50MS, comparison.Candidate.Latency.P50MS)),
			P95: roundPercent(finitePercentChange(comparison.Baseline.Latency.P95MS, comparison.Candidate.Latency.P95MS)),
			P99: roundPercent(finitePercentChange(comparison.Baseline.Latency.P99MS, comparison.Candidate.Latency.P99MS)),
			Max: roundPercent(finitePercentChange(comparison.Baseline.Latency.MaxMS, comparison.Candidate.Latency.MaxMS)),
		},
		Errors: compareCountJSON{
			Baseline:  comparison.Baseline.Errors,
			Candidate: comparison.Candidate.Errors,
		},
	})
}

func formatPercentChange(baseline, candidate float64) string {
	change, ok := percentChange(baseline, candidate)
	if !ok {
		return "n/a"
	}
	return fmt.Sprintf("%+.2f%%", change)
}

func finitePercentChange(baseline, candidate float64) float64 {
	change, ok := percentChange(baseline, candidate)
	if !ok {
		return 0
	}
	return change
}

func percentChange(baseline, candidate float64) (float64, bool) {
	if baseline == 0 {
		return 0, false
	}
	change := ((candidate - baseline) / baseline) * 100
	if math.IsInf(change, 0) || math.IsNaN(change) {
		return 0, false
	}
	return change, true
}

func formatMilliseconds(ms float64) string {
	return fmt.Sprintf("%.2fms", ms)
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}

func errorRate(result CatbenchResult) float64 {
	if result.Requests <= 0 {
		return 0
	}
	return (float64(result.Errors) / float64(result.Requests)) * 100
}

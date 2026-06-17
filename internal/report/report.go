package report

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/shanmeiliu/catbench/internal/runner"
)

func Print(w io.Writer, result runner.Result) {
	fmt.Fprintf(w, "Target: %s\n", result.URL)
	fmt.Fprintf(w, "Method: %s\n\n", result.Method)
	fmt.Fprintf(w, "Requests:     %d\n", result.Requests)
	fmt.Fprintf(w, "Concurrency:  %d\n", result.Concurrency)
	fmt.Fprintf(w, "Success:      %d\n", result.Success)
	fmt.Fprintf(w, "Errors:       %d\n", result.Errors)
	fmt.Fprintf(w, "Duration:     %s\n", formatDuration(result.Duration))
	fmt.Fprintf(w, "RPS:          %.2f\n\n", result.RPS)
	fmt.Fprintln(w, "Latency:")
	fmt.Fprintf(w, "p50:          %s\n", formatDuration(result.P50))
	fmt.Fprintf(w, "p95:          %s\n", formatDuration(result.P95))
	fmt.Fprintf(w, "p99:          %s\n", formatDuration(result.P99))
	fmt.Fprintf(w, "max:          %s\n\n", formatDuration(result.Max))
	fmt.Fprintln(w, "Status Codes:")

	codes := make([]int, 0, len(result.StatusCounts))
	for code := range result.StatusCounts {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	if len(codes) == 0 {
		fmt.Fprintln(w, "none:         0")
		return
	}

	for _, code := range codes {
		fmt.Fprintf(w, "%d:          %d\n", code, result.StatusCounts[code])
	}
}

func formatDuration(d time.Duration) string {
	return d.Round(time.Microsecond).String()
}

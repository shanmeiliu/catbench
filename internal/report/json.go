package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/shanmeiliu/catbench/internal/runner"
)

type CatbenchResult struct {
	Target      string         `json:"target"`
	Method      string         `json:"method"`
	Requests    int            `json:"requests"`
	Concurrency int            `json:"concurrency"`
	Success     int            `json:"success"`
	Errors      int            `json:"errors"`
	DurationMS  int64          `json:"duration_ms"`
	RPS         float64        `json:"rps"`
	Latency     JSONLatency    `json:"latency"`
	StatusCodes map[string]int `json:"status_codes"`
	Timestamp   string         `json:"timestamp"`
}

type JSONResult struct {
	Target      string         `json:"target"`
	Method      string         `json:"method"`
	Requests    int            `json:"requests"`
	Concurrency int            `json:"concurrency"`
	Success     int            `json:"success"`
	Errors      int            `json:"errors"`
	DurationMS  int64          `json:"duration_ms"`
	RPS         float64        `json:"rps"`
	Latency     JSONLatency    `json:"latency"`
	StatusCodes map[string]int `json:"status_codes"`
	Timestamp   string         `json:"timestamp"`
}

type JSONLatency struct {
	P50MS float64 `json:"p50_ms"`
	P95MS float64 `json:"p95_ms"`
	P99MS float64 `json:"p99_ms"`
	MaxMS float64 `json:"max_ms"`
}

func WriteJSON(w io.Writer, result runner.Result) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(toJSONResult(result))
}

func SaveJSON(path string, result runner.Result) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create result directory: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create result file: %w", err)
	}
	defer file.Close()

	if err := WriteJSON(file, result); err != nil {
		return fmt.Errorf("write result file: %w", err)
	}

	return nil
}

func LoadResultFile(path string) (CatbenchResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return CatbenchResult{}, err
	}
	defer file.Close()

	var result CatbenchResult
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&result); err != nil {
		return CatbenchResult{}, fmt.Errorf("invalid catbench JSON result: %w", err)
	}

	if err := validateCatbenchResult(result); err != nil {
		return CatbenchResult{}, err
	}

	return result, nil
}

func validateCatbenchResult(result CatbenchResult) error {
	if result.Target == "" {
		return fmt.Errorf("invalid catbench JSON result: missing target")
	}
	if result.Method == "" {
		return fmt.Errorf("invalid catbench JSON result: missing method")
	}
	if result.Requests <= 0 {
		return fmt.Errorf("invalid catbench JSON result: requests must be greater than 0")
	}
	if result.Concurrency <= 0 {
		return fmt.Errorf("invalid catbench JSON result: concurrency must be greater than 0")
	}
	if result.Success < 0 || result.Errors < 0 {
		return fmt.Errorf("invalid catbench JSON result: success and errors must be non-negative")
	}
	if result.Success+result.Errors != result.Requests {
		return fmt.Errorf("invalid catbench JSON result: success plus errors must equal requests")
	}
	if result.DurationMS < 0 {
		return fmt.Errorf("invalid catbench JSON result: duration_ms must be non-negative")
	}
	if result.RPS < 0 {
		return fmt.Errorf("invalid catbench JSON result: rps must be non-negative")
	}
	if result.Latency.P50MS < 0 || result.Latency.P95MS < 0 || result.Latency.P99MS < 0 || result.Latency.MaxMS < 0 {
		return fmt.Errorf("invalid catbench JSON result: latency values must be non-negative")
	}
	if result.StatusCodes == nil {
		return fmt.Errorf("invalid catbench JSON result: missing status_codes")
	}
	if result.Timestamp == "" {
		return fmt.Errorf("invalid catbench JSON result: missing timestamp")
	}

	return nil
}

func toJSONResult(result runner.Result) JSONResult {
	statusCodes := make(map[string]int, len(result.StatusCounts))
	for code, count := range result.StatusCounts {
		statusCodes[strconv.Itoa(code)] = count
	}

	return JSONResult{
		Target:      result.URL,
		Method:      result.Method,
		Requests:    result.Requests,
		Concurrency: result.Concurrency,
		Success:     result.Success,
		Errors:      result.Errors,
		DurationMS:  result.Duration.Milliseconds(),
		RPS:         result.RPS,
		Latency: JSONLatency{
			P50MS: durationMilliseconds(result.P50),
			P95MS: durationMilliseconds(result.P95),
			P99MS: durationMilliseconds(result.P99),
			MaxMS: durationMilliseconds(result.Max),
		},
		StatusCodes: statusCodes,
		Timestamp:   result.Timestamp.Format("2006-01-02T15:04:05Z"),
	}
}

func durationMilliseconds(d time.Duration) float64 {
	return float64(d) / float64(time.Millisecond)
}

package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Config struct {
	URL         string
	Method      string
	Requests    int
	Concurrency int
	Timeout     time.Duration
	Body        string
	Headers     map[string]string
}

type Result struct {
	URL          string
	Method       string
	Requests     int
	Concurrency  int
	Success      int
	Errors       int
	Duration     time.Duration
	RPS          float64
	P50          time.Duration
	P95          time.Duration
	P99          time.Duration
	Max          time.Duration
	StatusCounts map[int]int
}

type requestResult struct {
	latency time.Duration
	status  int
	success bool
}

func Run(cfg Config) (Result, error) {
	if err := validateConfig(&cfg); err != nil {
		return Result{}, err
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = cfg.Concurrency
	transport.MaxIdleConnsPerHost = cfg.Concurrency

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	jobs := make(chan struct{})
	results := make(chan requestResult, cfg.Requests)

	var wg sync.WaitGroup
	wg.Add(cfg.Concurrency)
	for i := 0; i < cfg.Concurrency; i++ {
		go func() {
			defer wg.Done()
			for range jobs {
				results <- executeRequest(client, cfg)
			}
		}()
	}

	start := time.Now()
	for i := 0; i < cfg.Requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(start)
	close(results)

	return summarize(cfg, results, duration), nil
}

func validateConfig(cfg *Config) error {
	cfg.Method = strings.ToUpper(strings.TrimSpace(cfg.Method))

	if cfg.URL == "" {
		return fmt.Errorf("--url is required")
	}
	if cfg.Method != http.MethodGet && cfg.Method != http.MethodPost {
		return fmt.Errorf("unsupported --method %q: only GET and POST are supported", cfg.Method)
	}
	if cfg.Requests <= 0 {
		return fmt.Errorf("--requests must be greater than 0")
	}
	if cfg.Concurrency <= 0 {
		return fmt.Errorf("--concurrency must be greater than 0")
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("--timeout must be greater than 0")
	}
	if cfg.Concurrency > cfg.Requests {
		cfg.Concurrency = cfg.Requests
	}

	return nil
}

func executeRequest(client *http.Client, cfg Config) requestResult {
	var body io.Reader
	if cfg.Method == http.MethodPost {
		body = bytes.NewBufferString(cfg.Body)
	}

	req, err := http.NewRequestWithContext(context.Background(), cfg.Method, cfg.URL, body)
	if err != nil {
		return requestResult{success: false}
	}

	for name, value := range cfg.Headers {
		req.Header.Set(name, value)
	}
	if cfg.Method == http.MethodPost && cfg.Body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return requestResult{latency: latency, success: false}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	success := resp.StatusCode >= 200 && resp.StatusCode < 400
	return requestResult{
		latency: latency,
		status:  resp.StatusCode,
		success: success,
	}
}

func summarize(cfg Config, results <-chan requestResult, duration time.Duration) Result {
	statusCounts := make(map[int]int)
	latencies := make([]time.Duration, 0, cfg.Requests)
	success := 0
	errors := 0

	for result := range results {
		latencies = append(latencies, result.latency)
		if result.status != 0 {
			statusCounts[result.status]++
		}
		if result.success {
			success++
		} else {
			errors++
		}
	}

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	rps := 0.0
	if duration > 0 {
		rps = float64(cfg.Requests) / duration.Seconds()
	}

	return Result{
		URL:          cfg.URL,
		Method:       cfg.Method,
		Requests:     cfg.Requests,
		Concurrency:  cfg.Concurrency,
		Success:      success,
		Errors:       errors,
		Duration:     duration,
		RPS:          rps,
		P50:          percentile(latencies, 50),
		P95:          percentile(latencies, 95),
		P99:          percentile(latencies, 99),
		Max:          maxLatency(latencies),
		StatusCounts: statusCounts,
	}
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	index := (p*len(sorted) + 99) / 100
	if index < 1 {
		index = 1
	}
	if index > len(sorted) {
		index = len(sorted)
	}

	return sorted[index-1]
}

func maxLatency(sorted []time.Duration) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	return sorted[len(sorted)-1]
}

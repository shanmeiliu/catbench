package runner

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRunRequestsMode(t *testing.T) {
	result, err := Run(Config{
		URL:         "http://example.test",
		Method:      http.MethodGet,
		Requests:    5,
		Concurrency: 2,
		Timeout:     time.Second,
		Client:      testClient(0, http.StatusOK),
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if result.Requests != 5 {
		t.Fatalf("Requests = %d, want 5", result.Requests)
	}
	if result.Success != 5 {
		t.Fatalf("Success = %d, want 5", result.Success)
	}
	if result.Errors != 0 {
		t.Fatalf("Errors = %d, want 0", result.Errors)
	}
	if result.StatusCounts[http.StatusOK] != 5 {
		t.Fatalf("StatusCounts[200] = %d, want 5", result.StatusCounts[http.StatusOK])
	}
}

func TestRunDurationMode(t *testing.T) {
	result, err := Run(Config{
		URL:           "http://example.test",
		Method:        http.MethodGet,
		DurationLimit: 30 * time.Millisecond,
		Concurrency:   2,
		Timeout:       time.Second,
		Client:        testClient(time.Millisecond, http.StatusOK),
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if result.Requests == 0 {
		t.Fatal("Requests = 0, want completed requests")
	}
	if result.Success != result.Requests {
		t.Fatalf("Success = %d, want %d", result.Success, result.Requests)
	}
	if result.Errors != 0 {
		t.Fatalf("Errors = %d, want 0", result.Errors)
	}
	if result.Duration < 30*time.Millisecond {
		t.Fatalf("Duration = %s, want at least 30ms", result.Duration)
	}
}

func TestRunRejectsRequestsAndDurationTogether(t *testing.T) {
	_, err := Run(Config{
		URL:           "http://example.test",
		Method:        http.MethodGet,
		Requests:      10,
		DurationLimit: time.Second,
		Concurrency:   1,
		Timeout:       time.Second,
	})
	if err == nil {
		t.Fatal("Run returned nil error, want validation error")
	}
	if !strings.Contains(err.Error(), "either --requests or --duration") {
		t.Fatalf("error = %q, want requests/duration validation", err)
	}
}

func TestRunRejectsMissingRequestsAndDuration(t *testing.T) {
	_, err := Run(Config{
		URL:         "http://example.test",
		Method:      http.MethodGet,
		Concurrency: 1,
		Timeout:     time.Second,
	})
	if err == nil {
		t.Fatal("Run returned nil error, want validation error")
	}
	if !strings.Contains(err.Error(), "exactly one of --requests or --duration") {
		t.Fatalf("error = %q, want missing mode validation", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func testClient(delay time.Duration, statusCode int) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if delay > 0 {
				time.Sleep(delay)
			}
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}),
	}
}

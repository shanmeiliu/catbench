package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shanmeiliu/catbench/internal/report"
	"github.com/shanmeiliu/catbench/internal/runner"
)

type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "run" {
		fmt.Fprintln(os.Stderr, "Usage: catbench run --url URL [--method GET|POST] [--requests N] [--concurrency N] [--timeout 10s] [--body JSON] [--header 'Name: Value'] [--output text|json] [--save path]")
		os.Exit(2)
	}

	runFlags := flag.NewFlagSet("run", flag.ExitOnError)
	var headers headerFlags
	var output string
	var savePath string
	cfg := runner.Config{}

	runFlags.StringVar(&cfg.URL, "url", "", "target URL")
	runFlags.StringVar(&cfg.Method, "method", "GET", "HTTP method: GET or POST")
	runFlags.IntVar(&cfg.Requests, "requests", 1000, "total requests to send")
	runFlags.IntVar(&cfg.Concurrency, "concurrency", 50, "number of concurrent workers")
	runFlags.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "request timeout")
	runFlags.StringVar(&cfg.Body, "body", "", "raw request body")
	runFlags.Var(&headers, "header", `request header in "Name: Value" format; can be repeated`)
	runFlags.StringVar(&output, "output", "text", "output format: text or json")
	runFlags.StringVar(&savePath, "save", "", "save benchmark result JSON to this path")

	if err := runFlags.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	parsedHeaders, err := parseHeaders(headers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	cfg.Headers = parsedHeaders

	output = strings.ToLower(strings.TrimSpace(output))
	if output != "text" && output != "json" {
		fmt.Fprintf(os.Stderr, "unsupported --output %q: expected text or json\n", output)
		os.Exit(2)
	}

	result, err := runner.Run(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if savePath != "" {
		if err := report.SaveJSON(savePath, result); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if output == "json" {
		if err := report.WriteJSON(os.Stdout, result); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	report.Print(os.Stdout, result)
}

func parseHeaders(values []string) (map[string]string, error) {
	headers := make(map[string]string, len(values))
	for _, value := range values {
		name, headerValue, found := strings.Cut(value, ":")
		if !found {
			return nil, fmt.Errorf("invalid --header %q: expected format \"Name: Value\"", value)
		}

		name = strings.TrimSpace(name)
		headerValue = strings.TrimSpace(headerValue)
		if name == "" {
			return nil, fmt.Errorf("invalid --header %q: header name cannot be empty", value)
		}

		headers[name] = headerValue
	}

	return headers, nil
}

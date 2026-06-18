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
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "run":
		run(os.Args[2:])
	case "compare":
		compare(os.Args[2:])
	default:
		printUsage()
		os.Exit(2)
	}
}

func run(args []string) {
	runFlags := flag.NewFlagSet("run", flag.ExitOnError)
	var headers headerFlags
	var output string
	var savePath string
	cfg := runner.Config{}

	runFlags.StringVar(&cfg.URL, "url", "", "target URL")
	runFlags.StringVar(&cfg.Method, "method", "GET", "HTTP method: GET or POST")
	runFlags.IntVar(&cfg.Requests, "requests", 0, "total requests to send")
	runFlags.DurationVar(&cfg.DurationLimit, "duration", 0, "benchmark duration")
	runFlags.IntVar(&cfg.Concurrency, "concurrency", 50, "number of concurrent workers")
	runFlags.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "request timeout")
	runFlags.StringVar(&cfg.Body, "body", "", "raw request body")
	runFlags.Var(&headers, "header", `request header in "Name: Value" format; can be repeated`)
	runFlags.StringVar(&output, "output", "text", "output format: text or json")
	runFlags.StringVar(&savePath, "save", "", "save benchmark result JSON to this path")

	if err := runFlags.Parse(args); err != nil {
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

func compare(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: catbench compare <baseline.json> <candidate.json> [--output text|json]")
		os.Exit(2)
	}

	baselinePath := args[0]
	candidatePath := args[1]
	compareFlags := flag.NewFlagSet("compare", flag.ExitOnError)
	var output string
	compareFlags.StringVar(&output, "output", "text", "output format: text or json")
	if err := compareFlags.Parse(args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if compareFlags.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "Usage: catbench compare <baseline.json> <candidate.json> [--output text|json]")
		os.Exit(2)
	}

	output = strings.ToLower(strings.TrimSpace(output))
	if output != "text" && output != "json" {
		fmt.Fprintf(os.Stderr, "unsupported --output %q: expected text or json\n", output)
		os.Exit(2)
	}

	comparison, err := report.CompareFiles(baselinePath, candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if output == "json" {
		if err := report.WriteCompareJSON(os.Stdout, comparison); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	report.PrintCompare(os.Stdout, comparison)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  catbench run --url URL (--requests N | --duration 30s) [--method GET|POST] [--concurrency N] [--timeout 10s] [--body JSON] [--header 'Name: Value'] [--output text|json] [--save path]")
	fmt.Fprintln(os.Stderr, "  catbench compare <baseline.json> <candidate.json> [--output text|json]")
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

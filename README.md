# Catbench

Catbench is a lightweight HTTP load testing CLI for benchmarking one HTTP endpoint.

Version 0.2 focuses on a single endpoint with a fixed number of requests, configurable concurrency, machine-readable JSON output, and saved benchmark result files. It uses Go's standard library HTTP client and does not require external dependencies.

## Build

```bash
go build -o catbench ./cmd/catbench
```

## Example Usage

```bash
./catbench run \
  --url http://localhost:8080/products \
  --requests 10000 \
  --concurrency 100
```

POST requests can include a raw JSON body and repeated headers:

```bash
./catbench run \
  --url http://localhost:8080/products \
  --method POST \
  --body '{"name":"Notebook"}' \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer token"
```

Print JSON instead of the text report:

```bash
./catbench run \
  --url http://localhost:8080/products \
  --output json
```

Save a benchmark result as a baseline file:

```bash
./catbench run \
  --url http://localhost:8080/products \
  --requests 5000 \
  --concurrency 100 \
  --save results/products-baseline.json
```

Compare a baseline against a candidate result:

```bash
./catbench compare \
  results/products-baseline.json \
  results/products-cache.json
```

Print comparison output as JSON:

```bash
./catbench compare \
  results/products-baseline.json \
  results/products-cache.json \
  --output json
```

## Flags

```text
--url string          required target URL
--method string       HTTP method, GET or POST (default GET)
--requests int        total requests to send (default 1000)
--concurrency int     number of concurrent workers (default 50)
--timeout duration    request timeout (default 10s)
--body string         optional raw JSON body for POST
--header string       request header in "Name: Value" format; can be repeated
--output string       output format, text or json (default text)
--save string         save benchmark result JSON to this path
```

Compare command:

```text
catbench compare <baseline.json> <candidate.json> [--output text|json]
```

## Example Output

```text
Target: http://localhost:8080/products
Method: GET

Requests:     10000
Concurrency:  100
Success:      9998
Errors:       2
Duration:     1.23s
RPS:          8123.45

Latency:
p50:          2ms
p95:          8ms
p99:          24ms
max:          110ms

Status Codes:
200:          9998
500:          2
```

## JSON Output

```json
{
  "target": "http://localhost:8080/products",
  "method": "GET",
  "requests": 10000,
  "concurrency": 100,
  "success": 10000,
  "errors": 0,
  "duration_ms": 1250,
  "rps": 8000.12,
  "latency": {
    "p50_ms": 2.1,
    "p95_ms": 8.3,
    "p99_ms": 25.6,
    "max_ms": 110.4
  },
  "status_codes": {
    "200": 10000
  },
  "timestamp": "2026-06-16T18:00:00Z"
}
```

## Baseline Workflow

Capture a baseline before changing the service:

```bash
./catbench run \
  --url http://localhost:8080/products \
  --requests 5000 \
  --concurrency 100 \
  --save results/products-baseline.json
```

Run the same benchmark after an optimization, such as adding a cache:

```bash
./catbench run \
  --url http://localhost:8080/products \
  --requests 5000 \
  --concurrency 100 \
  --save results/products-cache.json
```

Compare:

```bash
./catbench compare \
  results/products-baseline.json \
  results/products-cache.json
```

Use the comparison report to check:

```text
RPS change
p50 latency change
p95 latency change
p99 latency change
max latency change
success count
error count
error rate difference
```

## Compare Output

```text
Catbench Compare

Target:
baseline:  http://localhost:8080/products
candidate: http://localhost:8080/products

RPS:
baseline:  2109.86
candidate: 8240.12
change:    +290.55%

Latency:
p50: baseline 2.42ms -> candidate 0.80ms (-66.94%)
p95: baseline 22.49ms -> candidate 4.10ms (-81.77%)
p99: baseline 23.77ms -> candidate 8.50ms (-64.24%)
max: baseline 23.83ms -> candidate 18.00ms (-24.46%)

Success:
baseline:  5000
candidate: 5000

Errors:
baseline:  0
candidate: 0
rate diff: +0.00 percentage points
```

JSON compare output includes RPS change, latency change percentages, and error counts.

## Status

Catbench v0.2 supports load testing one endpoint at a time, saving benchmark results as JSON, and comparing saved benchmark results. Distributed workers, HTML reports, CSV reports, charts, and spike/ramp/soak modes are intentionally out of scope for this release.

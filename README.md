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

```text
Baseline:   results/products-baseline.json
After cache: results/products-cache.json

Metrics:
p95 latency
p99 latency
RPS
```

## Status

Catbench v0.2 supports load testing one endpoint at a time and saving benchmark results as JSON. Distributed workers, HTML reports, CSV reports, and spike/ramp/soak modes are intentionally out of scope for this release.

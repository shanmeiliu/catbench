# Catbench

Catbench is a lightweight HTTP load testing CLI for benchmarking one HTTP endpoint.

Version 0.1 focuses on a single endpoint with a fixed number of requests and configurable concurrency. It uses Go's standard library HTTP client and does not require external dependencies.

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

## Flags

```text
--url string          required target URL
--method string       HTTP method, GET or POST (default GET)
--requests int        total requests to send (default 1000)
--concurrency int     number of concurrent workers (default 50)
--timeout duration    request timeout (default 10s)
--body string         optional raw JSON body for POST
--header string       request header in "Name: Value" format; can be repeated
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

## Status

Catbench v0.1 supports load testing one endpoint at a time. Distributed workers, HTML reports, CSV reports, and spike/ramp/soak modes are intentionally out of scope for this release.

# Catbench

Catbench is a lightweight HTTP load testing CLI for benchmarking APIs, measuring latency, and simulating traffic spikes.

It is designed as a learning-focused performance engineering tool that can be used across different projects, including Rust, Go, Python, Node, and AI/RAG APIs.

## Goals

Catbench helps answer questions like:

- How many requests per second can this endpoint handle?
- What are the p50, p95, and p99 latencies?
- How does performance change as concurrency increases?
- What happens during sudden traffic spikes?
- Does caching improve throughput?
- Does the database connection pool become a bottleneck?

## Example Usage

```bash
catbench run \
  --url http://localhost:8080/products \
  --requests 10000 \
  --concurrency 100
```

Example output:

```text
Target: http://localhost:8080/products

Requests:     10000
Concurrency:  100
Success:      10000
Errors:       0

Throughput:
RPS:          18432.5

Latency:
p50:          2ms
p95:          8ms
p99:          24ms
max:          110ms
```

## Planned Features

### Phase 1: Basic Load Testing

* GET requests
* POST requests
* Custom headers
* JSON request body
* Total request mode
* Concurrency control
* Timeout support
* p50, p95, p99 latency
* Requests per second
* Success/error counts

### Phase 2: Test Modes

* Fixed request count
* Fixed duration
* Spike test
* Ramp-up test
* Soak test

### Phase 3: Reports

* Console summary
* JSON output
* CSV output
* HTML report

### Phase 4: Advanced Performance Testing

* Connection reuse tuning
* Rate limiting
* Request body templates
* Multiple endpoints
* Scenario-based testing
* Distributed workers

## Why Build This?

Catbench is both a practical tool and a performance engineering learning project.

The goal is not only to generate traffic, but to understand the full optimization loop:

```text
Measure
  ↓
Find bottleneck
  ↓
Optimize
  ↓
Measure again
```

Example use case with Charmaine Cat Studio:

```text
Baseline:
GET /products → PostgreSQL query every request

Optimization:
GET /products → in-memory cache

Compare:
Before cache: p95 = 25ms, RPS = 3,000
After cache:  p95 = 2ms,  RPS = 30,000
```

## Performance Concepts Explored

Catbench is intended to help study:

* HTTP connection reuse
* Concurrency
* Goroutines
* Latency percentiles
* Tail latency
* Throughput
* Database bottlenecks
* Connection pooling
* Caching effects
* Traffic spikes
* Saturation points
* Error rates under load

## Tech Stack

* Go
* Standard library HTTP client
* CLI flags
* No heavy dependencies initially

## Example Targets

Catbench can be used to benchmark:

```bash
catbench run --url http://localhost:8080/products
catbench run --url http://localhost:8000/api/chat
catbench run --url http://localhost:3000/api/search
catbench run --url http://localhost:8080/health
```

## Project Status

Early design phase.

Initial goal:

```text
Build a simple single-binary CLI that can benchmark one HTTP endpoint with configurable request count and concurrency.
```

## License

MIT



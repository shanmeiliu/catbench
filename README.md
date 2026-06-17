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
# Distributed Systems Learnings

A production-focused sandbox showcasing low-level infrastructure primitives, network routing architectures, and concurrent systems design implemented from scratch in Go.

Rather than focusing on business-logic applications (CRUD), this repository isolates and implements core architectural patterns found in modern reverse proxies, service meshes, observability pipelines, and distributed databases. Each sub-project is treated as a minimal, high-throughput component designed to explore networking protocols, thread safety, and failure-domain boundaries.

---

# 🛠️ Projects Overview

| Project                           | Core Concepts Covered                                                                         | Folder Link                                                    |
| :-------------------------------- | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------- |
| **TLS-Terminating Sidecar Proxy** | Network proxies, TLS termination, graceful shutdown, context propagation                      | [`/tls-sidecar-proxy`](#1-tls-terminating-sidecar-proxy)       |
| **Consistent Hash Ambassador**    | Algorithmic routing, data partitioning, virtual nodes, high-availability failover             | [`/consistent-hash-ambassador`](#2-consistent-hash-ambassador) |
| **Prometheus Metrics Adapter**    | Interface normalization, pull-based observability, in-memory metrics state, ticker scheduling | [`/metrics-adapter`](#3-prometheus-metrics-adapter)            |

---

# 1. TLS-Terminating Sidecar Proxy

Instead of forcing a backend application to handle HTTPS/TLS directly, this sidecar proxy process intercepts secure HTTPS traffic, terminates the TLS layer, and forwards plain HTTP traffic to the underlying backend service.

> 📂 **Source Code:** [`/tls-sidecar-proxy`](https://github.com/vats-24/distributed-systems-learnings/tree/main/tls-sidecar-proxy)

## Architecture

```text
Client (HTTPS)
        |
        v
TLS Reverse Proxy (:8443)
        |
        v
HTTP
        |
        v
Backend Service (:8080)
```

The backend remains entirely decoupled from and unaware of the TLS layer, mimicking production ingress and sidecar patterns (such as Envoy or NGINX) at a micro-scale.

## Core Features

- **Minimal HTTP Backend (`:8080`)**

  - Reads configuration (`PORT`) from environment variables
  - Exposes thread-safe `/stats` and `/health` endpoints
  - Implements graceful shutdown via `SIGINT`
  - Waits for in-flight requests using Go synchronization primitives

- **TLS Reverse Proxy (`:8443`)**

  - Accepts HTTPS traffic via self-signed certificates
  - Routes traffic using Go’s `httputil.ReverseProxy`
  - Logs backend failures with request context
  - Returns structured JSON errors instead of generic `502 Bad Gateway` responses

---

# 2. Consistent Hash Ambassador

An exploration of data partitioning, routing determinism, and high-availability proxy design. This system implements an **Ambassador Proxy** that leverages a custom consistent hashing ring to route requests deterministically based on shard keys, mimicking routing layers found in systems like Cassandra or Memcached.

> 📂 **Source Code:** [`/consistent-hash-ambassador`](https://github.com/vats-24/distributed-systems-learnings/tree/main/consistent-hash-ambassador)

## Architecture

```text
Client
   |
   v
Ambassador Proxy
   |
   v
Consistent Hash Ring
   |
   v
Deterministic Backend Node
```

## Why Consistent Hashing?

Naïve routing algorithms such as `hash(key) % serverCount` break catastrophically when nodes scale up or down, forcing large-scale reshuffling of data.

By mapping both keys and nodes onto a logical circular ring structure and traversing clockwise, this system ensures that adding or removing a node only impacts a fraction of the keys (`K/N`). This enables:

- Cache locality
- Session affinity
- Multi-tenant partitioning
- Reduced redistribution overhead

## Core Features

- **Deterministic Routing**

  - Inspects incoming HTTP requests for an `X-Shard-Key` header
  - Calculates ring placement deterministically

- **Virtual Nodes (Vnodes)**

  - Distributes virtual replicas across the ring
  - Prevents hotspotting
  - Improves load balancing across physical nodes

- **High-Performance Traversal**

  - Uses binary-search-based (`sort.Search`) ring lookups
  - Resolves backend ownership in `O(log N)` time

- **Fault Tolerance & Failover**

  - Uses `sync.RWMutex` for thread-safe mutations
  - Dynamically retries alternate nodes clockwise along the ring
  - Intercepts backend failures through reverse-proxy failover logic

---

# 3. Prometheus Metrics Adapter

An exploration of interface normalization and telemetry translation layers inspired by observability infrastructure such as Prometheus exporters and the OpenTelemetry Collector.

The adapter decouples the application from destination-specific formats by running a long-lived concurrent background worker that periodically scrapes raw JSON statistics from a foundation service, maps them onto in-memory metrics state, and serializes them on demand into the standard Prometheus exposition format.

> 📂 **Source Code:** [`/metrics-adapter`](https://github.com/vats-24/distributed-systems-learnings/tree/main/metrics-adapter)

## Architecture

```text
Foundation Service (JSON)
            |
            v
     [10s Ticker Worker]
            |
            v
Metrics Adapter (Translation Layer)
            |
            v
/metrics (Prometheus Format)
            |
            v
Prometheus Server (Pull Model)
```

## Design Intent: Why HTTP Pull Instead of gRPC Streaming?

While gRPC is common in modern distributed systems, Prometheus fundamentally operates using a pull-based observability model.

Designing this adapter as an HTTP-scraping exporter aligns with the architecture of standard telemetry collector sidecars, making pull-based polling the structurally correct choice over persistent RPC streams.

## Core Features

- **Interface Normalization**

  - Keeps the application agnostic of downstream monitoring backends
  - Consumes a single JSON telemetry contract
  - Enables compatibility with systems such as Prometheus or Datadog

- **In-Memory Time-Series Modeling**

  - Models metrics as live, thread-safe state objects
  - Distinguishes counters and gauges structurally
  - Supports multi-dimensional metric labels

- **Periodic Worker Scheduling**

  - Uses a background goroutine driven by `time.Ticker`
  - Executes isolated scrape-and-normalize cycles every 10 seconds

- **On-Demand Serialization**

  - Exposes a `/metrics` endpoint
  - Captures instantaneous snapshots of runtime state
  - Dynamically converts metrics into Prometheus exposition format

---

# 🧠 Key Engineering Takeaways

Building these infrastructure components demonstrates that distributed systems complexity comes from reasoning about state, interfaces, concurrency, and failure boundaries—not simply writing more code.

- **Interface Normalization & Interoperability**

  - Producer applications should emit stable contracts
  - Adapter layers should handle transport and protocol conversion
  - Example: translating JSON telemetry into Prometheus exposition format

- **Concurrency vs. Performance**

  - Prevent race conditions using mutexes
  - Use `RWMutex` to optimize read-heavy lookup paths
  - Coordinate background workers using goroutines and tickers

- **State Snapshotting**

  - Separate live in-memory state from serialization layers
  - Treat transport payloads as runtime snapshots of evolving systems state

- **Asynchronous Coordination**

  - Use Go channels and `context.Context` for:

    - Cancellation propagation
    - Timeout handling
    - Lifecycle management
    - Graceful shutdown orchestration

---

# 📌 Repository Focus

This repository is intentionally centered around infrastructure engineering concepts rather than product-facing application development.

The emphasis is on understanding:

- Network boundaries
- Routing determinism
- Failure recovery
- Concurrent runtime behavior
- Protocol adaptation
- Distributed systems architecture patterns

Each project is designed as a minimal but production-inspired systems component that prioritizes architectural clarity, concurrency safety, and operational realism.

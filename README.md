Here is the updated, unified `README.md` containing your new **Prometheus Metrics Adapter**.

I have integrated it seamlessly, updated the master table, and consolidated your new takeaways into the core engineering section at the bottom.

---

````markdown
# Distributed Systems Learnings

A production-focused sandbox showcasing low-level infrastructure primitives, network routing architectures, and concurrent systems design implemented from scratch in Go.

Rather than focusing on business-logic applications (CRUD), this repository isolates and implements core architectural patterns found in modern reverse proxies, service meshes, observability pipelines, and distributed databases. Each sub-project is treated as a minimal, high-throughput component designed to explore networking protocols, thread safety, and failure domain boundaries.

---

## 🛠️ Projects Overview

| Project                           | Core Concepts Covered                                                                         | Folder Link                                                    |
| :-------------------------------- | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------- |
| **TLS-Terminating Sidecar Proxy** | Network proxies, TLS termination, Graceful shutdown, Context propagation                      | [`/tls-sidecar-proxy`](#1-tls-terminating-sidecar-proxy)       |
| **Consistent Hash Ambassador**    | Algorithmic routing, Data partitioning, Virtual nodes, High-availability failover             | [`/consistent-hash-ambassador`](#2-consistent-hash-ambassador) |
| **Prometheus Metrics Adapter**    | Interface normalization, Pull-based observability, In-memory metrics state, Ticker scheduling | [`/prometheus-metrics-adapter`](#3-prometheus-metrics-adapter) |

---

## 1. TLS-Terminating Sidecar Proxy

Instead of forcing a core backend application to handle HTTPS/TLS directly, this separate sidecar proxy process intercepts secure HTTPS traffic, terminates the TLS layer, and forwards plain HTTP traffic to the underlying backend service.

> 📂 **Source Code:** [`/tls-sidecar-proxy`](https://github.com/vats-24/distributed-systems-learnings/tree/main/tls-sidecar-proxy)

### Architecture

```text
Client (HTTPS) ---> TLS Reverse Proxy (:8443) ---> HTTP ---> Backend Service (:8080)
```
````

The backend remains entirely decoupled from and unaware of the TLS layer, mimicking production ingress and sidecar patterns (like Envoy or NGINX) at a micro-scale.

### Core Features

- **Minimal HTTP Backend (`:8080`):** Reads configuration (`PORT`) from environment variables, exposes thread-safe `/stats` and `/health` tracking, and implements a strict **graceful shutdown** via `SIGINT` (waiting for in-flight requests using Go `sync` primitives).
- **TLS Reverse Proxy (`:8443`):** Accepts HTTPS via self-signed certificates, routes traffic using Go’s `httputil.ReverseProxy`, logs backend failures with context, and returns structured JSON errors instead of generic `502 Bad Gateway` responses.

---

## 2. Consistent Hash Ambassador

An exploration of data partitioning, routing determinism, and high-availability proxy design. This system implements an **Ambassador Proxy** that leverages a custom consistent hashing ring to route requests deterministically based on shard keys, mimicking routing layers found in systems like Cassandra or Memcached.

> 📂 **Source Code:** [`/consistent-hash-ambassador`](https://github.com/vats-24/distributed-systems-learnings/tree/main/consistent-hash-ambassador)

### Architecture

```text
Client ---> Ambassador Proxy ---> Consistent Hash Ring ---> Deterministic Backend Node

```

### Why Consistent Hashing?

Naïve routing algorithms like `hash(key) % serverCount` break catastrophically when nodes scale up or down, forcing a massive reshuffle of data. By mapping both keys and nodes onto a logical circular ring structure and traversing clockwise, this system ensures that adding or removing a node only impacts a fraction of the keys ($K/N$), enabling **cache locality**, **session affinity**, and **multi-tenant partitioning**.

### Core Features

- **Deterministic Routing:** Inspects incoming HTTP requests for an `X-Shard-Key` header to calculate ring placement.
- **Virtual Nodes (Vnodes):** Implements virtual node distribution across the ring to prevent "hot spotting" and ensure uniform data distribution across physical backends.
- **High-Performance Traversal:** Uses a binary-search-based (`sort.Search`) ring lookup to resolve backend ownership efficiently in $O(\log N)$ time.
- **Fault Tolerance & Failover:** Uses a fine-grained `sync.RWMutex` for thread-safe mutations and integrates dynamic reverse proxy routing that intercepts backend failures to automatically retry alternate nodes clockwise along the ring.

---

## 3. Prometheus Metrics Adapter

An exploration of **interface normalization** and telemetry translation layers inspired by observability infrastructure like Prometheus exporters and the OpenTelemetry Collector.

The adapter decouples the application from destination-specific formats by running a long-running concurrent background worker that periodically scrapes raw JSON statistics from a foundation service, maps them onto in-memory metrics state, and serializes them on-demand into the standard Prometheus exposition format.

> 📂 **Source Code:** [`/metrics-adapter`](https://github.com/vats-24/distributed-systems-learnings/tree/main/metrics-adapter)

### Architecture

```text
Foundation Service (JSON) ---> [Every 10s Ticker] ---> Metrics Adapter (Translation) ---> /metrics (Prometheus Format) ---> Prometheus Server (Pull)

```

### Design Intent: Why HTTP Pull over gRPC Streaming?

While gRPC feels modern for systems engineering, Prometheus fundamentally utilizes a **pull-based observability model**. Designing this adapter as an HTTP-scraping exporter matches the structural design principles of standard telemetry collector sidecars, making pull-based polling the architecturally correct choice over persistent RPC streams.

### Core Features

- **Interface Normalization:** Keeps the core application completely agnostic of downstream monitoring backends (e.g., Prometheus, Datadog) by consuming a single JSON telemetry contract.
- **In-Memory Time-Series Modeling:** Rather than simply manipulating strings, metrics are modeled as live, thread-safe state objects (Counters vs. Gauges with structural multi-dimensional Labels) continuously updated over time.
- **Periodic Worker Scheduling:** Implements a background goroutine driven by a Go `time.Ticker` executing an isolated scrape-and-normalize loop every 10 seconds.
- **On-Demand Serialization:** Serves a `/metrics` endpoint that captures an instantaneous runtime snapshot of the in-memory metric store and dynamically converts it into the exact line-oriented Prometheus text exposition format.

---

## 🧠 Key Engineering Takeaways

Building these infrastructure pieces highlights that systems complexity comes from reasoning about state, interfaces, and failure boundaries rather than raw lines of code:

- **Interface Normalization & Interoperability:** Designing systems where producer applications emit single contracts, leaving transport format conversion (like JSON to Prometheus text protocol) to decoupled adapter layers.
- **Concurrency vs. Performance:** Eliminating race conditions using mutexes, implementing read-write locks (`RWMutex`) to keep data-lookup paths non-blocking, and maintaining background goroutine workers driven by time tickers.
- **State Snapshotting:** Reasoning about live, in-memory states (like a metrics registry or a hash ring) versus their serialization/transport layers (like an HTTP response or a network payload).
- **Asynchronous Coordination:** Utilizing Go channels and `context` for clean timeout, lifecycle management, and cancellation propagation.

```

```

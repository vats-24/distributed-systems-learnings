# Distributed Systems Learnings

A repository dedicated to exploring low-level backend infrastructure, networking primitives, and distributed systems architecture using Go.

Rather than building standard CRUD applications, the projects here focus on deep, production-grade infrastructure concepts: concurrency synchronization, request lifecycle management, routing determinism, and fault tolerance.

---

## 🛠️ Projects Overview

| Project                           | Core Concepts Covered                                                             | Folder Link                                                    |
| :-------------------------------- | :-------------------------------------------------------------------------------- | :------------------------------------------------------------- |
| **TLS-Terminating Sidecar Proxy** | Network proxies, TLS termination, Graceful shutdown, Context propagation          | [`/tls-sidecar-proxy`](#1-tls-terminating-sidecar-proxy)       |
| **Consistent Hash Ambassador**    | Algorithmic routing, Data partitioning, Virtual nodes, High-availability failover | [`/consistent-hash-ambassador`](#2-consistent-hash-ambassador) |

---

## 1. TLS-Terminating Sidecar Proxy

Instead of forcing a core backend application to handle HTTPS/TLS directly, this separate sidecar proxy process intercepts secure HTTPS traffic, terminates the TLS layer, and forwards plain HTTP traffic to the underlying backend service.

> 📂 **Source Code:** [`/tls-sidecar-proxy`](https://github.com/vats-24/distributed-systems-learnings/tree/main/tls-sidecar-proxy)

### Architecture

```text
Client (HTTPS) ---> TLS Reverse Proxy (:8443) ---> HTTP ---> Backend Service (:8080)

```

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

## 🧠 Key Engineering Takeaways

Building these infrastructure pieces in less than 200 lines of code each highlighted that systems complexity comes from reasoning about state and failure boundaries rather than code volume:

- **Concurrency vs. Performance:** Eliminating race conditions using mutexes and utilizing read-write locks (`RWMutex`) to keep data-lookup paths non-blocking.
- **Asynchronous Coordination:** Utilizing Go channels and `context` for clean timeout and cancellation propagation.
- **Algorithmic Infrastructure:** Shifting perspective from simply _using_ pre-built load balancers to understanding low-level proxy mechanics, header manipulation, and network coordination.

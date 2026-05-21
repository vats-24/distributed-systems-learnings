#A Proper guide to learn distributed systems for first principles

## TLS-Terminating Sidecar Proxy

A minimal, high-performance infrastructure project built in Go to deeply understand concurrency, networking, and request lifecycle management.

Instead of forcing the backend to handle HTTPS/TLS directly, this separate sidecar proxy process intercepts HTTPS traffic, terminates the TLS layer, and forwards plain HTTP traffic to the underlying backend service.

> 📂 **Location:** This project lives inside the [`tls-sidecar-proxy`](https://github.com/vats-24/distributed-systems-learnings/tree/main/tls-sidecar-proxy) folder of the `distributed-systems-learnings` repository.

---

## Architecture

```text
Client (HTTPS) ---> TLS Reverse Proxy (:8443) ---> HTTP ---> Backend Service (:8080)

```

The backend service remains entirely decoupled from and unaware of the TLS layer, mimicking production ingress and sidecar patterns (like Envoy or NGINX) at a micro-scale.

---

## Core Components & Features

### 1. Minimal HTTP Backend Service (`:8080`)

- **Dynamic Configuration:** Reads dynamic configuration (`PORT`) from environment variables.
- **Observability:** Exposes `/health` (uptime JSON) and `/stats` (thread-safe request counting).
- **Production-Grade Lifecycle:** Implements **graceful shutdown** on `SIGINT`, waiting for in-flight requests to complete before exiting using Go `sync` primitives.

### 2. TLS Reverse Proxy Sidecar (`:8443`)

- **TLS Termination:** Accepts secure HTTPS traffic using a self-signed certificate.
- **Smart Forwarding:** Utilizes Go’s `httputil.ReverseProxy` to route traffic efficiently.
- **Resilient Error Handling:** Logs backend failures with explicit request path context and returns structured JSON errors instead of generic, unhelpful `502 Bad Gateway` responses.

---

## Key Engineering Takeaways

Infrastructure engineering is deceptively compact; less than 200 lines of code unlocked deep, practical experience with:

- **Concurrency & Synchronization:** Managing shared mutable state and eliminating race conditions using `sync.Mutex`.
- **Asynchronous Coordination:** Utilizing Go channels and `context` for clean timeout and cancellation propagation.
- **Network Proxies:** Shifting perspective from simply _using_ ingress tools to understanding low-level reverse proxy mechanics, connection forwarding, and header manipulation.

```

```

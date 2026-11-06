# Distributed Rate Limiter System Architecture

**Author:** Minh Kha Truong  
**Status:** Proposed / Design Phase  
**Tech Stack:** Go (Golang), Redis, Kubernetes

---

## 1. System Overview

This document outlines the architecture for a distributed rate limiter built from scratch. The system is designed to protect backend APIs from excessive traffic, prevent abuse, and ensure fair usage among clients. It utilizes a Gateway pattern to intercept requests, Redis for distributed state management, and Go's concurrency primitives for high-performance request handling.

## 2. High-Level Architecture

The architecture separates the rate-limiting logic from the core business logic, ensuring that the Backend API only processes legitimate, permitted requests.

### 2.1. Components

- **Kubernetes (K8s) Ingress / Load Balancer:** Acts as the entry point, distributing incoming client traffic across multiple Gateway instances.
- **Gateway (Go):** The core rate-limiting enforcement point. It intercepts GET requests, checks limits against Redis, and either rejects the request (HTTP 429) or forwards it to the Backend API.
- **Redis:** The centralized, in-memory data store acting as the source of truth for rate-limit counters. It ensures all Gateway instances share a synchronized view of client usage.
- **Backend API:** The protected resource. Upon receiving a validated request from the Gateway, it performs the necessary computation (e.g., generating a random hash/token) and returns the payload.

### 2.2. Request Flow

1.  **Client Request:** A client sends an HTTP GET request.
2.  **Load Balancing:** Kubernetes routes the request to one of the available Gateway instances.
3.  **State Check (Gateway -> Redis):** The Gateway extracts the client identifier (e.g., IP address, API key) and executes an atomic check against Redis to verify if the client has exceeded their limit.
4.  **Decision:**
    - _Allow:_ If the limit is not exceeded, the Gateway increments the counter in Redis and forwards the request to the Backend API.
    - _Deny:_ If the limit is exceeded, the Gateway immediately terminates the request and returns an `HTTP 429 Too Many Requests` response.
5.  **Backend Processing:** The Backend API processes the allowed request, generates the required token/hash, and sends the response back through the Gateway to the client.

## 3. Concurrency Model

The Gateway is built in Go, leveraging its lightweight concurrency model to handle high throughput.

### Worker Pool vs. Native Goroutines

While a traditional worker pool pattern can be used to strictly bound resources, the primary implementation leverages Go's native `net/http` server behavior, which spawns a new goroutine for every incoming request.

Because the Gateway's primary tasks are I/O bound (network calls to Redis and the Backend API), native goroutines provide excellent performance without the overhead of managing a custom worker pool. If downstream connection limits become necessary, bounded channels acting as semaphores will be introduced to throttle concurrent outbound connections.

## 4. Redis Integration & Atomicity

To prevent race conditions in a highly concurrent, distributed environment, the Gateway interacts with Redis using **Lua scripts**.

### The Concurrency Problem

A naive approach using standard Redis commands (`GET` followed by `INCR`) is susceptible to race conditions. Multiple concurrent requests might read the same counter value before any request increments it, allowing traffic bursts to bypass the limit.

### The Lua Script Solution

By embedding the rate-limiting logic (e.g., Token Bucket or Fixed Window algorithm) within a Lua script, Redis executes the evaluation and incrementation as a single, isolated, atomic operation. This guarantees absolute accuracy of the rate limits regardless of the number of concurrent Gateway instances.

## 5. Resilience and Failure Handling

As the Gateway sits in the critical path, system resilience is paramount.

- **Context Propagation:** Go's `context` package is utilized across all layers. If a client disconnects prematurely, the cancellation signal propagates, immediately halting Redis checks or Backend API calls to conserve resources.
- **Circuit Breaking & Failover:** In the event of a Redis outage, the Gateway must decide how to handle traffic. The system incorporates a localized circuit breaker with a configurable strategy:
  - _Fail-Closed (Default for strict APIs):_ Blocks all requests to protect the backend, returning `HTTP 503 Service Unavailable`.
  - _Fail-Open (Default for high-availability systems):_ Temporarily bypasses the rate limiter, forwarding traffic directly to the backend while logging the state degradation.

## 6. Future Enhancements

- **Algorithm Pluggability:** Abstracting the rate-limit logic to easily swap between Fixed Window, Sliding Window Log, and Token Bucket algorithms based on endpoint requirements.
- **gRPC Integration:** Expanding the Gateway to handle gRPC streams alongside standard HTTP/REST endpoints.

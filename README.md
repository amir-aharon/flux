# Flux

**High-performance infrastructure primitives and system design patterns for Go.**

Flux is a collection of production-ready, thread-safe components that form the building blocks of distributed systems. Rather than building end-user applications, Flux focuses on the "plumbing" of the cloud, implementing the core algorithms and architectures used to manage concurrency, traffic, and data reliability.

Each module in Flux is a self-contained study in system design, built with zero external dependencies and a focus on idiomatic Go.

## Modules

| Module     | Name            | Pattern              | Status   | Description                                                                           |
| ---------- | --------------- | -------------------- | -------- | ------------------------------------------------------------------------------------- |
| **Tempo**  | Rate Limiter    | `Token Bucket`       | ✅ Ready | A middleware for throttling API traffic and preventing DoS attacks.                   |
| **Swarm**  | Worker Pool     | `Fan-Out/Fan-In`     | ✅ Ready | A concurrency engine for processing massive job queues with fixed resources.          |
| **Echo**   | Pub/Sub Broker  | `Observer`           | ✅ Ready | A generic, thread-safe message broadcasting system for decoupled communication.       |
| **Fuse**   | Circuit Breaker | `State Machine`      | ✅ Ready | A resilience wrapper that prevents cascading failures by failing fast during outages. |
| **Zen**    | Load Balancer   | `Reverse Proxy`      | ✅ Ready | A Layer 7 HTTP traffic distributor using Round-Robin strategy.                        |
| **Stash**  | LRU Cache       | `Doubly Linked List` | ✅ Ready | An O(1) generic in-memory storage engine with fixed capacity and automatic eviction.  |
| **Mirage** | CDN Edge Node   | `Caching Proxy`      | ✅ Ready | An intelligent HTTP accelerator that caches reads and actively invalidates on writes. |

## Philosophy

- **Zero Magic:** No heavy frameworks. Just the standard library and raw logic.
- **Concurrency First:** Heavy use of Channels, Mutexes, and Atomics to handle scale.
- **Resilience:** Designed to handle failures gracefully (Context cancellation, graceful shutdowns, timeout management).

## Getting Started

### Prerequisites

- Go 1.22+

### Installation

```bash
git clone https://github.com/amir-aharon/flux.git
cd flux

```

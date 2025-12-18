# Zen

**A lightweight, thread-safe Layer 7 Load Balancer for Go.**

Zen acts as a high-performance reverse proxy that distributes incoming HTTP traffic across multiple backend servers. It is designed to be extensible, efficient, and compatible with the standard Go `net/http` ecosystem.

## Features

- **Layer 7 Load Balancing:** operates at the application layer, fully parsing HTTP requests and forwarding them intelligently.
- **Round-Robin Strategy:** Implements a thread-safe, overflow-protected rotation algorithm to distribute load evenly.
- **High Performance:** Pre-initializes reverse proxies and URL parsing during startup to ensure zero allocation overhead during request handling.
- **Extensible Design:** Built on the Strategy Pattern, allowing easy implementation of custom balancing logic (e.g., Least Connections, Weighted).
- **Standard Compatibility:** Implements `http.Handler`, allowing it to serve as a standalone gateway or mount within a larger router.

## 🚀 Usage

### Import

```go
import "github.com/amir-aharon/flux/zen"

```

### Example

```go
package main

import (
    "log"
    "net/http"
    "github.com/amir-aharon/flux/zen"
)

func main() {
    // 1. Define your backend targets
    backends := []string{
        "http://localhost:8081",
        "http://localhost:8082",
        "http://localhost:8083",
    }

    // 2. Initialize the Load Balancer
    // Zen pre-compiles URLs and initializes proxies here for performance
    lb, err := zen.NewLoadBalancer(backends)
    if err != nil {
        log.Fatalf("Failed to init load balancer: %v", err)
    }

    // 3. Start the Gateway
    log.Println("Load Balancer listening on :8080")
    if err := http.ListenAndServe(":8080", lb); err != nil {
        log.Fatal(err)
    }
}

```

## System Design Notes

- **Reverse Proxying:** Zen leverages Go's robust `net/http/httputil` to handle the complexities of HTTP forwarding (header copying, streaming bodies, etc.).
- **Optimization (Hot Path):** Unlike naive implementations that parse URLs or create proxies on every request, Zen uses a "Constructor Initialization" pattern. It builds a map of `*url.URL` to `*httputil.ReverseProxy` once at startup, reducing the per-request cost to a simple pointer lookup.
- **Concurrency:** The Round-Robin strategy uses a `sync.Mutex` to protect the internal counter, ensuring that concurrent requests never result in race conditions or skipped backends.
- **Architecture:** The logic is split into two distinct components: the **Strategy** (Brain) which decides _where_ to go, and the **Balancer** (Body) which executes the forwarding. This decoupling makes the system highly testable and modular.

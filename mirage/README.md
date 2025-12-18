# Mirage

**An intelligent HTTP Caching Reverse Proxy (CDN Edge Node).**

Mirage is the integration layer of the Flux ecosystem. It sits between your users and your backend services, absorbing read-heavy traffic by caching responses in memory. It combines the Round-Robin Load Balancing of **Zen** with the O(1) LRU storage of **Stash**.

## Features

- **Smart Caching:** Automatically intercepts and caches `GET` requests while letting `POST`, `PATCH`, and other methods pass through.
- **Active Invalidation:** Implements a "Write-Purges-Read" strategy. If a `DELETE` or `PUT` request is detected for a specific resource, the cached `GET` entry is immediately destroyed to prevent stale data.
- **Response Interception:** Uses a "Spy" (Response Recorder) pattern to capture the backend's response stream on the fly without blocking the user.
- **Configurable Capacity:** Built on top of Stash, allowing strict memory limits (e.g., "Max 10,000 items").

## 🚀 Usage

### Import

```go
import "github.com/amir-aharon/flux/mirage"
```

### Example

```go
package main

import (
    "log"
    "net/http"
    "github.com/amir-aharon/flux/mirage"
)

func main() {
    // 1. Define your upstream backends
    backends := []string{
        "http://localhost:8081",
        "http://localhost:8082",
    }

    // 2. Initialize the CDN
    // - Backends: Where to forward traffic on a "Miss"
    // - Capacity: 1000 items in the cache
    cdn, err := mirage.NewCDN(backends, 1000)
    if err != nil {
        log.Fatal(err)
    }

    // 3. Start the Proxy
    log.Println("🚀 Mirage CDN running on :8080")
    http.ListenAndServe(":8080", cdn)
}
```

## System Design Notes

### The "Spy" Architecture

Mirage cannot read the standard `http.ResponseWriter` because it is a write-only stream. To solve this, Mirage uses a **Decorator Pattern** (`ResponseRecorder`):

1. **Cache Miss:** Mirage creates a "Spy" writer and passes it to the Load Balancer.
2. **Interception:** As the backend writes data to the user, the Spy copies those bytes into a local buffer.
3. **Promotion:** If the request succeeds (200 OK), the buffer is frozen and promoted to the **Stash** cache.

### Invalidation Strategy

Mirage prioritizes consistency over raw hit rate.

- **Rule:** `DELETE /users/1`
- **Action:** Mirage calculates the key for `GET:/users/1` and forcefully removes it from Stash.
- **Result:** The next read is guaranteed to hit the backend, ensuring the user never sees the "ghost" of a deleted item.

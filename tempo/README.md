# Tempo

**A thread-safe, token-bucket rate limiter middleware for Go HTTP services.**

Tempo provides a robust mechanism to throttle incoming API requests based on IP address. It uses a "Lazy Refill" algorithm to maintain high performance without needing expensive background timers for every user.

## Features

- **Token Bucket Algorithm:** Allows for bursts of traffic while enforcing a long-term rate limit.
- **Thread-Safe:** Uses fine-grained locking to ensure safe concurrent access.
- **Memory Efficient:** Implements "Lazy Initialization" (creates buckets only when needed) and a background Garbage Collector to remove stale users.
- **Zero Dependencies:** Built using only the Go standard library.

## 🚀 Usage

### Import

```go
import "https://github.com/amir-aharon/flux/tempo"
```

### Example

```go
package main

import (
    "net/http"
    "time"
    "https://github.com/amir-aharon/flux/tempo"
)

func main() {
    // 1. Create a Limiter: 5 requests/second, capacity of 10 requests "at rest"
    limit := tempo.NewRateLimiter(5.0, 10.0)

    // 2. Start Background Cleanup (Optional but recommended)
    // Check every minute, delete users inactive for 5 minutes
    limit.StartBackgroundCleanup(1*time.Minute, 5*time.Minute)

    // 3. Wrap your handler
    http.HandleFunc("/api", tempo.RateLimiterMiddleware(limit, myHandler))

    http.ListenAndServe(":8080", nil)

}

func myHandler(w http.ResponseWriter, r \*http.Request) {
    w.Write([]byte("Success!"))
}
```

## System Design Notes

- Concurrency: Tempo uses a two-level locking strategy. A Manager mutex protects the map of users, while individual Bucket mutexes protect the math calculations. This prevents a "stop-the-world" pause during high traffic.

- Garbage Collection: The cleanup process uses a "Snapshot" technique to iterate over keys safely without blocking incoming requests.

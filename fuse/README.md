# Fuse

**A generic, concurrent Circuit Breaker implementation.**

Fuse protects your application from cascading failures by temporarily stopping the execution of operations that are likely to fail. It wraps your service calls (like database queries or API requests) and enforces a "fail-fast" mechanism when the downstream service is unhealthy, allowing it time to recover.

## Features

- **State Machine:** Implements a robust 3-state logic (Closed, Open, Half-Open) to manage service health.
- **Atomic Probing:** Prevents "Thundering Herd" scenarios during recovery by allowing only one request to test the connection (The Probe).
- **Custom Error Classification:** Allows you to define which errors count as "crashes" and which are safe logic errors.
- **Thread-Safe:** Fully concurrent design using `sync.Mutex` to handle high-throughput environments safely.
- **Zero-Overhead:** Uses lazy evaluation for timeouts, requiring no background goroutines or tickers.

## 🚀 Usage

### Import

```go
import "github.com/amir-aharon/flux/fuse"
```

### Example

```go
package main

import (
    "errors"
    "fmt"
    "time"
    "github.com/amir-aharon/flux/fuse"
)

func main() {
    // 1. Define Config: Trip after 3 failures, wait 5 seconds to retry
    cfg := fuse.Config{
        Threshold: 3,
        Interval:  5 * time.Second,
        // Optional: Ignore "not found" errors, only count "db dead" as failures
        IsSuccessful: func(err error) bool {
            return err.Error() == "user not found"
        },
    }

    // 2. Initialize Breaker
    cb := fuse.NewBreaker(cfg)

    // 3. Execute a risky operation
    err := cb.Execute(func() error {
        // Simulate a database call
        return errors.New("db dead")
    })

    if err != nil {
        fmt.Printf("Result: %v\n", err)
    }
}
```

## System Design Notes

- **Pattern:** Implements the **Circuit Breaker** pattern to provide stability and prevent system overload.
- **State Transitions:**
- **Closed:** Normal operation. Failures are counted.
- **Open:** Fail fast. All requests are blocked immediately without executing the job.
- **Half-Open:** Recovery mode. A single "Probe" request is allowed through to test the upstream service.

- **Probe Mechanism:** To prevent a "Thundering Herd" when the timeout expires, Fuse atomically elects a single request to become the Probe. All concurrent requests continue to fail fast until the Probe returns successfully.
- **Concurrency:**
- Uses the **Lock/Unlock Pattern** strictly for state checks and updates.
- The actual job execution happens **outside the lock**, ensuring that a slow database call does not block other goroutines from checking the breaker status.

- **Failure Logic:** Uses "Consecutive Failure" counting. A single successful request in the Closed state resets the failure count to zero, preventing the breaker from tripping due to occasional, non-persistent glitches.

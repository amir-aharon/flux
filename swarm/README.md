# Swarm

**A generic, concurrent worker pool implementation using the Fan-Out/Fan-In pattern.**

Swarm allows you to process a massive stream of tasks using a fixed number of goroutines. It is designed to bound resource usage (CPU/Memory) when dealing with heavy workloads, ensuring your system remains stable under load.

## Features

- **Generics:** Type-safe processing of any input `TaskT` to any output `ResT`.
- **Bounded Concurrency:** Strict control over the number of active goroutines.
- **Non-Blocking:** Submission and processing are decoupled via buffered channels.
- **Graceful Shutdown:** Ensures all pending tasks are completed before closing the results channel.

## 🚀 Usage

### Import

```go
import "https://github.com/amir-aharon/flux/swarm"
```

### Example

```go
package main

import (
    "fmt"
    "https://github.com/amir-aharon/flux/swarm"
)

// 1. Define your processing function
func double(n int) (int, error) {
    return n * 2, nil
}

func main() {
    // 2. Initialize Pool: 3 workers, queue size of 10
    pool := swarm.NewPool(3, 10, double)

    // 3. Start the engine
    pool.Run()

    // 4. Submit tasks
    for i := 0; i < 10; i++ {
        pool.Submit(i, i) // ID: i, Data: i
    }

    // 5. Signal shutdown (stop accepting new work)
    pool.Shutdown()

    // 6. Collect results (blocks until all work is done)
    for res := range pool.Results() {
        fmt.Printf("Task %d: %d (Err: %v)\n", res.TaskID, res.Value, res.Err)
    }
}
```

## System Design Notes

- Pattern: Uses the Fan-Out pattern to distribute work and Fan-In to collect results.

- Deadlock Prevention: The Results() method returns a channel that is only closed once all workers have finished and the `sync.WaitGroup` counter reaches zero.

- Resource Safety: By fixing the worker count, Swarm prevents memory exhaustion attacks where an unlimited number of goroutines might be spawned.

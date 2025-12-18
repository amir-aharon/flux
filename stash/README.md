# Stash

**A high-performance, thread-safe LRU Cache with O(1) complexity.**

Stash is a generic in-memory storage engine designed for scenarios where data access patterns follow a "locality of reference." It enforces a fixed memory footprint by automatically evicting the Least Recently Used (LRU) items when capacity is reached, ensuring your application never runs out of memory due to unbounded caching.

## Features

- **Constant Time O(1):** Achieves instant lookups, insertions, and deletions regardless of cache size.
- **Generics:** Type-safe implementation allowing any `comparable` key and any value type (no `interface{}` casting).
- **Thread-Safe:** Protected by fine-grained locking for safe concurrent access in high-throughput environments.
- **Fixed Capacity:** Strictly enforces a maximum item count to predict and control memory usage.

## 🚀 Usage

### Import

```go
import "github.com/amir-aharon/flux/stash"
```

### Example

```go
package main

import (
    "fmt"
    "github.com/amir-aharon/flux/stash"
)

type User struct {
    Name string
    Role string
}

func main() {
    // 1. Initialize a cache for Users with a capacity of 100
    cache := stash.NewCache[string, User](100)

    // 2. Add items (O(1))
    cache.Put("u_123", User{Name: "Alice", Role: "Admin"})

    // 3. Retrieve items (O(1))
    // This automatically promotes "u_123" to the Most Recently Used position
    if user, ok := cache.Get("u_123"); ok {
        fmt.Printf("Found user: %s\n", user.Name)
    }

    // 4. Automatic Eviction
    // If the cache exceeds 100 items, the user who hasn't been
    // accessed for the longest time is dropped automatically.
}
```

## System Design Notes

- **Dual Data Structures:** Stash combines a Go `map` with a custom `Doubly Linked List`.
- **The Map** provides O(1) access to specific nodes.
- **The Linked List** maintains the temporal order (Age) of items.

- **Pointer Arithmetic:** Unlike standard slice-based caches, Stash performs 0 iterations. Moving an item to the front involves changing exactly 4 pointers, making it extremely CPU efficient even under heavy load.
- **Collision Handling:** Relies on Go's native map hashing for key distribution, inheriting its speed and collision resistance.

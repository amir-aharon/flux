# Echo

**A generic, concurrent Pub/Sub (Publish/Subscribe) broker.**

Echo enables real-time, one-to-many message broadcasting within your Go application. It is designed to decouple producers from consumers, allowing independent scaling and strictly typed communication between different components of your system.

## Features

- **Generics:** Type-safe broadcasting of any message type.
- **Topic-Based Routing:** Messages are routed to specific topics, allowing granular subscriptions.
- **Fan-Out:** A single published message is efficiently delivered to all active subscribers.
- **Backpressure Support:** Configurable buffer sizes to manage slow consumers.
- **Thread-Safe:** Fully concurrent design using `RWMutex` to handle high-throughput environments.
- **Resource Management:** Automatic cleanup of subscriptions and empty topics to prevent memory leaks.

## 🚀 Usage

### Import

```go
import "github.com/amir-aharon/flux/echo"
```

### Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/amir-aharon/flux/echo"
)

func main() {
    // 1. Initialize Broker: Buffer size of 10 to handle bursts
    broker := echo.NewBroker[string](10)

    // 2. Subscribe to a topic
    sub := broker.Subscribe("alerts")
    defer sub.Unsubscribe() // Ensure cleanup

    // 3. Start a listener (Consumer)
    go func() {
        for msg := range sub.C {
            fmt.Printf("Received: %s\n", msg)
        }
    }()

    // 4. Publish messages (Producer)
    broker.Publish("alerts", "System is starting...")
    broker.Publish("alerts", "CPU usage at 90%")

    // Give time for the message to propagate in this example
    time.Sleep(100 * time.Millisecond)
}

```

## System Design Notes

- **Pattern:** Implements the **Fan-Out** pattern where one input (Publish) triggers writes to N output channels (Subscribers).
- **Topology:** Uses a "Map of Maps" (`map[Topic]map[*Subscription]chan`) to allow O(1) addition and removal of subscribers.
- **Coupling:** Employs the "Handle Pattern" where `Subscribe` returns a struct linked to the Broker, enabling an intuitive `sub.Unsubscribe()` API.
- **Concurrency:**
  - `Publish` uses a Read Lock (`RLock`) to allow parallel publishing.
  - `Subscribe`/`Unsubscribe` use a Write Lock (`Lock`) to ensure map integrity.
- **Blocking Policy:** This implementation enforces a strict delivery guarantee. If a subscriber's buffer is full, the Publisher will block until space is available, preventing data loss at the cost of potential latency for the producer.

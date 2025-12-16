package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/amir-aharon/flux/echo"
)

func main() {
	// Buffer 0: Publisher blocks until ALL subscribers receive the message
	broker := echo.NewBroker[string](0)
	var wg sync.WaitGroup

	// Subscriber 1: Fast
	sub1 := broker.Subscribe("updates")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			msg := <-sub1.C
			fmt.Printf("[Sub 1] Received: %s\n", msg)
		}
	}()

	// Subscriber 2: Slow (Simulates 1s latency)
	sub2 := broker.Subscribe("updates")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			msg := <-sub2.C
			time.Sleep(1 * time.Second) // Blocks the publisher
			fmt.Printf("[Sub 2] Received: %s\n", msg)
		}
	}()

	// Give subscribers time to register
	time.Sleep(100 * time.Millisecond)

	// Publisher
	for i := 1; i <= 3; i++ {
		msg := fmt.Sprintf("payload-%d", i)
		fmt.Printf("Publishing %s...\n", msg)

		broker.Publish("updates", msg) // Blocks here for 1s

		fmt.Println("Done publishing.")
	}

	wg.Wait()
}

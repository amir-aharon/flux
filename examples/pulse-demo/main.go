package main

import (
	"fmt"
	"time"

	"github.com/amir-aharon/flux/pulse"
)

func main() {
	mon := pulse.NewMonitor(2 * time.Second)
	url := "www.google.com/search"
	mon.AddURL(url)
	fmt.Printf("monitor: %v\n", mon)
	for stat := range mon.GetStatus() {
		fmt.Printf("record: %v\n", stat)
	}
}

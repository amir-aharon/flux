package main

import (
	"fmt"
	"time"

	"github.com/amir-aharon/flux/swarm"
)

func slowSquare(input int) (int, error) {
	time.Sleep(1 * time.Second)
	retval := input * input

	return retval, nil
}

func main() {
	workers := swarm.NewPool(5, 10, slowSquare)
	workers.Run()
	for i := 1; i <= 10; i++ {
		workers.Submit(i, i)
	}
	workers.Shutdown()

	for res := range workers.Results() {
		fmt.Println(res)
	}
}

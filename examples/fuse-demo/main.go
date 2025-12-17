package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/amir-aharon/flux/fuse"
)

func job() error {
	return errors.New("db dead")
}

func main() {
	cfg := fuse.Config{
		Threshold: 4,
		Interval:  time.Second,
	}
	breaker := fuse.NewBreaker(cfg)

	for i := range 15 {
		err := breaker.Execute(job)
		fmt.Printf("%d\t%s\t%v\n", i, breaker.State.String(), err)
		time.Sleep(500 * time.Millisecond)
	}
}

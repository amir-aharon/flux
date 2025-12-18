package main

import (
	"fmt"

	"github.com/amir-aharon/flux/stash"
)

func main() {
	cache := stash.NewCache[string, int](3)
	cache.Put("A", 1)
	cache.Put("B", 2)
	cache.Put("C", 3)
	cache.Get("A")
	cache.Put("D", 4)
	v, exists := cache.Get("B")
	fmt.Println(v, exists)
}

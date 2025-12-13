package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/amir-aharon/flux/tempo"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	rl := tempo.NewRateLimiter(1.0, 3.0)
	http.HandleFunc("/hello", tempo.RateLimiterMiddleware(rl, helloHandler))
	rl.StartBackgroundCleanup(time.Minute*1, time.Minute*5)
	http.ListenAndServe(":8090", nil)
}

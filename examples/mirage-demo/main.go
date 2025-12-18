package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/amir-aharon/flux/mirage"
)

func main() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("processing heavy request...")
			time.Sleep(2 * time.Second)
			w.Write([]byte("finished processing!"))
		})
		http.ListenAndServe(":8081", mux)
	}()

	time.Sleep(1 * time.Second)

	backends := []string{"http://localhost:8081"}
	cdn, err := mirage.NewCDN(backends, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CDN running on :8080")
	if err := http.ListenAndServe(":8080", cdn); err != nil {
		log.Fatal(err)
	}
}

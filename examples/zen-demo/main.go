package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amir-aharon/flux/zen"
)

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "hello\n")
}

func main() {

	backendPorts := []string{":8081", ":8082"}
	var backendUrls []string
	for _, port := range backendPorts {
		backendUrls = append(backendUrls, fmt.Sprintf("http://localhost%s", port))
		go func(p string) {
			mux := http.NewServeMux()
			mux.HandleFunc("/hello", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(rw, "hello from port %s", p)
			})
			fmt.Printf("Starting backend on %s\n", p)
			http.ListenAndServe(p, mux)
		}(port)
	}

	lb, err := zen.NewLoadBalancer(backendUrls)
	if err != nil {
		log.Fatalf("%v", err)
	}

	http.ListenAndServe(":8080", lb)
}

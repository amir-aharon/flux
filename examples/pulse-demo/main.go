package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/amir-aharon/flux/pulse"
)

type Server struct {
	healthy bool
	mu      sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		healthy: true,
	}
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	healthy := s.healthy
	s.mu.RUnlock()

	if healthy {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Service Unavailable")
	}
}

func (s *Server) toggleHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.healthy = !s.healthy
	newState := s.healthy
	s.mu.Unlock()

	status := "unhealthy"
	if newState {
		status = "healthy"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Toggled to: %s\n", status)
	log.Printf("Health status toggled to: %s", status)
}

func main() {
	mon := pulse.NewMonitor(2 * time.Second)
	server := NewServer()
	http.HandleFunc("/healthz", server.healthzHandler)
	http.HandleFunc("/toggle", server.toggleHandler)
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	mon.AddEndpoint("github.com/amir-aharon")
	mon.AddEndpoint("http://localhost:8080/healthz")

	go mon.Start()

	for {
		fmt.Printf("%v\n", mon.GetStatus())
		time.Sleep(2 * time.Second)
	}

}

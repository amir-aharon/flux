package zen

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend *url.URL

type LoadBalancingStrategy interface {
	SelectBackend() Backend
}

type RoundRobin struct {
	Backends     []Backend
	BackendIndex int
	mu           sync.Mutex
}

func (strat *RoundRobin) SelectBackend() Backend {
	strat.mu.Lock()
	defer strat.mu.Unlock()

	backend := strat.Backends[strat.BackendIndex]
	strat.BackendIndex = (strat.BackendIndex + 1) % len(strat.Backends)
	return backend
}

type LoadBalancer struct {
	strat   LoadBalancingStrategy
	proxies map[Backend]*httputil.ReverseProxy
}

func NewLoadBalancer(backends []string) (*LoadBalancer, error) {
	if len(backends) == 0 {
		return nil, errors.New("LB has to have at least one backend")
	}

	parsedBackends := make([]Backend, len(backends))
	proxies := make(map[Backend]*httputil.ReverseProxy)
	for i, backend := range backends {
		parsedBackend, err := url.Parse(backend)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse backend %s. full error: %v", backend, err)
		}
		parsedBackends[i] = parsedBackend
		proxies[parsedBackend] = httputil.NewSingleHostReverseProxy(parsedBackend)
	}

	return &LoadBalancer{
		strat: &RoundRobin{
			Backends:     parsedBackends,
			BackendIndex: 0,
		},
		proxies: proxies,
	}, nil
}

func (lb *LoadBalancer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	target := lb.strat.SelectBackend()
	proxy := lb.proxies[target]
	proxy.ServeHTTP(rw, req)
}

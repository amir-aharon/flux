package pulse

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type EndpointStatus struct {
	URL         url.URL
	Up          bool
	Latency     time.Duration
	LastChecked time.Time
}

type Monitor struct {
	Interval time.Duration
	Status   map[url.URL]*EndpointStatus
	mu       sync.RWMutex
}

func NewMonitor(interval time.Duration) *Monitor {
	return &Monitor{
		Interval: interval,
		Status:   make(map[url.URL]*EndpointStatus),
	}
}

func (mon *Monitor) AddEndpoint(rawUrl string) error {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	mon.mu.Lock()
	defer mon.mu.Unlock()

	if _, exists := mon.Status[*parsedUrl]; exists {
		return nil
	}
	mon.Status[*parsedUrl] = &EndpointStatus{URL: *parsedUrl}
	return nil
}

func (mon *Monitor) Start() {
	for {
		time.Sleep(mon.Interval)

		mon.mu.RLock()
		endpoints := make([]url.URL, 0, len(mon.Status))
		for k := range mon.Status {
			endpoints = append(endpoints, k)
		}
		mon.mu.RUnlock()

		for _, e := range endpoints {
			status := mon.HealthCheck(e)

			mon.mu.Lock()
			mon.Status[e] = status
			mon.mu.Unlock()
		}
	}
}

func (mon *Monitor) HealthCheck(endpoint url.URL) *EndpointStatus {
	reqTime := time.Now()
	resp, err := http.Get(endpoint.String())
	if err != nil {
		return &EndpointStatus{
			URL:         endpoint,
			Up:          false,
			LastChecked: reqTime,
		}
	}

	latency := time.Since(reqTime)
	defer resp.Body.Close()

	return &EndpointStatus{
		URL:         endpoint,
		Up:          resp.StatusCode == http.StatusOK,
		Latency:     latency,
		LastChecked: reqTime,
	}
}

func (mon *Monitor) GetStatus() []*EndpointStatus {
	mon.mu.RLock()
	defer mon.mu.RUnlock()

	// Might need a thread-safe iterator for larger scale
	// Currently chose a 'Snapshot' for simplicity and to avoid edge cases
	results := make([]*EndpointStatus, 0, len(mon.Status))
	for _, status := range mon.Status {
		results = append(results, status)
	}
	return results
}

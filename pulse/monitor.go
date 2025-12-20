package pulse

import (
	"iter"
	"maps"
	"net/url"
	"time"
)

type URLStatus struct {
	URL         *url.URL
	Up          bool
	Latency     time.Duration
	LastChecked time.Time
}

type Monitor struct {
	Interval time.Duration
	Status   map[*url.URL]*URLStatus
}

func NewMonitor(interval time.Duration) *Monitor {
	return &Monitor{
		Interval: interval,
		Status:   make(map[*url.URL]*URLStatus),
	}
}

func (mon *Monitor) AddURL(rawUrl string) error {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	if _, exists := mon.Status[parsedUrl]; exists {
		return nil
	}
	mon.Status[parsedUrl] = &URLStatus{URL: parsedUrl}
	return nil
}

func (mon *Monitor) Start() {
	for {
		time.Sleep(mon.Interval)
	}
}

func (mon *Monitor) GetStatus() iter.Seq[*URLStatus] {
	return maps.Values(mon.Status)
}

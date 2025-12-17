package fuse

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	Open State = iota
	Closed
	HalfOpen
)

func (s State) String() string {
	switch s {
	case Open:
		return "Open"
	case Closed:
		return "Closed"
	case HalfOpen:
		return "HalfOpen"
	default:
		return "Unknown"
	}
}

type Config struct {
	Threshold    int
	Interval     time.Duration
	IsSuccessful func(err error) bool
}

type Breaker struct {
	Config
	State        State
	FailureCount int
	OpenExpiry   time.Time
	mu           sync.Mutex
}

func NewBreaker(cfg Config) *Breaker {
	if cfg.IsSuccessful == nil {
		cfg.IsSuccessful = func(err error) bool { return false }
	}
	return &Breaker{
		Config: cfg,
		State:  Closed,
	}
}

func (b *Breaker) Execute(job func() error) error {
	b.mu.Lock()
	allowedToRun := false
	switch b.State {
	case Closed:
		allowedToRun = true
	case Open:
		if time.Now().After(b.OpenExpiry) {
			b.State = HalfOpen
			allowedToRun = true
		}
	}
	b.mu.Unlock()

	if !allowedToRun {
		return errors.New("service unavailable")
	}

	err := job()
	crashed := err != nil && !b.IsSuccessful(err)

	b.mu.Lock()
	defer b.mu.Unlock()

	if !crashed {
		b.FailureCount = 0
		b.State = Closed
	} else {
		if b.State == Closed {
			b.FailureCount++
			if b.FailureCount < b.Threshold {
				return err
			}
		}

		b.State = Open
		b.OpenExpiry = time.Now().Add(b.Interval)
	}
	return err
}

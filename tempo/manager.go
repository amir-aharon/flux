package tempo

import (
	"sync"
	"time"
)

type RateLimiter struct {
	clients  map[string]*TokenBucket
	mu       sync.Mutex
	rate     float64
	capacity float64
}

func NewRateLimiter(rate, capacity float64) *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string]*TokenBucket),
		mu:       sync.Mutex{},
		rate:     rate,
		capacity: capacity,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	bucket, exists := rl.clients[ip]
	if !exists {
		bucket = NewTokenBucket(rl.rate, rl.capacity)
		rl.clients[ip] = bucket
	}
	rl.mu.Unlock()
	return bucket.Allow()
}

func (rl *RateLimiter) Cleanup(ttl time.Duration) {
	rl.mu.Lock()
	ips := make([]string, 0, len(rl.clients))
	for ip := range rl.clients {
		ips = append(ips, ip)
	}
	rl.mu.Unlock()

	for _, ip := range ips {
		rl.mu.Lock()
		bucket, exists := rl.clients[ip]
		rl.mu.Unlock()

		if !exists {
			continue
		}

		bucket.mu.Lock()
		isStale := time.Since(bucket.lastRefill) > ttl
		bucket.mu.Unlock()
		if isStale {
			rl.mu.Lock()
			delete(rl.clients, ip)
			rl.mu.Unlock()
		}
	}
}

func (rl *RateLimiter) StartBackgroundCleanup(interval time.Duration, ttl time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			<-ticker.C
			rl.Cleanup(ttl)
		}
	}()
}

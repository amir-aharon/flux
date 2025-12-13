package tempo

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: rate,
		lastRefill: time.Now(),
		mu:         sync.Mutex{},
	}
}

// refill is NOT thread safe
func (tb *TokenBucket) refill() {
	elapsedSeconds := time.Since(tb.lastRefill).Seconds()
	tokenDebt := tb.refillRate * elapsedSeconds
	tb.tokens += tokenDebt
	tb.tokens = min(tb.tokens, tb.capacity)
	tb.lastRefill = time.Now()
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

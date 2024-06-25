package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate         float64
	capacity     float64
	tokens       float64
	lastRefill   time.Time
	mu           sync.Mutex
	timeProvider TimeProvider
}

// NewTokenBucket create a new Rate Limiter using Token Bucket algorithm.
// rate represents the number of tokens added per second.
// capacity is the total number of tokens allowed anytime.
// tp represents the time provider for getting the current time. helpful for mocking time when testing.
func NewTokenBucket(rate float64, capacity float64, tp TimeProvider) *TokenBucket {
	return &TokenBucket{
		rate:         rate,
		capacity:     capacity,
		tokens:       capacity,
		lastRefill:   tp.Now().UTC(),
		timeProvider: tp,
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 1. Calculate the number of tokens to add based on the time elapsed since last refill
	now := tb.timeProvider.Now()
	elapsed := now.Sub(tb.lastRefill)

	// 2. Add tokens to the bucket (up to capacity)
	tokensToAdd := elapsed.Seconds() * tb.rate
	tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
	tb.lastRefill = now

	// 3. If there's at least one token, consume it and return true
	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}

	// 4. Otherwise, return false
	return false
}

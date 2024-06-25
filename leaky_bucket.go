package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	rate         float64
	capacity     float64
	water        float64
	lastLeak     time.Time
	mu           sync.Mutex
	timeProvider TimeProvider
}

// NewLeakyBucket create a new Rate Limiter using Leaky Bucket algorithm.
// rate represents the rate at which the water is leaked from the bucket.
// capacity is the total number of water the bucket can hold.
// tp represents the time provider for getting the current UTC time. helpful for mocking time when testing.
func NewLeakyBucket(rate float64, capacity float64, tp TimeProvider) *LeakyBucket {
	return &LeakyBucket{
		rate:         rate,
		capacity:     capacity,
		water:        0,
		lastLeak:     tp.Now(),
		timeProvider: tp,
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 1. Calculate the amount of water leaked since last check
	now := lb.timeProvider.Now()
	elapsed := now.Sub(lb.lastLeak)
	leakedAmount := elapsed.Seconds() * lb.rate

	// 2. Remove the leaked water from the bucket
	lb.water = max(lb.water-leakedAmount, 0)
	lb.lastLeak = now

	// 3. If there's room in the bucket, add water and return true
	if lb.water < lb.capacity {
		// If requests have impacts on the rate limiter, instead of adding 1, we can add
		// an argument to Allow() that specifies the amount of water to add for each request.
		lb.water += 1
		return true
	}

	// 4. Otherwise, return false
	return false
}

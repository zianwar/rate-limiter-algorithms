package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	rate          int
	interval      time.Duration
	currentCount  int
	previousCount int
	windowStart   time.Time
	mu            sync.Mutex
}

func NewSlidingWindowCounter(rate int, interval time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		rate:          rate,
		interval:      interval,
		currentCount:  0,
		previousCount: 0,
		windowStart:   time.Now(),
	}
}

func (swc *SlidingWindowCounter) Allow() bool {
	// Implement Sliding Window Counter algorithm here
	// 1. Check if we're in a new window, if so shift the counts
	// 2. Calculate the weighted count based on the position in the current window
	// 3. If the weighted count is less than the rate, increment the current count and return true
	// 4. Otherwise, return false
	return false
}

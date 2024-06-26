package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	rate          float64
	interval      time.Duration
	currentCount  int
	previousCount int
	windowStart   time.Time
	mu            sync.Mutex
	timeProvider  TimeProvider
}

func NewSlidingWindowCounter(rate float64, interval time.Duration, tp TimeProvider) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		rate:          rate,
		interval:      interval,
		currentCount:  0,
		previousCount: 0,
		windowStart:   tp.Now(),
		timeProvider:  tp,
	}
}

func (swc *SlidingWindowCounter) Allow() bool {
	swc.mu.Lock()
	defer swc.mu.Unlock()

	// 1. Check if we're in a new window, if so shift the counts
	now := swc.timeProvider.Now()
	elapsed := now.Sub(swc.windowStart)

	if elapsed >= swc.interval {
		// Move currentCount to previousCount, reset currentCount, adjust windowStart
		swc.previousCount += swc.currentCount
		swc.currentCount = 0
		swc.windowStart = now
	}

	// 2. Calculate the weighted count based on the position in the current window
	elapsedFraction := elapsed.Seconds() / swc.interval.Seconds()
	weightedCount := max(0, float64(swc.previousCount)*(1.0-elapsedFraction)+float64(swc.currentCount))

	// 3. If the weighted count is less than the rate, increment the current count and return true
	if weightedCount < swc.rate {
		swc.currentCount++
		return true
	}

	// 4. Otherwise, return false
	return false
}

package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowCounter struct {
	rate         int
	interval     time.Duration
	count        int
	windowStart  time.Time
	mu           sync.Mutex
	timeProvider TimeProvider
}

// NewFixedWindowCounter creates a new FixedWindowCounter with a specified rate and interval.
// rate specifies the maximum number of requests per interval
// interval specifies the length of the time window
// tp is the TimeProvider that defines how current time is determined.
func NewFixedWindowCounter(rate int, interval time.Duration, tp TimeProvider) *FixedWindowCounter {
	return &FixedWindowCounter{
		rate:         rate,
		interval:     interval,
		count:        0,
		windowStart:  time.Now().UTC(),
		timeProvider: tp,
	}
}

func (fwc *FixedWindowCounter) Allow() bool {
	fwc.mu.Lock()
	defer fwc.mu.Unlock()

	now := fwc.timeProvider.Now()

	// 1. Check if we're in a new window, if so reset the count
	//
	//  windowStart     now
	//      v            v
	//      [----------------]
	//           interval
	//
	if now.Sub(fwc.windowStart) >= fwc.interval {
		// Reset the window
		fwc.count = 0
		fwc.windowStart = now
	}

	// 2. If the count is less than the rate, increment and return true
	if fwc.count < fwc.rate {
		fwc.count++
		return true
	}
	// 3. Otherwise, return false
	return false
}

package ratelimiter

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	rate := 1.0
	capacity := 5.0
	mockTimeProvider := &MockTimeProvider{currentTime: time.Now().UTC()}
	tb := NewTokenBucket(rate, capacity, mockTimeProvider)

	// Should allow the first request
	if !tb.Allow() {
		t.Error("Expected first request to be allowed")
	}

	// Should allow up to capacity requests in quick succession
	for i := 0; i < int(capacity)-1; i++ {
		if !tb.Allow() {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Next request should be denied
	if tb.Allow() {
		t.Error("Expected request to be denied after capacity is reached")
	}

	// Test refill over time

	// Waiting for 2 seconds should add 2 tokens.
	mockTimeProvider.Advance(2 * time.Second)

	for i := 0; i < 2; i++ {
		if !tb.Allow() {
			t.Error("Expected to allow request after waiting, tokens should have refilled")
		}
	}

	// Should be empty again
	if tb.Allow() {
		t.Error("Expected to throttle request, after two allowed requests")
	}
}

func TestLeakyBucket(t *testing.T) {
	rate := 1.0
	capacity := 5.0
	mockTimeProvider := &MockTimeProvider{currentTime: time.Now().UTC()}
	lb := NewLeakyBucket(rate, capacity, mockTimeProvider)

	// Test: Fill the bucket to capacity
	for i := 0; i < int(capacity); i++ {
		if !lb.Allow() {
			t.Errorf("Bucket should allow filling up to capacity, failed at %d", i+1)
		}
	}

	// Test: Ensure the bucket is full and additional request is denied
	if lb.Allow() {
		t.Errorf("Bucket should be full and deny any additional requests")
	}

	// Test: Wait for half the capacity to leak out
	mockTimeProvider.Advance(3 * time.Second) // waiting for 3 seconds, 3 units should leak out

	// Test: Bucket should allow 3 more requests
	for i := 0; i < 3; i++ {
		if !lb.Allow() {
			t.Errorf("Bucket should have allowed after waiting, failed at request %d", i+1)
		}
	}

	// Test: Again, check if the bucket denies after refilling just allowed units
	if lb.Allow() {
		t.Errorf("Bucket should be full again and deny any additional requests")
	}

	// Test: Ensure correct behavior over a longer leak period
	mockTimeProvider.Advance(5 * time.Second) // all water should leak out, as it's the total capacity

	// Test: Bucket should be empty and allow filling to capacity again
	for i := 0; i < int(capacity); i++ {
		if !lb.Allow() {
			t.Errorf("Bucket should be empty and allow refilling, failed at %d", i+1)
		}
	}
}

func TestFixedWindowCounter(t *testing.T) {
	rate := 5
	interval := 1 * time.Second
	mockTimeProvider := &MockTimeProvider{currentTime: time.Now().UTC()}
	fwc := NewFixedWindowCounter(rate, interval, mockTimeProvider)

	// Test exceeding rate in the first window
	for i := 0; i < rate; i++ {
		if !fwc.Allow() {
			t.Errorf("Request %d was unexpectedly denied", i+1)
		}
	}

	// Next request should be denied
	if fwc.Allow() {
		t.Error("Request should be denied as the rate limit has been reached")
	}

	// Wait for the next window
	mockTimeProvider.Advance(interval)

	// Test rate limit in the new window
	for i := 0; i < rate; i++ {
		if !fwc.Allow() {
			t.Errorf("Request %d in new window was unexpectedly denied", i+1)
		}
	}

	// Next request in the second window should also be denied
	if fwc.Allow() {
		t.Error("Request in the second window should be denied as the rate limit has been reached again")
	}
}

func TestSlidingWindowCounter(t *testing.T) {
	rate := 5.0
	interval := 1 * time.Minute
	mockTimeProvider := &MockTimeProvider{currentTime: time.Now().UTC()}
	swc := NewSlidingWindowCounter(rate, interval, mockTimeProvider)

	// Should allow first request (R0) at 0:00
	if !swc.Allow() {
		t.Error("First request was unexpectedly denied")
	}

	// Move time to 0:30
	mockTimeProvider.Advance(30 * time.Second)

	// Send 4 requests (R1, R2, R3, R4) to fill up the rate, they should be allowed.
	for i := 0; i < 4; i++ {
		if !swc.Allow() {
			t.Errorf("Request %d was unexpectedly denied", i+1)
		}
	}

	// Move time to 1:10, shift into a new window.
	mockTimeProvider.Advance(40 * time.Second)

	// Send R5, it should reset the counters and start new window.r
	if !swc.Allow() {
		t.Error("Request R5 was unexpectedly denied at the start of a new window.")
	}

	// Move time to 1:20, still in the same window.
	mockTimeProvider.Advance(10 * time.Second)

	// Next request should be allowed since we started calculating a new window.
	if swc.Allow() {
		t.Error("Request should be denied as the sliding window rate limit has been reached")
	}
}

package ratelimiter

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (rtp *RealTimeProvider) Now() time.Time {
	return time.Now().UTC()
}

type MockTimeProvider struct {
	currentTime time.Time
}

func (mtp *MockTimeProvider) Now() time.Time {
	return mtp.currentTime
}

func (mtp *MockTimeProvider) Advance(d time.Duration) {
	mtp.currentTime = mtp.currentTime.Add(d)
}

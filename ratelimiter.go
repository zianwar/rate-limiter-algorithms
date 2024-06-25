package ratelimiter

// RateLimiter interface defines the common method for all rate limiting algorithms
type RateLimiter interface {
	Allow() bool
}

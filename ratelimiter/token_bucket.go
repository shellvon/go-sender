package ratelimiter

import (
	"context"

	"golang.org/x/time/rate"
)

// TokenBucketRateLimiter is a token bucket rate limiter based on golang.org/x/time/rate.
type TokenBucketRateLimiter struct {
	limiter *rate.Limiter
}

// NewTokenBucketRateLimiter creates a new token bucket rate limiter
// qps: requests per second allowed
// burst: burst capacity
func NewTokenBucketRateLimiter(qps float64, burst int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(qps), burst),
	}
}

// Allow checks if a request is allowed to pass.
func (r *TokenBucketRateLimiter) Allow() bool {
	return r.limiter.Allow()
}

// Wait waits until a request can pass.
func (r *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	return r.limiter.Wait(ctx)
}

// Close shuts down the rate limiter (no-op).
func (r *TokenBucketRateLimiter) Close() error {
	return nil
}

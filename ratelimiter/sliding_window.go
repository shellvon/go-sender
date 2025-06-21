package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrInvalidParameters = errors.New("invalid parameters")
)

// SlidingWindowRateLimiter is a rate limiter based on a sliding window
type SlidingWindowRateLimiter struct {
	mutex           sync.RWMutex
	windowSize      time.Duration
	maxRequests     int
	requests        []time.Time
	lastCleanup     time.Time
	cleanupInterval time.Duration
}

// NewSlidingWindowRateLimiter creates a new sliding window rate limiter
func NewSlidingWindowRateLimiter(windowSize time.Duration, maxRequests int) (*SlidingWindowRateLimiter, error) {
	if windowSize <= 0 {
		return nil, fmt.Errorf("%w: windowSize must be positive", ErrInvalidParameters)
	}
	if maxRequests <= 0 {
		return nil, fmt.Errorf("%w: maxRequests must be positive", ErrInvalidParameters)
	}

	return &SlidingWindowRateLimiter{
		windowSize:      windowSize,
		maxRequests:     maxRequests,
		requests:        make([]time.Time, 0, maxRequests), // Pre-allocate capacity
		lastCleanup:     time.Now(),
		cleanupInterval: windowSize / 10, // Periodic cleanup interval
	}, nil
}

// Allow checks if a request is allowed to pass
func (r *SlidingWindowRateLimiter) Allow(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// Clean up expired requests
	r.cleanup(now)

	// Check if the limit is exceeded
	if len(r.requests) >= r.maxRequests {
		return fmt.Errorf("%w: %d requests in %v", ErrRateLimitExceeded, len(r.requests), r.windowSize)
	}

	// Record the current request
	r.requests = append(r.requests, now)
	return nil
}

// Wait waits until a request can be made or context is cancelled
func (r *SlidingWindowRateLimiter) Wait(ctx context.Context) error {
	for {
		if err := r.Allow(ctx); err == nil {
			return nil
		} else if !errors.Is(err, ErrRateLimitExceeded) {
			return err
		}

		// Calculate how long to wait before the next request can be made
		waitTime := r.calculateWaitTime()
		if waitTime <= 0 {
			continue // Immediate retry
		}

		// Wait or be cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Continue trying
		}
	}
}

// GetStats returns current statistics
func (r *SlidingWindowRateLimiter) GetStats() (currentRequests int, maxRequests int, windowSize time.Duration) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-r.windowSize)

	// Calculate the number of requests in the current window
	count := 0
	for _, reqTime := range r.requests {
		if reqTime.After(cutoff) {
			count++
		}
	}

	return count, r.maxRequests, r.windowSize
}

// Reset clears all request history
func (r *SlidingWindowRateLimiter) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.requests = r.requests[:0] // Reset slice but retain capacity
	r.lastCleanup = time.Now()
}

// cleanup removes expired requests from the sliding window
func (r *SlidingWindowRateLimiter) cleanup(now time.Time) {
	// Only clean up when necessary to avoid frequent cleanup
	if now.Sub(r.lastCleanup) < r.cleanupInterval && len(r.requests) < r.maxRequests*2 {
		return
	}

	cutoff := now.Add(-r.windowSize)

	// Find the first unexpired request
	validIdx := len(r.requests) // Default all requests are expired
	for i, reqTime := range r.requests {
		if reqTime.After(cutoff) {
			validIdx = i
			break
		}
	}

	// Remove expired requests
	if validIdx > 0 {
		// Move valid requests to the front of the slice
		copy(r.requests, r.requests[validIdx:])
		r.requests = r.requests[:len(r.requests)-validIdx]
	}

	// If the slice is too large, reallocate to release memory if it is too large
	if cap(r.requests) > r.maxRequests*4 && len(r.requests) < r.maxRequests {
		newRequests := make([]time.Time, len(r.requests), r.maxRequests)
		copy(newRequests, r.requests)
		r.requests = newRequests
	}

	r.lastCleanup = now
}

// calculateWaitTime calculates how long to wait before the next request can be made
func (r *SlidingWindowRateLimiter) calculateWaitTime() time.Duration {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.requests) < r.maxRequests {
		return 0
	}

	// Find the oldest request and calculate when it expires
	now := time.Now()
	oldestRequest := r.requests[0]
	expireTime := oldestRequest.Add(r.windowSize)

	if expireTime.After(now) {
		return expireTime.Sub(now)
	}

	return 0
}

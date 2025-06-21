package core

import (
	"context"
)

// Handler is the core function that middleware wraps. It represents the
// final action, like sending a notification.
type Handler func(ctx context.Context, msg Message) error

// Middleware is a function that wraps a Handler to add functionality.
type Middleware func(next Handler) Handler

// SenderMiddleware holds configurations for sender middlewares.
type SenderMiddleware struct {
	RateLimiter    RateLimiter
	Retry          *RetryPolicy
	Queue          Queue
	CircuitBreaker CircuitBreaker
	Metrics        MetricsCollector
}

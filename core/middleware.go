package core

import (
	"context"
)

// SenderMiddleware holds configurations for sender middlewares.
type SenderMiddleware struct {
	RateLimiter    RateLimiter
	Retry          *RetryPolicy
	Queue          Queue
	CircuitBreaker CircuitBreaker
	Metrics        MetricsCollector

	// beforeHooks are executed BEFORE each send. Returning a non-nil error aborts the send.
	beforeHooks []BeforeHook
	// afterHooks are executed AFTER the send completes. They never affect the original result.
	afterHooks []AfterHook
}

// BeforeHook is invoked right before the message is sent. The implementation may inspect
// or mutate *SendOptions. Returning an error aborts the whole send operation.
type BeforeHook func(ctx context.Context, msg Message, opts *SendOptions) error

// AfterHook is invoked after the send attempt finishes (regardless of success).
// It should be side-effect only: logging, metrics, masking, etc. It must not change
// the outcome, hence has no return value. `result` can be nil on failure; `err` carries
// the send error when present.
type AfterHook func(ctx context.Context, msg Message, opts *SendOptions, result *SendResult, err error)

// UseBeforeHook adds a BeforeHook to the middleware chain.
// It is executed right before the message is sent.
// The implementation may inspect or mutate *SendOptions. Returning an error aborts the whole send operation.
func (sm *SenderMiddleware) UseBeforeHook(h BeforeHook) {
	if h == nil {
		return
	}
	sm.beforeHooks = append(sm.beforeHooks, h)
}

// UseAfterHook adds an AfterHook to the middleware chain.
// It is executed after the send attempt finishes (regardless of success).
// It should be side-effect only: logging, metrics, masking, etc. It must not change
// the outcome, hence has no return value. `result` can be nil on failure; `err` carries
// the send error when present.
func (sm *SenderMiddleware) UseAfterHook(h AfterHook) {
	if h == nil {
		return
	}
	sm.afterHooks = append(sm.afterHooks, h)
}

// WithBeforeHooks overrides the beforeHooks slice.
// It is used for one-time configuration during construction.
func (sm *SenderMiddleware) WithBeforeHooks(hooks ...BeforeHook) {
	sm.beforeHooks = hooks
}

// WithAfterHooks overrides the afterHooks slice.
// It is used for one-time configuration during construction.
func (sm *SenderMiddleware) WithAfterHooks(hooks ...AfterHook) {
	sm.afterHooks = hooks
}

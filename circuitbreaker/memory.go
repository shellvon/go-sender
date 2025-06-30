package circuitbreaker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shellvon/go-sender/core"
)

// MemoryCircuitBreaker implements a circuit breaker pattern using in-memory state.
type MemoryCircuitBreaker struct {
	name         string
	maxFailures  int64
	resetTimeout time.Duration
	logger       core.Logger

	mu              sync.RWMutex
	state           CircuitState
	failureCount    int64
	lastFailureTime time.Time
	nextRetryTime   time.Time
}

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	// StateClosed represents the closed state of the circuit breaker.
	StateClosed CircuitState = iota
	// StateOpen represents the open state of the circuit breaker.
	StateOpen
	// StateHalfOpen represents the half-open state of the circuit breaker.
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// NewMemoryCircuitBreaker creates a new in-memory circuit breaker.
func NewMemoryCircuitBreaker(name string, maxFailures int64, resetTimeout time.Duration) *MemoryCircuitBreaker {
	return &MemoryCircuitBreaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
		logger:       &core.NoOpLogger{},
	}
}

// SetLogger sets the logger for the circuit breaker.
func (cb *MemoryCircuitBreaker) SetLogger(logger core.Logger) {
	cb.logger = logger
}

// Execute executes the given function with circuit breaker protection.
func (cb *MemoryCircuitBreaker) Execute(_ context.Context, fn func() error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn()
	cb.afterRequest(err)

	return err
}

func (cb *MemoryCircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateOpen:
		if now.After(cb.nextRetryTime) {
			cb.state = StateHalfOpen
			return nil
		}
		return fmt.Errorf("circuit breaker %s is OPEN", cb.name)

	case StateHalfOpen:
		// Half-open state allows one request to pass
		return nil

	case StateClosed:
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state: %v", cb.state)
	}
}

func (cb *MemoryCircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}
}

func (cb *MemoryCircuitBreaker) onSuccess() {
	//nolint:exhaustive // intentionally not all cases handled, default covers the rest
	switch cb.state {
	case StateHalfOpen:
		cb.state = StateClosed
		cb.failureCount = 0
		_ = cb.logger.Log(
			core.LevelInfo,
			"message",
			"circuit breaker half open to closed",
			"circuit_breaker",
			cb.name,
			"state",
			"HALF_OPEN -> CLOSED",
		)

	case StateClosed:
		cb.failureCount = 0
	}
}

func (cb *MemoryCircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	//nolint:exhaustive // intentionally not all cases handled, default covers the rest
	switch cb.state {
	case StateClosed:
		if cb.failureCount >= cb.maxFailures {
			cb.state = StateOpen
			cb.nextRetryTime = time.Now().Add(cb.resetTimeout)
			_ = cb.logger.Log(
				core.LevelWarn,
				"message",
				"circuit breaker closed to open",
				"circuit_breaker",
				cb.name,
				"state",
				"CLOSED -> OPEN",
				"failures",
				cb.failureCount,
			)
		}

	case StateHalfOpen:
		cb.state = StateOpen
		cb.nextRetryTime = time.Now().Add(cb.resetTimeout)
		_ = cb.logger.Log(
			core.LevelWarn,
			"message",
			"circuit breaker half open to open",
			"circuit_breaker",
			cb.name,
			"state",
			"HALF_OPEN -> OPEN",
		)
	}
}

// GetState returns the current state.
func (cb *MemoryCircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailureCount returns the failure count.
func (cb *MemoryCircuitBreaker) GetFailureCount() int64 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failureCount
}

// Reset resets the circuit breaker.
func (cb *MemoryCircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	_ = cb.logger.Log(
		core.LevelInfo,
		"message",
		"circuit breaker reset",
		"circuit_breaker",
		cb.name,
		"state",
		"RESET -> CLOSED",
	)
	cb.state = StateClosed
	cb.failureCount = 0
	cb.lastFailureTime = time.Time{}
}

func (cb *MemoryCircuitBreaker) Close() error {
	return nil
}

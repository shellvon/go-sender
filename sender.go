package gosender

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shellvon/go-sender/core"
)

// Sender is the main entry point for the go-sender framework
type Sender struct {
	providers  map[core.ProviderType]*core.ProviderDecorator
	middleware *core.SenderMiddleware
	logger     core.Logger
	mu         sync.RWMutex
	closed     bool
}

// NewSender creates a new Sender instance
func NewSender(logger core.Logger) *Sender {
	if logger == nil {
		logger = &core.NoOpLogger{}
	}

	return &Sender{
		providers:  make(map[core.ProviderType]*core.ProviderDecorator),
		middleware: &core.SenderMiddleware{},
		logger:     logger,
	}
}

// RegisterProvider registers a provider with the sender
func (s *Sender) RegisterProvider(providerType core.ProviderType, provider core.Provider, middleware *core.SenderMiddleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if middleware == nil {
		copy := *s.middleware
		middleware = &copy
	}

	s.providers[providerType] = core.NewProviderDecorator(provider, middleware, s.logger)
	s.logger.Log(core.LevelInfo, "message", "provider registered", "provider", provider.Name(), "type", providerType)
}

// UnregisterProvider removes a provider from the sender
func (s *Sender) UnregisterProvider(providerType core.ProviderType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("sender is closed, cannot unregister provider")
	}

	if _, exists := s.providers[providerType]; !exists {
		return fmt.Errorf("provider type %s not found", providerType)
	}

	delete(s.providers, providerType)
	s.logger.Log(core.LevelInfo, "message", "provider unregistered", "type", providerType)
	return nil
}

// Send sends a message using the appropriate provider
func (s *Sender) Send(ctx context.Context, message core.Message, opts ...core.SendOption) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("sender is closed")
	}

	providerType := message.ProviderType()
	provider, exists := s.providers[providerType]
	if !exists {
		return fmt.Errorf("no provider registered for type %s", providerType)
	}

	return provider.Send(ctx, message, opts...)
}

// GetProvider retrieves a provider by type
func (s *Sender) GetProvider(providerType core.ProviderType) (*core.ProviderDecorator, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, exists := s.providers[providerType]
	return provider, exists
}

// SendVia sends a message via a specific channel (provider/bot) identified by channel.
// The channel should be a string (provider name, bot name).
// This method goes through all middleware (rate limiting, retry, circuit breaker, etc.).
// It's equivalent to:
//
//	ctx = core.WithCtxItemName(ctx, channel);
//	return s.Send(ctx, message, opts...)
func (s *Sender) SendVia(ctx context.Context, channel string, message core.Message, opts ...core.SendOption) error {
	ctx = core.WithCtxItemName(ctx, channel)
	return s.Send(ctx, message, opts...)
}

// SetRateLimiter sets the rate limiter for the sender
func (s *Sender) SetRateLimiter(rateLimiter core.RateLimiter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.RateLimiter = rateLimiter
}

// SetRetryPolicy sets the retry policy for the sender
func (s *Sender) SetRetryPolicy(retryPolicy *core.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if retryPolicy != nil {
		if err := retryPolicy.Validate(); err != nil {
			return fmt.Errorf("invalid retry policy: %w", err)
		}
	}

	s.middleware.Retry = retryPolicy
	return nil
}

// SetQueue sets the queue for the sender
func (s *Sender) SetQueue(queue core.Queue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.Queue = queue
}

// SetCircuitBreaker sets the circuit breaker for the sender
func (s *Sender) SetCircuitBreaker(circuitBreaker core.CircuitBreaker) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.CircuitBreaker = circuitBreaker
}

// SetMetrics sets the metrics collector for the sender
func (s *Sender) SetMetrics(metrics core.MetricsCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.Metrics = metrics
}

// HealthCheck performs a health check on the sender and all its components
func (s *Sender) HealthCheck(ctx context.Context) *core.SenderHealth {
	s.mu.RLock()
	defer s.mu.RUnlock()

	health := &core.SenderHealth{
		Status:    core.HealthStatusHealthy,
		Timestamp: time.Now(),
		Providers: make(map[core.ProviderType]*core.ProviderHealth),
	}

	// Check providers
	for providerType, provider := range s.providers {
		providerHealth := &core.ProviderHealth{
			ProviderType: providerType,
			Status:       core.HealthStatusHealthy,
			LastCheck:    time.Now(),
		}

		// Check if the underlying provider implements HealthChecker
		if healthChecker, ok := provider.Provider.(core.HealthChecker); ok {
			check := healthChecker.HealthCheck(ctx)
			if check != nil {
				providerHealth.Status = check.Status
				providerHealth.Message = check.Message
			}
		}

		health.Providers[providerType] = providerHealth

		// Update overall status
		if providerHealth.Status == core.HealthStatusUnhealthy {
			health.Status = core.HealthStatusUnhealthy
		} else if providerHealth.Status == core.HealthStatusDegraded && health.Status == core.HealthStatusHealthy {
			health.Status = core.HealthStatusDegraded
		}
	}

	// Check queue health
	if s.middleware.Queue != nil {
		if queueHealthChecker, ok := s.middleware.Queue.(core.HealthChecker); ok {
			health.Queue = queueHealthChecker.HealthCheck(ctx)
		}
	}

	// Check metrics health
	if s.middleware.Metrics != nil {
		if metricsHealthChecker, ok := s.middleware.Metrics.(core.HealthChecker); ok {
			health.Metrics = metricsHealthChecker.HealthCheck(ctx)
		}
	}

	return health
}

// Close gracefully shuts down the sender and all its components
func (s *Sender) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	var errs []error

	// Close all providers
	for providerType, provider := range s.providers {
		if err := provider.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close provider %s: %w", providerType, err))
		}
	}

	// Close middleware components
	if s.middleware != nil {
		if s.middleware.RateLimiter != nil {
			if err := s.middleware.RateLimiter.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close rate limiter: %w", err))
			}
		}
		if s.middleware.Queue != nil {
			if err := s.middleware.Queue.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close queue: %w", err))
			}
		}
		if s.middleware.CircuitBreaker != nil {
			if err := s.middleware.CircuitBreaker.Close(); err != nil {
				errs = append(errs, fmt.Errorf("failed to close circuit breaker: %w", err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}
	return nil
}

// IsClosed returns true if the sender has been closed
func (s *Sender) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

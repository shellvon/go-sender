package gosender

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shellvon/go-sender/core"
)

// Sender is the main entry point for the go-sender framework.
type Sender struct {
	providers  map[core.ProviderType]*core.ProviderDecorator
	middleware *core.SenderMiddleware
	logger     core.Logger
	mu         sync.RWMutex
	closed     bool
	// defaultHTTPClient is the global default HTTP client for all HTTP-based providers. SMTP/email is not affected.
	defaultHTTPClient *http.Client
}

// Option defines a function type for configuring Sender.
type Option func(*Sender)

// WithLogger sets a custom logger for the Sender.
//
//   - logger: a custom logger implementing core.Logger interface.
//   - If nil, the default logger is retained (NoOpLogger).
func WithLogger(logger core.Logger) Option {
	return func(s *Sender) {
		if logger != nil {
			s.logger = logger
		}
	}
}

// NewSender creates a new Sender instance with optional configurations.
//
//   - opts: optional configuration options (e.g., WithLogger, WithMiddleware, etc.)
//   - The default logger is NoOpLogger. Use WithLogger to set a custom logger.
func NewSender(opts ...Option) *Sender {
	s := &Sender{
		providers:         make(map[core.ProviderType]*core.ProviderDecorator),
		middleware:        &core.SenderMiddleware{},
		logger:            &core.NoOpLogger{}, // Default to NoOpLogger. Use WithLogger for custom logging.
		defaultHTTPClient: core.DefaultHTTPClient(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// RegisterProvider registers a provider with the sender.
func (s *Sender) RegisterProvider(
	providerType core.ProviderType,
	provider core.Provider,
	middleware *core.SenderMiddleware,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if middleware == nil {
		copyMiddleware := *s.middleware
		middleware = &copyMiddleware
	}

	s.providers[providerType] = core.NewProviderDecorator(provider, middleware, s.logger)
	_ = s.logger.Log(
		core.LevelInfo,
		"message",
		"provider registered",
		"provider",
		provider.Name(),
		"type",
		providerType,
	) // ignore log error
}

// UnregisterProvider removes a provider from the sender.
func (s *Sender) UnregisterProvider(providerType core.ProviderType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return errors.New("sender is closed, cannot unregister provider")
	}

	if _, exists := s.providers[providerType]; !exists {
		return fmt.Errorf("provider type %s not found", providerType)
	}

	delete(s.providers, providerType)
	_ = s.logger.Log(core.LevelInfo, "message", "provider unregistered", "type", providerType) // ignore log error
	return nil
}

// Send sends a message and only returns error information.
// It is a thin wrapper around Send that discards the detailed SendResult.
func (s *Sender) Send(ctx context.Context, message core.Message, opts ...core.SendOption) error {
	_, err := s.SendWithResult(ctx, message, opts...)
	return err
}

// SendWithResult sends a message and returns the detailed SendResult.
func (s *Sender) SendWithResult(
	ctx context.Context,
	message core.Message,
	opts ...core.SendOption,
) (*core.SendResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, errors.New("sender is closed")
	}

	providerType := message.ProviderType()
	provider, exists := s.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("no provider registered for type %s", providerType)
	}

	allOpts := append([]core.SendOption{core.WithSendHTTPClient(s.defaultHTTPClient)}, opts...)
	return provider.Send(ctx, message, allOpts...)
}

// GetProvider retrieves a provider by type.
func (s *Sender) GetProvider(providerType core.ProviderType) (*core.ProviderDecorator, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	provider, exists := s.providers[providerType]
	return provider, exists
}

// SetRateLimiter sets the rate limiter for the sender.
//
// NOTE: The change only applies to providers that are registered **after** this
// method is called. Each `ProviderDecorator` receives a *copy* of
// `Sender.middleware` at registration time; already-registered providers keep
// their own copy and will NOT pick up the new rate-limiter.
func (s *Sender) SetRateLimiter(rateLimiter core.RateLimiter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.RateLimiter = rateLimiter
}

// SetRetryPolicy sets the retry policy for the sender.
//
// NOTE: The update only affects providers registered *after* the call. See the
// comment on `SetRateLimiter` for rationale.
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

// SetQueue sets the queue for the sender.
//
// NOTE: This only influences providers registered **after** the setter is
// invoked. If you need the queue for earlier providers, call `SetQueue` *before*
// `RegisterProvider`, or re-register the provider.
func (s *Sender) SetQueue(queue core.Queue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.Queue = queue
}

// SetCircuitBreaker sets the circuit breaker for the sender.
//
// NOTE: Only providers registered after this call will use the circuit breaker.
func (s *Sender) SetCircuitBreaker(circuitBreaker core.CircuitBreaker) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.CircuitBreaker = circuitBreaker
}

// SetMetrics sets the metrics collector for the sender.
//
// NOTE: Like other middleware setters, this affects *future* providers only.
func (s *Sender) SetMetrics(metrics core.MetricsCollector) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.middleware.Metrics = metrics
}

// SetDefaultHTTPClient sets the global default HTTP client for all HTTP-based providers.
// This only affects HTTP/REST providers; SMTP/email providers are not affected.
func (s *Sender) SetDefaultHTTPClient(client *http.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.defaultHTTPClient = client
}

// HealthCheck performs a health check on the sender and all its components.
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

// closeComponent safely closes a component and appends any error to errs.
func closeComponent(closer interface{ Close() error }, errs *[]error, desc string) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			*errs = append(*errs, fmt.Errorf("failed to close %s: %w", desc, err))
		}
	}
}

// Close gracefully shuts down the sender and all its components.
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
		closeComponent(s.middleware.RateLimiter, &errs, "rate limiter")
		closeComponent(s.middleware.Queue, &errs, "queue")
		closeComponent(s.middleware.CircuitBreaker, &errs, "circuit breaker")
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}
	return nil
}

// IsClosed returns true if the sender has been closed.
func (s *Sender) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

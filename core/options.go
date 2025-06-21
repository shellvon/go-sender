package core

import (
	"errors"
	"log"
	"sync"
	"time"
)

// SendOptions manages all configurations related to sending notifications.
type SendOptions struct {
	Async                 bool                   // Whether to send asynchronously.
	Priority              int                    // Priority (1-10, 10 is highest).
	DelayUntil            *time.Time             // Time until which sending should be delayed.
	Timeout               time.Duration          // Send operation timeout.
	Metadata              map[string]interface{} // Additional metadata.
	DisableCircuitBreaker bool                   // Disable circuit breaker middleware.
	DisableRateLimiter    bool                   // Disable rate limiter
	Callback              func(error)            // Callback executed after message processing (only effective for local/in-memory queue or async goroutine, not called in distributed queue like Redis)
	RetryPolicy           *RetryPolicy           // Custom retry policy (overrides global), deserializable from queue will lost the filter function.
}

// NotificationMiddleware holds configurations for notification middlewares.
type NotificationMiddleware struct {
	RateLimiter    RateLimiter      // Rate limiting configuration.
	Retry          *RetryPolicy     // Retry configuration.
	Queue          Queue            // Queue configuration.
	CircuitBreaker CircuitBreaker   // Circuit breaker configuration.
	Metrics        MetricsCollector // Metrics collection configuration.
}

// RetryPolicy manages unified retry settings.
type RetryPolicy struct {
	MaxAttempts   int           // Maximum number of retry attempts.
	InitialDelay  time.Duration // Initial delay before the first retry.
	MaxDelay      time.Duration // Maximum delay between retries.
	BackoffFactor float64       // Factor by which the delay increases with each attempt.
	Filter        RetryFilter   // Custom filter function to determine if retry should occur.
	// Internal state for managing retry attempts.
	currentAttempt int
	mu             sync.RWMutex
}

// Reset resets the retry counter for a new operation.
func (r *RetryPolicy) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentAttempt = 0
}

// ShouldRetry determines if a retry should be attempted based on the current attempt count and error.
func (r *RetryPolicy) ShouldRetry(attempt int, err error) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if we've exceeded the maximum number of attempts.
	if attempt >= r.MaxAttempts {
		return false
	}

	// If no filter is provided, don't retry
	if r.Filter == nil {
		return false
	}

	// Use the custom filter function to determine if retry should occur.
	return r.Filter(attempt, err)
}

// NextDelay calculates the delay before the next retry attempt using exponential backoff.
func (r *RetryPolicy) NextDelay(attempt int, err error) time.Duration {
	delay := time.Duration(float64(r.InitialDelay) * float64(attempt) * r.BackoffFactor)
	if delay > r.MaxDelay {
		delay = r.MaxDelay
	}
	return delay
}

// RetryOption is a function type for configuring retry behavior.
type RetryOption func(*RetryPolicy)

// WithRetryMaxAttempts sets the maximum number of retry attempts.
func WithRetryMaxAttempts(attempts int) RetryOption {
	return func(c *RetryPolicy) {
		c.MaxAttempts = attempts
	}
}

// WithRetryInitialDelay sets the initial delay before the first retry.
func WithRetryInitialDelay(delay time.Duration) RetryOption {
	return func(c *RetryPolicy) {
		c.InitialDelay = delay
	}
}

// WithRetryMaxDelay sets the maximum delay between retries.
func WithRetryMaxDelay(delay time.Duration) RetryOption {
	return func(c *RetryPolicy) {
		c.MaxDelay = delay
	}
}

// WithRetryBackoffFactor sets the factor by which the delay increases with each attempt.
func WithRetryBackoffFactor(factor float64) RetryOption {
	return func(c *RetryPolicy) {
		c.BackoffFactor = factor
	}
}

// WithRetryFilter sets the custom filter function to determine if retry should occur.
func WithRetryFilter(filter RetryFilter) RetryOption {
	return func(c *RetryPolicy) {
		c.Filter = filter
	}
}

// DefaultRetryFilter creates a default retry filter that uses retryable errors and optional classifier fallback.
func DefaultRetryFilter(retryableErrors []error, fallbackToClassifier bool) RetryFilter {
	return func(attempt int, err error) bool {
		if err == nil {
			return false
		}

		// Check if error is in retryable errors list
		if len(retryableErrors) > 0 {
			for _, retryableErr := range retryableErrors {
				if errors.Is(err, retryableErr) {
					return true
				}
			}
		}

		// Fallback to classifier if enabled
		if fallbackToClassifier {
			classifier := NewDefaultErrorClassifier()
			return classifier.IsRetryableError(err)
		}

		return false
	}
}

// SendOption defines the function type for configuring SendOptions.
type SendOption func(*SendOptions)

// WithSendAsync sets the option to send the message asynchronously.
// If no argument is provided, it defaults to true.
func WithSendAsync(async ...bool) SendOption {
	return func(opts *SendOptions) {
		if len(async) > 0 {
			opts.Async = async[0]
		} else {
			opts.Async = true
		}
	}
}

// WithSendPriority sets the priority for the message.
func WithSendPriority(priority int) SendOption {
	return func(opts *SendOptions) {
		opts.Priority = priority
	}
}

// WithSendDelay sets a delay for sending the message.
func WithSendDelay(delay time.Duration) SendOption {
	return func(opts *SendOptions) {
		delayUntil := time.Now().Add(delay)
		opts.DelayUntil = &delayUntil
	}
}

// WithSendTimeout sets a timeout for the send operation.
func WithSendTimeout(timeout time.Duration) SendOption {
	return func(opts *SendOptions) {
		opts.Timeout = timeout
	}
}

// WithSendMetadata adds a key-value pair to the message's metadata.
func WithSendMetadata(key string, value interface{}) SendOption {
	return func(opts *SendOptions) {
		if opts.Metadata == nil {
			opts.Metadata = make(map[string]interface{})
		}
		opts.Metadata[key] = value
	}
}

// WithSendDisableCircuitBreaker sets the DisableCircuitBreaker option.
func WithSendDisableCircuitBreaker(disable bool) SendOption {
	return func(o *SendOptions) {
		o.DisableCircuitBreaker = disable
	}
}

// WithSendDisableRateLimiter sets the DisableRateLimiter option.
func WithSendDisableRateLimiter(disable bool) SendOption {
	return func(o *SendOptions) {
		o.DisableRateLimiter = disable
	}
}

// WithSendCallback sets the callback function to be executed after message processing.
// Note: callback is only effective for local/in-memory queue or async goroutine scenarios, and will not be called in distributed queues (such as Redis).
func WithSendCallback(callback func(error)) SendOption {
	return func(o *SendOptions) {
		o.Callback = callback
	}
}

// WithSendRetryPolicy sets a custom retry policy for this send operation.
func WithSendRetryPolicy(policy *RetryPolicy) SendOption {
	return func(o *SendOptions) {
		o.RetryPolicy = policy
	}
}

// frameworkMetadataKey is the namespaced key for storing SendOptions in Metadata.
const frameworkMetadataKey = "_framework_send_options"

// defaultSerializer is the default SendOptionsSerializer instance
var defaultSerializer SendOptionsSerializer = &DefaultSendOptionsSerializer{}

// serializeSendOptions serializes relevant SendOptions fields to JSON for storage in Metadata.
// It logs a warning if the key already exists to alert about potential user conflicts.
func serializeSendOptions(opts *SendOptions, metadata map[string]interface{}) (map[string]interface{}, error) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	if _, exists := metadata[frameworkMetadataKey]; exists {
		log.Print("metadata key conflict detected", "key", frameworkMetadataKey, "overwriting", "true")
	}

	data, err := defaultSerializer.Serialize(opts)
	if err != nil {
		return nil, NewSenderError(ErrCodeQueueSerializationFailed, "failed to serialize SendOptions", err)
	}
	metadata[frameworkMetadataKey] = data
	return metadata, nil
}

// deserializeSendOptions deserializes SendOptions from Metadata.
func deserializeSendOptions(metadata map[string]interface{}) (*SendOptions, error) {
	opts := &SendOptions{
		Metadata: metadata, // Preserve original Metadata
	}
	if metadata == nil {
		return opts, nil
	}

	data, ok := metadata[frameworkMetadataKey]
	if !ok {
		return opts, nil // Use defaults if key is missing
	}

	dataBytes, ok := data.([]byte)
	if !ok {
		log.Print("invalid type for metadata key", "key", frameworkMetadataKey, "expected", "[]byte")
		return opts, nil
	}

	deserializedOpts, err := defaultSerializer.Deserialize(dataBytes)
	if err != nil {
		return nil, NewSenderError(ErrCodeQueueDeserializationFailed, "failed to deserialize SendOptions", err)
	}

	// Merge deserialized options with preserved metadata
	opts.Priority = deserializedOpts.Priority
	opts.Timeout = deserializedOpts.Timeout
	opts.DisableCircuitBreaker = deserializedOpts.DisableCircuitBreaker
	opts.DisableRateLimiter = deserializedOpts.DisableRateLimiter
	opts.RetryPolicy = deserializedOpts.RetryPolicy

	return opts, nil
}

// ValidateRetryPolicy validates the retry policy configuration
func (r *RetryPolicy) Validate() error {
	if r.MaxAttempts < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "max attempts cannot be negative", nil)
	}
	if r.InitialDelay < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "initial delay cannot be negative", nil)
	}
	if r.MaxDelay < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "max delay cannot be negative", nil)
	}
	if r.BackoffFactor <= 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "backoff factor must be positive", nil)
	}
	if r.InitialDelay > r.MaxDelay {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "initial delay cannot be greater than max delay", nil)
	}
	return nil
}

// NewRetryPolicy creates a new RetryPolicy with the given options.
func NewRetryPolicy(opts ...RetryOption) *RetryPolicy {
	policy := &RetryPolicy{
		MaxAttempts:   3,                                   // Default to 3 attempts
		InitialDelay:  time.Second,                         // Default to 1 second initial delay
		MaxDelay:      30 * time.Second,                    // Default to 30 seconds max delay
		BackoffFactor: 2.0,                                 // Default to exponential backoff
		Filter:        DefaultRetryFilter([]error{}, true), // Use default filter with fallback enabled
	}

	// Apply all options
	for _, opt := range opts {
		opt(policy)
	}

	return policy
}

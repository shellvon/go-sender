package core

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	// DefaultUserAgent is the default User-Agent for HTTP requests.
	DefaultUserAgent = "go-sender/1.0.0"
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout   = 30 * time.Second
	defaultMaxAttempts   = 3
	defaultInitialDelay  = 100 * time.Millisecond
	defaultMaxDelay      = 30 * time.Second
	defaultBackoffFactor = 2.0
	DefaultTimeout       = 30 * time.Second
	maxIdleConns         = 100
	maxIdleConnsPerHost  = 10
	idleConnTimeout      = 90 * time.Second
)

type BodyType string

const (
	BodyTypeJSON BodyType = "json"
	BodyTypeForm BodyType = "form"
	BodyTypeText BodyType = "text"
	BodyTypeXML  BodyType = "xml"
	BodyTypeRaw  BodyType = "raw"
	BodyTypeNone BodyType = "none"
)

// DefaultHTTPClient returns a default HTTP client with proper settings.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultHTTPTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        maxIdleConns,
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			IdleConnTimeout:     idleConnTimeout,
		},
	}
}

// EnsureHTTPClient ensures that the HTTP client has a default User-Agent.
func EnsureHTTPClient(client *http.Client) *http.Client {
	if client == nil {
		client = DefaultHTTPClient()
	}

	// Ensure User-Agent is set
	if client.Transport == nil {
		client.Transport = &http.Transport{}
	}

	return client
}

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
	// HTTPClient allows per-send custom HTTP client. Only affects HTTP-based providers; SMTP/email is not affected.
	HTTPClient *http.Client // Optional: custom HTTP client for this send
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
	Jitter         bool
}

// NewRetryPolicy creates a new RetryPolicy with the given options.
func NewRetryPolicy(opts ...RetryOption) *RetryPolicy {
	policy := &RetryPolicy{
		MaxAttempts:   defaultMaxAttempts,            // Default to 3 attempts
		InitialDelay:  defaultInitialDelay,           // Default to 100ms
		MaxDelay:      defaultMaxDelay,               // Default to 30 seconds max delay
		BackoffFactor: defaultBackoffFactor,          // Default to exponential backoff
		Jitter:        true,                          // Default to jitter enabled
		Filter:        DefaultRetryFilter(nil, true), // Default to retry on all errors
	}
	for _, opt := range opts {
		opt(policy)
	}
	return policy
}

// Reset resets the retry policy's current attempt.
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
func (r *RetryPolicy) NextDelay(attempt int, _ error) time.Duration {
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
	return func(_ int, err error) bool {
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

// WithSendHTTPClient sets a custom HTTP client for this send operation.
// Only affects HTTP-based providers; SMTP/email providers are not affected.
func WithSendHTTPClient(client *http.Client) SendOption {
	return func(opts *SendOptions) {
		opts.HTTPClient = EnsureHTTPClient(client)
	}
}

// frameworkMetadataKey is the namespaced key for storing SendOptions in Metadata.
const frameworkMetadataKey = "__gosender_framework_send_options"

// internalContextItemNameKey is the internal key for storing context item name in queue metadata.
const internalContextItemNameKey = "__gosender_internal_ctx_item_name__"

// defaultSerializer is the default SendOptionsSerializer instance.
//
//nolint:gochecknoglobals // Reason: defaultSerializer is a global default for SendOptions serialization
var defaultSerializer SendOptionsSerializer = &DefaultSendOptionsSerializer{}

// serializeSendOptions serializes relevant SendOptions fields to JSON for storage in Metadata.
// It logs a warning if the key already exists to alert about potential user conflicts.
// It also preserves context information (like specified item names) for queue recovery.
func serializeSendOptions(
	ctx context.Context,
	opts *SendOptions,
	metadata map[string]interface{},
) (map[string]interface{}, error) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	// Note: metadata key conflict detection removed for lint compliance

	// Preserve context information for queue recovery
	if itemName := GetItemNameFromCtx(ctx); itemName != "" {
		metadata[internalContextItemNameKey] = itemName
	}

	data, err := defaultSerializer.Serialize(opts)
	if err != nil {
		return nil, NewSenderError(ErrCodeQueueSerializationFailed, "failed to serialize SendOptions", err)
	}
	metadata[frameworkMetadataKey] = data
	return metadata, nil
}

// deserializeSendOptions deserializes SendOptions from Metadata and restores context information.
// Returns the restored context and SendOptions.
func deserializeSendOptions(
	ctx context.Context,
	metadata map[string]interface{},
) (context.Context, *SendOptions, error) {
	opts := &SendOptions{
		Metadata: metadata, // Preserve original Metadata
	}
	if metadata == nil {
		return ctx, opts, nil
	}

	// Restore context information from metadata
	if itemName, ok := metadata[internalContextItemNameKey].(string); ok && itemName != "" {
		ctx = WithCtxItemName(ctx, itemName)
	}

	data, ok := metadata[frameworkMetadataKey]
	if !ok {
		return ctx, opts, nil // Use defaults if key is missing
	}

	dataBytes, ok := data.([]byte)
	if !ok {
		// log.Print("invalid type for metadata key", "key", frameworkMetadataKey, "expected", "[]byte")
		return ctx, opts, nil
	}

	deserializedOpts, err := defaultSerializer.Deserialize(dataBytes)
	if err != nil {
		return nil, nil, NewSenderError(ErrCodeQueueDeserializationFailed, "failed to deserialize SendOptions", err)
	}

	// Merge deserialized options with preserved metadata
	opts.Priority = deserializedOpts.Priority
	opts.Timeout = deserializedOpts.Timeout
	opts.DisableCircuitBreaker = deserializedOpts.DisableCircuitBreaker
	opts.DisableRateLimiter = deserializedOpts.DisableRateLimiter
	opts.RetryPolicy = deserializedOpts.RetryPolicy

	return ctx, opts, nil
}

// Validate validates the retry policy configuration.
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

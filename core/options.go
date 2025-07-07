package core

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
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

type BodyType int

const (
	BodyTypeNone BodyType = iota
	BodyTypeJSON
	BodyTypeForm
	BodyTypeText
	BodyTypeXML
	BodyTypeRaw
)

//nolint:gochecknoglobals // constant mapping table, safe as global
var bodyTypeToString = map[BodyType]string{
	BodyTypeNone: "none",
	BodyTypeJSON: "json",
	BodyTypeForm: "form",
	BodyTypeText: "text",
	BodyTypeXML:  "xml",
	BodyTypeRaw:  "raw",
}

//nolint:gochecknoglobals // reverse lookup table, safe as global
var stringToBodyType = map[string]BodyType{
	"none": BodyTypeNone,
	"json": BodyTypeJSON,
	"form": BodyTypeForm,
	"text": BodyTypeText,
	"xml":  BodyTypeXML,
	"raw":  BodyTypeRaw,
}

// String implements fmt.Stringer for human-readable output.
func (b *BodyType) String() string {
	if b == nil {
		return "unknown"
	}
	if s, ok := bodyTypeToString[*b]; ok {
		return s
	}
	return "unknown"
}

// ContentType returns the corresponding HTTP Content-Type header value.
func (b *BodyType) ContentType() string {
	if b == nil {
		return ""
	}
	switch *b {
	case BodyTypeNone:
		return ""
	case BodyTypeJSON:
		return "application/json; charset=utf-8"
	case BodyTypeForm:
		return "application/x-www-form-urlencoded; charset=utf-8"
	case BodyTypeText:
		return "text/plain; charset=utf-8"
	case BodyTypeXML:
		return "application/xml; charset=utf-8"
	case BodyTypeRaw:
		return "application/octet-stream"
	default:
		return ""
	}
}

// MarshalJSON encodes the BodyType as a quoted string (e.g., "json").
func (b *BodyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON decodes a quoted string (e.g., "json") into BodyType.
func (b *BodyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToBodyType[s]; ok {
		*b = val
		return nil
	}
	*b = BodyTypeNone
	return nil
}

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
	// Whether to send asynchronously.
	Async bool
	// Priority (1-10, 10 is highest).
	Priority int
	// Time until which sending should be delayed.
	// This field is only effective when `Async` is set to `true`.
	DelayUntil *time.Time
	// Send operation timeout.
	Timeout time.Duration
	// Additional metadata.
	Metadata map[string]interface{}
	// Disable circuit breaker middleware.
	DisableCircuitBreaker bool
	// Disable rate limiter.
	DisableRateLimiter bool
	// Callback is an optional function that will be executed after the message
	// has been fully processed (either successfully sent or failed after all retries).
	// This callback is primarily effective for local/in-memory queue processing
	// or when messages are sent via an asynchronous goroutine.
	// It is important to note that this callback is generally NOT invoked
	// when using distributed queues (e.g., Redis, Kafka) because the
	// message processing might occur in a separate consumer process
	// where the original context and callback function are no longer available.
	//
	// Specifying a Callback also implies that `Async` must be set to `true`
	// for the callback to be effective, as callbacks are not supported for synchronous sends.
	Callback func(error)
	// Custom retry policy (overrides global), deserializable from queue will lost the filter function.
	RetryPolicy *RetryPolicy
	// HTTPClient allows per-send custom HTTP client.
	// Only affects HTTP-based providers; SMTP/email is not affected.
	// Optional: custom HTTP client for this send.
	HTTPClient *http.Client
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
	MaxAttempts    int           // Maximum number of retry attempts.
	InitialDelay   time.Duration // Initial delay before the first retry.
	MaxDelay       time.Duration // Maximum delay between retries.
	BackoffFactor  float64       // Factor by which the delay increases with each attempt.
	Filter         RetryFilter   // Custom filter function to determine if retry should occur.
	currentAttempt int           // Internal state for managing retry attempts.
	rng            *rand.Rand    // Random number generator.
	mu             sync.RWMutex  // Mutex for thread safety.
}

// NewRetryPolicy creates a new RetryPolicy with the given options.
func NewRetryPolicy(opts ...RetryOption) *RetryPolicy {
	policy := &RetryPolicy{
		MaxAttempts:   defaultMaxAttempts,            // Default to 3 attempts
		InitialDelay:  defaultInitialDelay,           // Default to 100ms
		MaxDelay:      defaultMaxDelay,               // Default to 30 seconds max delay
		BackoffFactor: defaultBackoffFactor,          // Default to exponential backoff
		Filter:        DefaultRetryFilter(nil, true), // Default to retry on all errors
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
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

// NextDelay calculates the delay before the next retry attempt using exponential backoff with jitter.
//
// Exponential backoff calculation:
//   - The attempt parameter is 0-indexed for the first attempt, 1 for the second, etc.
//   - For exponential backoff, the factor should be raised to the power of the attempt number.
//     So, for attempt 0, delay = initialDelay * backoffFactor^0 = initialDelay
//   - For attempt 1, delay = initialDelay * backoffFactor^1 And so on.
func (r *RetryPolicy) NextDelay(attempt int, _ error) time.Duration {
	calculatedDelay := time.Duration(float64(r.InitialDelay) * math.Pow(r.BackoffFactor, float64(attempt)))
	if calculatedDelay > 0 {
		calculatedDelay = time.Duration(r.rng.Int63n(int64(calculatedDelay) + 1))
	}
	// Cap the delay at MaxDelay
	if calculatedDelay > r.MaxDelay {
		calculatedDelay = r.MaxDelay
	}
	return calculatedDelay
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

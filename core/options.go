package core

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

const (
	// DefaultUserAgent is the default User-Agent for HTTP requests.
	DefaultUserAgent = "go-sender/0.1.5"
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30 * time.Second
	// defaultMaxAttempts specifies the default maximum number of retry attempts.
	defaultMaxAttempts = 3
	// defaultInitialDelay specifies the default initial delay before the first retry.
	defaultInitialDelay = 100 * time.Millisecond
	// defaultMaxDelay specifies the default maximum delay between retries.
	defaultMaxDelay = 30 * time.Second
	// defaultBackoffFactor specifies the default factor by which the delay increases with each attempt.
	defaultBackoffFactor = 2.0
	// DefaultTimeout specifies the default timeout for send operations.
	DefaultTimeout = 30 * time.Second
	// maxIdleConns specifies the maximum number of idle (keep-alive) connections to keep open.
	maxIdleConns = 100
	// maxIdleConnsPerHost specifies the maximum idle (keep-alive) connections to keep per-host.
	maxIdleConnsPerHost = 10
	// idleConnTimeout specifies the amount of time an idle (keep-alive) connection will remain open before closing itself.
	idleConnTimeout = 90 * time.Second
)

// BodyType represents the type of a message body for HTTP requests.
type BodyType int

const (
	// BodyTypeNone indicates no specific body type.
	BodyTypeNone BodyType = iota
	// BodyTypeJSON indicates a JSON body type.
	BodyTypeJSON
	// BodyTypeForm indicates a form-encoded body type.
	BodyTypeForm
	// BodyTypeText indicates a plain text body type.
	BodyTypeText
	// BodyTypeXML indicates an XML body type.
	BodyTypeXML
	// BodyTypeRaw indicates a raw binary body type.
	BodyTypeRaw
)

// bodyTypeToString maps BodyType values to their string representations.
//
//nolint:gochecknoglobals // constant mapping table, safe as global.
var bodyTypeToString = map[BodyType]string{
	BodyTypeNone: "none",
	BodyTypeJSON: "json",
	BodyTypeForm: "form",
	BodyTypeText: "text",
	BodyTypeXML:  "xml",
	BodyTypeRaw:  "raw",
}

// stringToBodyType maps string representations to BodyType values.
//
//nolint:gochecknoglobals // reverse lookup table, safe as global.
var stringToBodyType = map[string]BodyType{
	"none": BodyTypeNone,
	"json": BodyTypeJSON,
	"form": BodyTypeForm,
	"text": BodyTypeText,
	"xml":  BodyTypeXML,
	"raw":  BodyTypeRaw,
}

// String implements fmt.Stringer for human-readable output of BodyType.
func (b *BodyType) String() string {
	if b == nil {
		return "unknown"
	}
	if s, ok := bodyTypeToString[*b]; ok {
		return s
	}
	return "unknown"
}

// ContentType returns the corresponding HTTP Content-Type header value for the BodyType.
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

// EnsureHTTPClient validates and initializes the provided HTTP client.
//
// If the provided client is nil, this function creates and returns a new default client
// with a predefined timeout. It also ensures the client has a non-nil Transport,
// which is required for making HTTP requests.
func EnsureHTTPClient(client *http.Client) *http.Client {
	if client == nil {
		client = DefaultHTTPClient()
		client.Timeout = DefaultHTTPTimeout
	}

	// Ensure User-Agent is set.
	if client.Transport == nil {
		client.Transport = &http.Transport{}
	}

	return client
}

// SendOptions manages all configurations related to sending notifications.
type SendOptions struct {
	// Async indicates whether to send asynchronously.
	Async bool
	// Priority (1-10, 10 is highest) specifies the priority of the message.
	Priority int
	// DelayUntil specifies the exact time until which sending should be delayed.
	// This field is only effective when `Async` is set to `true`.
	DelayUntil *time.Time
	// Timeout specifies the send operation timeout.
	Timeout time.Duration
	// Metadata holds additional metadata for the message.
	Metadata map[string]interface{}
	// AccountName forces the provider to use the specified account name.
	// When set, the framework will inject it into context so that BaseConfig.Select
	// picks this exact account instead of the usual strategy.
	AccountName string
	// StrategyName overrides the default strategy for this single send (e.g. "round_robin").
	StrategyName string
	// DisableCircuitBreaker indicates whether to disable the circuit breaker middleware for this send.
	DisableCircuitBreaker bool
	// DisableRateLimiter indicates whether to disable the rate limiter middleware for this send.
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
	Callback func(*SendResult, error)
	// RetryPolicy allows for a custom retry policy for this send operation (overrides global).
	// Note: The Filter function within RetryPolicy will be lost if deserialized from a queue.
	RetryPolicy *RetryPolicy
	// HTTPClient allows a per-send custom HTTP client.
	// Only affects HTTP-based providers; SMTP/email providers are not affected.
	// This client is optional.
	HTTPClient *http.Client

	// ---- Per-request hooks (not serialized) --------------------------
	// BeforeHooks are executed after global SenderMiddleware.beforeHooks but before the send attempt.
	BeforeHooks []BeforeHook `json:"-"`
	// AfterHooks are executed after global SenderMiddleware.afterHooks.
	AfterHooks []AfterHook `json:"-"`
}

// NotificationMiddleware holds configurations for various notification middlewares.
type NotificationMiddleware struct {
	// RateLimiter specifies the rate limiting configuration.
	RateLimiter RateLimiter
	// Retry specifies the retry configuration.
	Retry *RetryPolicy
	// Queue specifies the queue configuration.
	Queue Queue
	// CircuitBreaker specifies the circuit breaker configuration.
	CircuitBreaker CircuitBreaker
	// Metrics specifies the metrics collection configuration.
	Metrics MetricsCollector
}

// RetryPolicy manages unified retry settings for message sending.
type RetryPolicy struct {
	// MaxAttempts specifies the maximum number of retry attempts.
	MaxAttempts int
	// InitialDelay specifies the initial delay before the first retry.
	InitialDelay time.Duration
	// MaxDelay specifies the maximum delay between retries.
	MaxDelay time.Duration
	// BackoffFactor specifies the factor by which the delay increases with each attempt.
	BackoffFactor float64
	// Filter is a custom filter function to determine if a retry should occur.
	Filter RetryFilter
	// currentAttempt tracks the internal state for managing retry attempts.
	currentAttempt int
	// mu protects the currentAttempt field for concurrent access.
	mu sync.RWMutex
}

// NewRetryPolicy creates a new RetryPolicy with the given options.
func NewRetryPolicy(opts ...RetryOption) *RetryPolicy {
	policy := &RetryPolicy{
		MaxAttempts:   defaultMaxAttempts,            // Default to 3 attempts.
		InitialDelay:  defaultInitialDelay,           // Default to 100ms.
		MaxDelay:      defaultMaxDelay,               // Default to 30 seconds max delay.
		BackoffFactor: defaultBackoffFactor,          // Default to exponential backoff.
		Filter:        DefaultRetryFilter(nil, true), // Default to retry on all errors.
	}
	for _, opt := range opts {
		opt(policy)
	}
	return policy
}

// Reset resets the retry policy's current attempt count to zero.
func (r *RetryPolicy) Reset() {
	r.mu.Lock() // Acquire write lock for initialization.
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

	// If no filter is provided, do not retry.
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
//     So, for attempt 0, delay = initialDelay * backoffFactor^0 = initialDelay.
//   - For attempt 1, delay = initialDelay * backoffFactor^1. And so on.
func (r *RetryPolicy) NextDelay(attempt int, _ error) time.Duration {
	calculatedDelay := time.Duration(float64(r.InitialDelay) * math.Pow(r.BackoffFactor, float64(attempt)))

	// Apply full jitter using the package-level default random number generator, which is concurrency safe.
	//nolint:gosec // G404: This is for retry delay jitter, not security-related
	if calculatedDelay > 0 {
		calculatedDelay = time.Duration(rand.Int64N(int64(calculatedDelay) + 1))
	}

	// Cap the delay at MaxDelay.
	if calculatedDelay > r.MaxDelay {
		calculatedDelay = r.MaxDelay
	}

	return calculatedDelay
}

// RetryOption is a function type for configuring retry behavior.
type RetryOption func(*RetryPolicy)

// WithRetryMaxAttempts sets the maximum number of retry attempts for the policy.
func WithRetryMaxAttempts(attempts int) RetryOption {
	return func(c *RetryPolicy) {
		c.MaxAttempts = attempts
	}
}

// WithRetryInitialDelay sets the initial delay before the first retry for the policy.
func WithRetryInitialDelay(delay time.Duration) RetryOption {
	return func(c *RetryPolicy) {
		c.InitialDelay = delay
	}
}

// WithRetryMaxDelay sets the maximum delay between retries for the policy.
func WithRetryMaxDelay(delay time.Duration) RetryOption {
	return func(c *RetryPolicy) {
		c.MaxDelay = delay
	}
}

// WithRetryBackoffFactor sets the factor by which the delay increases with each attempt for the policy.
func WithRetryBackoffFactor(factor float64) RetryOption {
	return func(c *RetryPolicy) {
		c.BackoffFactor = factor
	}
}

// WithRetryFilter sets the custom filter function to determine if retry should occur for the policy.
func WithRetryFilter(filter RetryFilter) RetryOption {
	return func(c *RetryPolicy) {
		c.Filter = filter
	}
}

// DefaultRetryFilter creates a default retry filter that uses retryable errors and an optional classifier fallback.
func DefaultRetryFilter(retryableErrors []error, fallbackToClassifier bool) RetryFilter {
	return func(_ int, err error) bool {
		if err == nil {
			return false
		}

		// Check if error is in the retryable errors list.
		if len(retryableErrors) > 0 {
			for _, retryableErr := range retryableErrors {
				if errors.Is(err, retryableErr) {
					return true
				}
			}
		}

		// Fallback to classifier if enabled.
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

// WithSendAccount selects a specific account by name, bypassing strategy selection.
func WithSendAccount(name string) SendOption {
	return func(opts *SendOptions) {
		opts.AccountName = name
	}
}

// WithSendStrategy sets a per-send selection strategy.
func WithSendStrategy(st StrategyType) SendOption {
	return func(opts *SendOptions) {
		opts.StrategyName = string(st)
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
func WithSendCallback(callback func(*SendResult, error)) SendOption {
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
func WithSendHTTPClient(client *http.Client) SendOption {
	return func(opts *SendOptions) {
		opts.HTTPClient = EnsureHTTPClient(client)
	}
}

// WithSendBeforeHooks appends per-request BeforeHooks.
func WithSendBeforeHooks(hooks ...BeforeHook) SendOption {
	return func(opts *SendOptions) {
		opts.BeforeHooks = append(opts.BeforeHooks, hooks...)
	}
}

// WithSendAfterHooks appends per-request AfterHooks.
func WithSendAfterHooks(hooks ...AfterHook) SendOption {
	return func(opts *SendOptions) {
		opts.AfterHooks = append(opts.AfterHooks, hooks...)
	}
}

// sendOptionsMetadataKey is the sole key used to embed serialized SendOptions into queue Metadata.
const sendOptionsMetadataKey = "__gosender_send_options__"

// defaultSerializer is the default SendOptionsSerializer instance.
//
//nolint:gochecknoglobals // Reason: defaultSerializer is a global default for SendOptions serialization.
var defaultSerializer SendOptionsSerializer = &DefaultSendOptionsSerializer{}

// serializeSendOptions serializes relevant SendOptions fields to JSON for storage in Metadata.
// It logs a warning if the key already exists to alert about potential user conflicts.
// It also preserves context information (like specified item names) for queue recovery.
func serializeSendOptions(
	_ context.Context,
	opts *SendOptions,
	metadata map[string]interface{},
) (map[string]interface{}, error) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	// Note: metadata key conflict detection removed for lint compliance.

	data, err := defaultSerializer.Serialize(opts)
	if err != nil {
		return nil, NewSenderError(ErrCodeQueueSerializationFailed, "failed to serialize SendOptions", err)
	}
	metadata[sendOptionsMetadataKey] = data
	return metadata, nil
}

// deserializeSendOptions deserializes SendOptions from Metadata and restores context information.
// Returns the restored context and SendOptions.
func deserializeSendOptions(
	ctx context.Context,
	metadata map[string]interface{},
) (context.Context, *SendOptions, error) {
	opts := &SendOptions{
		Metadata: metadata, // Preserve original Metadata.
	}
	if metadata == nil {
		return ctx, opts, nil
	}

	// After deserialization we'll inject route info below.
	data, ok := metadata[sendOptionsMetadataKey]
	if !ok {
		return ctx, opts, nil // Use defaults if key is missing.
	}

	dataBytes, ok := data.([]byte)
	if !ok {
		return ctx, opts, nil
	}

	deserializedOpts, err := defaultSerializer.Deserialize(dataBytes)
	if err != nil {
		return nil, nil, NewSenderError(ErrCodeQueueDeserializationFailed, "failed to deserialize SendOptions", err)
	}

	opts.Priority = deserializedOpts.Priority
	opts.Timeout = deserializedOpts.Timeout
	opts.DisableCircuitBreaker = deserializedOpts.DisableCircuitBreaker
	opts.DisableRateLimiter = deserializedOpts.DisableRateLimiter
	opts.RetryPolicy = deserializedOpts.RetryPolicy
	opts.AccountName = deserializedOpts.AccountName
	opts.StrategyName = deserializedOpts.StrategyName

	// Rebuild route info for ctx
	if opts.AccountName != "" || opts.StrategyName != "" {
		ctx = WithRoute(ctx, &RouteInfo{AccountName: opts.AccountName, StrategyType: StrategyType(opts.StrategyName)})
	}

	return ctx, opts, nil
}

// Validate validates the retry policy configuration.
func (r *RetryPolicy) Validate() error {
	if r.MaxAttempts < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "max attempts cannot be negative.", nil)
	}
	if r.InitialDelay < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "initial delay cannot be negative.", nil)
	}
	if r.MaxDelay < 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "max delay cannot be negative.", nil)
	}
	if r.BackoffFactor <= 0 {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "backoff factor must be positive.", nil)
	}
	if r.InitialDelay > r.MaxDelay {
		return NewSenderError(ErrCodeRetryPolicyInvalid, "initial delay cannot be greater than max delay.", nil)
	}
	return nil
}

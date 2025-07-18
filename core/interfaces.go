package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// ProviderType defines the type of a notification provider.
type ProviderType string

const (
	// ProviderTypeSMS represents an SMS notification provider.
	ProviderTypeSMS ProviderType = "sms"
	// ProviderTypeEmail represents an email notification provider.
	ProviderTypeEmail ProviderType = "email"
	// ProviderTypeWecombot represents a WeCom bot notification provider.
	ProviderTypeWecombot ProviderType = "wecombot"
	// ProviderTypeWebhook represents a generic webhook notification provider.
	ProviderTypeWebhook ProviderType = "webhook"
	// ProviderTypeTelegram represents a Telegram bot notification provider.
	ProviderTypeTelegram ProviderType = "telegram"
	// ProviderTypeDingtalk represents a DingTalk bot notification provider.
	ProviderTypeDingtalk ProviderType = "dingtalk"
	// ProviderTypeLark represents a Lark/Feishu bot notification provider.
	ProviderTypeLark ProviderType = "lark"
	// ProviderTypeServerChan represents a ServerChan notification provider.
	ProviderTypeServerChan ProviderType = "serverchan"
	// ProviderTypeEmailAPI represents an emailapi notification provider.
	ProviderTypeEmailAPI ProviderType = "emailapi"
)

// HealthStatus represents the health status of a component.
type HealthStatus string

const (
	// HealthStatusHealthy represents a healthy status.
	HealthStatusHealthy HealthStatus = "healthy"
	// HealthStatusUnhealthy represents an unhealthy status.
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	// HealthStatusDegraded represents a degraded status.
	HealthStatusDegraded HealthStatus = "degraded"
	// HealthStatusUnknown represents an unknown status.
	HealthStatusUnknown HealthStatus = "unknown"
)

// HealthCheck represents a health check result.
type HealthCheck struct {
	Status    HealthStatus            `json:"status"`
	Message   string                  `json:"message,omitempty"`
	Details   map[string]interface{}  `json:"details,omitempty"`
	Timestamp time.Time               `json:"timestamp"`
	Checks    map[string]*HealthCheck `json:"checks,omitempty"`
}

// HealthChecker defines the interface for health checking.
type HealthChecker interface {
	// HealthCheck performs a health check and returns the result
	HealthCheck(ctx context.Context) *HealthCheck
}

// ProviderHealth represents the health status of a provider.
type ProviderHealth struct {
	ProviderType ProviderType  `json:"provider_type"`
	Status       HealthStatus  `json:"status"`
	Message      string        `json:"message,omitempty"`
	LastCheck    time.Time     `json:"last_check"`
	ErrorRate    float64       `json:"error_rate"`
	Latency      time.Duration `json:"latency"`
}

// SenderHealth represents the overall health status of the sender.
type SenderHealth struct {
	Status    HealthStatus                     `json:"status"`
	Message   string                           `json:"message,omitempty"`
	Timestamp time.Time                        `json:"timestamp"`
	Providers map[ProviderType]*ProviderHealth `json:"providers"`
	Queue     *HealthCheck                     `json:"queue,omitempty"`
	Metrics   *HealthCheck                     `json:"metrics,omitempty"`
}

// Message is an interface that all specific message types must implement.
// It ensures that messages can validate themselves before being sent.
type Message interface {
	// Validate checks if the message content is valid
	Validate() error
	// ProviderType returns the type of provider this message is intended for.
	ProviderType() ProviderType

	// GetSubProvider returns the sub-provider name (empty when not applicable).
	GetSubProvider() string

	// MsgID returns a unique id for this message (default: uuid, overridable)
	MsgID() string
}

// DefaultMessage provides a base implementation for Message with a unique id.
type DefaultMessage struct {
	msgID  string                 // 可选，允许自定义消息ID，未设置时自动生成
	Extras map[string]interface{} `json:"extras,omitempty"`
}

// LoggerAware is implemented by types that can set a logger.
type LoggerAware interface {
	SetLogger(Logger)
}

// MsgID returns the unique id of the message.
func (m *DefaultMessage) MsgID() string {
	if m.msgID == "" {
		m.msgID = uuid.NewString()
	}
	return m.msgID
}

// ProviderSendOptions defines per-request parameters for Provider.Send.
type ProviderSendOptions struct {
	HTTPClient *http.Client
}

// Provider is the interface that all providers must implement.
// It returns a detailed SendResult along with error (if any).
type Provider interface {
	Send(ctx context.Context, msg Message, opts *ProviderSendOptions) (*SendResult, error)
	// Name returns the unique name of the provider.
	Name() string
}

// RateLimiter is an interface for controlling the rate of operations.
type RateLimiter interface {
	// Allow checks if an operation is permitted without blocking.
	Allow() bool
	// Wait blocks until an operation is permitted or the context is cancelled.
	Wait(ctx context.Context) error
	// Close shuts down the rate limiter, releasing any resources.
	Close() error
}

// Comparable defines an interface for types that can be compared.
type Comparable[T any] interface {
	// Compare returns true if the current item should come before 'other' (higher priority),
	// and false otherwise (lower priority or equal).
	Compare(other T) bool
}

// Schedulable is an optional interface for items that can be scheduled for a future time.
type Schedulable interface {
	// SetScheduledAt sets the time when the item should be processed.
	SetScheduledAt(t time.Time)
	// GetScheduledAt returns the scheduled processing time.
	GetScheduledAt() *time.Time
}

// QueueItem represents an item to be processed within a notification queue.
type QueueItem struct {
	ID          string // Unique message id, recommended to use message.MsgID()
	Provider    string
	Message     Message
	Priority    int
	ScheduledAt *time.Time
	Metadata    map[string]interface{}
	CreatedAt   time.Time
	// Callback is executed after message processing (success or failure)
	Callback func(*SendResult, error)
}

// Compare determines the priority order for QueueItem.
// A smaller Priority value indicates higher priority.
// If priorities are equal, items created earlier have higher precedence.
func (q *QueueItem) Compare(other *QueueItem) bool {
	if other == nil {
		return false
	}
	if q.Priority == other.Priority {
		return q.CreatedAt.Before(other.CreatedAt)
	}
	return q.Priority < other.Priority
}

// SetScheduledAt implements the Schedulable interface for QueueItem.
func (q *QueueItem) SetScheduledAt(t time.Time) {
	q.ScheduledAt = &t
}

// GetScheduledAt implements the Schedulable interface for QueueItem.
func (q *QueueItem) GetScheduledAt() *time.Time {
	return q.ScheduledAt
}

// Queue is an interface for a message queuing system.
type Queue interface {
	// Enqueue adds an item to the queue for immediate processing.
	Enqueue(ctx context.Context, item *QueueItem) error
	// EnqueueDelayed adds an item to the queue to be processed after a specified delay.
	EnqueueDelayed(ctx context.Context, item *QueueItem, delay time.Duration) error
	// Dequeue retrieves an item from the queue for processing.
	Dequeue(ctx context.Context) (*QueueItem, error)
	// Size returns the current number of items in the queue.
	Size() int
	// Close shuts down the queue, releasing any resources.
	Close() error
}

// PerformanceMetrics represents detailed performance metrics.
type PerformanceMetrics struct {
	SendLatency         time.Duration `json:"send_latency"`
	QueueLatency        time.Duration `json:"queue_latency,omitempty"`
	RetryCount          int           `json:"retry_count,omitempty"`
	ErrorRate           float64       `json:"error_rate"`
	Throughput          float64       `json:"throughput"` // messages per second
	QueueSize           int           `json:"queue_size,omitempty"`
	CircuitBreakerState string        `json:"circuit_breaker_state,omitempty"`
	RateLimitRemaining  int           `json:"rate_limit_remaining,omitempty"`
}

// MetricsData represents structured metrics data.
type MetricsData struct {
	Provider     string                 `json:"provider"`
	Success      bool                   `json:"success"`
	Duration     time.Duration          `json:"duration"`
	Operation    string                 `json:"operation,omitempty"`
	ErrorType    string                 `json:"error_type,omitempty"`
	RetryCount   int                    `json:"retry_count,omitempty"`
	QueueSize    int                    `json:"queue_size,omitempty"`
	QueueLatency time.Duration          `json:"queue_latency,omitempty"`
	Performance  *PerformanceMetrics    `json:"performance,omitempty"`
	Custom       map[string]interface{} `json:"custom,omitempty"`
}

// MetricsCollector is an interface for collecting performance and outcome metrics.
type MetricsCollector interface {
	// RecordSendResult records the result of a send operation.
	RecordSendResult(data MetricsData)
}

// CircuitBreaker is an interface for implementing the Circuit Breaker pattern.
type CircuitBreaker interface {
	// Execute attempts to run a function, applying circuit breaker logic.
	// It returns an error if the circuit is open or the function execution fails.
	Execute(ctx context.Context, fn func() error) error
	// Close shuts down the circuit breaker.
	Close() error
}

// RetryFilter is a function type that determines whether an error should trigger a retry.
type RetryFilter func(attempt int, err error) bool

// SendOptionsSerializer defines the interface for serializing/deserializing SendOptions.
type SendOptionsSerializer interface {
	Serialize(opts *SendOptions) ([]byte, error)
	Deserialize(data []byte) (*SendOptions, error)
}

// DefaultSendOptionsSerializer is the default implementation of SendOptionsSerializer.
type DefaultSendOptionsSerializer struct{}

// Serialize serializes SendOptions to JSON bytes.
func (s *DefaultSendOptionsSerializer) Serialize(opts *SendOptions) ([]byte, error) {
	if opts == nil {
		return nil, errors.New("send options cannot be nil")
	}

	data := sendOptionsData{
		Priority:              opts.Priority,
		Timeout:               int64(opts.Timeout),
		DisableCircuitBreaker: opts.DisableCircuitBreaker,
		DisableRateLimiter:    opts.DisableRateLimiter,
		Metadata:              opts.Metadata,
		AccountName:           opts.AccountName,
		StrategyName:          opts.StrategyName,
	}

	// Convert RetryPolicy to serializable format if present
	if opts.RetryPolicy != nil {
		data.RetryPolicy = &serializableRetryPolicy{
			MaxAttempts:   opts.RetryPolicy.MaxAttempts,
			InitialDelay:  int64(opts.RetryPolicy.InitialDelay),
			MaxDelay:      int64(opts.RetryPolicy.MaxDelay),
			BackoffFactor: opts.RetryPolicy.BackoffFactor,
		}
	}

	return json.Marshal(data)
}

// Deserialize deserializes JSON bytes to SendOptions.
func (s *DefaultSendOptionsSerializer) Deserialize(data []byte) (*SendOptions, error) {
	if len(data) == 0 {
		return &SendOptions{}, nil
	}

	var dataStruct sendOptionsData
	if err := json.Unmarshal(data, &dataStruct); err != nil {
		return nil, fmt.Errorf("failed to deserialize send options: %w", err)
	}

	opts := &SendOptions{
		Priority:              dataStruct.Priority,
		Timeout:               time.Duration(dataStruct.Timeout),
		DisableCircuitBreaker: dataStruct.DisableCircuitBreaker,
		DisableRateLimiter:    dataStruct.DisableRateLimiter,
		Metadata:              dataStruct.Metadata,
		AccountName:           dataStruct.AccountName,
		StrategyName:          dataStruct.StrategyName,
	}

	// Convert serializable RetryPolicy back to RetryPolicy if present
	if dataStruct.RetryPolicy != nil {
		opts.RetryPolicy = NewRetryPolicy(
			WithRetryMaxAttempts(dataStruct.RetryPolicy.MaxAttempts),
			WithRetryInitialDelay(time.Duration(dataStruct.RetryPolicy.InitialDelay)),
			WithRetryMaxDelay(time.Duration(dataStruct.RetryPolicy.MaxDelay)),
			WithRetryBackoffFactor(dataStruct.RetryPolicy.BackoffFactor),
			WithRetryFilter(DefaultRetryFilter([]error{}, true)),
		)
	}

	return opts, nil
}

// sendOptionsData represents the serializable data structure for SendOptions.
type sendOptionsData struct {
	Priority              int                    `json:"priority"`
	Timeout               int64                  `json:"timeout_ns"`
	DisableCircuitBreaker bool                   `json:"disable_circuit_breaker"`
	DisableRateLimiter    bool                   `json:"disable_rate_limiter"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
	AccountName           string                 `json:"account_name,omitempty"`
	StrategyName          string                 `json:"strategy_name,omitempty"`
	// Serializable retry policy (without Filter function)
	RetryPolicy *serializableRetryPolicy `json:"retry_policy,omitempty"`
}

// serializableRetryPolicy represents a serializable version of RetryPolicy.
type serializableRetryPolicy struct {
	MaxAttempts   int     `json:"max_attempts"`
	InitialDelay  int64   `json:"initial_delay_ns"`
	MaxDelay      int64   `json:"max_delay_ns"`
	BackoffFactor float64 `json:"backoff_factor"`
	// Note: Filter function cannot be serialized, will use default filter on deserialization
}

// ConfigProvider defines the interface for configuration providers.
type ConfigProvider interface {
	GetStrategy() StrategyType
}

const (
	// OperationEnqueue represents the enqueue operation.
	OperationEnqueue = "enqueue"
	// OperationDequeue represents the dequeue operation.
	OperationDequeue = "dequeue"
	OperationSent    = "sent"
)

// GetExtraString retrieves a string value from extras.
func (m *DefaultMessage) GetExtraString(key string) (string, bool) {
	if m.Extras == nil {
		return "", false
	}
	if value, ok := m.Extras[key]; ok {
		if str, okStr := value.(string); okStr {
			return str, true
		}
	}
	return "", false
}

// GetExtraStringOrDefault retrieves a string value from extras with a default fallback.
func (m *DefaultMessage) GetExtraStringOrDefault(key, defaultValue string) string {
	if value, ok := m.GetExtraString(key); ok && value != "" {
		return value
	}
	return defaultValue
}

// GetExtraInt retrieves an integer value from extras.
func (m *DefaultMessage) GetExtraInt(key string) (int, bool) {
	if m.Extras == nil {
		return 0, false
	}
	if value, ok := m.Extras[key]; ok {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i, true
			}
		}
	}
	return 0, false
}

// GetExtraIntOrDefault retrieves an integer value from extras with a default fallback.
func (m *DefaultMessage) GetExtraIntOrDefault(key string, defaultValue int) int {
	if value, ok := m.GetExtraInt(key); ok {
		return value
	}
	return defaultValue
}

// GetExtraBool retrieves a boolean value from extras.
func (m *DefaultMessage) GetExtraBool(key string) (bool, bool) {
	if m.Extras == nil {
		return false, false
	}
	if value, ok := m.Extras[key]; ok {
		if b, okBool := value.(bool); okBool {
			return b, true
		}
	}
	return false, false
}

// GetExtraBoolOrDefault retrieves a boolean value from extras with a default fallback.
func (m *DefaultMessage) GetExtraBoolOrDefault(key string, defaultValue bool) bool {
	if value, ok := m.GetExtraBool(key); ok {
		return value
	}
	return defaultValue
}

// GetExtraFloat retrieves a float64 value from extras.
func (m *DefaultMessage) GetExtraFloat(key string) (float64, bool) {
	if m.Extras == nil {
		return 0, false
	}
	if value, ok := m.Extras[key]; ok {
		switch v := value.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

// GetExtraFloatOrDefault retrieves a float64 value from extras with a default fallback.
func (m *DefaultMessage) GetExtraFloatOrDefault(key string, defaultValue float64) float64 {
	if value, ok := m.GetExtraFloat(key); ok {
		return value
	}
	return defaultValue
}

func (m *DefaultMessage) GetSubProvider() string {
	return ""
}

// SendResult represents the result of a send operation.
type SendResult struct {
	Config     interface{} // 发送时的配置信息
	StatusCode int         // HTTP状态码
	Headers    http.Header // 响应头
	Body       []byte      // 响应体
}

package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	// DefaultSendTimeout is the default timeout for send operations
	DefaultSendTimeout = 30 * time.Second
)

// ProviderDecorator is a decorated Provider that includes middleware for various concerns.
type ProviderDecorator struct {
	Provider
	middleware *SenderMiddleware
	workers    sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	logger     Logger
	tracer     trace.Tracer
	meter      metric.Meter
}

// Global callbackRegistry, only used for local/in-memory queue/testing scenarios
var callbackRegistry sync.Map // map[msgID]func(error)

// NewProviderDecorator creates a new ProviderDecorator instance.
func NewProviderDecorator(provider Provider, middleware *SenderMiddleware, logger Logger) *ProviderDecorator {
	ctx, cancel := context.WithCancel(context.Background())

	pd := &ProviderDecorator{
		Provider:   provider,
		middleware: middleware,
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		tracer:     otel.Tracer("provider-decorator"),
		meter:      otel.Meter("provider-decorator"),
	}

	// Check if provider supports logger injection
	if loggerAware, ok := provider.(LoggerAwareProvider); ok {
		loggerAware.SetLogger(pd.logger)
	}

	// Start the queue processor if a queue is configured.
	if middleware != nil && middleware.Queue != nil {
		pd.startQueueProcessor()
	}

	return pd
}

// Send applies middleware in a layered fashion.
func (pd *ProviderDecorator) Send(ctx context.Context, message Message, opts ...SendOption) error {
	sendOpts := &SendOptions{}
	for _, opt := range opts {
		opt(sendOpts)
	}

	if sendOpts.Async {
		// Enqueue chain: only enqueue and record metrics
		err := pd.sendAsync(ctx, message, sendOpts)
		pd.recordEnqueue(ctx, message, err)
		return err
	}

	// Synchronous chain: rate limiting -> circuit breaker -> retry -> send -> metrics
	return pd.executeWithMiddleware(ctx, message, sendOpts)
}

// Unified logic for consumer and synchronous chains
func (pd *ProviderDecorator) executeWithMiddleware(ctx context.Context, message Message, opts *SendOptions) error {
	startTime := time.Now()
	ctx, span := pd.tracer.Start(ctx, "provider.send",
		trace.WithAttributes(attribute.String("provider", pd.Provider.Name())),
	)
	defer span.End()

	pd.logger.Log(LevelDebug, "message", "provider send start", "message_id", message.MsgID())

	// Rate limiting
	if pd.middleware != nil && pd.middleware.RateLimiter != nil && !opts.DisableRateLimiter {
		if !pd.middleware.RateLimiter.Allow() {
			err := NewSenderError(ErrCodeRateLimitExceeded, "rate limit exceeded", nil)
			pd.afterSend(ctx, message, err, startTime)
			return err
		}
	}
	// Circuit breaker
	if pd.middleware != nil && pd.middleware.CircuitBreaker != nil && !opts.DisableCircuitBreaker {
		return pd.middleware.CircuitBreaker.Execute(ctx, func() error {
			return pd.doSendWithRetry(ctx, message, opts, startTime)
		})
	}
	return pd.doSendWithRetry(ctx, message, opts, startTime)
}

func (pd *ProviderDecorator) doSendWithRetry(ctx context.Context, message Message, opts *SendOptions, startTime time.Time) error {
	var err error
	if pd.middleware != nil && pd.middleware.Retry != nil {
		err = pd.sendWithRetry(ctx, message, opts)
	} else {
		err = pd.executeSend(ctx, message, opts)
	}
	pd.afterSend(ctx, message, err, startTime)
	return err
}

// Observability: record metrics and logs after sending
func (pd *ProviderDecorator) afterSend(ctx context.Context, message Message, err error, start time.Time) {
	if pd.middleware != nil && pd.middleware.Metrics != nil {
		pd.middleware.Metrics.RecordSendResult(MetricsData{
			Provider: string(pd.Provider.Name()),
			Success:  err == nil,
			Duration: time.Since(start),
		})
	}
	pd.logger.Log(LevelInfo, "message", "provider send end", "message_id", message.MsgID(), "success", fmt.Sprintf("%v", err == nil), "duration", time.Since(start))
}

// Observability: record metrics and logs when enqueueing
func (pd *ProviderDecorator) recordEnqueue(ctx context.Context, message Message, err error) {
	if pd.middleware != nil && pd.middleware.Metrics != nil {
		pd.middleware.Metrics.RecordSendResult(MetricsData{
			Provider:  string(pd.Provider.Name()),
			Success:   err == nil,
			Duration:  0,
			Operation: "enqueue",
		})
	}
	pd.logger.Log(LevelInfo, "message", "provider enqueue", "message_id", message.MsgID(), "success", fmt.Sprintf("%v", err == nil), "error", fmt.Sprintf("%v", err))
}

// The queue consumer worker uses the same logic as the synchronous chain
func (pd *ProviderDecorator) processQueueItem(ctx context.Context, item *QueueItem) {
	// Deserialize SendOptions and restore context
	restoredCtx, opts, err := deserializeSendOptions(ctx, item.Metadata)
	if err != nil {
		pd.logger.Log(LevelWarn, "message", "deserialize send options failed", "error", err.Error())
		opts = &SendOptions{} // fallback
		restoredCtx = ctx
	}
	err = pd.executeWithMiddleware(restoredCtx, item.Message, opts)
	// Find and invoke callback (only effective for local/in-memory queue)
	if cb, ok := callbackRegistry.LoadAndDelete(item.Message.MsgID()); ok {
		if callback, ok := cb.(func(error)); ok {
			callback(err)
		}
	}
}

// sendAsync sends the message asynchronously, using a queue if available, otherwise a goroutine.
func (pd *ProviderDecorator) sendAsync(ctx context.Context, message Message, opts *SendOptions) error {
	metadata, err := serializeSendOptions(ctx, opts, opts.Metadata)
	if err != nil {
		return fmt.Errorf("failed to serialize SendOptions for queue: %w", err)
	}

	item := &QueueItem{
		ID:        message.MsgID(),
		Provider:  pd.Provider.Name(),
		Message:   message,
		Priority:  opts.Priority,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	if opts.Callback != nil {
		callbackRegistry.Store(message.MsgID(), opts.Callback)
	}

	if pd.middleware != nil && pd.middleware.Queue != nil {
		return pd.middleware.Queue.Enqueue(ctx, item)
	}

	// Fallback to goroutine if no queue is configured
	go func() {
		err := pd.executeSend(context.Background(), message, opts)
		if opts.Callback != nil {
			opts.Callback(err)
		}
		if err != nil {
			pd.logger.Log(LevelError, "message", "async send failed", "message_id", message.MsgID(), "error", fmt.Sprintf("%v", err))
		}
	}()

	return nil
}

// sendWithRetry attempts to send the message with retry logic.
func (pd *ProviderDecorator) sendWithRetry(ctx context.Context, message Message, opts *SendOptions) error {
	retryPolicy := opts.RetryPolicy
	if retryPolicy == nil {
		retryPolicy = pd.middleware.Retry
	}

	if retryPolicy == nil {
		return pd.executeSend(ctx, message, opts)
	}

	var lastErr error
	for attempt := 0; attempt <= retryPolicy.MaxAttempts; attempt++ {
		err := pd.executeSend(ctx, message, opts)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry based on the retry filter
		if retryPolicy.Filter != nil && !retryPolicy.Filter(attempt, err) {
			pd.logger.Log(LevelWarn, "message", "retry filtered", "attempt", attempt, "error", err.Error())
			break
		}

		// If this is the last attempt, don't wait
		if attempt == retryPolicy.MaxAttempts {
			break
		}

		// Wait before retrying using NextDelay method
		delay := retryPolicy.NextDelay(attempt, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			continue
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", retryPolicy.MaxAttempts+1, lastErr)
}

// executeSend executes the actual send operation with timeout and metrics.
func (pd *ProviderDecorator) executeSend(ctx context.Context, message Message, opts *SendOptions) error {
	var (
		err     error
		timeout = opts.Timeout
	)

	if timeout == 0 {
		timeout = DefaultSendTimeout
	}

	if opts.Metadata != nil {
		ctx = WithCtxSendMetadata(ctx, opts.Metadata)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute the actual send operation
	err = pd.Provider.Send(ctx, message)

	return err
}

// startQueueProcessor starts a goroutine to continuously dequeue and process messages.
func (pd *ProviderDecorator) startQueueProcessor() {
	if pd.middleware == nil || pd.middleware.Queue == nil {
		pd.logger.Log(LevelInfo, "message", "queue processor not started", "message", "Queue processor not started: queue is not configured.")
		return
	}

	pd.workers.Add(1)
	go func() {
		defer pd.workers.Done()
		pd.logger.Log(LevelInfo, "message", "queue processor started", "message", "Queue processor started")

		for {
			select {
			case <-pd.ctx.Done():
				pd.logger.Log(LevelInfo, "message", "queue processor shutting down", "message", "Queue processor shutting down")
				return
			default:
				item, err := pd.middleware.Queue.Dequeue(pd.ctx)
				if err != nil {
					if err == context.Canceled {
						pd.logger.Log(LevelInfo, "message", "queue dequeue cancelled", "message", "Queue dequeue cancelled")
						return
					}
					pd.logger.Log(LevelError, "message", "queue dequeue failed", "message", "Queue dequeue failed", "error", fmt.Sprintf("%v", err))
					time.Sleep(time.Second)
					continue
				}

				// Process the item
				pd.processQueueItem(pd.ctx, item)
			}
		}
	}()
}

// Close gracefully shuts down the ProviderDecorator and its associated middleware.
func (pd *ProviderDecorator) Close() error {
	// Signal all goroutines to shut down.
	if pd.cancel != nil {
		pd.cancel()
	}

	// Wait for all workers to finish
	pd.workers.Wait()

	pd.logger.Log(LevelInfo, "message", "All provider workers shut down")

	// Close the underlying provider if it implements io.Closer
	if closer, ok := pd.Provider.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("failed to close underlying provider: %w", err)
		}
	}

	pd.logger.Log(LevelInfo, "message", "ProviderDecorator closed successfully")
	return nil
}

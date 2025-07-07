package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// DefaultSendTimeout is the default timeout for send operations.
	DefaultSendTimeout = 30 * time.Second
	queueBackoff       = 10 * time.Millisecond
)

// ProviderDecorator is a decorated Provider that includes middleware for various concerns.
type ProviderDecorator struct {
	Provider

	middleware *SenderMiddleware
	workers    sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	logger     Logger
}

// Global callbackRegistry, only used for local/in-memory queue/testing scenarios.
//
//nolint:gochecknoglobals // Reason: callbackRegistry is a global callback registry for async message handling
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
	}

	// Check if provider supports logger injection
	if loggerAware, ok := provider.(LoggerAware); ok {
		loggerAware.SetLogger(logger)
	}

	// Start the queue processor if a queue is configured.
	if middleware != nil && middleware.Queue != nil {
		pd.startQueueProcessor()
	}

	return pd
}

// Send applies middleware in a layered fashion.
func (pd *ProviderDecorator) Send(ctx context.Context, message Message, opts ...SendOption) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	sendOpts := &SendOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(sendOpts)
		}
	}

	if sendOpts.Async {
		// Enqueue chain: only enqueue and record metrics
		err := pd.sendAsync(ctx, message, sendOpts)
		pd.recordMetric(OperationEnqueue, message, err == nil, 0, 0)
		return err
	}

	// Synchronous chain: rate limiting -> circuit breaker -> retry -> send -> metrics
	startTime := time.Now()
	err := pd.executeWithMiddleware(ctx, message, sendOpts)
	pd.recordMetric(OperationSent, message, err == nil, time.Since(startTime), 0)
	return err
}

// Unified logic for consumer and synchronous chains.
func (pd *ProviderDecorator) executeWithMiddleware(ctx context.Context, message Message, opts *SendOptions) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if pd.logger != nil {
		_ = pd.logger.Log(
			LevelDebug,
			"message",
			"provider send start",
			"message_id",
			message.MsgID(),
			"provider",
			pd.Provider.Name(),
		)
	}

	// Rate limiting
	if pd.middleware != nil && pd.middleware.RateLimiter != nil && !opts.DisableRateLimiter {
		if !pd.middleware.RateLimiter.Allow() {
			return NewSenderError(ErrCodeRateLimitExceeded, "rate limit exceeded", nil)
		}
	}
	// Circuit breaker
	if pd.middleware != nil && pd.middleware.CircuitBreaker != nil && !opts.DisableCircuitBreaker {
		return pd.middleware.CircuitBreaker.Execute(ctx, func() error {
			return pd.doSendWithRetry(ctx, message, opts)
		})
	}
	return pd.doSendWithRetry(ctx, message, opts)
}

func (pd *ProviderDecorator) doSendWithRetry(ctx context.Context, message Message, opts *SendOptions) error {
	var err error
	if pd.middleware != nil && pd.middleware.Retry != nil {
		err = pd.sendWithRetry(ctx, message, opts)
	} else {
		err = pd.executeSend(ctx, message, opts)
	}
	return err
}

// The queue consumer worker uses the same logic as the synchronous chain.
func (pd *ProviderDecorator) processQueueItem(ctx context.Context, item *QueueItem) {
	queueLatency := time.Duration(0)
	if !item.CreatedAt.IsZero() {
		queueLatency = time.Since(item.CreatedAt)
	}
	pd.recordMetric(OperationDequeue, item.Message, true, 0, queueLatency)
	// Deserialize SendOptions and restore context
	restoredCtx, opts, err := deserializeSendOptions(ctx, item.Metadata)
	if err != nil {
		if pd.logger != nil {
			_ = pd.logger.Log(LevelWarn, "message", "deserialize send options failed", "error", err.Error())
		}
		opts = &SendOptions{} // fallback
		restoredCtx = ctx
	}

	// If ScheduledAt is set and is in the future, wait until that time
	if item.ScheduledAt != nil && item.ScheduledAt.After(time.Now()) {
		pd.logInfo(
			fmt.Sprintf(
				"Message %s scheduled for future processing, waiting until %s",
				item.ID,
				item.ScheduledAt.Format(time.RFC3339),
			),
		)
		select {
		case <-time.After(time.Until(*item.ScheduledAt)):
			// Waited successfully, continue processing
		case <-ctx.Done():
			// Context cancelled while waiting, do not process
			pd.logWarn(fmt.Sprintf("Message %s processing cancelled during scheduled wait: %v", item.ID, ctx.Err()))
			return
		}
	}

	err = pd.executeWithMiddleware(restoredCtx, item.Message, opts)
	// Find and invoke callback (only effective for local/in-memory queue)
	if cb, ok := callbackRegistry.LoadAndDelete(item.Message.MsgID()); ok {
		if callback, okCallback := cb.(func(error)); okCallback {
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
		ID:          message.MsgID(),
		Provider:    pd.Provider.Name(),
		Message:     message,
		Priority:    opts.Priority,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		ScheduledAt: opts.DelayUntil,
	}

	if opts.Callback != nil {
		callbackRegistry.Store(message.MsgID(), opts.Callback)
	}

	if pd.middleware != nil && pd.middleware.Queue != nil {
		return pd.middleware.Queue.Enqueue(ctx, item)
	}

	// Fallback to goroutine if no queue is configured
	go func() {
		// If DelayUntil is set, wait until that time
		if opts.DelayUntil != nil && opts.DelayUntil.After(time.Now()) {
			select {
			case <-time.After(time.Until(*opts.DelayUntil)):
				// Waited successfully
			case <-context.Background().Done(): // Use background context to prevent premature cancellation
				// Context cancelled while waiting, do not send
				if opts.Callback != nil {
					opts.Callback(context.Background().Err())
				}
				return
			}
		}
		errSend := pd.executeSend(context.Background(), message, opts)
		if opts.Callback != nil {
			opts.Callback(errSend)
		}
		if errSend != nil && pd.logger != nil {
			_ = pd.logger.Log(
				LevelError,
				"message",
				"async send failed",
				"message_id",
				message.MsgID(),
				"error",
				fmt.Sprintf("%v", errSend),
			)
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
		if retryPolicy.Filter == nil || !retryPolicy.Filter(attempt, err) {
			break
		}
		if pd.logger != nil {
			_ = pd.logger.Log(LevelWarn, "message", "retry filtered", "attempt", attempt, "error", err.Error())
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

	// Convert SendOptions to ProviderSendOptions
	providerOpts := &ProviderSendOptions{
		HTTPClient: EnsureHTTPClient(opts.HTTPClient),
	}

	// Execute the actual send operation
	err = pd.Provider.Send(ctx, message, providerOpts)

	return err
}

// startQueueProcessor starts a goroutine to continuously dequeue and process messages.
func (pd *ProviderDecorator) startQueueProcessor() {
	if pd.middleware == nil || pd.middleware.Queue == nil {
		pd.logInfo("Queue processor not started: queue is not configured.")
		return
	}

	pd.workers.Add(1)
	go pd.queueProcessorLoop()
}

func (pd *ProviderDecorator) queueProcessorLoop() {
	defer pd.workers.Done()
	pd.logInfo("Queue processor started")

	for {
		select {
		case <-pd.ctx.Done():
			pd.logInfo("Queue processor shutting down")
			return
		default:
			item, err := pd.middleware.Queue.Dequeue(pd.ctx)
			if err != nil {
				if pd.handleDequeueError(err) {
					return
				}
				continue
			}
			if item == nil {
				// Queue is empty; back off briefly to avoid busy loop
				time.Sleep(queueBackoff)
				continue
			}
			pd.processQueueItem(pd.ctx, item)
		}
	}
}

func (pd *ProviderDecorator) handleDequeueError(err error) bool {
	if errors.Is(err, context.Canceled) {
		pd.logInfo("Queue dequeue cancelled")
		return true
	}
	pd.logError("Queue dequeue failed", err)
	time.Sleep(time.Second)
	return false
}

func (pd *ProviderDecorator) logInfo(msg string) {
	if pd.logger != nil {
		_ = pd.logger.Log(LevelInfo, "message", msg)
	}
}

func (pd *ProviderDecorator) logWarn(msg string) {
	if pd.logger != nil {
		_ = pd.logger.Log(LevelWarn, "message", msg)
	}
}

func (pd *ProviderDecorator) logError(msg string, err error) {
	if pd.logger != nil {
		_ = pd.logger.Log(LevelError, "message", msg, "error", fmt.Sprintf("%v", err))
	}
}

// Close gracefully shuts down the ProviderDecorator and its associated middleware.
func (pd *ProviderDecorator) Close() error {
	// Signal all goroutines to shut down.
	if pd.cancel != nil {
		pd.cancel()
	}

	// Wait for all workers to finish
	pd.workers.Wait()

	if pd.logger != nil {
		_ = pd.logger.Log(LevelInfo, "message", "All provider workers shut down")
	}

	// Close the underlying provider if it implements io.Closer
	if closer, ok := pd.Provider.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("failed to close underlying provider: %w", err)
		}
	}

	if pd.logger != nil {
		_ = pd.logger.Log(LevelInfo, "message", "ProviderDecorator closed successfully")
	}
	return nil
}

func (pd *ProviderDecorator) recordMetric(
	operation string,
	_ Message,
	success bool,
	duration time.Duration,
	queueLatency time.Duration,
) {
	queueSize := 0
	if pd.middleware != nil && pd.middleware.Queue != nil {
		queueSize = pd.middleware.Queue.Size()
	}
	if pd.middleware != nil && pd.middleware.Metrics != nil {
		pd.middleware.Metrics.RecordSendResult(MetricsData{
			Provider:     pd.Provider.Name(),
			Success:      success,
			Duration:     duration,
			Operation:    operation,
			QueueLatency: queueLatency,
			QueueSize:    queueSize,
		})
	}
}

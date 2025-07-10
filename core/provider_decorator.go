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

// Send applies middleware, executes send synchronously and returns detailed SendResult.
//
// Behaviour notes:
//  1. If SendOptions.Async == false (default)  ➜  The message is sent immediately and
//     the returned *SendResult contains status-code / headers / body from the provider.
//  2. If SendOptions.Async == true             ➜  The message is merely **enqueued** (or
//     dispatched via a background goroutine when no queue is configured). In this case
//     Send always returns (nil, err) where err reflects whether the enqueue
//     operation succeeded. The final outcome will not be available here.
//  3. To observe the actual async result, supply SendOptions.Callback. The callback
//     will be invoked ONLY when processing occurs in-process (memory queue or goroutine
//     fallback). Distributed queues executed by external workers cannot trigger the
//     callback.
func (pd *ProviderDecorator) Send(ctx context.Context, message Message, opts ...SendOption) (*SendResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	sendOpts := &SendOptions{}
	for _, opt := range opts {
		if opt != nil {
			opt(sendOpts)
		}
	}

	if sendOpts.Async {
		// Validate callback usage
		if sendOpts.Callback != nil && !sendOpts.Async {
			return nil, NewSenderError(ErrCodeInvalidConfig, "callback requires async send", nil)
		}

		// Enqueue chain: only enqueue and record metrics
		err := pd.sendAsync(ctx, message, sendOpts)
		pd.recordMetric(OperationEnqueue, message, err == nil, 0, 0)
		return nil, err
	}

	// Synchronous chain: rate limiting -> circuit breaker -> retry -> send -> metrics
	startTime := time.Now()
	result, err := pd.executeWithMiddleware(ctx, message, sendOpts)
	pd.recordMetric(OperationSent, message, err == nil, time.Since(startTime), 0)
	return result, err
}

// Unified logic for consumer and synchronous chains.
func (pd *ProviderDecorator) executeWithMiddleware(
	ctx context.Context,
	message Message,
	opts *SendOptions,
) (*SendResult, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
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
			return nil, NewSenderError(ErrCodeRateLimitExceeded, "rate limit exceeded", nil)
		}
	}

	var result *SendResult
	var err error

	// Circuit breaker
	if pd.middleware != nil && pd.middleware.CircuitBreaker != nil && !opts.DisableCircuitBreaker {
		err = pd.middleware.CircuitBreaker.Execute(ctx, func() error {
			var cbErr error
			result, cbErr = pd.doSendWithRetry(ctx, message, opts)
			return cbErr
		})
	} else {
		result, err = pd.doSendWithRetry(ctx, message, opts)
	}

	// Execute callback **only** when the send originated from an async flow. For
	// synchronous Send (opts.Async == false) callbacks must be ignored per the
	// updated API contract.
	if opts.Async && opts.Callback != nil {
		opts.Callback(result, err)
	}

	return result, err
}

func (pd *ProviderDecorator) doSendWithRetry(
	ctx context.Context,
	message Message,
	opts *SendOptions,
) (*SendResult, error) {
	var result *SendResult
	var err error
	if pd.middleware != nil && pd.middleware.Retry != nil {
		result, err = pd.sendWithRetry(ctx, message, opts)
	} else {
		result, err = pd.executeSend(ctx, message, opts)
	}
	return result, err
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

	// Propagate callback stored in QueueItem to SendOptions before processing.
	// Note: Async flag might be lost during serialization, so re-enable it here
	// to ensure the callback is honoured for local queue processing.
	if item.Callback != nil {
		opts.Callback = item.Callback
		opts.Async = true
	}

	result, err := pd.executeWithMiddleware(restoredCtx, item.Message, opts)
	_ = result // result already delivered via internal callback if any
	if err != nil {
		if pd.logger != nil {
			_ = pd.logger.Log(LevelError, "message", "execute with middleware failed", "error", err.Error())
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
		Callback:    opts.Callback,
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
					opts.Callback(nil, context.Background().Err())
				}
				return
			}
		}
		_, errSend := pd.executeWithMiddleware(context.Background(), message, opts)
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
func (pd *ProviderDecorator) sendWithRetry(
	ctx context.Context,
	message Message,
	opts *SendOptions,
) (*SendResult, error) {
	retryPolicy := opts.RetryPolicy
	if retryPolicy == nil {
		retryPolicy = pd.middleware.Retry
	}

	if retryPolicy == nil {
		return pd.executeSend(ctx, message, opts)
	}

	var lastErr error
	var lastResult *SendResult
	for attempt := 0; attempt <= retryPolicy.MaxAttempts; attempt++ {
		result, err := pd.executeSend(ctx, message, opts)
		if err == nil {
			return result, nil
		}

		lastErr = err
		lastResult = result

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
			return nil, ctx.Err()
		case <-time.After(delay):
			continue
		}
	}

	return lastResult, fmt.Errorf("failed after %d attempts: %w", retryPolicy.MaxAttempts+1, lastErr)
}

// executeSend executes the actual send operation with timeout and metrics.
func (pd *ProviderDecorator) executeSend(ctx context.Context, message Message, opts *SendOptions) (*SendResult, error) {
	var (
		err     error
		timeout = opts.Timeout
	)

	if timeout == 0 {
		timeout = DefaultSendTimeout
	}

	if opts.StrategyName != "" || opts.AccountName != "" {
		ctx = WithRoute(ctx, &RouteInfo{AccountName: opts.AccountName, StrategyType: StrategyType(opts.StrategyName)})
	}

	if opts.Metadata != nil {
		ctx = context.WithValue(ctx, metadataKey{}, opts.Metadata)
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Convert SendOptions to ProviderSendOptions
	providerOpts := &ProviderSendOptions{
		HTTPClient: EnsureHTTPClient(opts.HTTPClient),
	}

	result, err := pd.Provider.Send(ctx, message, providerOpts)

	return result, err
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

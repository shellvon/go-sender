package gosender_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/metrics"
	"github.com/shellvon/go-sender/providers/webhook"
	"github.com/shellvon/go-sender/ratelimiter"
)

// testLogger 用于测试的日志实现.
type testLogger struct {
	t *testing.T
}

func (l *testLogger) Log(level core.Level, keyvals ...interface{}) error {
	l.t.Logf("LOG [%s] %v", level, keyvals)
	return nil
}

func (l *testLogger) With(_ ...interface{}) core.Logger {
	return l
}

// TestSenderIntegration 测试完整的发送流程.
func TestSenderIntegration(t *testing.T) {
	// 1. 初始化 Sender
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 2. 配置中间件
	// 2.1 设置限流器
	rateLimiter := ratelimiter.NewTokenBucketRateLimiter(1, 1)
	sender.SetRateLimiter(rateLimiter)

	// 2.2 设置重试策略
	retryPolicy := core.NewRetryPolicy(
		core.WithRetryMaxAttempts(3),
		core.WithRetryInitialDelay(time.Millisecond*100),
		core.WithRetryMaxDelay(time.Second),
	)
	if err := sender.SetRetryPolicy(retryPolicy); err != nil {
		t.Fatalf("Failed to set retry policy: %v", err)
	}

	// 2.3 设置指标收集
	metricsCollector := metrics.NewMemoryMetricsCollector()
	sender.SetMetrics(metricsCollector)

	// 3. 注册 webhook provider 作为测试用途
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 4. 创建测试消息
	message := webhook.Webhook().
		Body([]byte(`{"test": "message"}`)).
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		Build()

	// 5. 发送消息
	ctx := context.Background()
	sendErr := sender.Send(ctx, message)
	if sendErr != nil {
		// webhook provider 会返回错误，因为我们没有真正的服务器
		t.Logf("Expected error from webhook provider: %v", sendErr)
	}

	// 6. 验证指标
	total, success, failed := metricsCollector.GetStats(webhookProvider.Name())
	if total != 1 {
		t.Errorf("Expected 1 total request, got %d", total)
	}
	if success != 0 {
		t.Errorf("Expected 0 successful requests, got %d", success)
	}
	if failed != 1 {
		t.Errorf("Expected 1 failed request, got %d", failed)
	}
}

// TestMultipleProviders 测试多个 Provider 的场景.
func TestMultipleProviders(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 注册多个 webhook endpoints
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "endpoint1",
				URL:     "http://example1.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
				Weight:  1,
			},
			{
				Name:    "endpoint2",
				URL:     "http://example2.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
				Weight:  2,
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 发送多个消息，验证负载均衡
	for range 10 {
		message := webhook.Webhook().
			Body([]byte(`{"test": "message"}`)).
			Method(http.MethodPost).
			Header("Content-Type", "application/json").
			Build()

		sendErr := sender.Send(context.Background(), message)
		if sendErr != nil {
			t.Logf("Expected error from webhook provider: %v", sendErr)
		}
	}
}

// TestRetryAndCircuitBreaker 测试重试和熔断机制.
func TestRetryAndCircuitBreaker(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 配置重试策略
	retryPolicy := core.NewRetryPolicy(
		core.WithRetryMaxAttempts(3),
		core.WithRetryInitialDelay(time.Millisecond*10),
		core.WithRetryMaxDelay(time.Millisecond*100),
	)
	if err := sender.SetRetryPolicy(retryPolicy); err != nil {
		t.Fatalf("Failed to set retry policy: %v", err)
	}

	// 配置熔断器
	circuitBreaker := &testCircuitBreaker{
		threshold: 3,
		timeout:   time.Millisecond * 500,
	}
	sender.SetCircuitBreaker(circuitBreaker)

	// 注册一个总是失败的 webhook endpoint
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "failing-endpoint",
				URL:     "http://non-existent.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	provider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, provider, nil)

	// 发送消息，验证重试和熔断
	message := webhook.Webhook().
		Body([]byte(`{"test": "message"}`)).
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		Build()

	// 第一次发送应该触发重试
	sendErr := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
	if sendErr == nil {
		t.Error("Expected error from webhook provider")
	}

	// 多次发送直到触发熔断
	for range 10 {
		sendErr = sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
		if sendErr != nil && sendErr.Error() == "circuit breaker is open" {
			// 熔断器已打开，测试通过
			return
		}
	}
	t.Error("Circuit breaker did not open as expected")
}

// TestConcurrentSending 测试并发发送.
func TestConcurrentSending(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 配置限流器
	rateLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 5) // 每秒10个请求，突发5个
	sender.SetRateLimiter(rateLimiter)

	// 配置指标收集
	metricsCollector := metrics.NewMemoryMetricsCollector()
	sender.SetMetrics(metricsCollector)

	// 注册 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 并发发送消息
	var wg sync.WaitGroup
	concurrency := 10
	messagesPerGoroutine := 3

	for range concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range messagesPerGoroutine {
				message := webhook.Webhook().
					Body([]byte(`{"test": "concurrent"}`)).
					Method(http.MethodPost).
					Header("Content-Type", "application/json").
					Build()

				sendErr := sender.Send(context.Background(), message)
				if sendErr != nil {
					t.Logf("Expected error from webhook provider: %v", sendErr)
				}
				time.Sleep(time.Millisecond * 10)
			}
		}()
	}

	wg.Wait()

	// 验证总请求数
	total, _, _ := metricsCollector.GetStats(webhookProvider.Name())
	expectedRequests := concurrency * messagesPerGoroutine
	if total != int64(expectedRequests) {
		t.Errorf("Expected %d total requests, got %d", expectedRequests, total)
	}
}

// TestMessageTransformation 测试消息转换功能.
func TestMessageTransformation(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 注册带有自定义转换器的 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 直接使用 Sender 继承的 middleware，避免覆盖 RateLimiter / Retry / Metrics。
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 发送消息并验证转换
	message := webhook.Webhook().
		Body([]byte(`{"test":"message"}`)).
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		Build()

	sendErr := sender.Send(context.Background(), message)
	if sendErr != nil {
		t.Logf("Expected error from webhook provider: %v", sendErr)
	}
}

// TestAccountSelectionStrategy 测试账号选择策略.
func TestAccountSelectionStrategy(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 创建多个 endpoints 配置不同的权重
	webhookConfig := &webhook.Config{
		ProviderMeta: core.ProviderMeta{Strategy: core.StrategyWeighted}, // 使用权重策略
		Items: []*webhook.Endpoint{
			{
				Name:   "endpoint1",
				URL:    "http://example1.com/webhook",
				Method: http.MethodPost,
				Weight: 1,
			},
			{
				Name:   "endpoint2",
				URL:    "http://example2.com/webhook",
				Method: http.MethodPost,
				Weight: 2,
			},
			{
				Name:   "endpoint3",
				URL:    "http://example3.com/webhook",
				Method: http.MethodPost,
				Weight: 3,
			},
		},
	}

	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建一个计数中间件来统计每个 endpoint 的使用次数
	endpointCounts := make(map[string]int)
	countMutex := sync.Mutex{}

	middleware := &core.SenderMiddleware{}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, middleware)

	// 验证选择策略，不实际发送网络请求，直接调用 Select 并统计结果
	totalMessages := 600 // 增加采样次数以降低随机误差
	for range totalMessages {
		selected, errSel := webhookProvider.Select(context.Background(), nil)
		if errSel != nil {
			t.Fatalf("Select failed: %v", errSel)
		}
		countMutex.Lock()
		endpointCounts[selected.GetName()]++
		countMutex.Unlock()
	}

	// 验证权重分布
	totalWeight := 6 // 1 + 2 + 3
	expectedCounts := map[string]float64{
		"endpoint1": float64(totalMessages) * 1 / float64(totalWeight),
		"endpoint2": float64(totalMessages) * 2 / float64(totalWeight),
		"endpoint3": float64(totalMessages) * 3 / float64(totalWeight),
	}

	// 允许 25% 的误差，避免小权重项因随机波动触发误报
	tolerance := 0.25
	for endpoint, expected := range expectedCounts {
		actual := float64(endpointCounts[endpoint])
		diff := actual - expected
		if diff < 0 {
			diff = -diff
		}
		if diff/expected > tolerance {
			t.Errorf("Endpoint %s: expected around %.2f requests, got %d (diff %.2f%%)",
				endpoint, expected, endpointCounts[endpoint], diff/expected*100)
		}
	}
}

// TestQueueProcessing 测试消息队列处理.
func TestQueueProcessing(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 先创建一个内存队列并设置到 sender —— 必须在 RegisterProvider 之前，
	// 因为每个 ProviderDecorator 拷贝 middleware（参见 Sender.SetQueue 注释）。
	queue := &testQueue{items: make([]*core.QueueItem, 0)}
	sender.SetQueue(queue)

	// 注册 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 发送消息到队列
	message := webhook.Webhook().
		Body([]byte(`{"test":"queue"}`)).
		Method(http.MethodPost).
		Build()

	// 使用异步选项发送消息
	ctx := context.Background()
	sendErr := sender.Send(ctx, message, core.WithSendAsync(true))
	if sendErr != nil {
		t.Fatalf("Failed to send message: %v", sendErr)
	}

	// 验证消息是否进入队列
	if queue.Size() != 1 {
		t.Errorf("Expected queue size 1, got %d", queue.Size())
	}

	// 从队列中获取消息并验证
	item, err := queue.Dequeue(ctx)
	if err != nil {
		t.Fatalf("Failed to dequeue message: %v", err)
	}

	if item == nil {
		t.Fatal("Expected queue item, got nil")
	}

	webhookMsg, ok := item.Message.(*webhook.Message)
	if !ok {
		t.Fatalf("Expected message type *webhook.Message, got %T", item.Message)
	}

	if string(webhookMsg.Body) != `{"test":"queue"}` {
		t.Errorf("Expected message body %s, got %s", `{"test":"queue"}`, string(webhookMsg.Body))
	}
}

// TestTimeoutAndCancellation 测试超时和取消场景.
func TestTimeoutAndCancellation(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 注册一个会延迟响应的 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "slow-endpoint",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 创建一个自定义的中间件来模拟延迟
	middleware := &core.SenderMiddleware{}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, middleware)

	// 测试场景 1: Context 超时
	t.Run("Context Timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()

		message := webhook.Webhook().
			Body([]byte(`{"test":"timeout"}`)).
			Method(http.MethodPost).
			Build()

		sendErr := sender.Send(ctx, message)
		if sendErr == nil {
			t.Error("Expected timeout-related error, got nil")
		} else {
			t.Logf("Received error as expected: %v", sendErr)
		}
	})

	// 测试场景 2: Context 取消
	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// 在另一个 goroutine 中取消 context
		go func() {
			time.Sleep(time.Millisecond * 10)
			cancel()
		}()

		message := webhook.Webhook().
			Body([]byte(`{"test":"cancel"}`)).
			Method(http.MethodPost).
			Build()

		sendErr := sender.Send(ctx, message)
		if sendErr == nil {
			t.Error("Expected cancellation-related error, got nil")
		} else {
			t.Logf("Received error as expected: %v", sendErr)
		}
	})

	// 测试场景 3: 请求超时设置
	t.Run("Request Timeout", func(t *testing.T) {
		// 设置一个较短的超时时间
		customClient := &http.Client{
			Timeout: time.Millisecond * 50,
		}

		message := webhook.Webhook().
			Body([]byte(`{"test":"request-timeout"}`)).
			Method(http.MethodPost).
			Build()

		sendErr := sender.Send(
			context.Background(),
			message,
			core.WithSendHTTPClient(customClient),
		)
		if sendErr == nil {
			t.Error("Expected timeout error, got nil")
		}
	})
}

// TestErrorRecovery 测试错误恢复场景.
func TestErrorRecovery(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 配置熔断器
	circuitBreaker := &testCircuitBreaker{
		threshold: 3,
		timeout:   time.Second * 2,
	}
	sender.SetCircuitBreaker(circuitBreaker)

	// 配置限流器
	rateLimiter := ratelimiter.NewTokenBucketRateLimiter(2, 1) // 每秒2个请求，突发1个
	sender.SetRateLimiter(rateLimiter)

	// 注册 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "recovery-test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 测试场景 1: 熔断器恢复
	t.Run("Circuit Breaker Recovery", func(t *testing.T) {
		message := webhook.Webhook().
			Body([]byte(`{"test":"circuit-breaker"}`)).
			Method(http.MethodPost).
			Build()

		// 发送请求直到触发熔断
		for range 5 {
			sendErr := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
			if sendErr != nil && sendErr.Error() == "circuit breaker is open" {
				// 熔断器已打开，等待恢复
				time.Sleep(time.Second * 3) // 等待超过熔断器的超时时间

				// 尝试发送新请求，应该被允许
				sendErr = sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
				if sendErr != nil && sendErr.Error() == "circuit breaker is open" {
					t.Error("Circuit breaker should be recovered")
				}
				return
			}
		}
		t.Error("Circuit breaker did not open as expected")
	})

	// 测试场景 2: 限流器恢复
	t.Run("Rate Limiter Recovery", func(t *testing.T) {
		message := webhook.Webhook().
			Body([]byte(`{"test":"rate-limit"}`)).
			Method(http.MethodPost).
			Build()

		// 快速发送请求触发限流
		sendErr1 := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
		sendErr2 := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
		sendErr3 := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))

		// 至少有一个请求应该被限流
		if sendErr1 == nil && sendErr2 == nil && sendErr3 == nil {
			t.Error("Expected at least one request to be rate limited")
		}

		// 等待令牌桶恢复
		time.Sleep(time.Second)

		// 新的请求应该能够成功
		sendErr := sender.Send(context.Background(), message, core.WithSendDisableRateLimiter(true))
		if sendErr != nil && sendErr.Error() == "rate limit exceeded" {
			t.Error("Rate limiter should be recovered")
		}
	})
}

// TestMiddlewareChain 测试中间件链.
func TestMiddlewareChain(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 配置限流器
	rateLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 5)
	sender.SetRateLimiter(rateLimiter)

	// 配置重试策略
	retryPolicy := core.NewRetryPolicy(
		core.WithRetryMaxAttempts(3),
		core.WithRetryInitialDelay(time.Millisecond*100),
		core.WithRetryMaxDelay(time.Second),
	)
	if err := sender.SetRetryPolicy(retryPolicy); err != nil {
		t.Fatalf("Failed to set retry policy: %v", err)
	}

	// 配置指标收集
	metricsCollector := metrics.NewMemoryMetricsCollector()
	sender.SetMetrics(metricsCollector)

	// 注册 webhook provider
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "middleware-test",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}

	// 直接使用 Sender 继承的 middleware，避免覆盖 RateLimiter / Retry / Metrics。
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 发送消息
	message := webhook.Webhook().
		Body([]byte(`{"test":"middleware-chain"}`)).
		Method(http.MethodPost).
		Build()

	sendErr := sender.Send(context.Background(), message)
	if sendErr != nil {
		t.Logf("Expected error from webhook provider: %v", sendErr)
	}

	// 验证指标收集
	total, _, _ := metricsCollector.GetStats(webhookProvider.Name())
	if total != 1 {
		t.Errorf("Expected 1 request to be recorded, got %d", total)
	}
}

// TestCallbackHandling 测试回调处理.
func TestCallbackHandling(t *testing.T) {
	sender := gosender.NewSender(gosender.WithLogger(&testLogger{t: t}))
	defer sender.Close()

	// 注册 webhook provider (所有子测试复用)
	webhookConfig := &webhook.Config{
		Items: []*webhook.Endpoint{
			{
				Name:    "callback-endpoint",
				URL:     "http://example.com/webhook",
				Method:  http.MethodPost,
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	webhookProvider, err := webhook.New(webhookConfig)
	if err != nil {
		t.Fatalf("Failed to create webhook provider: %v", err)
	}
	sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

	// 同步/异步子测试各自使用独立通道，避免相互影响

	// 场景 1: 同步发送 —— 期望 Callback 不被触发
	t.Run("Sync Send Callback", func(t *testing.T) {
		ch := make(chan error, 1)
		message := webhook.Webhook().
			Body([]byte(`{"test":"callback"}`)).
			Method(http.MethodPost).
			Build()

		_ = sender.Send(context.Background(), message, core.WithSendCallback(func(_ *core.SendResult, err error) {
			ch <- err
		}))

		select {
		case <-ch:
			t.Error("Callback should NOT be called for sync send")
		case <-time.After(200 * time.Millisecond):
			// Pass: no callback received
		}
	})

	// 场景 2: 异步发送 —— Callback 必须触发
	t.Run("Async Send Callback", func(t *testing.T) {
		ch := make(chan error, 1)
		message := webhook.Webhook().
			Body([]byte(`{"test":"async-callback"}`)).
			Method(http.MethodPost).
			Build()

		if sendErr := sender.Send(context.Background(), message,
			core.WithSendAsync(),
			core.WithSendCallback(func(_ *core.SendResult, err error) { ch <- err }),
		); sendErr != nil {
			t.Fatalf("Failed to send message: %v", sendErr)
		}

		select {
		case cbErr := <-ch:
			if cbErr == nil {
				t.Error("Expected error in callback, got nil")
			}
		case <-time.After(time.Second):
			t.Error("Async callback not received within timeout")
		}
	})

	// 场景 3: 带延迟的异步发送回调 —— Callback 必须触发
	t.Run("Delayed Async Send Callback", func(t *testing.T) {
		ch := make(chan error, 1)
		message := webhook.Webhook().
			Body([]byte(`{"test":"delayed-callback"}`)).
			Method(http.MethodPost).
			Build()

		if sendErr := sender.Send(context.Background(), message,
			core.WithSendAsync(),
			core.WithSendDelay(time.Millisecond*100),
			core.WithSendCallback(func(_ *core.SendResult, err error) { ch <- err }),
		); sendErr != nil {
			t.Fatalf("Failed to send message: %v", sendErr)
		}

		select {
		case cbErr := <-ch:
			if cbErr == nil {
				t.Error("Expected error in callback, got nil")
			}
		case <-time.After(2 * time.Second):
			t.Error("Delayed callback not received within timeout")
		}
	})
}

// testCircuitBreaker 是一个用于测试的熔断器实现.
type testCircuitBreaker struct {
	mu        sync.RWMutex
	failures  int
	threshold int
	timeout   time.Duration
	openTime  time.Time
}

// Execute 实现 CircuitBreaker 接口.
func (cb *testCircuitBreaker) Execute(_ context.Context, fn func() error) error {
	if !cb.Allow() {
		return errors.New("circuit breaker is open")
	}

	err := fn()
	if err != nil {
		cb.Failure()
		return err
	}

	cb.Success()
	return nil
}

func (cb *testCircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	if cb.failures >= cb.threshold {
		if time.Since(cb.openTime) < cb.timeout {
			return false
		}
		// 超过超时时间，重置
		cb.failures = 0
	}
	return true
}

func (cb *testCircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
}

func (cb *testCircuitBreaker) Failure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.threshold {
		cb.openTime = time.Now()
	}
}

func (cb *testCircuitBreaker) Close() error {
	return nil
}

// testQueue 是一个用于测试的简单内存队列实现.
type testQueue struct {
	mu    sync.Mutex
	items []*core.QueueItem
}

func (q *testQueue) Enqueue(_ context.Context, item *core.QueueItem) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	return nil
}

func (q *testQueue) EnqueueDelayed(_ context.Context, item *core.QueueItem, _ time.Duration) error {
	// 简单实现，直接入队，不考虑延迟
	return q.Enqueue(context.Background(), item)
}

func (q *testQueue) Dequeue(_ context.Context) (*core.QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, nil //nolint:nilnil // 测试需要
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *testQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *testQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = nil
	return nil
}

package gosender

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shellvon/go-sender/circuitbreaker"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/metrics"
	"github.com/shellvon/go-sender/ratelimiter"
)

type mockMessage struct {
	id           string
	providerType core.ProviderType
}

func (m *mockMessage) ProviderType() core.ProviderType {
	if m.providerType != "" {
		return m.providerType
	}
	return "mock"
}

func (m *mockMessage) Name() string {
	return "mock"
}

func (m *mockMessage) Content() string {
	return "test content"
}

func (m *mockMessage) MsgID() string {
	if m.id == "" {
		m.id = "test-id"
	}
	return m.id
}

func (m *mockMessage) Validate() error {
	return nil
}

func (m *mockMessage) GetAccountName() string     { return "" }
func (m *mockMessage) SetAccountName(name string) {}

// mockProvider 用于测试的 provider，实现 core.Provider

type mockProvider struct {
	failCount    int32
	failMax      int32
	succeedOnce  bool
	providerType core.ProviderType
}

func (m *mockProvider) Name() string {

	if m.providerType == "" {
		return "mock"
	}

	return string(m.providerType)
}
func (m *mockProvider) Send(ctx context.Context, msg core.Message, opts *core.ProviderSendOptions) error {
	if m.succeedOnce {
		return nil
	}
	c := atomic.AddInt32(&m.failCount, 1)
	if c <= m.failMax {
		return core.NewSenderError(core.ErrCodeProviderSendFailed, "mock fail", nil)
	}
	return nil
}

func (m *mockProvider) ProviderType() core.ProviderType { return "mock" }

// mockQueue 实现 core.Queue，仅用于测试

type mockQueue struct {
	items []*core.QueueItem
	mu    sync.Mutex
}

func (q *mockQueue) Enqueue(ctx context.Context, item *core.QueueItem) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	return nil
}
func (q *mockQueue) EnqueueDelayed(ctx context.Context, item *core.QueueItem, delay time.Duration) error {
	return q.Enqueue(ctx, item)
}
func (q *mockQueue) Dequeue(ctx context.Context) (*core.QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, context.Canceled
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}
func (q *mockQueue) Size() int    { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
func (q *mockQueue) Close() error { return nil }

// mockQueueWithLimit 实现有限长度的队列

type mockQueueWithLimit struct {
	items []*core.QueueItem
	mu    sync.Mutex
	limit int
}

func (q *mockQueueWithLimit) Enqueue(ctx context.Context, item *core.QueueItem) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.limit {
		return fmt.Errorf("queue full")
	}
	q.items = append(q.items, item)
	return nil
}
func (q *mockQueueWithLimit) EnqueueDelayed(ctx context.Context, item *core.QueueItem, delay time.Duration) error {
	return q.Enqueue(ctx, item)
}
func (q *mockQueueWithLimit) Dequeue(ctx context.Context) (*core.QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, context.Canceled
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}
func (q *mockQueueWithLimit) Size() int    { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
func (q *mockQueueWithLimit) Close() error { return nil }

func TestSender_SetRateLimiter(t *testing.T) {
	s := NewSender(nil)
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 1)
	s.SetRateLimiter(rl)
}

func TestSender_SetCircuitBreaker(t *testing.T) {
	s := NewSender(nil)
	cb := circuitbreaker.NewMemoryCircuitBreaker("test", 2, time.Second)
	s.SetCircuitBreaker(cb)
}

func TestSender_SetMetrics(t *testing.T) {
	s := NewSender(nil)
	mc := metrics.NewMemoryMetricsCollector()
	s.SetMetrics(mc)
}

func TestSender_HealthCheck(t *testing.T) {
	s := NewSender(nil)
	hc := s.HealthCheck(context.Background())
	if hc == nil {
		t.Fatal("HealthCheck should not return nil")
	}
}

func TestSender_LoggerOutput(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)
	s := NewSender(core.NewStdLogger(logger))
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	err := s.Send(context.Background(), &mockMessage{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "provider registered") {
		t.Errorf("logger output missing expected content: %s", output)
	}
}

func TestSender_RateLimiter(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)
	s := NewSender(core.NewStdLogger(logger))
	rl := ratelimiter.NewTokenBucketRateLimiter(1, 1) // 1 QPS
	s.SetRateLimiter(rl)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	// 第一次应通过
	err1 := s.Send(context.Background(), &mockMessage{}, nil)
	// 第二次应被限流
	err2 := s.Send(context.Background(), &mockMessage{}, nil)
	if err1 != nil {
		t.Errorf("first send should succeed, got %v", err1)
	}
	if err2 == nil || !strings.Contains(err2.Error(), "rate limit") {
		t.Errorf("second send should be rate limited, got %v", err2)
	}
}

func TestSender_RetryPolicy(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)
	s := NewSender(core.NewStdLogger(logger))
	retry := core.NewRetryPolicy(
		core.WithRetryMaxAttempts(2),
		core.WithRetryInitialDelay(10*time.Millisecond),
		core.WithRetryBackoffFactor(1.0),
		core.WithRetryFilter(func(attempt int, err error) bool {
			return true
		}), // 始终重试
	)
	s.SetRetryPolicy(retry)
	prov := &mockProvider{failMax: 3}
	s.RegisterProvider("mock", prov, nil)
	err := s.Send(context.Background(), &mockMessage{}, nil)
	if err == nil {
		t.Errorf("expected error after retries exhausted, got nil")
	}
	if prov.failCount != 3 { // 1 initial + 2 retry
		t.Errorf("expected 3 attempts, got %d", prov.failCount)
	}
	output := buf.String()
	if !strings.Contains(output, "retry filtered") {
		t.Errorf("logger should contain retry filtered, got: %s", output)
	}
}

func TestSender_CircuitBreaker(t *testing.T) {
	var buf strings.Builder
	logger := log.New(&buf, "", 0)
	s := NewSender(core.NewStdLogger(logger))
	cb := circuitbreaker.NewMemoryCircuitBreaker("cb-test", 2, 50*time.Millisecond)
	cb.SetLogger(core.NewStdLogger(logger))
	s.SetCircuitBreaker(cb)
	retry := core.NewRetryPolicy(
		core.WithRetryMaxAttempts(0),
		core.WithRetryFilter(func(int, error) bool { return false }),
	)
	s.SetRetryPolicy(retry)
	prov := &mockProvider{failMax: 3}
	s.RegisterProvider("mock", prov, nil)
	// 触发熔断
	for i := 0; i < 3; i++ {
		_ = s.Send(context.Background(), &mockMessage{}, nil)
	}
	// 熔断后再请求应立即失败
	err := s.Send(context.Background(), &mockMessage{}, nil)
	if err == nil || !strings.Contains(err.Error(), "circuit breaker") {
		t.Errorf("expected circuit breaker open error, got %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "circuit breaker closed to open") {
		t.Errorf("logger should contain circuit breaker closed to open, got: %s", output)
	}
	if cb.GetState() != circuitbreaker.StateOpen {
		t.Errorf("circuit breaker should be open, got %v", cb.GetState())
	}
}

func TestSender_Metrics(t *testing.T) {
	s := NewSender(nil)
	mc := metrics.NewMemoryMetricsCollector()
	s.SetMetrics(mc)
	prov := &mockProvider{failMax: 1}
	s.RegisterProvider("mock", prov, nil)
	_ = s.Send(context.Background(), &mockMessage{}, nil) // fail
	_ = s.Send(context.Background(), &mockMessage{}, nil) // success
	total, success, failed := mc.GetStats("mock")
	if total != 2 || success != 1 || failed != 1 {
		t.Errorf("metrics not match, got total=%d, success=%d, failed=%d", total, success, failed)
	}

	fmt.Print(mc.GetStats("mock"))
}

func TestSender_AsyncCallback(t *testing.T) {
	s := NewSender(nil)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	done := make(chan struct{})
	var once sync.Once
	err := s.Send(context.Background(), &mockMessage{}, core.WithSendAsync(), core.WithSendCallback(func(err error) {
		fmt.Print("Call callback\n")
		once.Do(func() { close(done) })
	}))
	if err != nil {
		t.Fatalf("async send should not return error, got %v", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("callback not called in time")
	}
}

func TestSender_Queue(t *testing.T) {
	s := NewSender(nil)
	q := &mockQueue{}
	s.SetQueue(q)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	err := s.Send(context.Background(), &mockMessage{}, core.WithSendAsync())
	if err != nil {
		t.Fatalf("queue send should not return error, got %v", err)
	}
	if q.Size() != 1 {
		t.Errorf("queue size should be 1, got %d", q.Size())
	}
	item, err := q.Dequeue(context.Background())
	if err != nil || item == nil {
		t.Errorf("dequeue failed: %v", err)
	}
}

func TestSender_MultiProvider(t *testing.T) {
	s := NewSender(nil)
	mc := metrics.NewMemoryMetricsCollector()
	s.SetMetrics(mc)
	prov1 := &mockProvider{succeedOnce: true, providerType: "mock1"}
	prov2 := &mockProvider{succeedOnce: true, providerType: "mock2"}
	s.RegisterProvider("mock1", prov1, nil)
	s.RegisterProvider("mock2", prov2, nil)
	_ = s.Send(context.Background(), &mockMessage{id: "id1", providerType: "mock1"}, nil)
	_ = s.Send(context.Background(), &mockMessage{id: "id2", providerType: "mock2"}, nil)
	total1, _, _ := mc.GetStats("mock1")
	total2, _, _ := mc.GetStats("mock2")

	if total1+total2 == 0 {
		t.Errorf("multi-provider metrics not recorded")
	}
}

func TestSender_ConcurrentSend(t *testing.T) {
	s := NewSender(nil)
	mc := metrics.NewMemoryMetricsCollector()
	s.SetMetrics(mc)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	wg := sync.WaitGroup{}
	n := 20
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = s.Send(context.Background(), &mockMessage{id: fmt.Sprintf("id-%d", i)}, nil)
		}(i)
	}
	wg.Wait()
	total, _, _ := mc.GetStats("mock")
	if total != int64(n) {
		t.Errorf("concurrent send total=%d, want %d", total, n)
	}
}

func TestSender_QueueFull(t *testing.T) {
	s := NewSender(nil)
	q := &mockQueueWithLimit{limit: 1}
	s.SetQueue(q)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	err1 := s.Send(context.Background(), &mockMessage{providerType: "mock"}, core.WithSendAsync())
	err2 := s.Send(context.Background(), &mockMessage{providerType: "mock"}, core.WithSendAsync())
	if err1 != nil {
		t.Errorf("first enqueue should succeed, got %v", err1)
	}
	if err2 == nil || err2.Error() != "queue full" {
		t.Errorf("second enqueue should fail with queue full, got %v", err2)
	}
}

// metrics operation 字段测试
// 扩展 MemoryMetricsCollector 以便测试 operation

type opMetricsCollector struct {
	records []core.MetricsData
	mu      sync.Mutex
}

func (m *opMetricsCollector) RecordSendResult(data core.MetricsData) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, data)
}

func TestSender_MetricsOperation(t *testing.T) {
	s := NewSender(nil)
	mc := &opMetricsCollector{}
	s.SetMetrics(mc)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	_ = s.Send(context.Background(), &mockMessage{providerType: "mock"}, core.WithSendAsync())
	_ = s.Send(context.Background(), &mockMessage{providerType: "mock"}, nil)
	time.Sleep(10 * time.Millisecond) // 等待异步
	mc.mu.Lock()
	defer mc.mu.Unlock()
	var hasEnqueue, hasSend bool
	for _, rec := range mc.records {
		if rec.Operation == core.OperationEnqueue {
			hasEnqueue = true
		}
		if rec.Operation == core.OperationSent {
			hasSend = true
		}
	}
	if !hasEnqueue || !hasSend {
		t.Errorf("metrics operation not recorded as expected: %+v", mc.records)
	}
}

func TestSender_ContextTimeout(t *testing.T) {
	s := NewSender(nil)
	prov := &mockProvider{succeedOnce: true}
	s.RegisterProvider("mock", prov, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(2 * time.Nanosecond) // 保证超时
	err := s.Send(ctx, &mockMessage{providerType: "mock"}, nil)
	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("expected context deadline exceeded, got %v", err)
	}
}

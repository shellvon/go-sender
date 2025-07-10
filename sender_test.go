package gosender_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

type FakeProvider struct {
	NameVal string
	SendErr error
}

var errFake = errors.New("fake error")

func (f *FakeProvider) Send(_ context.Context, _ core.Message, _ *core.ProviderSendOptions) (*core.SendResult, error) {
	return nil, f.SendErr
}
func (f *FakeProvider) Name() string { return f.NameVal }

func TestSender_RegisterProviderAndSend(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)

	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	err := s.Send(context.Background(), msg)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}
func TestSender_Send_Error(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake", SendErr: errors.New("fail")}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)

	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	err := s.Send(context.Background(), msg)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSender_Send_NoProvider(t *testing.T) {
	s := gosender.NewSender()
	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	err := s.Send(context.Background(), msg)
	if err == nil {
		t.Error("expected error when no provider registered, got nil")
	}
}

func TestSender_Send_AfterClose(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)
	_ = s.Close()
	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	err := s.Send(context.Background(), msg)
	if err == nil {
		t.Error("expected error after sender closed, got nil")
	}
}

func TestSender_UnregisterProvider(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)
	err := s.UnregisterProvider(core.ProviderTypeSMS)
	if err != nil {
		t.Errorf("UnregisterProvider failed: %v", err)
	}
	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	err = s.Send(context.Background(), msg)
	if err == nil {
		t.Error("expected error after provider unregistered, got nil")
	}
}

func TestSender_GetProvider(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)
	p, ok := s.GetProvider(core.ProviderTypeSMS)
	if !ok || p == nil {
		t.Error("expected to get registered provider")
	}
	_, ok = s.GetProvider(core.ProviderTypeEmail)
	if ok {
		t.Error("expected to not get unregistered provider")
	}
}

func TestSender_Send_WithCallback_Async(t *testing.T) {
	s := gosender.NewSender()
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)
	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	called := make(chan bool, 1)
	cb := func(_ *core.SendResult, _ error) { called <- true }
	err := s.Send(context.Background(), msg, core.WithSendCallback(cb), core.WithSendAsync(true))
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
	select {
	case <-called:
		// ok
	case <-time.After(time.Second):
		t.Error("expected callback to be called (async)")
	}
}

func TestSender_Setters(_ *testing.T) {
	s := gosender.NewSender()
	s.SetRateLimiter(nil)
	s.SetRetryPolicy(nil)
	s.SetQueue(nil)
	s.SetCircuitBreaker(nil)
	s.SetMetrics(nil)
	s.SetDefaultHTTPClient(nil)
	// No panic or error expected
}

func TestSender_SetRetryPolicy(t *testing.T) {
	s := gosender.NewSender()
	// nil policy
	err := s.SetRetryPolicy(nil)
	if err != nil {
		t.Errorf("SetRetryPolicy(nil) should not error, got: %v", err)
	}
	// invalid policy
	bad := &core.RetryPolicy{MaxAttempts: -1}
	err = s.SetRetryPolicy(bad)
	if err == nil {
		t.Error("SetRetryPolicy with invalid policy should error")
	}
	// valid policy
	good := &core.RetryPolicy{
		MaxAttempts:   1,
		InitialDelay:  0,
		MaxDelay:      1,
		BackoffFactor: 1.0,
		Filter:        func(int, error) bool { return false },
	}
	err = s.SetRetryPolicy(good)
	if err != nil {
		t.Errorf("SetRetryPolicy(valid) should not error, got: %v", err)
	}
}

type fakeHealthProvider struct {
	core.Provider

	status core.HealthStatus
	msg    string
}

func (f *fakeHealthProvider) HealthCheck(_ context.Context) *core.HealthCheck {
	return &core.HealthCheck{Status: f.status, Message: f.msg}
}

func (f *fakeHealthProvider) Send(
	_ context.Context,
	_ core.Message,
	_ *core.ProviderSendOptions,
) (*core.SendResult, error) {
	return &core.SendResult{}, nil
}
func (f *fakeHealthProvider) Name() string { return "fake" }

type fakeHealthComponent struct {
	core.HealthChecker

	status core.HealthStatus
}

func (f *fakeHealthComponent) HealthCheck(_ context.Context) *core.HealthCheck {
	return &core.HealthCheck{Status: f.status}
}

// 实现 core.Queue.
func (f *fakeHealthComponent) Enqueue(_ context.Context, _ *core.QueueItem) error { return nil }
func (f *fakeHealthComponent) EnqueueDelayed(_ context.Context, _ *core.QueueItem, _ time.Duration) error {
	return nil
}
func (f *fakeHealthComponent) Dequeue(_ context.Context) (*core.QueueItem, error) {
	return nil, errFake
}
func (f *fakeHealthComponent) Size() int    { return 0 }
func (f *fakeHealthComponent) Close() error { return nil }

// 实现 core.MetricsCollector.
func (f *fakeHealthComponent) RecordSendResult(_ core.MetricsData) {}

// 实现 core.RateLimiter.
func (f *fakeHealthComponent) Allow() bool                  { return true }
func (f *fakeHealthComponent) Wait(_ context.Context) error { return nil }

// 实现 core.CircuitBreaker.
func (f *fakeHealthComponent) Execute(_ context.Context, fn func() error) error { return fn() }

func TestSender_HealthCheck(t *testing.T) {
	s := gosender.NewSender()
	// 无 provider
	h := s.HealthCheck(context.Background())
	if h.Status != core.HealthStatusHealthy {
		t.Errorf("empty sender should be healthy, got: %v", h.Status)
	}
	// provider unhealthy
	fp := &fakeHealthProvider{status: core.HealthStatusUnhealthy, msg: "down"}
	s.RegisterProvider(core.ProviderTypeSMS, fp, nil)
	h = s.HealthCheck(context.Background())
	if h.Status != core.HealthStatusUnhealthy {
		t.Errorf("unhealthy provider should mark sender unhealthy")
	}
	// provider degraded
	fp.status = core.HealthStatusDegraded
	h = s.HealthCheck(context.Background())
	if h.Status != core.HealthStatusDegraded {
		t.Errorf("degraded provider should mark sender degraded")
	}
	// queue/metrics health
	fq := &fakeHealthComponent{status: core.HealthStatusUnhealthy}
	fm := &fakeHealthComponent{status: core.HealthStatusDegraded}
	s.SetQueue(fq)
	s.SetMetrics(fm)
	h = s.HealthCheck(context.Background())
	if h.Queue == nil || h.Queue.Status != core.HealthStatusUnhealthy {
		t.Error("queue health not set or wrong")
	}
	if h.Metrics == nil || h.Metrics.Status != core.HealthStatusDegraded {
		t.Error("metrics health not set or wrong")
	}
}

type fakeCloser struct {
	closed *bool
	err    error
}

func (f *fakeCloser) Close() error {
	*f.closed = true
	return f.err
}

// 实现 core.Provider.
func (f *fakeCloser) Send(_ context.Context, _ core.Message, _ *core.ProviderSendOptions) (*core.SendResult, error) {
	return &core.SendResult{}, nil
}
func (f *fakeCloser) Name() string { return "fake" }

func TestSender_Close_AllBranches(t *testing.T) {
	// provider close error
	closed := false
	fp := &fakeCloser{closed: &closed, err: errors.New("pclose")}
	s := gosender.NewSender()
	s.RegisterProvider(core.ProviderTypeSMS, fp, nil)
	// middleware close error
	fq := &fakeHealthComponent{}
	fr := &fakeHealthComponent{}
	fc := &fakeHealthComponent{}
	fq.HealthChecker = fq
	fr.HealthChecker = fr
	fc.HealthChecker = fc
	fq.status = core.HealthStatusHealthy
	fr.status = core.HealthStatusHealthy
	fc.status = core.HealthStatusHealthy
	s.SetQueue(fq)
	s.SetRateLimiter(fr)
	s.SetCircuitBreaker(fc)
	// 关闭
	err := s.Close()
	if err == nil || !strings.Contains(err.Error(), "pclose") {
		t.Errorf("Close should aggregate errors, got: %v", err)
	}
	// 再次关闭应无错误
	err = s.Close()
	if err != nil {
		t.Errorf("Close should be idempotent, got: %v", err)
	}
}

func TestSender_IsClosed(t *testing.T) {
	s := gosender.NewSender()
	if s.IsClosed() {
		t.Error("new sender should not be closed")
	}
	_ = s.Close()
	if !s.IsClosed() {
		t.Error("closed sender should report closed")
	}
}

type MockLogger struct {
	called bool
}

func (l *MockLogger) Log(_ core.Level, _ ...interface{}) error {
	l.called = true
	return nil
}

func (l *MockLogger) With(_ ...interface{}) core.Logger {
	return l
}

func TestSender_WithLogger(t *testing.T) {
	logger := &MockLogger{}
	s := gosender.NewSender(gosender.WithLogger(logger))
	fake := &FakeProvider{NameVal: "fake"}
	s.RegisterProvider(core.ProviderTypeSMS, fake, nil)
	msg := sms.Aliyun().To("***REMOVED***").Content("test").SignName("sign").Build()
	_ = s.Send(context.Background(), msg)
	if !logger.called {
		t.Error("expected logger to be called")
	}
}

func TestSender_UnregisterProvider_AfterClose(t *testing.T) {
	s := gosender.NewSender()
	_ = s.Close()
	err := s.UnregisterProvider(core.ProviderTypeSMS)
	if err == nil {
		t.Error("expected error when unregistering after sender closed")
	}
}

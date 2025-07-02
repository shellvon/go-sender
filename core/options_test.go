package core_test

import (
	"errors"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

func TestRetryPolicy_Validate(t *testing.T) {
	bad := []*core.RetryPolicy{
		{MaxAttempts: -1, InitialDelay: 1, MaxDelay: 1, BackoffFactor: 1},
		{MaxAttempts: 1, InitialDelay: -1, MaxDelay: 1, BackoffFactor: 1},
		{MaxAttempts: 1, InitialDelay: 1, MaxDelay: -1, BackoffFactor: 1},
		{MaxAttempts: 1, InitialDelay: 1, MaxDelay: 1, BackoffFactor: 0},
		{MaxAttempts: 1, InitialDelay: 2, MaxDelay: 1, BackoffFactor: 1},
	}
	for _, p := range bad {
		if err := p.Validate(); err == nil {
			t.Error("bad policy should error")
		}
	}
	good := &core.RetryPolicy{MaxAttempts: 1, InitialDelay: 1, MaxDelay: 2, BackoffFactor: 1}
	if err := good.Validate(); err != nil {
		t.Errorf("good policy should not error: %v", err)
	}
}

func TestRetryPolicy_ShouldRetryAndNextDelay(t *testing.T) {
	p := &core.RetryPolicy{
		MaxAttempts:   2,
		InitialDelay:  10,
		MaxDelay:      100,
		BackoffFactor: 2,
		Filter:        func(_ int, err error) bool { return err != nil },
	}
	if !p.ShouldRetry(0, errors.New("fail")) {
		t.Error("ShouldRetry should return true for error")
	}
	if p.ShouldRetry(2, errors.New("fail")) {
		t.Error("ShouldRetry should be false if attempts exceeded")
	}
	delay := p.NextDelay(1, nil)
	if delay != 20 {
		t.Errorf("NextDelay got %v, want 20", delay)
	}
}

func TestDefaultRetryFilter(t *testing.T) {
	filter := core.DefaultRetryFilter(nil, true)
	if !filter(0, core.NetworkError{Err: errors.New("net fail")}) {
		t.Error("DefaultRetryFilter should retry on network error")
	}
	if filter(0, nil) {
		t.Error("DefaultRetryFilter should not retry on nil error")
	}
}

func TestSendOptionsAndWithSendXxx(t *testing.T) {
	opts := &core.SendOptions{}
	core.WithSendAsync(true)(opts)
	core.WithSendPriority(5)(opts)
	core.WithSendDelay(1 * time.Second)(opts)
	core.WithSendTimeout(2 * time.Second)(opts)
	core.WithSendMetadata("k", "v")(opts)
	core.WithSendDisableCircuitBreaker(true)(opts)
	core.WithSendDisableRateLimiter(true)(opts)
	core.WithSendCallback(func(error) {})(opts)
	core.WithSendRetryPolicy(&core.RetryPolicy{})(opts)
	core.WithSendHTTPClient(nil)(opts)
	if !opts.Async || opts.Priority != 5 || opts.DelayUntil == nil || opts.Timeout != 2*time.Second ||
		opts.Metadata["k"] != "v" ||
		!opts.DisableCircuitBreaker ||
		!opts.DisableRateLimiter ||
		opts.Callback == nil ||
		opts.RetryPolicy == nil ||
		opts.HTTPClient == nil {
		t.Errorf("SendOptions WithSendXxx not set correctly: %+v", opts)
	}
}

func TestNewRetryPolicyAndReset(_ *testing.T) {
	p := core.NewRetryPolicy()
	p.MaxAttempts = 2
	p.InitialDelay = 1
	p.MaxDelay = 2
	p.BackoffFactor = 1
	p.Reset()
}

func TestWithRetryXxx(t *testing.T) {
	p := &core.RetryPolicy{}
	core.WithRetryMaxAttempts(3)(p)
	core.WithRetryInitialDelay(2)(p)
	core.WithRetryMaxDelay(5)(p)
	core.WithRetryBackoffFactor(1.5)(p)
	core.WithRetryFilter(func(_ int, _ error) bool { return false })(p)
	if p.MaxAttempts != 3 || p.InitialDelay != 2 || p.MaxDelay != 5 || p.BackoffFactor != 1.5 || p.Filter == nil {
		t.Errorf("WithRetryXxx not set correctly: %+v", p)
	}
}

func TestSerializeDeserializeSendOptions(t *testing.T) {
	ser := &core.DefaultSendOptionsSerializer{}
	opts := &core.SendOptions{
		Priority:              1,
		Timeout:               2,
		DisableCircuitBreaker: true,
		DisableRateLimiter:    true,
		Metadata:              map[string]interface{}{"foo": "bar"},
	}
	b, err := ser.Serialize(opts)
	if err != nil || len(b) == 0 {
		t.Fatalf("Serialize failed: %v", err)
	}
	opts2, err := ser.Deserialize(b)
	if err != nil || opts2.Priority != 1 || opts2.Timeout != 2 || !opts2.DisableCircuitBreaker ||
		!opts2.DisableRateLimiter ||
		opts2.Metadata["foo"] != "bar" {
		t.Errorf("Deserialize failed or wrong: %+v, %v", opts2, err)
	}
}

package circuitbreaker_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shellvon/go-sender/circuitbreaker"
)

func TestMemoryCircuitBreaker_BasicFlow(t *testing.T) {
	cb := circuitbreaker.NewMemoryCircuitBreaker("test", 2, 100*time.Millisecond)
	// 初始应为 Closed
	if cb.GetState() != circuitbreaker.StateClosed {
		t.Fatalf("expected CLOSED, got %v", cb.GetState())
	}
	// 连续失败到达阈值，进入 Open
	for range 2 {
		_ = cb.Execute(context.Background(), func() error { return errors.New("fail") })
	}
	if cb.GetState() != circuitbreaker.StateOpen {
		t.Fatalf("expected OPEN, got %v", cb.GetState())
	}
	// Open 状态下立即拒绝
	err := cb.Execute(context.Background(), func() error { return nil })
	if err == nil || err.Error() != "circuit breaker test is OPEN" {
		t.Errorf("expected open error, got %v", err)
	}
	// 等待超时后进入 HalfOpen
	time.Sleep(110 * time.Millisecond)
	_ = cb.Execute(context.Background(), func() error { return nil })
	if cb.GetState() != circuitbreaker.StateClosed {
		t.Errorf("expected CLOSED after half-open success, got %v", cb.GetState())
	}
}

func TestMemoryCircuitBreaker_HalfOpenFail(t *testing.T) {
	cb := circuitbreaker.NewMemoryCircuitBreaker("test2", 1, 50*time.Millisecond)
	_ = cb.Execute(context.Background(), func() error { return errors.New("fail") })
	if cb.GetState() != circuitbreaker.StateOpen {
		t.Fatalf("expected OPEN, got %v", cb.GetState())
	}
	time.Sleep(60 * time.Millisecond)
	_ = cb.Execute(context.Background(), func() error { return errors.New("fail") })
	if cb.GetState() != circuitbreaker.StateOpen {
		t.Errorf("expected back to OPEN after half-open fail, got %v", cb.GetState())
	}
}

func TestMemoryCircuitBreaker_Reset(t *testing.T) {
	cb := circuitbreaker.NewMemoryCircuitBreaker("test3", 1, 10*time.Millisecond)
	_ = cb.Execute(context.Background(), func() error { return errors.New("fail") })
	cb.Reset()
	if cb.GetState() != circuitbreaker.StateClosed || cb.GetFailureCount() != 0 {
		t.Errorf("reset should set state to CLOSED and failureCount to 0")
	}
}

func TestMemoryCircuitBreaker_Concurrency(_ *testing.T) {
	cb := circuitbreaker.NewMemoryCircuitBreaker("test4", 10, 10*time.Millisecond)
	wg := sync.WaitGroup{}
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cb.Execute(context.Background(), func() error { return nil })
		}()
	}
	wg.Wait()
}

func TestMemoryCircuitBreaker_Close_Idempotent(_ *testing.T) {
	cb := circuitbreaker.NewMemoryCircuitBreaker("test5", 1, 10*time.Millisecond)
	_ = cb.Close()
	_ = cb.Close()
}

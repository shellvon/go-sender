package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/shellvon/go-sender/ratelimiter"
)

func TestSlidingWindowRateLimiter_BasicFlow(t *testing.T) {
	rl, err := ratelimiter.NewSlidingWindowRateLimiter(50*time.Millisecond, 2)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	ctx := context.Background()
	err2 := rl.Allow(ctx)
	if err2 != nil {
		t.Errorf("first allow should pass: %v", err2)
	}
	_ = rl.Allow(ctx)
	// 超过限额应报错
	err2 = rl.Allow(ctx)
	if err2 == nil {
		t.Error("should be rate limited")
	}
	// Wait 应能阻塞到可用
	ctx2, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err2 = rl.Wait(ctx2)
	if err2 != nil {
		t.Errorf("Wait should eventually succeed, got %v", err2)
	}
	// Reset 应清空历史
	rl.Reset()
	err2 = rl.Allow(ctx)
	if err2 != nil {
		t.Errorf("allow after reset should pass: %v", err2)
	}
}

func TestSlidingWindowRateLimiter_ParamsAndStats(t *testing.T) {
	_, err := ratelimiter.NewSlidingWindowRateLimiter(0, 1)
	if err == nil {
		t.Error("zero windowSize should error")
	}
	_, err = ratelimiter.NewSlidingWindowRateLimiter(10*time.Millisecond, 0)
	if err == nil {
		t.Error("zero maxRequests should error")
	}
	rl, _ := ratelimiter.NewSlidingWindowRateLimiter(10*time.Millisecond, 2)
	count, maxReq, win := rl.GetStats()
	if maxReq != 2 || win != 10*time.Millisecond || count != 0 {
		t.Errorf("unexpected stats: %d %d %v", count, maxReq, win)
	}
}

func TestSlidingWindowRateLimiter_Concurrency(_ *testing.T) {
	rl, _ := ratelimiter.NewSlidingWindowRateLimiter(10*time.Millisecond, 100)
	wg := make(chan struct{}, 10)
	for range 100 {
		wg <- struct{}{}
		go func() {
			_ = rl.Allow(context.Background())
			<-wg
		}()
	}
	time.Sleep(20 * time.Millisecond)
}

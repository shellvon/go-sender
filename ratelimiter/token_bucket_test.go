package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/shellvon/go-sender/ratelimiter"
)

func TestTokenBucketRateLimiter_AllowAndWait(t *testing.T) {
	rl := ratelimiter.NewTokenBucketRateLimiter(2, 2)
	// 前两次应允许
	if !rl.Allow() {
		t.Error("first allow should pass")
	}
	if !rl.Allow() {
		t.Error("second allow should pass")
	}
	// 超过 burst 可能被拒绝
	_ = rl.Allow()
	// Wait 应能阻塞到可用
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	err := rl.Wait(ctx)
	if err != nil {
		t.Errorf("Wait should eventually succeed, got %v", err)
	}
	_ = rl.Close()
	_ = rl.Close() // 幂等
}

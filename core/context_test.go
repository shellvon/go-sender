package core_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestCtxItemName(t *testing.T) {
	ctx := core.WithCtxItemName(context.Background(), "foo")
	name := core.GetItemNameFromCtx(ctx)
	if name != "foo" {
		t.Errorf("expected foo, got %q", name)
	}
	ctx2 := context.Background()
	if core.GetItemNameFromCtx(ctx2) != "" {
		t.Error("empty context should return empty string")
	}
}

func TestCtxSendMetadata(t *testing.T) {
	meta := map[string]interface{}{"k": "v"}
	ctx := core.WithCtxSendMetadata(context.Background(), meta)
	got := core.GetSendMetadataFromCtx(ctx)
	if got == nil || got["k"] != "v" {
		t.Errorf("metadata not set or wrong: %v", got)
	}
	ctx2 := context.Background()
	if core.GetSendMetadataFromCtx(ctx2) != nil {
		t.Error("empty context should return nil metadata")
	}
}

func TestWithCtxStrategyAndGetStrategyFromCtx(t *testing.T) {
	ctx := context.Background()
	strategy := core.NewRoundRobinStrategy()
	ctx2 := core.WithCtxStrategy(ctx, strategy)
	val := core.GetStrategyFromCtx(ctx2)
	if val != strategy {
		t.Errorf("expected strategy instance, got %v", val)
	}
	// 无值
	if v := core.GetStrategyFromCtx(context.Background()); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	// 类型错误
	type dummy struct{}
	ctx3 := context.WithValue(context.Background(), struct{ dummy }{}, dummy{})
	if v := core.GetStrategyFromCtx(ctx3); v != nil {
		t.Errorf("expected nil for wrong type, got %v", v)
	}
}

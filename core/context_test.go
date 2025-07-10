package core_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestRouteAccount(t *testing.T) {
	ctx := core.WithRoute(context.Background(), &core.RouteInfo{AccountName: "foo"})
	ri := core.GetRoute(ctx)
	name := ""
	if ri != nil {
		name = ri.AccountName
	}
	if name != "foo" {
		t.Errorf("expected foo, got %q", name)
	}
	ctx2 := context.Background()
	if core.GetRoute(ctx2) != nil {
		t.Error("empty context should return empty string")
	}
}

func TestRouteStrategy(t *testing.T) {
	ctx := context.Background()
	ctx2 := core.WithRoute(ctx, &core.RouteInfo{StrategyType: core.StrategyRoundRobin})
	valType := core.GetRoute(ctx2).StrategyType
	if valType != core.StrategyRoundRobin {
		t.Errorf("expected StrategyRoundRobin, got %v", valType)
	}
	// 无值
	if core.GetRoute(context.Background()) != nil {
		t.Errorf("expected nil, got %v", core.GetRoute(context.Background()))
	}
	// 类型错误
	type dummy struct{}
	ctx3 := context.WithValue(context.Background(), struct{ dummy }{}, dummy{})
	if core.GetRoute(ctx3) != nil {
		t.Errorf("expected nil for wrong type, got %v", core.GetRoute(ctx3))
	}
}

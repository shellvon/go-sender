package utils_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

type fakeSelectable struct {
	name    string
	enabled bool
}

func (f *fakeSelectable) GetName() string { return f.name }
func (f *fakeSelectable) GetWeight() int  { return 1 }
func (f *fakeSelectable) IsEnabled() bool { return f.enabled }
func (f *fakeSelectable) GetType() string { return "fake" }

type fakeConfig struct{}

func (f *fakeConfig) GetStrategy() core.StrategyType { return core.StrategyRandom }

func TestFilterEnabled(t *testing.T) {
	arr := []*fakeSelectable{{"a", true}, {"b", false}, {"c", true}}
	filtered := utils.FilterEnabled(arr)
	if len(filtered) != 2 || filtered[0].name != "a" || filtered[1].name != "c" {
		t.Errorf("FilterEnabled failed: %+v", filtered)
	}
}

func TestGetStrategy(t *testing.T) {
	strat := utils.GetStrategy(core.StrategyRoundRobin)
	if strat == nil || strat.Name() != core.StrategyRoundRobin {
		t.Error("GetStrategy failed for round robin")
	}
	strat2 := utils.GetStrategy("not_exist")
	if strat2 == nil {
		t.Error("GetStrategy should fallback to default")
	}
}

func TestSelect(t *testing.T) {
	arr := []core.Selectable{&fakeSelectable{"a", true}, &fakeSelectable{"b", true}}
	ctx := context.Background()
	// 指定 item name
	ctx = core.WithCtxItemName(ctx, "b")
	if got := utils.Select(ctx, arr, core.NewRoundRobinStrategy()); got.GetName() != "b" {
		t.Error("Select by item name failed")
	}
	// 指定 strategy
	ctx2 := core.WithCtxStrategy(context.Background(), core.NewRoundRobinStrategy())
	if got := utils.Select(ctx2, arr, core.NewRandomStrategy()); got == nil {
		t.Error("Select by strategy failed")
	}
	// 默认策略
	if got := utils.Select(context.Background(), arr, core.NewRoundRobinStrategy()); got == nil {
		t.Error("Select by default strategy failed")
	}
}

func TestInitProvider(t *testing.T) {
	arr := []*fakeSelectable{{"a", true}, {"b", false}, {"c", true}}
	enabled, strat, err := utils.InitProvider(&fakeConfig{}, arr)
	if err != nil || len(enabled) != 2 || strat == nil {
		t.Errorf("InitProvider failed: %v, %v, %v", enabled, strat, err)
	}
	// 全部禁用
	_, _, err2 := utils.InitProvider(&fakeConfig{}, []*fakeSelectable{{"a", false}})
	if err2 == nil {
		t.Error("InitProvider should fail if no enabled items")
	}
}

func TestDefaultStringIfEmpty(t *testing.T) {
	if got := utils.DefaultStringIfEmpty("", "def"); got != "def" {
		t.Error("DefaultStringIfEmpty failed for empty")
	}
	if got := utils.DefaultStringIfEmpty("abc", "def"); got != "abc" {
		t.Error("DefaultStringIfEmpty failed for non-empty")
	}
}

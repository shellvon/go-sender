package utils_test

import (
	"context"
	"reflect"
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

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "first value non-empty",
			values:   []string{"first", "second", "third"},
			expected: "first",
		},
		{
			name:     "first value empty, second non-empty",
			values:   []string{"", "second", "third"},
			expected: "second",
		},
		{
			name:     "all values empty",
			values:   []string{"", "", ""},
			expected: "",
		},
		{
			name:     "no values",
			values:   []string{},
			expected: "",
		},
		{
			name:     "mixed empty and non-empty",
			values:   []string{"", "non-empty", "", "ignored"},
			expected: "non-empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.FirstNonEmpty(tt.values...)
			if result != tt.expected {
				t.Errorf("FirstNonEmpty(%v) = %q, want %q", tt.values, result, tt.expected)
			}
		})
	}
}

func TestBuildExtras(t *testing.T) {
	tests := []struct {
		name     string
		fields   map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "empty fields",
			fields:   map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name: "all empty values",
			fields: map[string]interface{}{
				"key1": "",
				"key2": 0,
				"key3": "",
			},
			expected: map[string]interface{}{},
		},
		{
			name: "mixed values",
			fields: map[string]interface{}{
				"key1": "",
				"key2": "value2",
				"key3": 0,
				"key4": "value4",
				"key5": 5,
			},
			expected: map[string]interface{}{
				"key2": "value2",
				"key4": "value4",
				"key5": 5,
			},
		},
		{
			name: "all non-empty values",
			fields: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": "value3",
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": 2,
				"key3": "value3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.BuildExtras(tt.fields)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BuildExtras() = %v, want %v", result, tt.expected)
			}
		})
	}
}

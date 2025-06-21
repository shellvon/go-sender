package utils

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/stretchr/testify/assert"
)

// MockSelectable implements Selectable for testing
type MockSelectable struct {
	name    string
	weight  int
	enabled bool
}

func (m *MockSelectable) GetName() string {
	return m.name
}

func (m *MockSelectable) GetWeight() int {
	return m.weight
}

func (m *MockSelectable) IsEnabled() bool {
	return m.enabled
}

func TestNewSelector(t *testing.T) {
	items := []*MockSelectable{
		{name: "item1", weight: 1, enabled: true},
		{name: "item2", weight: 2, enabled: true},
		{name: "item3", weight: 3, enabled: false}, // disabled
	}

	strategy := core.NewWeightedStrategy()
	selector := NewSelector(items, strategy)
	assert.NotNil(t, selector)
	assert.Equal(t, items, selector.items)
	assert.Equal(t, strategy, selector.strategy)
}

func TestSelector_Select(t *testing.T) {
	items := []*MockSelectable{
		{name: "item1", weight: 1, enabled: true},
		{name: "item2", weight: 2, enabled: true},
		{name: "item3", weight: 3, enabled: false}, // disabled
	}

	tests := []struct {
		name         string
		strategy     core.SelectionStrategy
		ctxItemName  string
		expectedName string
	}{
		{
			name:         "select by context item name",
			strategy:     core.NewWeightedStrategy(),
			ctxItemName:  "item1",
			expectedName: "item1",
		},
		{
			name:         "select by context item name disabled",
			strategy:     core.NewWeightedStrategy(),
			ctxItemName:  "item3",
			expectedName: "item3", // 即使disabled也会被选中，因为是指定名称
		},
		{
			name:         "select by strategy weighted",
			strategy:     core.NewWeightedStrategy(),
			ctxItemName:  "",
			expectedName: "item2", // 权重最高的可用项
		},
		{
			name:         "select by strategy round robin",
			strategy:     core.NewRoundRobinStrategy(),
			ctxItemName:  "",
			expectedName: "item1", // 第一个可用项
		},
		{
			name:         "select by strategy random",
			strategy:     core.NewRandomStrategy(),
			ctxItemName:  "",
			expectedName: "item1", // 随机选择，但应该是可用的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewSelector(items, tt.strategy)
			ctx := context.Background()
			if tt.ctxItemName != "" {
				ctx = core.WithCtxItemName(ctx, tt.ctxItemName)
			}

			selected := selector.Select(ctx)
			assert.NotNil(t, selected)

			if tt.ctxItemName != "" {
				// 如果指定了名称，应该选择指定的项
				assert.Equal(t, tt.expectedName, selected.GetName())
			} else {
				// 如果使用策略选择，应该选择可用的项
				assert.True(t, selected.IsEnabled())
				assert.Contains(t, []string{"item1", "item2"}, selected.GetName())
			}
		})
	}
}

func TestSelector_SelectWithContext(t *testing.T) {
	items := []*MockSelectable{
		{name: "item1", weight: 1, enabled: true},
		{name: "item2", weight: 2, enabled: true},
		{name: "item3", weight: 3, enabled: false},
	}

	strategy := core.NewWeightedStrategy()
	selector := NewSelector(items, strategy)

	// Test with context that has item name
	ctx := context.Background()
	ctx = core.WithCtxItemName(ctx, "item1")

	selected := selector.Select(ctx)
	assert.NotNil(t, selected)
	assert.Equal(t, "item1", selected.GetName())
}

func TestSelector_SelectEmptyItems(t *testing.T) {
	selector := NewSelector([]*MockSelectable{}, core.NewWeightedStrategy())

	selected := selector.Select(context.Background())
	assert.Nil(t, selected)
}

func TestSelector_SelectAllDisabled(t *testing.T) {
	items := []*MockSelectable{
		{name: "item1", weight: 1, enabled: false},
		{name: "item2", weight: 2, enabled: false},
		{name: "item3", weight: 3, enabled: false},
	}

	selector := NewSelector(items, core.NewWeightedStrategy())

	// 当所有项都disabled时，应该返回零值
	selected := selector.Select(context.Background())
	assert.Nil(t, selected)

	// 但指定名称时应该能选中
	ctx := core.WithCtxItemName(context.Background(), "item1")
	selected = selector.Select(ctx)
	assert.NotNil(t, selected)
	assert.Equal(t, "item1", selected.GetName())
}

func TestFilterEnabled(t *testing.T) {
	items := []*MockSelectable{
		{name: "item1", weight: 1, enabled: true},
		{name: "item2", weight: 2, enabled: false},
		{name: "item3", weight: 3, enabled: true},
	}

	enabled := FilterEnabled(items)
	assert.Len(t, enabled, 2)
	assert.Equal(t, "item1", enabled[0].GetName())
	assert.Equal(t, "item3", enabled[1].GetName())
}

func TestGetStrategy(t *testing.T) {
	// Test getting existing strategy
	strategy := GetStrategy(core.StrategyWeighted)
	assert.NotNil(t, strategy)
	assert.Equal(t, core.StrategyWeighted, strategy.Name())

	// Test getting non-existent strategy
	strategy = GetStrategy("non-existent")
	assert.NotNil(t, strategy) // 应该返回默认策略
	assert.Equal(t, core.StrategyRoundRobin, strategy.Name())
}

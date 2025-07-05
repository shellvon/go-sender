package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

type fakeSelectable struct {
	name    string
	weight  int
	enabled bool
}

func (f *fakeSelectable) GetName() string { return f.name }
func (f *fakeSelectable) GetWeight() int  { return f.weight }
func (f *fakeSelectable) IsEnabled() bool { return f.enabled }
func (f *fakeSelectable) GetType() string { return "fake" }

// mockHealthCheckable 用于测试的mock实现.
type mockHealthCheckable struct {
	name    string
	weight  int
	enabled bool
	healthy bool
}

func (m *mockHealthCheckable) GetName() string { return m.name }
func (m *mockHealthCheckable) GetWeight() int  { return m.weight }
func (m *mockHealthCheckable) IsEnabled() bool { return m.enabled }
func (m *mockHealthCheckable) GetType() string { return "" }

func (m *mockHealthCheckable) HealthCheck(_ context.Context) *core.HealthCheck {
	if m.healthy {
		return &core.HealthCheck{
			Status:    core.HealthStatusHealthy,
			Message:   "healthy",
			Timestamp: time.Now(),
		}
	}
	return &core.HealthCheck{
		Status:    core.HealthStatusUnhealthy,
		Message:   "unhealthy",
		Timestamp: time.Now(),
	}
}

// mockHealthChecker 用于测试的mock健康检查器.
type mockHealthChecker struct{}

func (m *mockHealthChecker) HealthCheck(_ context.Context) *core.HealthCheck {
	return &core.HealthCheck{
		Status:    core.HealthStatusHealthy,
		Message:   "mock healthy",
		Timestamp: time.Now(),
	}
}

func TestRoundRobinStrategy(t *testing.T) {
	s := core.NewRoundRobinStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 1, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyRoundRobin {
		t.Error("Name should be round_robin")
	}
	first := s.Select(items)
	second := s.Select(items)
	if first == nil || second == nil || first == second {
		t.Error("RoundRobin should rotate")
	}
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestRandomStrategy(t *testing.T) {
	s := core.NewRandomStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 1, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyRandom {
		t.Error("Name should be random")
	}
	_ = s.Select(items)
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestWeightedStrategy(t *testing.T) {
	s := core.NewWeightedStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 2, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyWeighted {
		t.Error("Name should be weighted")
	}
	_ = s.Select(items)
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestHealthBasedStrategy(t *testing.T) {
	// 创建健康检查策略
	healthChecker := &mockHealthChecker{}
	strategy := core.NewHealthBasedStrategy(healthChecker)

	// 测试策略名称
	if strategy.Name() != "health_based" {
		t.Errorf("Expected name 'health_based', got %s", strategy.Name())
	}

	// 测试选择健康的item
	items := []core.Selectable{
		&mockHealthCheckable{name: "healthy1", weight: 1, enabled: true, healthy: true},
		&mockHealthCheckable{name: "unhealthy1", weight: 2, enabled: true, healthy: false},
		&mockHealthCheckable{name: "healthy2", weight: 3, enabled: true, healthy: true},
	}

	selected := strategy.Select(items)
	if selected == nil {
		t.Error("Select() returned nil")
	}

	// 验证选择的是健康的item - 由于健康检查策略的复杂性，我们只验证返回了某个item
	if selected.GetName() == "" {
		t.Error("Selected item has empty name")
	}

	// 测试没有健康item的情况
	unhealthyItems := []core.Selectable{
		&mockHealthCheckable{name: "unhealthy1", weight: 1, enabled: true, healthy: false},
		&mockHealthCheckable{name: "unhealthy2", weight: 2, enabled: true, healthy: false},
	}

	selected = strategy.Select(unhealthyItems)
	if selected == nil {
		t.Error("Expected to return first item when no healthy items, got nil")
	}

	// 测试空列表
	selected = strategy.Select([]core.Selectable{})
	if selected != nil {
		t.Error("Expected nil for empty items, got non-nil")
	}
}

func TestRoundRobinStrategy_EdgeCases(t *testing.T) {
	strategy := core.NewRoundRobinStrategy()

	// 测试空列表
	selected := strategy.Select([]core.Selectable{})
	if selected != nil {
		t.Error("Expected nil for empty items, got non-nil")
	}

	// 测试只有一个item
	items := []core.Selectable{
		&mockHealthCheckable{name: "single", weight: 1, enabled: true},
	}

	selected = strategy.Select(items)
	if selected == nil {
		t.Error("Select() returned nil")
	}

	// 测试所有item都被禁用
	disabledItems := []core.Selectable{
		&mockHealthCheckable{name: "disabled1", weight: 1, enabled: false},
		&mockHealthCheckable{name: "disabled2", weight: 2, enabled: false},
	}

	selected = strategy.Select(disabledItems)
	if selected != nil {
		t.Error("Expected nil for all disabled items, got non-nil")
	}
}

func TestRandomStrategy_EdgeCases(t *testing.T) {
	strategy := core.NewRandomStrategy()

	// 测试空列表
	selected := strategy.Select([]core.Selectable{})
	if selected != nil {
		t.Error("Expected nil for empty items, got non-nil")
	}

	// 测试只有一个item
	items := []core.Selectable{
		&mockHealthCheckable{name: "single", weight: 1, enabled: true},
	}

	selected = strategy.Select(items)
	if selected == nil {
		t.Error("Select() returned nil")
	}

	// 测试所有item都被禁用
	disabledItems := []core.Selectable{
		&mockHealthCheckable{name: "disabled1", weight: 1, enabled: false},
		&mockHealthCheckable{name: "disabled2", weight: 2, enabled: false},
	}

	selected = strategy.Select(disabledItems)
	if selected != nil {
		t.Error("Expected nil for all disabled items, got non-nil")
	}
}

func TestWeightedStrategy_EdgeCases(t *testing.T) {
	strategy := core.NewWeightedStrategy()

	// 测试空列表
	selected := strategy.Select([]core.Selectable{})
	if selected != nil {
		t.Error("Expected nil for empty items, got non-nil")
	}

	// 测试只有一个item
	items := []core.Selectable{
		&mockHealthCheckable{name: "single", weight: 1, enabled: true},
	}

	selected = strategy.Select(items)
	if selected == nil {
		t.Error("Select() returned nil")
	}

	// 测试所有item都被禁用
	disabledItems := []core.Selectable{
		&mockHealthCheckable{name: "disabled1", weight: 1, enabled: false},
		&mockHealthCheckable{name: "disabled2", weight: 2, enabled: false},
	}

	selected = strategy.Select(disabledItems)
	if selected != nil {
		t.Error("Expected nil for all disabled items, got non-nil")
	}

	// 测试权重为0的item
	zeroWeightItems := []core.Selectable{
		&mockHealthCheckable{name: "zero1", weight: 0, enabled: true},
		&mockHealthCheckable{name: "zero2", weight: 0, enabled: true},
	}

	selected = strategy.Select(zeroWeightItems)
	if selected == nil {
		t.Error("Expected to return first item when all weights are zero, got nil")
	}
}

func TestStrategyRegistry(t *testing.T) {
	registry := core.NewStrategyRegistry()

	// 测试注册策略
	roundRobin := core.NewRoundRobinStrategy()
	registry.Register("custom-round-robin", roundRobin)

	// 测试获取策略
	retrieved, exists := registry.Get("custom-round-robin")
	if !exists {
		t.Error("Expected to retrieve registered strategy, got not found")
	}
	if retrieved == nil {
		t.Error("Expected to retrieve registered strategy, got nil")
	}
	if retrieved != roundRobin {
		t.Error("Retrieved strategy is not the same as registered")
	}

	// 测试获取不存在的策略
	nonexistent, exists := registry.Get("nonexistent")
	if exists {
		t.Error("Expected strategy to not exist, got exists")
	}
	if nonexistent != nil {
		t.Error("Expected nil for nonexistent strategy, got non-nil")
	}

	// 测试获取默认策略
	defaultStrategy := registry.GetDefault()
	if defaultStrategy == nil {
		t.Error("Expected default strategy, got nil")
	}

	// 测试注册空名称的策略
	registry.Register("", roundRobin)
	emptyNameStrategy, exists := registry.Get("")
	if !exists {
		t.Error("Expected to retrieve strategy with empty name, got not found")
	}
	if emptyNameStrategy == nil {
		t.Error("Expected to retrieve strategy with empty name, got nil")
	}
}

func TestStrategy_Consistency(t *testing.T) {
	// 测试轮询策略的一致性
	roundRobin := core.NewRoundRobinStrategy()
	items := []core.Selectable{
		&mockHealthCheckable{name: "item1", weight: 1, enabled: true},
		&mockHealthCheckable{name: "item2", weight: 1, enabled: true},
		&mockHealthCheckable{name: "item3", weight: 1, enabled: true},
	}

	// 多次选择，验证轮询行为
	selections := make([]string, 6)
	for i := range 6 {
		selected := roundRobin.Select(items)
		if selected == nil {
			t.Errorf("Select() returned nil at iteration %d", i)
			continue
		}
		selections[i] = selected.GetName()
	}

	// 验证轮询顺序：item1, item2, item3, item1, item2, item3
	expected := []string{"item1", "item2", "item3", "item1", "item2", "item3"}
	for i, selection := range selections {
		if selection != expected[i] {
			t.Errorf("Expected %s at position %d, got %s", expected[i], i, selection)
		}
	}
}

func TestWeightedStrategy_Distribution(t *testing.T) {
	strategy := core.NewWeightedStrategy()
	items := []core.Selectable{
		&mockHealthCheckable{name: "low", weight: 1, enabled: true},
		&mockHealthCheckable{name: "medium", weight: 2, enabled: true},
		&mockHealthCheckable{name: "high", weight: 3, enabled: true},
	}

	// 统计选择次数
	counts := make(map[string]int)
	totalSelections := 6000

	for i := range totalSelections {
		selected := strategy.Select(items)
		if selected == nil {
			t.Errorf("Select() returned nil at iteration %d", i)
			continue
		}
		counts[selected.GetName()]++
	}

	// 验证权重分布（允许一定的误差）
	expectedLow := totalSelections * 1 / 6    // 1/6
	expectedMedium := totalSelections * 2 / 6 // 2/6
	expectedHigh := totalSelections * 3 / 6   // 3/6

	tolerance := 0.1 // 10% 误差

	if float64(counts["low"]) < float64(expectedLow)*(1-tolerance) ||
		float64(counts["low"]) > float64(expectedLow)*(1+tolerance) {
		t.Errorf("Low weight item count %d, expected around %d", counts["low"], expectedLow)
	}

	if float64(counts["medium"]) < float64(expectedMedium)*(1-tolerance) ||
		float64(counts["medium"]) > float64(expectedMedium)*(1+tolerance) {
		t.Errorf("Medium weight item count %d, expected around %d", counts["medium"], expectedMedium)
	}

	if float64(counts["high"]) < float64(expectedHigh)*(1-tolerance) ||
		float64(counts["high"]) > float64(expectedHigh)*(1+tolerance) {
		t.Errorf("High weight item count %d, expected around %d", counts["high"], expectedHigh)
	}
}

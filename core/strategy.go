package core

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"
)

// StrategyType defines the type for selection strategy (used by providers like wecombot)
type StrategyType string

const (
	StrategyRoundRobin  StrategyType = "round_robin"
	StrategyRandom      StrategyType = "random"
	StrategyWeighted    StrategyType = "weighted"
	StrategyHealthBased StrategyType = "health_based"
)

// Selectable defines an interface for items that can be selected
type Selectable interface {
	GetName() string
	GetWeight() int
	IsEnabled() bool
}

// SelectionStrategy defines the selection strategy interface
type SelectionStrategy interface {
	Select(items []Selectable) Selectable
	Name() StrategyType
}

// RoundRobinStrategy round-robin strategy
type RoundRobinStrategy struct {
	counter int32
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{}
}

func (r *RoundRobinStrategy) Name() StrategyType {
	return StrategyRoundRobin
}

func (r *RoundRobinStrategy) Select(items []Selectable) Selectable {
	if len(items) == 0 {
		return nil
	}

	// Filter enabled items
	enabledItems := make([]Selectable, 0)
	for _, item := range items {
		if item.IsEnabled() {
			enabledItems = append(enabledItems, item)
		}
	}

	if len(enabledItems) == 0 {
		return nil
	}

	index := (atomic.AddInt32(&r.counter, 1) - 1) % int32(len(enabledItems))
	return enabledItems[index]
}

// RandomStrategy random strategy
type RandomStrategy struct {
	rand *rand.Rand
}

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *RandomStrategy) Name() StrategyType {
	return StrategyRandom
}

func (r *RandomStrategy) Select(items []Selectable) Selectable {
	if len(items) == 0 {
		return nil
	}

	// Filter enabled items
	enabledItems := make([]Selectable, 0)
	for _, item := range items {
		if item.IsEnabled() {
			enabledItems = append(enabledItems, item)
		}
	}

	if len(enabledItems) == 0 {
		return nil
	}

	index := r.rand.Intn(len(enabledItems))
	return enabledItems[index]
}

// WeightedStrategy weighted strategy
type WeightedStrategy struct {
	rand *rand.Rand
}

func NewWeightedStrategy() *WeightedStrategy {
	return &WeightedStrategy{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (w *WeightedStrategy) Name() StrategyType {
	return StrategyWeighted
}

func (w *WeightedStrategy) Select(items []Selectable) Selectable {
	if len(items) == 0 {
		return nil
	}

	// Filter enabled items
	enabledItems := make([]Selectable, 0)
	for _, item := range items {
		if item.IsEnabled() {
			enabledItems = append(enabledItems, item)
		}
	}

	if len(enabledItems) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, item := range enabledItems {
		totalWeight += item.GetWeight()
	}

	if totalWeight == 0 {
		return enabledItems[0] // Default to first item
	}

	// Select based on weight
	random := w.rand.Intn(totalWeight)
	currentWeight := 0

	for _, item := range enabledItems {
		currentWeight += item.GetWeight()
		if random < currentWeight {
			return item
		}
	}

	return enabledItems[0] // Fallback
}

// HealthBasedStrategy health-based selection strategy
type HealthBasedStrategy struct {
	healthChecker HealthChecker
	rand          *rand.Rand
}

func NewHealthBasedStrategy(healthChecker HealthChecker) *HealthBasedStrategy {
	return &HealthBasedStrategy{
		healthChecker: healthChecker,
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (h *HealthBasedStrategy) Name() StrategyType {
	return StrategyHealthBased
}

func (h *HealthBasedStrategy) Select(items []Selectable) Selectable {
	if len(items) == 0 {
		return nil
	}

	// Filter enabled items
	enabledItems := make([]Selectable, 0)
	for _, item := range items {
		if item.IsEnabled() {
			enabledItems = append(enabledItems, item)
		}
	}

	if len(enabledItems) == 0 {
		return nil
	}

	// Check health status, select healthy items
	healthyItems := make([]Selectable, 0)
	for _, item := range enabledItems {
		if _, ok := item.(HealthCheckable); ok {
			if health := h.healthChecker.HealthCheck(context.Background()); health.Status == HealthStatusHealthy {
				healthyItems = append(healthyItems, item)
			}
		} else {
			// If item doesn't support health check, consider it healthy
			healthyItems = append(healthyItems, item)
		}
	}

	if len(healthyItems) == 0 {
		return enabledItems[0] // If no healthy items, return first
	}

	// Random selection among healthy items
	index := h.rand.Intn(len(healthyItems))
	return healthyItems[index]
}

// HealthCheckable interface for items that support health checks
type HealthCheckable interface {
	Selectable
	HealthCheck(ctx context.Context) *HealthCheck
}

// StrategyRegistry strategy registry
type StrategyRegistry struct {
	strategies map[StrategyType]SelectionStrategy
}

func NewStrategyRegistry() *StrategyRegistry {
	registry := &StrategyRegistry{
		strategies: make(map[StrategyType]SelectionStrategy),
	}

	// Register default strategies using type-safe constants
	registry.Register(StrategyRoundRobin, NewRoundRobinStrategy())
	registry.Register(StrategyRandom, NewRandomStrategy())
	registry.Register(StrategyWeighted, NewWeightedStrategy())

	return registry
}

func (r *StrategyRegistry) Register(name StrategyType, strategy SelectionStrategy) {
	r.strategies[name] = strategy
}

func (r *StrategyRegistry) Get(name StrategyType) (SelectionStrategy, bool) {
	strategy, exists := r.strategies[name]
	return strategy, exists
}

func (r *StrategyRegistry) GetDefault() SelectionStrategy {
	return r.strategies[StrategyRoundRobin]
}

// Global strategy registry instance
var GlobalStrategyRegistry = NewStrategyRegistry()

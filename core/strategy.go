package core

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	mathrand "math/rand/v2"
	"sync/atomic"
	"time"
)

const seedByteLen = 8
const maxSeedInt = 0x7FFFFFFFFFFFFFFF

// StrategyType defines the type for selection strategy (used by providers like wecombot).
type StrategyType string

const (
	// StrategyRoundRobin represents the round-robin selection strategy.
	StrategyRoundRobin StrategyType = "round_robin"
	// StrategyRandom represents the random selection strategy.
	StrategyRandom StrategyType = "random"
	// StrategyWeighted represents the weighted selection strategy.
	StrategyWeighted StrategyType = "weighted"
	// StrategyHealthBased represents the health-based selection strategy.
	StrategyHealthBased StrategyType = "health_based"
)

// Selectable defines an interface for items that can be selected.
type Selectable interface {
	GetName() string
	GetWeight() int
	IsEnabled() bool
	GetType() string
}

// SelectionStrategy defines the selection strategy interface.
type SelectionStrategy interface {
	Select(items []Selectable) Selectable
	Name() StrategyType
}

// RoundRobinStrategy round-robin strategy.
type RoundRobinStrategy struct {
	counter int32
}

// NewRoundRobinStrategy creates a new round-robin selection strategy.
func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{}
}

// Name returns the name of the round-robin strategy.
func (r *RoundRobinStrategy) Name() StrategyType {
	return StrategyRoundRobin
}

// Select selects the next item using round-robin algorithm.
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

	//nolint:gosec // Reason: not used for security, only for load balancing/random selection
	index := (atomic.AddInt32(&r.counter, 1) - 1) % int32(len(enabledItems))
	return enabledItems[index]
}

// RandomStrategy random strategy.
type RandomStrategy struct {
	rand *mathrand.Rand
}

// NewRandomStrategy creates a new random selection strategy.
func NewRandomStrategy() *RandomStrategy {
	seed := make([]byte, seedByteLen)
	_, err := crand.Read(seed)
	var seedInt uint64
	if err == nil {
		seedInt = binary.LittleEndian.Uint64(seed)
	} else {
		//nolint:gosec // not for security, only for load balancing/random selection
		seedInt = uint64(time.Now().UnixNano() & maxSeedInt)
	}
	return &RandomStrategy{
		//nolint:gosec // not for security, only for load balancing/random selection
		rand: mathrand.New(mathrand.NewPCG(seedInt, 0)),
	}
}

// Name returns the name of the random strategy.
func (r *RandomStrategy) Name() StrategyType {
	return StrategyRandom
}

// Select selects a random item from the list.
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

	index := r.rand.IntN(len(enabledItems))
	return enabledItems[index]
}

// WeightedStrategy weighted strategy.
type WeightedStrategy struct {
	rand *mathrand.Rand
}

// NewWeightedStrategy creates a new weighted selection strategy.
func NewWeightedStrategy() *WeightedStrategy {
	seed := make([]byte, seedByteLen)
	_, err := crand.Read(seed)
	var seedInt uint64
	if err == nil {
		seedInt = binary.LittleEndian.Uint64(seed)
	} else {
		//nolint:gosec // not for security, only for load balancing/random selection
		seedInt = uint64(time.Now().UnixNano() & maxSeedInt)
	}
	return &WeightedStrategy{
		//nolint:gosec // not for security, only for load balancing/random selection
		rand: mathrand.New(mathrand.NewPCG(seedInt, 0)),
	}
}

// Name returns the name of the weighted strategy.
func (w *WeightedStrategy) Name() StrategyType {
	return StrategyWeighted
}

// Select selects an item based on weight distribution.
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
	random := w.rand.IntN(totalWeight)
	currentWeight := 0

	for _, item := range enabledItems {
		currentWeight += item.GetWeight()
		if random < currentWeight {
			return item
		}
	}

	return enabledItems[0] // Fallback
}

// HealthBasedStrategy health-based selection strategy.
type HealthBasedStrategy struct {
	healthChecker HealthChecker
	rand          *mathrand.Rand
}

// NewHealthBasedStrategy creates a new health-based selection strategy.
func NewHealthBasedStrategy(healthChecker HealthChecker) *HealthBasedStrategy {
	seed := make([]byte, seedByteLen)
	_, err := crand.Read(seed)
	var seedInt uint64
	if err == nil {
		seedInt = binary.LittleEndian.Uint64(seed)
	} else {
		//nolint:gosec // not for security, only for load balancing/random selection
		seedInt = uint64(time.Now().UnixNano() & maxSeedInt)
	}
	return &HealthBasedStrategy{
		healthChecker: healthChecker,
		//nolint:gosec // not for security, only for load balancing/random selection
		rand: mathrand.New(mathrand.NewPCG(seedInt, 0)),
	}
}

// Name returns the name of the health-based strategy.
func (h *HealthBasedStrategy) Name() StrategyType {
	return StrategyHealthBased
}

// Select selects the healthiest item from the list.
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
		var health *HealthCheck
		// 1. 若 item 自身实现了 HealthCheckable，则直接调用
		if hcItem, ok := item.(HealthCheckable); ok {
			health = hcItem.HealthCheck(context.Background())
		} else if h.healthChecker != nil {
			// 2. 否则退回到全局/外部注入的 healthChecker（如果存在）
			health = h.healthChecker.HealthCheck(context.Background())
		}

		// Treat nil/unknown/degraded as unhealthy; 仅保留 healthy
		if health == nil || health.Status == HealthStatusHealthy {
			healthyItems = append(healthyItems, item)
		}
	}

	if len(healthyItems) == 0 {
		return enabledItems[0] // If no healthy items, return first
	}

	// Random selection among healthy items
	index := h.rand.IntN(len(healthyItems))
	return healthyItems[index]
}

// HealthCheckable interface for items that support health checks.
type HealthCheckable interface {
	Selectable
	HealthCheck(ctx context.Context) *HealthCheck
}

// StrategyRegistry strategy registry.
type StrategyRegistry struct {
	strategies map[StrategyType]SelectionStrategy
}

// NewStrategyRegistry creates a new strategy registry.
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

// Register registers a strategy with the given name.
func (r *StrategyRegistry) Register(name StrategyType, strategy SelectionStrategy) {
	r.strategies[name] = strategy
}

// Get retrieves a strategy by name.
func (r *StrategyRegistry) Get(name StrategyType) (SelectionStrategy, bool) {
	strategy, exists := r.strategies[name]
	return strategy, exists
}

// GetDefault returns the default strategy (round-robin).
func (r *StrategyRegistry) GetDefault() SelectionStrategy {
	return r.strategies[StrategyRoundRobin]
}

// GlobalStrategyRegistry is the global strategy registry instance.
//
//nolint:gochecknoglobals // Reason: GlobalStrategyRegistry is a global registry for selection strategies
var GlobalStrategyRegistry = NewStrategyRegistry()

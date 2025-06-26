package utils

import (
	"context"
	"errors"
	"sync"

	"github.com/shellvon/go-sender/core"
)

// FilterEnabled filters enabled items from a slice
func FilterEnabled[T core.Selectable](items []T) []T {
	enabled := make([]T, 0, len(items))
	for i := range items {
		if items[i].IsEnabled() {
			enabled = append(enabled, items[i])
		}
	}
	return enabled
}

// GetStrategy gets strategy from registry or returns default
func GetStrategy(strategyType core.StrategyType) core.SelectionStrategy {
	strategy, exists := core.GlobalStrategyRegistry.Get(strategyType)
	if !exists {
		strategy = core.GlobalStrategyRegistry.GetDefault()
	}
	return strategy
}

// Selector provides a generic way to select items from a collection
// based on context, message, and strategy
type Selector[T core.Selectable] struct {
	items    []T
	strategy core.SelectionStrategy
	mu       sync.RWMutex
}

// NewSelector creates a new selector with the given items and strategy
func NewSelector[T core.Selectable](items []T, strategy core.SelectionStrategy) *Selector[T] {
	return &Selector[T]{
		items:    items,
		strategy: strategy,
	}
}

// Select selects an item based on the following priority:
// 1. Item name from context (if specified)
// 2. Strategy from context (if specified)
// 3. Default strategy
func (s *Selector[T]) Select(ctx context.Context) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. Use item name from context if specified
	if itemName := core.GetItemNameFromCtx(ctx); itemName != "" {
		for i := range s.items {
			if s.items[i].GetName() == itemName {
				return s.items[i]
			}
		}
	}

	// 2. Use strategy from context if specified
	if strategy := core.GetStrategyFromCtx(ctx); strategy != nil {
		selectables := make([]core.Selectable, len(s.items))
		for i, v := range s.items {
			selectables[i] = v
		}
		if selected := strategy.Select(selectables); selected != nil {
			name := selected.GetName()
			for i := range s.items {
				if s.items[i].GetName() == name {
					return s.items[i]
				}
			}
		}
	}

	// 3. Use default strategy
	selectables := make([]core.Selectable, len(s.items))
	for i, v := range s.items {
		selectables[i] = v
	}
	if selected := s.strategy.Select(selectables); selected != nil {
		name := selected.GetName()
		for i := range s.items {
			if s.items[i].GetName() == name {
				return s.items[i]
			}
		}
	}

	var zero T
	return zero
}

// InitProvider initializes a provider with common setup logic
func InitProvider[T core.Selectable](config core.ConfigProvider, items []T) ([]T, *Selector[T], error) {
	// Filter enabled items
	enabledItems := FilterEnabled(items)

	// Check if any enabled items exist
	if len(enabledItems) == 0 {
		return nil, nil, errors.New("no enabled items found")
	}

	// Get strategy and create selector
	strategy := GetStrategy(config.GetStrategy())
	selector := NewSelector(enabledItems, strategy)

	return enabledItems, selector, nil
}

// DefaultStringIfEmpty returns def if s is empty, otherwise returns s.
func DefaultStringIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

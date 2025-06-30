package utils

import (
	"context"
	"errors"

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

// Select selects an item from the provided items based on the following priority:
// 1. Item name from context (if specified)
// 2. Strategy from context (if specified)
// 3. Default strategy
func Select(ctx context.Context, items []core.Selectable, strategy core.SelectionStrategy) core.Selectable {
	// 1. Use item name from context if specified
	if itemName := core.GetItemNameFromCtx(ctx); itemName != "" {
		for i := range items {
			if items[i].GetName() == itemName {
				return items[i]
			}
		}
	}

	// 2. Use strategy from context if specified
	if ctxStrategy := core.GetStrategyFromCtx(ctx); ctxStrategy != nil {
		if selected := ctxStrategy.Select(items); selected != nil {
			return selected
		}
	}

	// 3. Use default strategy
	if selected := strategy.Select(items); selected != nil {
		return selected
	}

	return nil
}

// InitProvider initializes a provider with common setup logic
func InitProvider[T core.Selectable](config core.ConfigProvider, items []T) ([]T, core.SelectionStrategy, error) {
	// Filter enabled items
	enabledItems := FilterEnabled(items)

	// Check if any enabled items exist
	if len(enabledItems) == 0 {
		return nil, nil, errors.New("no enabled items found")
	}

	// Get strategy
	strategy := GetStrategy(config.GetStrategy())

	return enabledItems, strategy, nil
}

// DefaultStringIfEmpty returns def if s is empty, otherwise returns s.
func DefaultStringIfEmpty(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

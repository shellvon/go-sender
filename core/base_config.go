package core

import (
	"context"
	"errors"
	"strconv"
	"sync"
)

// Validatable is implemented by types that can validate their own fields.
// A nil/empty error means the value is valid.
// It is intentionally kept very small so that concrete items (Account, Endpoint…)
// can opt-in to extra validation without forcing every implementation to do so.
type Validatable interface {
	Validate() error
}

// BaseConfig provides a generic implementation of common provider configuration
// concerns: disabled flag, items slice non-empty, duplicate names, at least one
// enabled item, per-item Validate, CRUD helpers and strategy access.
//
// T must satisfy core.Selectable so it can be used by the core load-balancing layer.
// If T also implements Validatable, its Validate method will be invoked as part
// of the global validation.
type BaseConfig[T Selectable] struct {
	ProviderMeta `json:",inline" yaml:",inline"`

	// Items holds all accounts / endpoints / whatever a provider selects from.
	Items []T `json:"items" yaml:"items"`

	mu sync.RWMutex
}

// Compile-time assertion: BaseConfig[T] implements ProviderConfig[T].
type _dummySelectable struct{}

func (_dummySelectable) GetName() string { return "" }
func (_dummySelectable) GetWeight() int  { return 1 }
func (_dummySelectable) IsEnabled() bool { return true }
func (_dummySelectable) GetType() string { return "" }

var _ ProviderConfig[_dummySelectable] = (*BaseConfig[_dummySelectable])(nil)

// Validate performs provider-level and item-level sanity checks.
//
//  1. provider must not be disabled
//  2. slice must be non-empty
//  3. each selectable must have a unique, non-empty name
//  4. at least one selectable must be enabled
//  5. if item implements Validatable it must pass its own validation
func (c *BaseConfig[T]) Validate() error {
	if c.Disabled {
		return errors.New("provider is disabled")
	}
	if len(c.Items) == 0 {
		return errors.New("no items found")
	}

	nameSet := make(map[string]struct{}, len(c.Items))
	hasEnabled := false

	for idx, item := range c.Items {
		name := item.GetName()
		if name == "" {
			return errors.New("item name cannot be empty at index " + strconv.Itoa(idx))
		}
		if _, exists := nameSet[name]; exists {
			return errors.New("duplicate item name: " + name)
		}
		nameSet[name] = struct{}{}

		if item.IsEnabled() {
			hasEnabled = true
		}

		if v, ok := any(item).(Validatable); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	if !hasEnabled {
		return errors.New("all items are disabled")
	}
	return nil
}

// GetItems returns the underlying items slice.
func (c *BaseConfig[T]) GetItems() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]T, len(c.Items))
	copy(out, c.Items)
	return out
}

// GetStrategy proxies to embedded ProviderMeta with sane default.
func (c *BaseConfig[T]) GetStrategy() StrategyType { return c.ProviderMeta.GetStrategy() }

// ---- CRUD helpers (thread-safe) -------------------------------------------

// Add inserts a new item after passing validation for duplicates etc.
func (c *BaseConfig[T]) Add(item T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// cheap duplicate check
	for _, it := range c.Items {
		if it.GetName() == item.GetName() {
			return errors.New("duplicate item name: " + item.GetName())
		}
	}
	// per-item validation if implemented
	if v, ok := any(item).(Validatable); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	c.Items = append(c.Items, item)
	return nil
}

// Delete removes an item by name. It is a no-op if the name is unknown.
func (c *BaseConfig[T]) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, it := range c.Items {
		if it.GetName() == name {
			// delete while preserving order
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

// Update replaces an existing item with the same name.
func (c *BaseConfig[T]) Update(item T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, it := range c.Items {
		if it.GetName() == item.GetName() {
			if v, ok := any(item).(Validatable); ok {
				if err := v.Validate(); err != nil {
					return err
				}
			}
			c.Items[i] = item
			return nil
		}
	}
	return errors.New("item not found: " + item.GetName())
}

// Select selects an item from the config based on the following priority:
// 0. Filter items by filter function to get a subset of items
// 1. Item name from context (if specified)
// 2. Strategy from context (if specified)
// 3. Default strategy.
func (c *BaseConfig[T]) Select(ctx context.Context, filter func(T) bool) (T, error) {
	c.mu.RLock()
	itemsCopy := make([]T, len(c.Items))
	copy(itemsCopy, c.Items)
	c.mu.RUnlock()

	var filtered []T
	for _, item := range itemsCopy {
		if filter == nil || filter(item) {
			filtered = append(filtered, item)
		}
	}
	var zero T
	if len(filtered) == 0 {
		return zero, errors.New("no available config after filtering")
	}
	if len(filtered) == 1 {
		if !filtered[0].IsEnabled() {
			return zero, errors.New("selected item is disabled")
		}
		return filtered[0], nil
	}
	ri := GetRoute(ctx)
	stType := c.GetStrategy()
	// 1. 指定账号
	if ri != nil && ri.AccountName != "" {
		return c.findEnabledByName(ri.AccountName, filtered)
	}
	if ri != nil && ri.StrategyType != "" {
		stType = ri.StrategyType
	}

	// 3. 必须有策略
	strategy, _ := GlobalStrategyRegistry.Get(stType)
	if strategy == nil {
		return zero, errors.New("no strategy specified or unknown strategy")
	}

	selected := strategy.Select(toSelectables(filtered))
	if selected != nil {
		return c.findEnabledByName(selected.GetName(), filtered)
	}
	return zero, errors.New("no config selected by strategy")
}

// findEnabledByName returns the first enabled item with matching name.
func (c *BaseConfig[T]) findEnabledByName(name string, items []T) (T, error) {
	var zero T
	for _, item := range items {
		if item.GetName() == name {
			if !item.IsEnabled() {
				return zero, errors.New("selected item is disabled")
			}
			return item, nil
		}
	}
	return zero, errors.New("no item found with the given name")
}

// toSelectables converts []T to []Selectable.
func toSelectables[T Selectable](items []T) []Selectable {
	selectables := make([]Selectable, len(items))
	for i, item := range items {
		selectables[i] = item
	}
	return selectables
}

// ProviderConfig is a generic interface for provider configuration.
type ProviderConfig[T any] interface {
	Validate() error
	GetItems() []T
	GetStrategy() StrategyType
}

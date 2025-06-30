package core

import (
	"testing"
	"time"

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

func (m *MockSelectable) GetType() string {
	return ""
}

func TestRoundRobinStrategy(t *testing.T) {
	items := []Selectable{
		&MockSelectable{name: "item1", weight: 1, enabled: true},
		&MockSelectable{name: "item2", weight: 1, enabled: true},
		&MockSelectable{name: "item3", weight: 1, enabled: true},
	}

	strategy := NewRoundRobinStrategy()

	// Test round robin selection
	selected1 := strategy.Select(items)
	assert.Equal(t, "item1", selected1.GetName())

	selected2 := strategy.Select(items)
	assert.Equal(t, "item2", selected2.GetName())

	selected3 := strategy.Select(items)
	assert.Equal(t, "item3", selected3.GetName())

	// Should cycle back to first
	selected4 := strategy.Select(items)
	assert.Equal(t, "item1", selected4.GetName())
}

func TestRandomStrategy(t *testing.T) {
	items := []Selectable{
		&MockSelectable{name: "item1", weight: 1, enabled: true},
		&MockSelectable{name: "item2", weight: 1, enabled: true},
		&MockSelectable{name: "item3", weight: 1, enabled: true},
	}

	strategy := NewRandomStrategy()

	// Test random selection (should select one of the items)
	selected := strategy.Select(items)
	assert.NotNil(t, selected)
	assert.Contains(t, []string{"item1", "item2", "item3"}, selected.GetName())
}

func TestWeightedStrategy(t *testing.T) {
	items := []Selectable{
		&MockSelectable{name: "item1", weight: 1, enabled: true},
		&MockSelectable{name: "item2", weight: 2, enabled: true},
		&MockSelectable{name: "item3", weight: 3, enabled: true},
	}

	strategy := NewWeightedStrategy()

	// Test weighted selection
	selected := strategy.Select(items)
	assert.NotNil(t, selected)
	assert.Contains(t, []string{"item1", "item2", "item3"}, selected.GetName())
}

func TestHealthBasedStrategy(t *testing.T) {
	items := []Selectable{
		&MockSelectable{name: "item1", weight: 1, enabled: true},
		&MockSelectable{name: "item2", weight: 1, enabled: true},
		&MockSelectable{name: "item3", weight: 1, enabled: true},
	}

	strategy := NewHealthBasedStrategy(nil)

	// Test health-based selection
	selected := strategy.Select(items)
	assert.NotNil(t, selected)
	assert.Contains(t, []string{"item1", "item2", "item3"}, selected.GetName())
}

func TestStrategyRegistry(t *testing.T) {
	registry := NewStrategyRegistry()

	// Register strategies
	registry.Register(StrategyRoundRobin, NewRoundRobinStrategy())
	registry.Register(StrategyRandom, NewRandomStrategy())
	registry.Register(StrategyWeighted, NewWeightedStrategy())
	registry.Register(StrategyHealthBased, NewHealthBasedStrategy(nil))

	// Test getting strategies
	roundRobin, exists := registry.Get(StrategyRoundRobin)
	assert.True(t, exists)
	assert.NotNil(t, roundRobin)

	random, exists := registry.Get(StrategyRandom)
	assert.True(t, exists)
	assert.NotNil(t, random)

	weighted, exists := registry.Get(StrategyWeighted)
	assert.True(t, exists)
	assert.NotNil(t, weighted)

	healthBased, exists := registry.Get(StrategyHealthBased)
	assert.True(t, exists)
	assert.NotNil(t, healthBased)

	// Test getting non-existent strategy
	nonExistent, exists := registry.Get("non-existent")
	assert.False(t, exists)
	assert.Nil(t, nonExistent)
}

func TestRetryPolicy(t *testing.T) {
	policy := &RetryPolicy{
		MaxAttempts:   3,
		InitialDelay:  100,
		MaxDelay:      1000,
		BackoffFactor: 2.0,
		Filter: func(attempt int, err error) bool {
			return attempt < 3
		},
	}

	// Test validation
	err := policy.Validate()
	assert.NoError(t, err)

	// Test reset
	policy.Reset()
	assert.Equal(t, 0, policy.currentAttempt)

	// Test should retry
	assert.True(t, policy.ShouldRetry(1, assert.AnError))
	assert.True(t, policy.ShouldRetry(2, assert.AnError))
	assert.False(t, policy.ShouldRetry(3, assert.AnError))

	// Test next delay
	delay1 := policy.NextDelay(1, assert.AnError)
	assert.Equal(t, time.Duration(200), delay1) // 100 * 1 * 2.0

	delay2 := policy.NextDelay(2, assert.AnError)
	assert.Equal(t, time.Duration(400), delay2) // 100 * 2 * 2.0
}

func TestRetryPolicyValidation(t *testing.T) {
	tests := []struct {
		name    string
		policy  *RetryPolicy
		wantErr bool
	}{
		{
			name: "valid policy",
			policy: &RetryPolicy{
				MaxAttempts:   3,
				InitialDelay:  100,
				MaxDelay:      1000,
				BackoffFactor: 2.0,
			},
			wantErr: false,
		},
		{
			name: "negative max attempts",
			policy: &RetryPolicy{
				MaxAttempts:   -1,
				InitialDelay:  100,
				MaxDelay:      1000,
				BackoffFactor: 2.0,
			},
			wantErr: true,
		},
		{
			name: "negative initial delay",
			policy: &RetryPolicy{
				MaxAttempts:   3,
				InitialDelay:  -100,
				MaxDelay:      1000,
				BackoffFactor: 2.0,
			},
			wantErr: true,
		},
		{
			name: "negative max delay",
			policy: &RetryPolicy{
				MaxAttempts:   3,
				InitialDelay:  100,
				MaxDelay:      -1000,
				BackoffFactor: 2.0,
			},
			wantErr: true,
		},
		{
			name: "negative backoff factor",
			policy: &RetryPolicy{
				MaxAttempts:   3,
				InitialDelay:  100,
				MaxDelay:      1000,
				BackoffFactor: -2.0,
			},
			wantErr: true,
		},
		{
			name: "max delay less than initial delay",
			policy: &RetryPolicy{
				MaxAttempts:   3,
				InitialDelay:  1000,
				MaxDelay:      100,
				BackoffFactor: 2.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewRetryPolicy(t *testing.T) {
	policy := NewRetryPolicy(
		WithRetryMaxAttempts(5),
		WithRetryInitialDelay(200),
		WithRetryMaxDelay(2000),
		WithRetryBackoffFactor(3.0),
	)

	assert.Equal(t, 5, policy.MaxAttempts)
	assert.Equal(t, time.Duration(200), policy.InitialDelay)
	assert.Equal(t, time.Duration(2000), policy.MaxDelay)
	assert.Equal(t, 3.0, policy.BackoffFactor)
}

func TestDefaultRetryFilter(t *testing.T) {
	retryableErrors := []error{assert.AnError}

	filter := DefaultRetryFilter(retryableErrors, true)

	// Test with retryable error
	assert.True(t, filter(1, assert.AnError))

	// Test with non-retryable error
	assert.False(t, filter(1, nil))
}

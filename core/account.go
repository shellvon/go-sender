package core

// BaseConfig provides common configuration fields for all providers
type BaseConfig struct {
	Disabled bool         `json:"disabled"` // Whether the provider is enabled
	Strategy StrategyType `json:"strategy"` // Selection strategy
}

// GetStrategy returns the selection strategy, defaulting to round_robin
func (c *BaseConfig) GetStrategy() StrategyType {
	if c.Strategy == "" {
		return StrategyRoundRobin
	}
	return c.Strategy
}

// IsDisabled returns whether the provider is disabled
func (c *BaseConfig) IsDisabled() bool {
	return c.Disabled
}

// Account represents a unified account configuration that uses a single key field
type Account struct {
	Name     string `json:"name" yaml:"name"`
	Weight   int    `json:"weight" yaml:"weight"`
	Disabled bool   `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Key      string `json:"key" yaml:"key"`
	Webhook  string `json:"webhook,omitempty" yaml:"webhook,omitempty"`
}

func (a *Account) IsEnabled() bool {
	return !a.Disabled
}

func (a *Account) GetName() string {
	return a.Name
}

func (a *Account) GetWeight() int {
	if a.Weight <= 0 {
		return 1
	}
	return a.Weight
}

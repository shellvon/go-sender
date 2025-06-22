package serverchan

import "github.com/shellvon/go-sender/core"

// Config holds configuration for the Server酱 provider
type Config struct {
	Disabled bool              `json:"disabled"` // Whether the provider is disabled
	Accounts []Account         `json:"accounts"` // Multiple account configuration
	Strategy core.StrategyType `json:"strategy"` // Selection strategy
}

type Account struct {
	Name     string `json:"name"`               // Account name
	Key      string `json:"key"`                // SendKey
	Weight   int    `json:"weight"`             // Weight for weighted strategy
	Disabled bool   `json:"disabled,omitempty"` // Whether this account is disabled
}

func (c *Config) GetStrategy() core.StrategyType {
	if c.Strategy == "" {
		return core.StrategyRoundRobin
	}
	return c.Strategy
}

func (c Config) IsConfigured() bool {
	return !c.Disabled && len(c.Accounts) > 0
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

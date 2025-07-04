package email

import "github.com/shellvon/go-sender/core"

// Account is defined in account.go

// Config holds the configuration for the email provider.
type Config struct {
	core.ProviderMeta

	Accounts []*Account `json:"accounts"` // list of SMTP accounts
}

// IsConfigured checks if the email provider has a valid configuration.
func (c *Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

// GetStrategy returns the selection strategy, defaulting to round_robin.
func (c *Config) GetStrategy() core.StrategyType {
	if c.Strategy == "" {
		return core.StrategyRoundRobin
	}
	return c.Strategy
}

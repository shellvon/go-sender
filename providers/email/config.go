package email

import "github.com/shellvon/go-sender/core"

// Account represents a single email account configuration.
type Account struct {
	Name     string `json:"name"`               // A unique name for the account.
	Type     string `json:"type,omitempty"`     // Subprovider type identifier for account selection filtering
	Host     string `json:"host"`               // SMTP server host.
	Port     int    `json:"port"`               // SMTP server port.
	Username string `json:"username"`           // SMTP username.
	Password string `json:"password"`           // SMTP password.
	From     string `json:"from"`               // Default 'From' address.
	Weight   int    `json:"weight"`             // Weight for the weighted selection strategy.
	Disabled bool   `json:"disabled,omitempty"` // Whether this account is disabled (default: false, i.e., enabled by default)
}

// Config holds the configuration for the email provider.
type Config struct {
	core.BaseConfig

	Accounts []Account `json:"accounts"` // A list of email accounts to use.
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

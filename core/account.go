// Package core provides the core functionality for the go-sender library.
package core

// BaseConfig provides common configuration fields for all providers.
type BaseConfig struct {
	Disabled bool         `json:"disabled"` // Whether the provider is enabled
	Strategy StrategyType `json:"strategy"` // Selection strategy
}

// GetStrategy returns the selection strategy, defaulting to round_robin.
func (c *BaseConfig) GetStrategy() StrategyType {
	if c.Strategy == "" {
		return StrategyRoundRobin
	}
	return c.Strategy
}

// IsDisabled returns whether the provider is disabled.
func (c *BaseConfig) IsDisabled() bool {
	return c.Disabled
}

// Account represents a provider account configuration.
type Account struct {
	Name         string `json:"name"                    yaml:"name"`                    // Unique identifier for the account, used for selection and identification
	Type         string `json:"type,omitempty"          yaml:"type,omitempty"`          // Subprovider type identifier for account selection filtering (e.g., aliyun, tencent, emailjs, etc.)
	Weight       int    `json:"weight"                  yaml:"weight"`                  // Weight value for weighted selection strategy (higher weight = higher priority)
	Disabled     bool   `json:"disabled,omitempty"      yaml:"disabled,omitempty"`      // Whether this account is disabled (disabled accounts won't be selected)
	Key          string `json:"key"                     yaml:"key"`                     // Primary authentication key (API key, App ID, username, etc.)
	Secret       string `json:"secret,omitempty"        yaml:"secret,omitempty"`        // Secondary authentication secret (API secret, App Secret, password, etc.)
	Webhook      string `json:"webhook,omitempty"       yaml:"webhook,omitempty"`       // Webhook URL for callbacks or notifications
	Endpoint     string `json:"endpoint,omitempty"      yaml:"endpoint,omitempty"`      // Custom API endpoint URL (overrides default provider endpoint)
	IntlEndpoint string `json:"intl_endpoint,omitempty" yaml:"intl_endpoint,omitempty"` // Custom international API endpoint URL (for international services)
	From         string `json:"from,omitempty"          yaml:"from,omitempty"`          // Sender identifier (email address, phone number, bot token, etc.)
}

// IsEnabled returns whether the account is enabled.
func (a *Account) IsEnabled() bool {
	return !a.Disabled
}

// GetName returns the account name.
func (a *Account) GetName() string {
	return a.Name
}

// GetWeight returns the account weight for load balancing.
func (a *Account) GetWeight() int {
	if a.Weight <= 0 {
		return 1
	}
	return a.Weight
}

// GetType returns the subprovider type of this account.
func (a *Account) GetType() string {
	return a.Type
}

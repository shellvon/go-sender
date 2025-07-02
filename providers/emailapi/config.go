package emailapi

import "github.com/shellvon/go-sender/core"

// Config holds configuration for the emailapi provider and all supported API email services.
type Config struct {
	// BaseConfig contains common configuration like strategy and disabled flag
	core.BaseConfig

	Accounts []*core.Account `json:"accounts"`
}

// IsConfigured checks if the EmailAPI configuration is valid and ready to use.
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

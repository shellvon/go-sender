package telegram

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Telegram provider
type Config struct {
	core.BaseConfig
	Accounts []core.Account `json:"accounts"` // Multiple account configurations
}

// IsConfigured checks if the Telegram configuration is valid
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

package serverchan

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the ServerChan provider.
type Config struct {
	core.ProviderMeta

	Accounts []*Account `json:"accounts"` // Multiple account configuration
}

// IsConfigured checks if the ServerChan configuration is valid.
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

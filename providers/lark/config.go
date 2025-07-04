package lark

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Lark provider.
type Config struct {
	core.ProviderMeta

	Accounts []*Account `json:"accounts"` // Multiple accounts configuration
}

// IsConfigured checks if the Lark configuration is valid.
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

package dingtalk

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the DingTalk provider
type Config struct {
	core.BaseConfig
	Accounts []core.Account `json:"accounts"` // Multiple accounts configuration
}

// IsConfigured checks if the DingTalk configuration is valid
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Accounts) > 0
}

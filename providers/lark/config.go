package lark

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Lark/Feishu Robot provider
type Config struct {
	Disabled bool              `json:"disabled"` // Whether the lark provider is enabled (default: false)
	Bots     []Bot             `json:"bots"`     // Multiple bots configuration
	Strategy core.StrategyType `json:"strategy"` // Selection strategy: round_robin, random, weighted
}

// Bot represents a single Lark/Feishu bot configuration
type Bot struct {
	Name     string `json:"name"`               // Bot name
	Webhook  string `json:"webhook"`            // Lark webhook URL
	Secret   string `json:"secret,omitempty"`   // Lark webhook secret (optional)
	Weight   int    `json:"weight"`             // Weight for weighted strategy
	Disabled bool   `json:"disabled,omitempty"` // Whether this bot is disabled (default: false)
}

// GetStrategy returns the selection strategy, defaulting to round_robin.
func (c *Config) GetStrategy() core.StrategyType {
	if c.Strategy == "" {
		return core.StrategyRoundRobin
	}
	return c.Strategy
}

// IsConfigured checks if the Lark configuration is valid
func (c Config) IsConfigured() bool {
	return !c.Disabled && len(c.Bots) > 0
}

func (b *Bot) IsEnabled() bool {
	return !b.Disabled
}

func (b *Bot) GetName() string {
	return b.Name
}

func (b *Bot) GetWeight() int {
	if b.Weight <= 0 {
		return 1
	}
	return b.Weight
}

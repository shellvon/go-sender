package telegram

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Telegram provider
type Config struct {
	Disabled          bool              `json:"disabled"`           // Whether the telegram provider is enabled (default: false)
	Bots              []Bot             `json:"bots"`               // Multiple bot configurations
	Strategy          core.StrategyType `json:"strategy"`           // Selection strategy: round_robin, random, weighted
	DisableMiddleware bool              `json:"disable_middleware"` // Whether to disable middleware (for embedded usage)
}

// Bot represents a single Telegram bot configuration
type Bot struct {
	Name     string `json:"name"`     // Bot name for selection
	Token    string `json:"token"`    // Telegram bot token
	Weight   int    `json:"weight"`   // Weight for weighted strategy
	Disabled bool   `json:"disabled"` // Whether this bot is disabled (default: false)
}

// GetStrategy returns the selection strategy, defaulting to round_robin.
func (c *Config) GetStrategy() core.StrategyType {
	if c.Strategy == "" {
		return core.StrategyRoundRobin
	}
	return c.Strategy
}

// IsConfigured checks if the Telegram configuration is valid
func (c Config) IsConfigured() bool {
	return !c.Disabled && len(c.Bots) > 0
}

// IsEnabled checks if the bot is enabled
func (b *Bot) IsEnabled() bool {
	return !b.Disabled
}

// GetName returns the bot name for strategy selection
func (b *Bot) GetName() string {
	return b.Name
}

// GetWeight returns the bot weight for strategy selection
func (b *Bot) GetWeight() int {
	if b.Weight <= 0 {
		return 1
	}
	return b.Weight
}

// Package telegram provides telegram notification support for go-sender.
//
// This package implements the go-sender Provider interface for telegram
// messaging service, enabling seamless integration with the go-sender
// notification system.
//
// Basic usage:
//
//	provider, err := telegram.NewProvider([]*telegram.Account{account})
//	msg := telegram.Text().Content("Hello World").Build()
//	provider.Send(context.Background(), msg, nil)
package telegram

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

// Provider implements the Telegram provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var (
	_ core.Provider = (*Provider)(nil)
)

// New creates a new Telegram provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeTelegram),
		newTelegramTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies Telegram Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new Telegram provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := telegram.NewProvider([]*telegram.Account{account1, account2},
//	    telegram.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeTelegram,
		func(meta core.ProviderMeta, items []*Account) *Config {
			return &Config{
				ProviderMeta: meta,
				Items:        items,
			}
		},
		New,
		opts...,
	)
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return string(core.ProviderTypeTelegram)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: telegram.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*telegram.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)

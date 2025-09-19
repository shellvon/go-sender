// Package lark provides Lark/Feishu notification support for go-sender.
//
// This package enables sending rich messages to Lark (Feishu) groups and users
// via webhook bots, supporting text, post, image, interactive, and share chat messages.
// It provides comprehensive support for Lark's messaging capabilities with
// automatic message type routing and formatting.
//
// Basic usage:
//
//	account := lark.NewAccount("your-webhook-token")
//	provider, err := lark.NewProvider([]*lark.Account{account})
//	msg := lark.Text().Content("Hello from go-sender!").Build()
//	provider.Send(context.Background(), msg, nil)
//
// For more examples, see the package README and examples directory.
package lark

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

// Provider implements the Lark provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new Lark provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeLark),
		newLarkTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies Lark Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new Lark provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := lark.NewProvider([]*lark.Account{account1, account2},
//	    lark.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeLark,
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
	return string(core.ProviderTypeLark)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: lark.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*lark.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)

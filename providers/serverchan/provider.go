// Package serverchan provides ServerChan notification support for go-sender.
//
// This package enables sending notifications via ServerChan (Server酱), a popular
// WeChat notification service for developers. It supports sending text messages
// with optional markdown formatting to WeChat through ServerChan's API.
//
// Basic usage:
//
//	account := serverchan.NewAccount("your-sckey")
//	provider, err := serverchan.NewProvider([]*serverchan.Account{account})
//	msg := serverchan.NewMessage("Alert", "Something happened!")
//	provider.Send(context.Background(), msg, nil)
//
// For more examples, see the package README and examples directory.
package serverchan

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

// Provider implements the ServerChan provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new ServerChan provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeServerChan),
		newServerChanTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies ServerChan Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new ServerChan provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := serverchan.NewProvider([]*serverchan.Account{account1, account2},
//	    serverchan.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeServerChan,
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
	return string(core.ProviderTypeServerChan)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: serverchan.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*serverchan.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)

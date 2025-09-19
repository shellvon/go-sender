// Package dingtalk provides DingTalk notification support for go-sender.
//
// This package enables sending rich messages to DingTalk groups via custom robots,
// supporting text, markdown, link cards, action cards, and feed cards with optional
// security signature verification.
//
// Basic usage:
//
//	provider, err := dingtalk.NewProvider([]*dingtalk.Account{
//	    dingtalk.NewAccount("your-access-token"),
//	})
//	msg := dingtalk.Text().Content("Hello World").Build()
//	provider.Send(context.Background(), msg, nil)
//
// For more examples, see the package README and examples directory.
package dingtalk

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

// Provider implements the DingTalk provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new DingTalk provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeDingtalk),
		newDingTalkTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies DingTalk Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new DingTalk provider with the given accounts and options.
//
// At least one account is required.
//
// Example:
//
//	provider, err := dingtalk.NewProvider([]*dingtalk.Account{account1, account2},
//	    dingtalk.Strategy(core.StrategyWeighted))
func NewProvider(accounts []*Account, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		accounts,
		core.ProviderTypeDingtalk,
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
	return string(core.ProviderTypeDingtalk)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: dingtalk.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*dingtalk.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)

// Package webhook provides generic webhook notification support for go-sender.
//
// This package enables sending HTTP requests to arbitrary webhook endpoints,
// supporting various HTTP methods, headers, and body formats. It provides
// flexible configuration for custom integrations and third-party services
// that accept webhook notifications.
//
// Basic usage:
//
//	endpoint := webhook.NewEndpoint("https://hooks.example.com/webhook", "POST")
//	provider, err := webhook.NewProvider([]*webhook.Endpoint{endpoint})
//	msg := webhook.Webhook().Body([]byte(`{"message": "Hello!"}`)).Build()
//	provider.Send(context.Background(), msg, nil)
//
// For more examples, see the package README and examples directory.
package webhook

import (
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// Provider implements the Webhook provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Endpoint]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new Webhook provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeWebhook),
		newWebhookTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// ProviderOption represents a function that modifies Webhook Provider configuration.
type ProviderOption func(*Config)

// NewProvider creates a new Webhook provider with the given endpoints and options.
//
// At least one endpoint is required.
//
// Example:
//
//	provider, err := webhook.NewProvider([]*webhook.Endpoint{endpoint1, endpoint2},
//	    webhook.Strategy(core.StrategyWeighted))
func NewProvider(endpoints []*Endpoint, opts ...ProviderOption) (*Provider, error) {
	return core.CreateProvider(
		endpoints,
		core.ProviderTypeWebhook,
		func(meta core.ProviderMeta, items []*Endpoint) *Config {
			return &Config{
				ProviderMeta: meta,
				Items:        items,
			}
		},
		New,
		opts...,
	)
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

// Re-exported core provider options for cleaner API
// These provide convenient aliases: webhook.Strategy(core.StrategyWeighted) instead of core.WithStrategy[*webhook.Config](core.StrategyWeighted).
var (
	Strategy         = core.WithStrategy[*Config]
	ProviderDisabled = core.WithProviderDisabled[*Config]
)

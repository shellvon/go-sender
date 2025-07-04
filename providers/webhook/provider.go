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
		&webhookTransformer{},
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

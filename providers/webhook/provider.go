package webhook

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Webhook provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Endpoint]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new Webhook provider instance.
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("webhook provider is not configured or is disabled")
	}

	enabledEndpoints, _, err := utils.InitProvider(&config, config.Endpoints)
	if err != nil {
		return nil, errors.New("no enabled webhook endpoints found")
	}

	// 获取策略
	strategy := utils.GetStrategy(config.GetStrategy())

	// 创建泛型 provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeWebhook),
		enabledEndpoints,
		&webhookTransformer{},
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the name of the provider.
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

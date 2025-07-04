package dingtalk

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the DingTalk provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new DingTalk provider instance.
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("dingtalk provider is not configured or is disabled")
	}
	strategy := utils.GetStrategy(config.GetStrategy())

	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeDingtalk),
		config.Accounts,
		newDingTalkTransformer(),
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return string(core.ProviderTypeDingtalk)
}

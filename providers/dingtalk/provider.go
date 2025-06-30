package dingtalk

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the DingTalk provider using generic base
type Provider struct {
	*providers.HTTPProvider[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new DingTalk provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("dingtalk provider is not configured or is disabled")
	}

	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled dingtalk accounts found")
	}

	// Get strategy
	strategy := utils.GetStrategy(config.GetStrategy())

	// Create generic provider with transformer from transformer.go
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeDingtalk),
		enabledAccounts,
		newDingTalkTransformer(),
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeDingtalk)
}

package lark

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Lark provider using generic base
type Provider struct {
	*providers.HTTPProvider[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new Lark provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("lark provider is not configured or is disabled")
	}

	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled lark accounts found")
	}

	// Get strategy
	strategy := utils.GetStrategy(config.GetStrategy())

	// Create generic provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeLark),
		enabledAccounts,
		newLarkTransformer(),
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeLark)
}

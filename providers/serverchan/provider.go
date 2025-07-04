package serverchan

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the ServerChan provider using generic base.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new ServerChan provider instance.
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("serverchan provider is not configured or is disabled")
	}

	accounts := config.Accounts

	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled serverchan accounts found")
	}

	// Get strategy
	strategy := utils.GetStrategy(config.GetStrategy())

	// Create generic provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeServerChan),
		enabledAccounts,
		newTransformer(),
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return string(core.ProviderTypeServerChan)
}

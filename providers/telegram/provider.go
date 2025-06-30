package telegram

import (
	"errors"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the Telegram provider using generic base
type Provider struct {
	*providers.HTTPProvider[*core.Account]
}

var (
	_ core.Provider = (*Provider)(nil)
)

// New creates a new Telegram provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("telegram provider is not configured or is disabled")
	}

	accounts := make([]*core.Account, len(config.Accounts))
	for i := range config.Accounts {
		accounts[i] = &config.Accounts[i]
	}

	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled telegram accounts found")
	}

	// 获取策略
	strategy := utils.GetStrategy(config.GetStrategy())

	// 创建泛型 provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeTelegram),
		enabledAccounts,
		&telegramTransformer{},
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeTelegram)
}

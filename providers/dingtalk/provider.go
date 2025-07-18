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

// Name returns the provider name.
func (p *Provider) Name() string {
	return string(core.ProviderTypeDingtalk)
}

package emailapi

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/utils"
)

// Provider is the main emailapi provider, supporting multiple API-based email services.
type Provider struct {
	*providers.HTTPProvider[*core.Account]
}

var _ core.Provider = (*Provider)(nil)

// transformerRegistry global transformer registry.
//
//nolint:gochecknoglobals // Reason: transformerRegistry is a global registry for emailapi transformers
var transformerRegistry = make(map[string]core.HTTPTransformer[*core.Account])

// registryMutex global mutex for transformerRegistry.
//
//nolint:gochecknoglobals // Reason: registryMutex is a global mutex for transformerRegistry
var registryMutex sync.RWMutex

// RegisterTransformer registers transformer to global registry.
func RegisterTransformer(subProvider string, transformer core.HTTPTransformer[*core.Account]) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	transformerRegistry[subProvider] = transformer
}

// GetTransformer gets transformer from registry.
func GetTransformer(subProvider string) (core.HTTPTransformer[*core.Account], bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	transformer, exists := transformerRegistry[subProvider]
	return transformer, exists
}

// emailAPITransformer implements core.HTTPTransformer[*core.Account], selects specific transformer based on SubProvider.
type emailAPITransformer struct{}

// CanTransform checks if this is an EmailAPI message.
func (t *emailAPITransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeEmailAPI
}

// Transform gets specific transformer from registry based on SubProvider for conversion.
func (t *emailAPITransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for emailapi transformer: %T", msg)
	}

	// Get transformer from registry
	transformer, exists := GetTransformer(emailMsg.SubProvider)
	if !exists {
		return nil, nil, fmt.Errorf("unsupported EmailAPI sub-provider: %s", emailMsg.SubProvider)
	}

	return transformer.Transform(ctx, msg, account)
}

// New creates a new emailapi Provider with the given config.
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("emailapi provider is not configured or is disabled")
	}

	// Convert config accounts to core.Account
	accounts := make([]*core.Account, len(config.Accounts))
	for i, acc := range config.Accounts {
		accounts[i] = &core.Account{
			Name:         acc.Name,
			Key:          acc.APIKey,
			Secret:       acc.APISecret,
			From:         acc.From,
			Endpoint:     acc.Domain,
			IntlEndpoint: "",               // EmailAPI doesn't have international endpoints
			Webhook:      "",               // EmailAPI doesn't have webhook field
			Type:         string(acc.Type), // Sub-provider type
			Weight:       acc.Weight,
			Disabled:     acc.Disabled,
		}
	}

	enabledAccounts, _, err := utils.InitProvider(&config, accounts)
	if err != nil {
		return nil, errors.New("no enabled emailapi accounts found")
	}

	// Create strategy
	strategy := utils.GetStrategy(config.Strategy)

	// Create generic provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeEmailAPI),
		enabledAccounts,
		&emailAPITransformer{},
		strategy,
	)

	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider type.
func (p *Provider) Name() string {
	return string(core.ProviderTypeEmailAPI)
}

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
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// transformerRegistry global transformer registry.
//
//nolint:gochecknoglobals // Reason: transformerRegistry is a global registry for emailapi transformers
var transformerRegistry = make(map[string]core.HTTPTransformer[*Account])

// registryMutex global mutex for transformerRegistry.
//
//nolint:gochecknoglobals // Reason: registryMutex is a global mutex for transformerRegistry
var registryMutex sync.RWMutex

// RegisterTransformer registers transformer to global registry.
func RegisterTransformer(subProvider string, transformer core.HTTPTransformer[*Account]) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	transformerRegistry[subProvider] = transformer
}

// GetTransformer gets transformer from registry.
func GetTransformer(subProvider string) (core.HTTPTransformer[*Account], bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	transformer, exists := transformerRegistry[subProvider]
	return transformer, exists
}

// emailAPITransformer implements core.HTTPTransformer[*Account], selects specific transformer based on SubProvider.
type emailAPITransformer struct{}

// CanTransform checks if this is an EmailAPI message.
func (t *emailAPITransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeEmailAPI
}

// Transform gets specific transformer from registry based on SubProvider for conversion.
func (t *emailAPITransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for emailapi transformer: %T", msg)
	}

	if emailMsg.SubProvider == "" {
		return nil, nil, errors.New("sub-provider is required for emailapi transformer")
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

	// Create generic provider
	httpProvider := providers.NewHTTPProvider(
		string(core.ProviderTypeEmailAPI),
		config.Accounts,
		&emailAPITransformer{},
		utils.GetStrategy(config.Strategy),
	)
	return &Provider{
		HTTPProvider: httpProvider,
	}, nil
}

// Name returns the provider type.
func (p *Provider) Name() string {
	return string(core.ProviderTypeEmailAPI)
}

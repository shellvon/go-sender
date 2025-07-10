package emailapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// Config is the type alias for core.BaseConfig[*Account].
type Config = core.BaseConfig[*Account]

// Provider is the main emailapi provider, supporting multiple API-based email services.
type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// emailAPIRegistry is a shared registry for emailapi sub-provider transformers.
//
//nolint:gochecknoglobals // Global registry is acceptable for package-level look-ups.
var emailAPIRegistry = providers.NewTransformerRegistry[*Account]()

// RegisterTransformer registers transformer to the package registry.
func RegisterTransformer(subProvider string, transformer core.HTTPTransformer[*Account]) {
	emailAPIRegistry.Register(subProvider, transformer)
}

// GetTransformer gets transformer from the package registry.
func GetTransformer(subProvider string) (core.HTTPTransformer[*Account], bool) {
	return emailAPIRegistry.Get(subProvider)
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
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
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

// newEmailAPITransformer constructs a new emailAPITransformer.
func newEmailAPITransformer() core.HTTPTransformer[*Account] {
	return &emailAPITransformer{}
}

// New creates a new emailapi Provider with the given config.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeEmailAPI),
		newEmailAPITransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

// Name returns the provider type.
func (p *Provider) Name() string {
	return string(core.ProviderTypeEmailAPI)
}

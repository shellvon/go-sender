package providers

import (
	"fmt"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
)

// ProviderBuilder defines the interface for building providers from configuration
type ProviderBuilder interface {
	// BuildProvider creates a provider instance from parsed accounts
	BuildProvider(accounts []any) (core.Provider, error)
	// GetProviderType returns the provider type this builder handles
	GetProviderType() core.ProviderType
}

// MessageBuilder defines the interface for building messages from CLI flags
type MessageBuilder interface {
	// BuildMessage creates a message from CLI flags
	BuildMessage(flags *cli.CLIFlags) (core.Message, error)
	// GetProviderType returns the provider type this builder handles
	GetProviderType() core.ProviderType
	// ValidateFlags validates CLI flags for this provider
	ValidateFlags(flags *cli.CLIFlags) error
}

// ProviderRegistry manages provider and message builders
type ProviderRegistry struct {
	providerBuilders map[core.ProviderType]ProviderBuilder
	messageBuilders  map[core.ProviderType]MessageBuilder
}

// NewProviderRegistry creates a new provider registry
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providerBuilders: make(map[core.ProviderType]ProviderBuilder),
		messageBuilders:  make(map[core.ProviderType]MessageBuilder),
	}
}

// RegisterProviderBuilder registers a provider builder
func (r *ProviderRegistry) RegisterProviderBuilder(builder ProviderBuilder) {
	r.providerBuilders[builder.GetProviderType()] = builder
}

// RegisterMessageBuilder registers a message builder
func (r *ProviderRegistry) RegisterMessageBuilder(builder MessageBuilder) {
	r.messageBuilders[builder.GetProviderType()] = builder
}

// BuildProviders builds all providers from configuration and registers them with sender
func (r *ProviderRegistry) BuildProviders(sender *gosender.Sender, config *cli.RootConfig) error {
	// Group accounts by provider type
	accountsByType := make(map[core.ProviderType][]any)

	for _, account := range config.Accounts {
		providerStr, ok := account["provider"].(string)
		if !ok {
			continue
		}

		providerType := core.ProviderType(providerStr)
		accountsByType[providerType] = append(accountsByType[providerType], account)
	}

	// Build and register each provider
	for providerType, accounts := range accountsByType {
		builder, exists := r.providerBuilders[providerType]
		if !exists {
			return fmt.Errorf("no provider builder registered for type: %s", providerType)
		}

		provider, err := builder.BuildProvider(accounts)
		if err != nil {
			return fmt.Errorf("failed to build provider %s: %w", providerType, err)
		}

		sender.RegisterProvider(providerType, provider, nil)
	}

	return nil
}

// BuildMessage builds a message for the specified provider type
func (r *ProviderRegistry) BuildMessage(providerType core.ProviderType, flags *cli.CLIFlags) (core.Message, error) {
	builder, exists := r.messageBuilders[providerType]
	if !exists {
		return nil, fmt.Errorf("no message builder registered for provider type: %s", providerType)
	}

	if err := builder.ValidateFlags(flags); err != nil {
		return nil, fmt.Errorf("flag validation failed for %s: %w", providerType, err)
	}

	return builder.BuildMessage(flags)
}

// GetSupportedProviders returns all supported provider types
func (r *ProviderRegistry) GetSupportedProviders() []core.ProviderType {
	var providers []core.ProviderType
	for providerType := range r.providerBuilders {
		providers = append(providers, providerType)
	}
	return providers
}

// ValidateFlags validates flags for the specified provider
func (r *ProviderRegistry) ValidateFlags(providerType core.ProviderType, flags *cli.CLIFlags) error {
	builder, exists := r.messageBuilders[providerType]
	if !exists {
		return fmt.Errorf("no message builder registered for provider type: %s", providerType)
	}

	return builder.ValidateFlags(flags)
}

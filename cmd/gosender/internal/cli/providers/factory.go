package providers

import (
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// DefaultProviderRegistry creates and configures the default provider registry
func DefaultProviderRegistry() *ProviderRegistry {
	registry := NewProviderRegistry()

	// Register email provider and message builders
	registry.RegisterProviderBuilder(&EmailProviderBuilder{})
	registry.RegisterMessageBuilder(&EmailMessageBuilder{})

	// Register wecombot provider and message builders
	registry.RegisterProviderBuilder(&WeComBotProviderBuilder{})
	registry.RegisterMessageBuilder(&WeComBotMessageBuilder{})

	// TODO: Add other providers here as they are implemented
	// registry.RegisterProviderBuilder(&SMSProviderBuilder{})
	// registry.RegisterMessageBuilder(&SMSMessageBuilder{})
	// registry.RegisterProviderBuilder(&DingTalkProviderBuilder{})
	// registry.RegisterMessageBuilder(&DingTalkMessageBuilder{})

	return registry
}

// GetProviderType converts string to ProviderType with validation
func GetProviderType(provider string) (core.ProviderType, error) {
	if provider == "" {
		return core.ProviderTypeSMS, nil // Default
	}

	providerType := core.ProviderType(provider)

	// Validate that it's a known provider type
	switch providerType {
	case core.ProviderTypeSMS,
		core.ProviderTypeEmail,
		core.ProviderTypeDingtalk,
		core.ProviderTypeWebhook,
		core.ProviderTypeTelegram,
		core.ProviderTypeLark,
		core.ProviderTypeWecombot,
		core.ProviderTypeServerChan,
		core.ProviderTypeEmailAPI:
		return providerType, nil
	default:
		return "", fmt.Errorf("unsupported provider type: %s", provider)
	}
}

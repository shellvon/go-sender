package emailapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
)

// Provider is the main emailapi provider, supporting multiple API-based email services.
type Provider struct {
	config  Config
	clients map[string]APIClient // key: provider type (e.g., "mailgun")
}

// APIClient defines the interface for all API email service clients.
type APIClient interface {
	Send(ctx context.Context, msg *Message) error
	Name() string
}

// Registry for all supported API email providers.
var providerRegistry = make(map[string]func(AccountConfig) APIClient)

// RegisterProvider registers a new API email provider implementation.
func RegisterProvider(name string, constructor func(AccountConfig) APIClient) {
	providerRegistry[name] = constructor
}

// New creates a new emailapi Provider with the given config.
func New(config Config) (*Provider, error) {
	if len(config.Accounts) == 0 {
		return nil, errors.New("at least one account is required")
	}
	clients := make(map[string]APIClient)
	for _, acc := range config.Accounts {
		ctor, ok := providerRegistry[acc.Name]
		if !ok {
			ctor, ok = providerRegistry[config.ProviderType]
			if !ok {
				return nil, fmt.Errorf("unsupported email API provider: %s", acc.Name)
			}
		}
		clients[acc.Name] = ctor(acc)
	}
	return &Provider{
		config:  config,
		clients: clients,
	}, nil
}

// getProviderFromMessage tries to get provider from message.Extras["provider"]
func (p *Provider) getProviderFromMessage(msg *Message) APIClient {
	if msg.Extras != nil {
		if providerName, ok := msg.Extras["provider"].(string); ok {
			if client, ok := p.clients[providerName]; ok {
				return client
			}
		}
	}
	return nil
}

// Send dispatches an email message using the selected API provider/account.
func (p *Provider) Send(ctx context.Context, msg core.Message) error {
	emailMsg, ok := msg.(*Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *emailapi.Message, got %T", msg))
	}
	if err := emailMsg.Validate(); err != nil {
		return err
	}
	client := p.getProviderFromMessage(emailMsg)
	if client == nil {
		// fallback: use the first enabled client
		for _, c := range p.clients {
			client = c
			break
		}
	}
	if client == nil {
		return errors.New("no enabled email API provider available")
	}
	return client.Send(ctx, emailMsg)
}

// Name returns the provider type.
func (p *Provider) Name() string {
	return p.config.ProviderType
}

package sms

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider supports multiple SMS service providers and strategy selection
type Provider struct {
	providers []*SMSProvider
	selector  *utils.Selector[*SMSProvider]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new SMS provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("SMS provider is not configured or is disabled")
	}

	// Convert to pointer slice
	providers := make([]*SMSProvider, len(config.Providers))
	for i := range config.Providers {
		providers[i] = &config.Providers[i]
	}

	// Use common initialization logic
	enabledProviders, selector, err := utils.InitProvider(&config, providers)
	if err != nil {
		return nil, errors.New("no enabled SMS providers found")
	}

	return &Provider{
		providers: enabledProviders,
		selector:  selector,
	}, nil
}

// Send sends an SMS message
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	smsMsg, ok := message.(*Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *sms.Message, got %T", message))
	}

	if err := smsMsg.Validate(); err != nil {
		return err
	}

	provider := p.selectProvider(ctx)
	if provider == nil {
		return errors.New("no available provider")
	}
	return p.doSendSMS(ctx, provider, smsMsg)
}

// selectProvider selects a provider based on context
func (p *Provider) selectProvider(ctx context.Context) *SMSProvider {
	return p.selector.Select(ctx)
}

// doSendSMS dispatches to the correct provider implementation
func (p *Provider) doSendSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	switch provider.Type {
	case ProviderTypeCl253:
		return sendCl253SMS(ctx, provider, msg)
	case ProviderTypeLuosimao:
		return sendLuosimaoSMS(ctx, provider, msg)
	case ProviderTypeSmsbao:
		return sendSmsbaoSMS(ctx, provider, msg)
	case ProviderTypeJuhe:
		return sendJuheSMS(ctx, provider, msg)
	case ProviderTypeHuawei:
		return sendHuaweiSMS(ctx, provider, msg)
	case ProviderTypeAliyun:
		return sendAliyunSMS(ctx, provider, msg)
	case ProviderTypeUcp:
		return sendUcpSMS(ctx, provider, msg)
	case ProviderTypeYunpian:
		return sendYunpianSMS(ctx, provider, msg)
	case ProviderTypeSubmail:
		return sendSubmailSMS(ctx, provider, msg)
	case ProviderTypeVolc:
		return sendVolcSMS(ctx, provider, msg)
	default:
		return fmt.Errorf("unsupported or unimplemented SMS provider type: %s", provider.Type)
	}
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeSMS)
}

package sms

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// SMSProviderInterface defines the interface for SMS providers
type SMSProviderInterface interface {
	// Send sends an SMS message
	Send(ctx context.Context, msg *Message) error

	// GetCapabilities returns the provider's capabilities
	GetCapabilities() *Capabilities

	// CheckCapability checks if a specific capability is supported
	CheckCapability(msg *Message) error

	// GetLimits returns the provider's limits for a specific message type
	GetLimits(msgType MessageType) Limits

	// GetName returns the provider name
	GetName() string

	// GetType returns the provider type
	GetType() string

	// IsEnabled returns if the provider is enabled
	IsEnabled() bool

	// GetWeight returns the provider weight for selection
	GetWeight() int

	// CheckConfigured checks if the provider is configured correctly
	CheckConfigured() error
}

// Provider supports multiple SMS service providers and strategy selection
type Provider struct {
	providers []SMSProviderInterface
	selector  *utils.Selector[SMSProviderInterface]
}

var _ core.Provider = (*Provider)(nil)

// ProviderFactories 动态发现的 Provider 工厂，供自动化工具使用
var ProviderFactories = make(map[ProviderType]func(SMSProvider) SMSProviderInterface)

// RegisterProviderConstructor 注册provider构造器（用于动态加载）
func RegisterProviderConstructor(typ ProviderType, constructor interface{}) {
	ProviderFactories[typ] = func(cfg SMSProvider) SMSProviderInterface {
		args := []reflect.Value{reflect.ValueOf(cfg)}
		results := reflect.ValueOf(constructor).Call(args)
		return results[0].Interface().(SMSProviderInterface)
	}
}

// New creates a new SMS provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("SMS provider is not configured or is disabled")
	}

	// Convert config providers to SMSProviderInterface
	providers := make([]SMSProviderInterface, 0, len(config.Providers))
	for _, cfg := range config.Providers {
		provider, err := createProvider(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", cfg.Name, err)
		}
		if err := provider.CheckConfigured(); err != nil {
			return nil, fmt.Errorf("provider conf failed %s: %w", cfg.Name, err)
		}
		if provider.IsEnabled() {
			providers = append(providers, provider)
		}
	}

	if len(providers) == 0 {
		return nil, errors.New("no enabled SMS providers found")
	}

	// Create selector
	selector := utils.NewSelector(providers, utils.GetStrategy(config.Strategy))

	return &Provider{
		providers: providers,
		selector:  selector,
	}, nil
}

// createProvider creates a specific SMS provider based on configuration
func createProvider(cfg SMSProvider) (SMSProviderInterface, error) {
	if factory, ok := ProviderFactories[cfg.Type]; ok {
		return factory(cfg), nil
	}

	// 如果动态发现失败，回退到硬编码的switch
	switch cfg.Type {
	case ProviderTypeAliyun:
		return NewAliyunProvider(cfg), nil
	case ProviderTypeCl253:
		return NewCl253Provider(cfg), nil
	case ProviderTypeYuntongxun:
		return NewYuntongxunProvider(cfg), nil
	case ProviderTypeJuhe:
		return NewJuheProvider(cfg), nil
	case ProviderTypeHuawei:
		return NewHuaweiProvider(cfg), nil
	case ProviderTypeVolc:
		return NewVolcProvider(cfg), nil
	case ProviderTypeSmsbao:
		return NewSmsbaoProvider(cfg), nil
	case ProviderTypeSubmail:
		return NewSubmailProvider(cfg), nil
	case ProviderTypeUcp:
		return NewUcpProvider(cfg), nil
	case ProviderTypeLuosimao:
		return NewLuosimaoProvider(cfg), nil
	case ProviderTypeYunpian:
		return NewYunpianProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported SMS provider type: %s", cfg.Type)
	}
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

	// Try to get provider from message context first
	provider := p.getProviderFromMessage(ctx, smsMsg)
	if provider == nil {
		// Fallback to selector
		provider = p.selector.Select(ctx)
	}

	if provider == nil {
		return errors.New("no available provider")
	}

	return provider.Send(ctx, smsMsg)
}

// getProviderFromMessage tries to get provider from message context or message itself
func (p *Provider) getProviderFromMessage(ctx context.Context, msg *Message) SMSProviderInterface {
	// Check if message has specific provider preference
	if msg.Extras != nil {
		if providerName, ok := msg.Extras["provider"].(string); ok {
			for _, provider := range p.providers {
				if provider.GetName() == providerName {
					return provider
				}
			}
		}
	}

	// Check context for provider preference
	if ctx != nil {
		if providerName, ok := ctx.Value("sms_provider").(string); ok {
			for _, provider := range p.providers {
				if provider.GetName() == providerName {
					return provider
				}
			}
		}
	}

	return nil
}

// GetCapabilities returns capabilities of the selected provider
func (p *Provider) GetCapabilities(ctx context.Context) (*Capabilities, error) {
	provider := p.selector.Select(ctx)
	if provider == nil {
		return nil, errors.New("no available provider")
	}

	return provider.GetCapabilities(), nil
}

// GetProviderCapabilities returns capabilities of a specific provider
func (p *Provider) GetProviderCapabilities(providerName string) (*Capabilities, error) {
	for _, provider := range p.providers {
		if provider.GetName() == providerName {
			return provider.GetCapabilities(), nil
		}
	}
	return nil, fmt.Errorf("provider not found: %s", providerName)
}

// CheckCapability checks if a specific capability is supported by the selected provider
func (p *Provider) CheckCapability(ctx context.Context, msg *Message) (bool, error) {
	provider := p.selector.Select(ctx)
	if provider == nil {
		return false, errors.New("no available provider")
	}

	return provider.CheckCapability(msg) == nil, nil
}

// GetAllProviders returns all available providers
func (p *Provider) GetAllProviders() []SMSProviderInterface {
	return p.providers
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeSMS)
}

// DefaultCheckCapability 统一校验 provider 能力
// 1. 判断国内/国际能力
// 2. 判断类型/类别能力
// 3. 判断手机号数量是否超限
func DefaultCheckCapability(p SMSProviderInterface, msg *Message) error {
	caps := p.GetCapabilities()
	var regionCap *RegionCapability
	var msgCap *MessageCapability
	if msg.IsIntl() {
		regionCap = &caps.SMS.International
		msgCap = &caps.SMS
	} else {
		regionCap = &caps.SMS.Domestic
		msgCap = &caps.SMS
	}
	// 1. 判断是否支持该区域
	if regionCap == nil || (!regionCap.Single && !regionCap.Batch) {
		return fmt.Errorf("provider %s does not support %s region", p.GetName(), map[bool]string{true: "international", false: "domestic"}[msg.IsIntl()])
	}
	// 2. 判断类型/类别能力
	if !caps.Supports(msg.Type, msg.Category) {
		return fmt.Errorf("provider %s does not support type %s, category %s", p.GetName(), msg.Type, msg.Category)
	}
	// 3. 判断手机号数量
	if len(msg.Mobiles) == 0 {
		return fmt.Errorf("no mobile numbers provided")
	}
	if len(msg.Mobiles) == 1 && !regionCap.Single {
		return fmt.Errorf("provider %s does not support single send for this region", p.GetName())
	}
	if len(msg.Mobiles) > 1 {
		if !regionCap.Batch {
			return fmt.Errorf("provider %s does not support batch send for this region", p.GetName())
		}
		limits := msgCap.Limits
		if limits.MaxBatchSize > 0 && len(msg.Mobiles) > limits.MaxBatchSize {
			return fmt.Errorf("provider %s batch send exceeds max limit: %d", p.GetName(), limits.MaxBatchSize)
		}
	}
	return nil
}

// ValidateForSend 统一校验 provider 配置、消息、能力
// 1. 配置校验（CheckConfigured）
// 2. 消息校验（msg.Validate）
// 3. 能力校验（CheckCapability）
func ValidateForSend(p SMSProviderInterface, msg *Message) error {
	if err := p.CheckConfigured(); err != nil {
		return err
	}
	if err := msg.Validate(); err != nil {
		return err
	}
	if err := p.CheckCapability(msg); err != nil {
		return err
	}
	return nil
}

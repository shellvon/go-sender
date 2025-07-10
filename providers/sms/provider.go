package sms

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// smsRegistry is a package-level registry leveraging the shared providers.TransformerRegistry
// to avoid duplicated code.
//
//nolint:gochecknoglobals // Global registry is acceptable for package-wide look-ups.
var smsRegistry = providers.NewTransformerRegistry[*Account]()

// RegisterTransformer registers a transformer for a given SMS sub-provider.
func RegisterTransformer(subProvider string, transformer core.HTTPTransformer[*Account]) {
	smsRegistry.Register(subProvider, transformer)
}

// GetTransformer retrieves a transformer for a given SMS sub-provider.
func GetTransformer(subProvider string) (core.HTTPTransformer[*Account], bool) {
	return smsRegistry.Get(subProvider)
}

// smsTransformer 实现 core.HTTPTransformer[*Account]，根据SubProvider选择具体的transformer.
type smsTransformer struct{}

// CanTransform 判断是否为 SMS 消息.
func (t *smsTransformer) CanTransform(msg core.Message) bool {
	return msg.ProviderType() == core.ProviderTypeSMS
}

// Transform 根据SubProvider从注册表获取具体的transformer进行转换.
func (t *smsTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for sms transformer: %T", msg)
	}

	if smsMsg.SubProvider == "" {
		return nil, nil, errors.New("sub-provider is required for sms transformer")
	}

	// 从注册表获取transformer
	transformer, exists := GetTransformer(smsMsg.SubProvider)
	if !exists {
		return nil, nil, fmt.Errorf("unsupported SMS sub-provider: %s", smsMsg.SubProvider)
	}

	return transformer.Transform(ctx, msg, account)
}

// newSMSTransformer returns a new instance of the package-local smsTransformer.
func newSMSTransformer() core.HTTPTransformer[*Account] {
	return &smsTransformer{}
}

// New creates a new SMS provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeSMS),
		newSMSTransformer(),
		config,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{HTTPProvider: httpProvider}, nil
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeSMS)
}

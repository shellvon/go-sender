package sms

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type Config = core.BaseConfig[*Account]

type Provider struct {
	*providers.HTTPProvider[*Account]
}

var _ core.Provider = (*Provider)(nil)

// transformerRegistry global transformer registry.
//
//nolint:gochecknoglobals // Reason: transformerRegistry is a global registry for sms transformers
var transformerRegistry = make(map[string]core.HTTPTransformer[*Account])

// registryMutex global mutex for transformerRegistry.
//
//nolint:gochecknoglobals // Reason: registryMutex is a global mutex for transformerRegistry
var registryMutex sync.RWMutex

// RegisterTransformer 注册transformer到全局注册表.
func RegisterTransformer(subProvider string, transformer core.HTTPTransformer[*Account]) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	transformerRegistry[subProvider] = transformer
}

// GetTransformer 从注册表获取transformer.
func GetTransformer(subProvider string) (core.HTTPTransformer[*Account], bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	transformer, exists := transformerRegistry[subProvider]
	return transformer, exists
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
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
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

// New creates a new SMS provider instance.
func New(config *Config) (*Provider, error) {
	httpProvider, err := providers.NewHTTPProvider(
		string(core.ProviderTypeSMS),
		&smsTransformer{},
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

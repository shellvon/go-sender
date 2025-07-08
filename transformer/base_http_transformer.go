package transformer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
)

type HandlerFunc[M core.Message, C core.Selectable] func(ctx context.Context, msg M, cfg C) (*core.HTTPRequestSpec, core.ResponseHandler, error)

type BeforeHook[M core.Message, C core.Selectable] func(ctx context.Context, msg M, cfg C) error

type AfterHook[M core.Message, C core.Selectable] func(ctx context.Context, msg M, cfg C, handlerErr error) error

// ConfigResolver allows dynamically adjusting the config (e.g. account) before executing handler.
// It receives the ctx + message + current cfg, and returns the cfg that should be used.
type ConfigResolver[M core.Message, C core.Selectable] func(ctx context.Context, msg M, current C) (C, error)

type Option[M core.Message, C core.Selectable] func(*BaseHTTPTransformer[M, C])

type BaseHTTPTransformer[M core.Message, C core.Selectable] struct {
	providerType core.ProviderType // sms / webhook / emailapi …
	subProvider  string            // aliyun / resend / ""(default)

	responseCfg *core.ResponseHandlerConfig

	handler       HandlerFunc[M, C]
	beforeHooks   []BeforeHook[M, C]
	afterHooks    []AfterHook[M, C]
	resolveConfig ConfigResolver[M, C]
}

func NewBaseHTTPTransformer[M core.Message, C core.Selectable](
	providerType core.ProviderType,
	subProvider string,
	cfg *core.ResponseHandlerConfig,
	opts ...Option[M, C],
) *BaseHTTPTransformer[M, C] {
	if cfg == nil {
		cfg = &core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: false,
		}
	}

	defaultSize := 4
	bt := &BaseHTTPTransformer[M, C]{
		providerType: providerType,
		subProvider:  subProvider,
		responseCfg:  cfg,
		beforeHooks:  make([]BeforeHook[M, C], 0, defaultSize),
		afterHooks:   make([]AfterHook[M, C], 0, defaultSize),
	}

	for _, o := range opts {
		o(bt)
	}

	return bt
}

func (t *BaseHTTPTransformer[M, C]) CanTransform(msg core.Message) bool {
	m, ok := msg.(M)
	if !ok {
		return false
	}
	if m.ProviderType() != t.providerType {
		return false
	}
	if t.subProvider != "" && m.GetSubProvider() != t.subProvider {
		return false
	}
	return true
}

func (t *BaseHTTPTransformer[M, C]) Transform(
	ctx context.Context,
	msg core.Message,
	cfg C,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	m, ok := msg.(M)
	if !ok {
		return nil, nil, fmt.Errorf("invalid message type for %s provider: %T", t.providerType, msg)
	}

	if t.resolveConfig != nil {
		var err error
		cfg, err = t.resolveConfig(ctx, m, cfg)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, bh := range t.beforeHooks {
		if err := bh(ctx, m, cfg); err != nil {
			return nil, nil, err
		}
	}

	if t.handler == nil {
		return nil, nil, fmt.Errorf("%s provider has no handler registered", t.fullName())
	}

	reqSpec, respHandler, err := t.handler(ctx, m, cfg)
	if respHandler == nil {
		respHandler = t.buildResponseHandler()
	}

	for _, ah := range t.afterHooks {
		if aftErr := ah(ctx, m, cfg, err); aftErr != nil {
			err = aftErr
		}
	}

	return reqSpec, respHandler, err
}

func (t *BaseHTTPTransformer[M, C]) buildResponseHandler() core.ResponseHandler {
	generic := core.NewResponseHandler(t.responseCfg)
	fullName := t.fullName()

	return func(resp *http.Response) error {
		if err := generic(resp); err != nil {
			// 仅增加前缀信息，避免依赖 sms.NewProviderError
			return fmt.Errorf("[%s] %w", fullName, err)
		}
		return nil
	}
}

// fullName 返回 providerType[.subProvider] 形式，用于日志 / 错误标识。
func (t *BaseHTTPTransformer[M, C]) fullName() string {
	if t.subProvider == "" {
		return string(t.providerType)
	}
	return fmt.Sprintf("%s.%s", t.providerType, t.subProvider)
}

func AddBeforeHook[M core.Message, C core.Selectable](h BeforeHook[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.beforeHooks = append(b.beforeHooks, h)
	}
}

func SetBeforeHooks[M core.Message, C core.Selectable](hooks ...BeforeHook[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.beforeHooks = hooks
	}
}

func AddAfterHook[M core.Message, C core.Selectable](h AfterHook[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.afterHooks = append(b.afterHooks, h)
	}
}

func SetAfterHooks[M core.Message, C core.Selectable](hooks ...AfterHook[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.afterHooks = hooks
	}
}

// WithHandler 注册唯一的 Handler。 (保持原名以避免修改过多调用).
func WithHandler[M core.Message, C core.Selectable](h HandlerFunc[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.handler = h
	}
}

// WithConfigResolver 设置动态配置解析器。
func WithConfigResolver[M core.Message, C core.Selectable](r ConfigResolver[M, C]) Option[M, C] {
	return func(b *BaseHTTPTransformer[M, C]) {
		b.resolveConfig = r
	}
}

// NewSimpleHTTPTransformer creates a new BaseHTTPTransformer with default options.
// It validates the message and sets the handler.
//
// The returned transformer is ready to be used with providers.NewHTTPProvider.
func NewSimpleHTTPTransformer[M core.Message, C core.Selectable](
	providerType core.ProviderType,
	subProvider string,
	respCfg *core.ResponseHandlerConfig,
	handler HandlerFunc[M, C],
	extraOpts ...Option[M, C],
) *BaseHTTPTransformer[M, C] {
	opts := []Option[M, C]{
		AddBeforeHook(func(_ context.Context, msg M, _ C) error {
			return msg.Validate()
		}),
		WithHandler(handler),
	}
	opts = append(opts, extraOpts...)

	return NewBaseHTTPTransformer(providerType, subProvider, respCfg, opts...)
}

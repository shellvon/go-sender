package sms

import (
	"context"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// BaseTransformer is the base transformer for SMS providers.
// It supports sending text message, voice message, and mms message.
// It also supports adding before and after hooks to the transformer.

type (
	HandlerFunc func(ctx context.Context, msg *Message, account *Account) (*core.HTTPRequestSpec, core.ResponseHandler, error)
	Option      func(*BaseTransformer)
	HTTPOption  = transformer.Option[*Message, *Account]
	HTTPOptions = []HTTPOption
)

func WithSMSHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[SMSText] = fn }
}
func WithVoiceHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[Voice] = fn }
}
func WithMMSHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[MMS] = fn }
}
func WithHandler(mt MessageType, fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[mt] = fn }
}

func AddBeforeHook(h transformer.BeforeHook[*Message, *Account]) HTTPOption {
	return transformer.AddBeforeHook(h)
}

func SetBeforeHooks(hooks ...transformer.BeforeHook[*Message, *Account]) HTTPOption {
	return transformer.SetBeforeHooks(hooks...)
}

func AddAfterHook(h transformer.AfterHook[*Message, *Account]) HTTPOption {
	return transformer.AddAfterHook(h)
}

func SetAfterHooks(hooks ...transformer.AfterHook[*Message, *Account]) HTTPOption {
	return transformer.SetAfterHooks(hooks...)
}

func WithConfigResolver(r transformer.ConfigResolver[*Message, *Account]) HTTPOption {
	return transformer.WithConfigResolver(r)
}

type BaseTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]

	handlers map[MessageType]HandlerFunc
}

// NewBaseTransformer constructs a BaseTransformer.
func NewBaseTransformer(
	subProvider string,
	cfg *core.ResponseHandlerConfig,
	httpOpts HTTPOptions,
	opts ...Option,
) *BaseTransformer {
	bt := &BaseTransformer{
		handlers: make(map[MessageType]HandlerFunc),
	}

	for _, o := range opts {
		o(bt)
	}

	baseOpts := []transformer.Option[*Message, *Account]{
		transformer.AddBeforeHook(func(_ context.Context, m *Message, acc *Account) error {
			m.ApplyCommonDefaults(acc)
			return m.Validate()
		}),
		transformer.WithHandler(bt.dispatchHandler),
	}

	baseOpts = append(baseOpts, httpOpts...)

	bt.BaseHTTPTransformer = transformer.NewBaseHTTPTransformer(
		core.ProviderTypeSMS,
		subProvider,
		cfg,
		baseOpts...,
	)

	return bt
}

func (t *BaseTransformer) dispatchHandler(
	ctx context.Context,
	msg *Message,
	acc *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	handler, ok := t.handlers[msg.Type]
	if !ok {
		return nil, nil, fmt.Errorf("sms.%s unsupported message type %v", msg.GetSubProvider(), msg.Type)
	}

	return handler(ctx, msg, acc)
}

package sms

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shellvon/go-sender/core"
)

// HandlerFunc is the common signature for an SMS request builder.
type HandlerFunc func(ctx context.Context, msg *Message, account *Account) (*core.HTTPRequestSpec, core.ResponseHandler, error)

// Option configures a BaseTransformer at construction time.
type Option func(*BaseTransformer)

// WithSMSHandler registers a plain-text SMS handler.
func WithSMSHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[SMSText] = fn }
}

// WithVoiceHandler registers a voice SMS handler.
func WithVoiceHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[Voice] = fn }
}

// WithMMSHandler registers an MMS handler.
func WithMMSHandler(fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[MMS] = fn }
}

// WithHandler registers a handler for an arbitrary MessageType.
func WithHandler(mt MessageType, fn HandlerFunc) Option {
	return func(b *BaseTransformer) { b.handlers[mt] = fn }
}

// ---------------- Before / After hook types ----------------

// BeforeHook is executed before the concrete handler; returning error aborts processing.
type BeforeHook func(ctx context.Context, msg *Message, account *Account) error

// AfterHook is executed after the concrete handler; it receives the original error (if any).
// If it returns a non-nil error, that error is propagated instead.
type AfterHook func(ctx context.Context, msg *Message, account *Account, handlerErr error) error

// WithBeforeHook appends a before-hook.
func WithBeforeHook(h BeforeHook) Option {
	return func(b *BaseTransformer) { b.befores = append(b.befores, h) }
}

// WithAfterHook appends an after-hook.
func WithAfterHook(h AfterHook) Option {
	return func(b *BaseTransformer) { b.afters = append(b.afters, h) }
}

// BaseTransformer centralises common functionality for all SMS sub-providers:
//  1. CanTransform – filter by sub-provider.
//  2. Unified HTTP response handling via core.ResponseHandlerConfig.
//  3. Generic validation (mobiles, etc.).
//
// A concrete provider only needs to:
//   - Embed *BaseTransformer.
//   - Register one or more handlers via WithSMSHandler / WithVoiceHandler / WithMMSHandler
//     (or call Register manually).
//   - If the concrete handler returns a nil ResponseHandler, BaseTransformer will
//     automatically substitute buildResponseHandler(), created from the
//     ResponseHandlerConfig passed to NewBaseTransformer.
//
// This eliminates repetitive boilerplate across individual SMS providers.
type BaseTransformer struct {
	providerName string
	subProvider  string

	// Response validation configuration.
	responseCfg *core.ResponseHandlerConfig

	// handlers maps MessageType → concrete handler.
	handlers map[MessageType]HandlerFunc

	// hook chains
	befores []BeforeHook
	afters  []AfterHook
}

// NewBaseTransformer constructs a BaseTransformer.
// providerName is typically "sms"; subProvider identifies the concrete platform ("cl253" etc.).
// If cfg is nil the response handler defaults to JSON with status-code only validation.
func NewBaseTransformer(
	providerName, subProvider string,
	cfg *core.ResponseHandlerConfig,
	opts ...Option,
) *BaseTransformer {
	if cfg == nil {
		cfg = &core.ResponseHandlerConfig{
			ResponseType:     core.BodyTypeJSON,
			ValidateResponse: false,
		}
	}

	defaultSize := 4
	bt := &BaseTransformer{
		providerName: providerName,
		subProvider:  subProvider,
		responseCfg:  cfg,
		handlers:     make(map[MessageType]HandlerFunc),
		befores:      make([]BeforeHook, 0, defaultSize),
		afters:       make([]AfterHook, 0, defaultSize),
	}

	// Built-in before-hook: core defaults + base validation.
	bt.befores = append(bt.befores, func(_ context.Context, m *Message, acc *Account) error {
		// default apply
		m.ApplyCommonDefaults(acc)
		// base validation
		if len(m.Mobiles) == 0 {
			return NewProviderError(providerName, "MISSING_MOBILE", "at least one mobile number is required")
		}
		return nil
	})

	// Apply functional options to register provider-specific handlers and hooks.
	for _, o := range opts {
		o(bt)
	}

	return bt
}

// Register adds/overrides a handler at runtime. It is rarely needed outside tests.
func (t *BaseTransformer) Register(mt MessageType, fn HandlerFunc) {
	if t.handlers == nil {
		t.handlers = make(map[MessageType]HandlerFunc)
	}
	t.handlers[mt] = fn
}

// CanTransform implements core.HTTPTransformer.CanTransform – it matches the configured sub-provider.
func (t *BaseTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == t.subProvider
}

// Transform is the main entry of core.HTTPTransformer. It delegates to the
// registered handler based on MessageType; returns an error if unsupported.
func (t *BaseTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("invalid message type for %s provider: %T", t.providerName, msg)
	}

	// 1) Run before-hooks (built-in + provider-specific).
	for _, bh := range t.befores {
		if err := bh(ctx, smsMsg, account); err != nil {
			return nil, nil, err
		}
	}

	// 2) Lookup concrete handler for the message type and execute.
	if handler, exists := t.handlers[smsMsg.Type]; exists {
		reqSpec, respHandler, err := handler(ctx, smsMsg, account)

		// If handler didn't provide a ResponseHandler, fall back to the generic one.
		if respHandler == nil {
			respHandler = t.buildResponseHandler()
		}

		// 3) Run after-hooks (if any). First hook receiving error may alter it.
		for _, ah := range t.afters {
			if aftErr := ah(ctx, smsMsg, account, err); aftErr != nil {
				err = aftErr
			}
		}
		return reqSpec, respHandler, err
	}
	return nil, nil, fmt.Errorf("%s provider does not support message type %v", t.subProvider, smsMsg.Type)
}

// buildResponseHandler builds a provider-specific ResponseHandler and wraps
// generic errors into ProviderError for unified error reporting.
func (t *BaseTransformer) buildResponseHandler() core.ResponseHandler {
	generic := core.NewResponseHandler(t.responseCfg)
	providerName := t.providerName
	return func(resp *http.Response) error {
		if err := generic(resp); err != nil {
			return NewProviderError(providerName, "API_ERROR", err.Error())
		}
		return nil
	}
}

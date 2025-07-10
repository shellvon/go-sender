package providers

import (
	"context"
	"fmt"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// HTTPProvider is a unified HTTP Provider base class using generic design
// T must implement the core.Selectable interface, typically *core.Account.
type HTTPProvider[T core.Selectable] struct {
	name        string
	config      *core.BaseConfig[T]
	transformer core.HTTPTransformer[T]
}

// NewHTTPProvider creates a new HTTP Provider from a config object.
// The config must implement Validate, GetItems, and GetStrategy.
func NewHTTPProvider[T core.Selectable](
	name string,
	transformer core.HTTPTransformer[T],
	config *core.BaseConfig[T],
) (*HTTPProvider[T], error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil for provider %s", name)
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &HTTPProvider[T]{
		name:        name,
		config:      config,
		transformer: transformer,
	}, nil
}

func (p *HTTPProvider[T]) Select(ctx context.Context, filter func(T) bool) (T, error) {
	return p.config.Select(ctx, filter)
}

// Send implements core.Provider interface and returns detailed SendResult.
func (p *HTTPProvider[T]) Send(
	ctx context.Context,
	msg core.Message,
	opts *core.ProviderSendOptions,
) (*core.SendResult, error) {
	if opts == nil {
		opts = &core.ProviderSendOptions{}
	}

	// Filter selection by sub-provider (same as Send).
	filter := func(item T) bool {
		sub := ""
		if subProviderMsg, ok := msg.(interface{ GetSubProvider() string }); ok {
			sub = subProviderMsg.GetSubProvider()
		}
		return sub == "" || item.GetType() == sub
	}

	selectedConfig, err := p.config.Select(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Transform request
	reqSpec, handler, err := p.transformer.Transform(ctx, msg, selectedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to transform message: %w", err)
	}

	// Execute HTTP request capturing detailed response
	result, err := p.executeHTTPRequest(ctx, reqSpec, handler, opts)
	if result != nil {
		result.Config = selectedConfig // attach config for observability
	}
	return result, err
}

func (p *HTTPProvider[T]) executeHTTPRequest(
	ctx context.Context,
	reqSpec *core.HTTPRequestSpec,
	handler core.SendResultHandler,
	opts *core.ProviderSendOptions,
) (*core.SendResult, error) {
	// Prepare headers map (ensure non-nil so we can mutate)
	headers := make(map[string]string)
	for k, v := range reqSpec.Headers {
		headers[k] = v
	}

	httpOpts := utils.HTTPRequestOptions{
		Method:  reqSpec.Method,
		Headers: headers,
		Client:  opts.HTTPClient,
		Timeout: reqSpec.Timeout,
		Query:   reqSpec.QueryParams,
	}
	httpOpts.Raw = reqSpec.Body
	if ct := reqSpec.BodyType.ContentType(); ct != "" {
		httpOpts.Headers["Content-Type"] = ct
	}

	resp, err := utils.SendRequest(ctx, reqSpec.URL, httpOpts)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Read body bytes and reset resp.Body so handler can consume without network I/O again.
	bodyBytes, respHeaders, readErr := utils.ReadAndClose(resp)
	if readErr != nil {
		return nil, readErr
	}

	if handler == nil {
		handler = core.NewSendResultHandler(&core.ResponseHandlerConfig{CheckBody: false})
	}

	// Build SendResult early so we can return it even on handler error.
	result := &core.SendResult{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       bodyBytes,
	}

	if errHandler := handler(result); errHandler != nil {
		return result, errHandler
	}

	return result, nil
}

// Name returns the provider name.
func (p *HTTPProvider[T]) Name() string {
	return p.name
}

// GetConfigs returns all configurations.
func (p *HTTPProvider[T]) GetConfigs() []T {
	return p.config.GetItems()
}

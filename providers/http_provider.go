package providers

import (
	"context"
	"fmt"
	"net/url"

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

// Send implements the core.Provider interface.
func (p *HTTPProvider[T]) Send(ctx context.Context, msg core.Message, opts *core.ProviderSendOptions) error {
	if opts == nil {
		opts = &core.ProviderSendOptions{}
	}
	// 过滤逻辑：如果 msg 有 GetSubProvider，则只选 type 匹配的
	filter := func(item T) bool {
		sub := ""
		if subProviderMsg, ok := msg.(interface{ GetSubProvider() string }); ok {
			sub = subProviderMsg.GetSubProvider()
		}
		return sub == "" || item.GetType() == sub
	}
	selectedConfig, err := p.config.Select(ctx, filter)
	if err != nil {
		return err
	}

	// Transform request
	reqSpec, handler, err := p.transformer.Transform(ctx, msg, selectedConfig)
	if err != nil {
		return fmt.Errorf("failed to transform message: %w", err)
	}
	// Execute HTTP request
	return p.executeHTTPRequest(ctx, reqSpec, handler, opts)
}

// Name returns the provider name.
func (p *HTTPProvider[T]) Name() string {
	return p.name
}

// GetConfigs returns all configurations.
func (p *HTTPProvider[T]) GetConfigs() []T {
	return p.config.GetItems()
}

// executeHTTPRequest executes HTTP request.
func (p *HTTPProvider[T]) executeHTTPRequest(
	ctx context.Context,
	reqSpec *core.HTTPRequestSpec,
	handler core.ResponseHandler,
	opts *core.ProviderSendOptions,
) error {
	// Build URL (including query parameters)
	requestURL := reqSpec.URL
	if len(reqSpec.QueryParams) > 0 {
		parsedURL, err := url.Parse(reqSpec.URL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		query := parsedURL.Query()
		for k, v := range reqSpec.QueryParams {
			query.Set(k, v)
		}
		parsedURL.RawQuery = query.Encode()
		requestURL = parsedURL.String()
	}

	// Prepare HTTP request options
	httpOpts := utils.HTTPRequestOptions{
		Method:  reqSpec.Method,
		Headers: reqSpec.Headers,
		Client:  opts.HTTPClient,
		Timeout: reqSpec.Timeout,
	}

	// Set request body based on body type
	switch reqSpec.BodyType {
	case core.BodyTypeJSON:
		httpOpts.JSON = reqSpec.Body
	case core.BodyTypeForm:
		httpOpts.Raw = reqSpec.Body
		httpOpts.Headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	case core.BodyTypeRaw:
		httpOpts.Raw = reqSpec.Body
	case core.BodyTypeNone:
		httpOpts.Raw = nil
	case core.BodyTypeText:
		httpOpts.Raw = reqSpec.Body
		httpOpts.Headers["Content-Type"] = "text/plain; charset=utf-8"
	case core.BodyTypeXML:
		httpOpts.Raw = reqSpec.Body
		httpOpts.Headers["Content-Type"] = "application/xml; charset=utf-8"
	default:
		// Auto-detect, default to JSON
		httpOpts.JSON = reqSpec.Body
	}

	// Execute request
	body, statusCode, err := utils.DoRequest(ctx, requestURL, httpOpts)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// Use custom handler to process response
	if handler != nil {
		return handler(statusCode, body)
	}

	// Default response handling
	return p.defaultResponseHandler(statusCode, body)
}

// defaultResponseHandler is the default response handler.
func (p *HTTPProvider[T]) defaultResponseHandler(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status code %d. Response body: %s", statusCode, string(body))
	}
	return nil
}

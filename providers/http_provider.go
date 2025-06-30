package providers

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// HTTPProvider is a unified HTTP Provider base class using generic design
// T must implement the core.Selectable interface, typically *core.Account
type HTTPProvider[T core.Selectable] struct {
	name        string
	configs     []T
	strategy    core.SelectionStrategy
	transformer core.HTTPTransformer[T]
}

// NewHTTPProvider creates a new HTTP Provider
func NewHTTPProvider[T core.Selectable](name string, configs []T, transformer core.HTTPTransformer[T], strategy core.SelectionStrategy) *HTTPProvider[T] {
	return &HTTPProvider[T]{
		name:        name,
		configs:     configs,
		strategy:    strategy,
		transformer: transformer,
	}
}

// Send implements the core.Provider interface
func (p *HTTPProvider[T]) Send(ctx context.Context, msg core.Message, opts *core.ProviderSendOptions) error {
	if opts == nil {
		opts = &core.ProviderSendOptions{}
	}

	// Select configuration
	var selectedConfig T
	if len(p.configs) == 1 {
		selectedConfig = p.configs[0]
	} else if len(p.configs) > 1 {
		// Filter configurations based on message's SubProvider
		availableConfigs := p.filterConfigsByMessage(msg)
		if len(availableConfigs) == 0 {
			return fmt.Errorf("no suitable account found for the specified provider type")
		}

		// Convert to Selectable interface
		selectables := make([]core.Selectable, len(availableConfigs))
		for i, config := range availableConfigs {
			selectables[i] = config
		}

		selected := utils.Select(ctx, selectables, p.strategy)
		if selected == nil {
			return fmt.Errorf("no suitable account selected")
		}

		// Find the corresponding configuration
		for _, config := range availableConfigs {
			if config.GetName() == selected.GetName() {
				selectedConfig = config
				break
			}
		}
	} else {
		return errors.New("no available config")
	}

	if !selectedConfig.IsEnabled() {
		return fmt.Errorf("the selected account is disabled")
	}

	// Transform request
	reqSpec, handler, err := p.transformer.Transform(ctx, msg, selectedConfig)
	if err != nil {
		return fmt.Errorf("failed to transform message: %w", err)
	}

	// Execute HTTP request
	return p.executeHTTPRequest(ctx, reqSpec, handler, opts)
}

// filterConfigsByMessage filters configurations based on message
func (p *HTTPProvider[T]) filterConfigsByMessage(msg core.Message) []T {
	// Try to get SubProvider from message
	var subProvider string
	if subProviderMsg, ok := msg.(interface{ GetSubProvider() string }); ok {
		subProvider = subProviderMsg.GetSubProvider()
	}

	if subProvider == "" {
		return p.configs
	}

	// Filter configurations
	filtered := make([]T, 0, len(p.configs))
	for _, config := range p.configs {
		if config.GetType() == subProvider {
			filtered = append(filtered, config)
		}
	}
	return filtered
}

// Name returns the provider name
func (p *HTTPProvider[T]) Name() string {
	return p.name
}

// GetConfigs returns all configurations
func (p *HTTPProvider[T]) GetConfigs() []T {
	return p.configs
}

// SelectConfig selects a configuration (for special methods like UploadMedia)
func (p *HTTPProvider[T]) SelectConfig(ctx context.Context) T {
	if len(p.configs) == 1 {
		return p.configs[0]
	} else if len(p.configs) > 1 {
		// Convert to Selectable interface
		selectables := make([]core.Selectable, len(p.configs))
		for i, config := range p.configs {
			selectables[i] = config
		}

		selected := utils.Select(ctx, selectables, p.strategy)
		if selected == nil {
			var zero T
			return zero
		}

		// Find the corresponding configuration
		for _, config := range p.configs {
			if config.GetName() == selected.GetName() {
				return config
			}
		}
	}
	var zero T
	return zero
}

// executeHTTPRequest executes HTTP request
func (p *HTTPProvider[T]) executeHTTPRequest(ctx context.Context, reqSpec *core.HTTPRequestSpec, handler core.ResponseHandler, opts *core.ProviderSendOptions) error {
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
	case "json":
		httpOpts.JSON = reqSpec.Body
	case "form":
		// Convert []byte to form data, simplified handling
		httpOpts.Raw = reqSpec.Body
	case "raw":
		httpOpts.Raw = reqSpec.Body
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

// defaultResponseHandler is the default response handler
func (p *HTTPProvider[T]) defaultResponseHandler(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}
	return nil
}

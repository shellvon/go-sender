package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider implements the webhook provider
type Provider struct {
	endpoints []*Endpoint
	selector  *utils.Selector[*Endpoint]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new webhook provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("webhook provider is not configured or is disabled")
	}

	// Convert to pointer slice
	endpoints := make([]*Endpoint, len(config.Endpoints))
	for i := range config.Endpoints {
		endpoints[i] = &config.Endpoints[i]
	}

	// Use common initialization logic
	enabledEndpoints, selector, err := utils.InitProvider(&config, endpoints)
	if err != nil {
		return nil, errors.New("no enabled webhook endpoints found")
	}

	return &Provider{
		endpoints: enabledEndpoints,
		selector:  selector,
	}, nil
}

// Send sends a webhook message
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	webhookMsg, ok := message.(*Message)
	if !ok {
		return fmt.Errorf("unsupported message type for webhook provider: %T", message)
	}

	if err := webhookMsg.Validate(); err != nil {
		return err
	}

	selectedEndpoint := p.selector.Select(ctx)
	if selectedEndpoint == nil {
		return errors.New("no available endpoint")
	}

	return p.doSendWebhook(ctx, selectedEndpoint, webhookMsg)
}

// doSendWebhook performs the actual webhook request
func (p *Provider) doSendWebhook(ctx context.Context, endpoint *Endpoint, message *Message) error {
	// Build the request URL with query parameters
	requestURL := endpoint.URL
	if len(endpoint.QueryParams) > 0 {
		parsedURL, err := url.Parse(endpoint.URL)
		if err != nil {
			return fmt.Errorf("invalid endpoint URL: %w", err)
		}

		query := parsedURL.Query()
		for k, v := range endpoint.QueryParams {
			query.Set(k, v)
		}
		parsedURL.RawQuery = query.Encode()
		requestURL = parsedURL.String()
	}

	// Determine HTTP method
	method := endpoint.Method
	if method == "" {
		method = "POST"
	}

	// Prepare headers
	headers := make(map[string]string)
	if endpoint.Headers != nil {
		for k, v := range endpoint.Headers {
			headers[k] = v
		}
	}

	// Add message headers
	if message.Headers != nil {
		for k, v := range message.Headers {
			headers[k] = v
		}
	}

	// Set default Content-Type if not provided
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	// Send the request
	_, statusCode, err := utils.DoRequest(ctx, requestURL, utils.RequestOptions{
		Method:      method,
		Body:        message.Body,
		Headers:     headers,
		ContentType: headers["Content-Type"],
	})

	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}

	// Check if the status code indicates success (2xx range)
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("webhook request failed with status code: %d", statusCode)
	}

	return nil
}

// Name returns the name of the provider
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

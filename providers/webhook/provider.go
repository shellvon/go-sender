package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider represents a generic webhook provider
type Provider struct {
	endpoints []*Endpoint
	selector  *utils.Selector[*Endpoint]
}

var (
	_ core.Provider = (*Provider)(nil)
)

// New creates a new webhook provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("invalid webhook configuration: no endpoints configured or provider disabled")
	}

	// Convert to pointer slice
	endpoints := make([]*Endpoint, len(config.Endpoints))
	for i := range config.Endpoints {
		endpoints[i] = &config.Endpoints[i]
	}

	// Use common initialization logic
	enabledEndpoints, selector, err := utils.InitProvider(&config, endpoints)
	if err != nil {
		return nil, errors.New("no enabled endpoints")
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
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *webhook.Message, got %T", message))
	}

	if err := webhookMsg.Validate(); err != nil {
		return err
	}

	endpoint := p.selectEndpoint(ctx)
	if endpoint == nil {
		return errors.New("no available endpoint")
	}
	return p.doSendWebhook(ctx, endpoint, webhookMsg)
}

// selectEndpoint selects an endpoint based on context
func (p *Provider) selectEndpoint(ctx context.Context) *Endpoint {
	return p.selector.Select(ctx)
}

// doSendWebhook performs the actual HTTP request
func (p *Provider) doSendWebhook(ctx context.Context, endpoint *Endpoint, msg *Message) error {
	// 1. Build complete URL with query parameters
	finalURL := p.buildURLWithQueryParams(endpoint, msg)

	// 2. Merge headers
	headers := p.mergeHeaders(endpoint.Headers, msg.Headers)

	// 3. Prepare request body
	body, err := json.Marshal(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 4. Send HTTP request
	_, _, err = utils.DoRequest(ctx, finalURL, utils.RequestOptions{
		Method:      endpoint.Method,
		Headers:     headers,
		Body:        body,
		ContentType: "application/json",
	})

	return err
}

// buildURLWithQueryParams builds the complete URL with path variables and query parameters
func (p *Provider) buildURLWithQueryParams(endpoint *Endpoint, msg *Message) string {
	// 1. Start with base URL
	baseURL := endpoint.URL

	// 2. Replace path variables if any
	if msg.PathVars != nil {
		for key, value := range msg.PathVars {
			placeholder := "{" + key + "}"
			baseURL = strings.ReplaceAll(baseURL, placeholder, value)
		}
	}

	// 3. Add query parameters
	queryParams := p.mergeQueryParams(endpoint.QueryParams, msg.QueryParams)
	if len(queryParams) > 0 {
		// Parse existing URL to add query parameters
		parsedURL, err := url.Parse(baseURL)
		if err == nil {
			q := parsedURL.Query()
			for k, v := range queryParams {
				q.Set(k, v)
			}
			parsedURL.RawQuery = q.Encode()
			baseURL = parsedURL.String()
		}
	}

	return baseURL
}

// mergeHeaders merges endpoint headers with message headers
func (p *Provider) mergeHeaders(endpointHeaders, messageHeaders map[string]string) map[string]string {
	merged := make(map[string]string)

	// Add endpoint headers first
	for k, v := range endpointHeaders {
		merged[k] = v
	}

	// Override with message headers
	for k, v := range messageHeaders {
		merged[k] = v
	}

	return merged
}

// mergeQueryParams merges endpoint query parameters with message query parameters
func (p *Provider) mergeQueryParams(endpointParams, messageParams map[string]string) map[string]string {
	merged := make(map[string]string)

	// Add endpoint parameters first
	for k, v := range endpointParams {
		merged[k] = v
	}

	// Override with message parameters
	for k, v := range messageParams {
		merged[k] = v
	}

	return merged
}

// Name returns the provider name
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

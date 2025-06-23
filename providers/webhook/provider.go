package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

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
	// Build the request URL with path variables and query parameters
	requestURL := endpoint.URL

	// Use message's URL building functionality if path params or query params are provided
	if len(message.PathParams) > 0 || len(message.QueryParams) > 0 {
		builtURL, err := message.buildURL(endpoint.URL)
		if err != nil {
			return fmt.Errorf("failed to build URL: %w", err)
		}
		requestURL = builtURL
	} else if len(endpoint.QueryParams) > 0 {
		// Fallback to endpoint query params if no message query params
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

	// Determine HTTP method (message method overrides endpoint method)
	method := endpoint.Method
	if message.Method != "" {
		method = message.Method
	}
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

	// Send request
	responseBody, statusCode, err := utils.DoRequest(ctx, requestURL, utils.RequestOptions{
		Method:  method,
		Headers: headers,
		Raw:     message.Body,
	})

	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}

	// Handle response based on configuration
	return p.handleResponse(endpoint, statusCode, responseBody)
}

// handleResponse processes the webhook response based on endpoint configuration
func (p *Provider) handleResponse(endpoint *Endpoint, statusCode int, responseBody []byte) error {
	// Check status code first
	if !p.isSuccessStatusCode(endpoint, statusCode) {
		return fmt.Errorf("webhook request failed with status code: %d", statusCode)
	}

	// If no response validation is configured, consider it successful
	if endpoint.ResponseConfig == nil || !endpoint.ResponseConfig.ValidateResponse {
		return nil
	}

	// Validate response body based on configuration
	return p.validateResponseBody(endpoint.ResponseConfig, responseBody)
}

// isSuccessStatusCode checks if the status code indicates success
func (p *Provider) isSuccessStatusCode(endpoint *Endpoint, statusCode int) bool {
	// Use custom success status codes if configured
	if endpoint.ResponseConfig != nil && len(endpoint.ResponseConfig.SuccessStatusCodes) > 0 {
		for _, code := range endpoint.ResponseConfig.SuccessStatusCodes {
			if statusCode == code {
				return true
			}
		}
		return false
	}

	// Default: 2xx range indicates success
	return statusCode >= 200 && statusCode < 300
}

// validateResponseBody validates the response body based on configuration
func (p *Provider) validateResponseBody(config *ResponseConfig, responseBody []byte) error {
	switch config.ResponseType {
	case "json":
		return p.validateJSONResponse(config, responseBody)
	case "text":
		return p.validateTextResponse(config, responseBody)
	case "xml":
		return p.validateXMLResponse(config, responseBody)
	case "none":
		// No validation needed
		return nil
	default:
		// Default to no validation
		return nil
	}
}

// validateJSONResponse validates JSON response
func (p *Provider) validateJSONResponse(config *ResponseConfig, responseBody []byte) error {
	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Check success field if configured
	if config.SuccessField != "" {
		if successValue, exists := response[config.SuccessField]; exists {
			successStr := fmt.Sprintf("%v", successValue)
			if successStr != config.SuccessValue {
				// Extract error message if available
				errorMsg := p.extractErrorMessage(config, response)
				return fmt.Errorf("webhook returned failure: %s", errorMsg)
			}
		}
	}

	return nil
}

// validateTextResponse validates text response using regex patterns
func (p *Provider) validateTextResponse(config *ResponseConfig, responseBody []byte) error {
	responseText := string(responseBody)

	// Check error pattern first
	if config.ErrorPattern != "" {
		matched, err := regexp.MatchString(config.ErrorPattern, responseText)
		if err != nil {
			return fmt.Errorf("invalid error pattern: %w", err)
		}
		if matched {
			return fmt.Errorf("webhook returned error response: %s", responseText)
		}
	}

	// Check success pattern if configured
	if config.SuccessPattern != "" {
		matched, err := regexp.MatchString(config.SuccessPattern, responseText)
		if err != nil {
			return fmt.Errorf("invalid success pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("webhook response does not match success pattern: %s", responseText)
		}
	}

	return nil
}

// validateXMLResponse validates XML response (placeholder for future implementation)
func (p *Provider) validateXMLResponse(config *ResponseConfig, responseBody []byte) error {
	// TODO: Implement XML response validation
	return fmt.Errorf("XML response validation not yet implemented")
}

// extractErrorMessage extracts error message from JSON response
func (p *Provider) extractErrorMessage(config *ResponseConfig, response map[string]interface{}) string {
	if config.ErrorField != "" {
		if errorValue, exists := response[config.ErrorField]; exists {
			return fmt.Sprintf("%v", errorValue)
		}
	}
	if config.MessageField != "" {
		if messageValue, exists := response[config.MessageField]; exists {
			return fmt.Sprintf("%v", messageValue)
		}
	}
	return "unknown error"
}

// Name returns the name of the provider
func (p *Provider) Name() string {
	return string(core.ProviderTypeWebhook)
}

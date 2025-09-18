package webhook

import (
	"errors"

	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Webhook provider.
// Webhook provider uses Endpoint instead of Account because:
//   - Webhook is essentially calling different endpoints, not managing accounts
//   - Endpoint directly contains URL, Method, Headers and other webhook-specific configurations
//   - This design better fits the webhook usage scenario
//   - Endpoint implements core.Selectable interface to maintain architectural compatibility
type Config = core.BaseConfig[*Endpoint]

// Endpoint represents a single webhook endpoint configuration.
type Endpoint struct {
	Name        string            `json:"name"`         // Endpoint name for selection
	URL         string            `json:"url"`          // Webhook URL
	Method      string            `json:"method"`       // HTTP method, defaults to POST
	Headers     map[string]string `json:"headers"`      // Fixed request headers
	QueryParams map[string]string `json:"query_params"` // Fixed query parameters
	Weight      int               `json:"weight"`       // Weight for weighted strategy
	Disabled    bool              `json:"disabled"`     // Whether this endpoint is disabled (default: false)

	// Response handling configuration
	ResponseConfig *ResponseConfig `json:"response_config,omitempty"` // Response handling configuration
}

// Endpoint implements core.Selectable and core.Validatable.
var _ core.Selectable = (*Endpoint)(nil)
var _ core.Validatable = (*Endpoint)(nil)

// ResponseConfig defines how to handle webhook responses.
type ResponseConfig = core.ResponseHandlerConfig

// EndpointOption represents a function that modifies Webhook Endpoint configuration.
type EndpointOption func(*Endpoint)

// NewEndpoint creates a new webhook endpoint with the given URL and options.
//
// Example:
//
//	endpoint := webhook.NewEndpoint("https://api.example.com/webhook",
//	    webhook.Name("primary-webhook"),
//	    webhook.Weight(2),
//	    webhook.WithMethod("POST"),
//	    webhook.WithHeader("Content-Type", "application/json"))
func NewEndpoint(url string, opts ...EndpointOption) *Endpoint {
	endpoint := &Endpoint{
		Name:        "webhook-default",
		URL:         url,
		Method:      "POST", // Default to POST
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
		Weight:      1,
		Disabled:    false,
	}

	// Apply all options
	for _, opt := range opts {
		opt(endpoint)
	}

	return endpoint
}

// IsEnabled checks if the endpoint is enabled.
func (e *Endpoint) IsEnabled() bool {
	return !e.Disabled
}

// GetName returns the endpoint name for strategy selection.
func (e *Endpoint) GetName() string {
	return e.Name
}

// GetWeight returns the endpoint weight for strategy selection.
func (e *Endpoint) GetWeight() int {
	if e.Weight <= 0 {
		return 1
	}
	return e.Weight
}

// GetType returns the subprovider type of this endpoint.
func (e *Endpoint) GetType() string {
	return "" // Webhook endpoints don't have subprovider types
}

// Validate validates the endpoint configuration.
func (e *Endpoint) Validate() error {
	if e.Name == "" {
		return errors.New("name is required")
	}
	if e.URL == "" {
		return errors.New("url is required")
	}
	return nil
}

// Webhook-specific endpoint options

// Name sets the endpoint name.
func Name(name string) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Name = name
	}
}

// Weight sets the endpoint weight for load balancing.
func Weight(weight int) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Weight = weight
	}
}

// Disabled disables the endpoint.
func Disabled() EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Disabled = true
	}
}

// WithMethod sets the HTTP method (GET, POST, PUT, DELETE, etc.)
func WithMethod(method string) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.Method = method
	}
}

// WithHeader adds a header to the webhook request.
func WithHeader(key, value string) EndpointOption {
	return func(endpoint *Endpoint) {
		if endpoint.Headers == nil {
			endpoint.Headers = make(map[string]string)
		}
		endpoint.Headers[key] = value
	}
}

// WithHeaders sets multiple headers at once.
func WithHeaders(headers map[string]string) EndpointOption {
	return func(endpoint *Endpoint) {
		if endpoint.Headers == nil {
			endpoint.Headers = make(map[string]string)
		}
		for k, v := range headers {
			endpoint.Headers[k] = v
		}
	}
}

// WithQueryParam adds a query parameter to the webhook URL.
func WithQueryParam(key, value string) EndpointOption {
	return func(endpoint *Endpoint) {
		if endpoint.QueryParams == nil {
			endpoint.QueryParams = make(map[string]string)
		}
		endpoint.QueryParams[key] = value
	}
}

// WithQueryParams sets multiple query parameters at once.
func WithQueryParams(params map[string]string) EndpointOption {
	return func(endpoint *Endpoint) {
		if endpoint.QueryParams == nil {
			endpoint.QueryParams = make(map[string]string)
		}
		for k, v := range params {
			endpoint.QueryParams[k] = v
		}
	}
}

// WithResponseConfig sets the response handling configuration.
func WithResponseConfig(config *ResponseConfig) EndpointOption {
	return func(endpoint *Endpoint) {
		endpoint.ResponseConfig = config
	}
}

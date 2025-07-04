package webhook

import (
	"github.com/shellvon/go-sender/core"
)

// Webhook provider uses Endpoint instead of Account because:
//   - Webhook is essentially calling different endpoints, not managing accounts
//   - Endpoint directly contains URL, Method, Headers and other webhook-specific configurations
//   - This design better fits the webhook usage scenario
//   - Endpoint implements core.Selectable interface to maintain architectural compatibility
//
// Config holds configuration for the Webhook provider.
type Config struct {
	core.ProviderMeta

	Endpoints []*Endpoint `json:"endpoints"` // Multiple webhook endpoints configuration
}

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

// Endpoint implements core.Selectable.
var _ core.Selectable = (*Endpoint)(nil)

// ResponseConfig defines how to handle webhook responses.
type ResponseConfig struct {
	// Success criteria
	SuccessStatusCodes []int `json:"success_status_codes,omitempty"` // Custom success status codes (default: 2xx)

	// Response validation
	ValidateResponse bool `json:"validate_response,omitempty"` // Whether to validate response body

	// Response parsing
	ResponseType core.BodyType `json:"response_type,omitempty"` // "json", "text", "xml", "none"

	// JSON response validation (when ResponseType is "json")
	SuccessField string `json:"success_field,omitempty"` // Field name indicating success (e.g., "success", "ok")
	SuccessValue string `json:"success_value,omitempty"` // Expected value for success (e.g., "true", "ok")
	ErrorField   string `json:"error_field,omitempty"`   // Field name containing error message
	MessageField string `json:"message_field,omitempty"` // Field name containing response message

	// Text response validation (when ResponseType is "text")
	SuccessPattern string `json:"success_pattern,omitempty"` // Regex pattern for success response
	ErrorPattern   string `json:"error_pattern,omitempty"`   // Regex pattern for error response
}

// IsConfigured checks if the Webhook configuration is valid.
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Endpoints) > 0
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

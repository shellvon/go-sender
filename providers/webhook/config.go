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

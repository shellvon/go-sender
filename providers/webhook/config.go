package webhook

import (
	"github.com/shellvon/go-sender/core"
)

// Config holds configuration for the Webhook provider
type Config struct {
	core.BaseConfig
	Endpoints         []Endpoint `json:"endpoints"`          // Multiple webhook endpoints configuration
	DisableMiddleware bool       `json:"disable_middleware"` // Whether to disable middleware (for embedded usage)
}

// Endpoint represents a single webhook endpoint configuration
type Endpoint struct {
	Name        string            `json:"name"`         // Endpoint name for selection
	URL         string            `json:"url"`          // Webhook URL
	Method      string            `json:"method"`       // HTTP method, defaults to POST
	Headers     map[string]string `json:"headers"`      // Fixed request headers
	QueryParams map[string]string `json:"query_params"` // Fixed query parameters
	Weight      int               `json:"weight"`       // Weight for weighted strategy
	Disabled    bool              `json:"disabled"`     // Whether this endpoint is disabled (default: false)
}

// IsConfigured checks if the Webhook configuration is valid
func (c Config) IsConfigured() bool {
	return !c.IsDisabled() && len(c.Endpoints) > 0
}

// IsEnabled checks if the endpoint is enabled
func (e *Endpoint) IsEnabled() bool {
	return !e.Disabled
}

// GetName returns the endpoint name for strategy selection
func (e *Endpoint) GetName() string {
	return e.Name
}

// GetWeight returns the endpoint weight for strategy selection
func (e *Endpoint) GetWeight() int {
	if e.Weight <= 0 {
		return 1
	}
	return e.Weight
}

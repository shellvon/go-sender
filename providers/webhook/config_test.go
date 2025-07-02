package webhook_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/webhook"
)

func TestConfig_IsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   webhook.Config
		expected bool
	}{
		{
			name: "valid config",
			config: webhook.Config{
				Endpoints: []webhook.Endpoint{
					{Name: "test", URL: "https://example.com/webhook"},
				},
			},
			expected: true,
		},
		{
			name:     "empty config",
			config:   webhook.Config{},
			expected: false,
		},
		{
			name: "disabled config",
			config: webhook.Config{
				BaseConfig: core.BaseConfig{Disabled: true},
				Endpoints: []webhook.Endpoint{
					{Name: "test", URL: "https://example.com/webhook"},
				},
			},
			expected: false,
		},
		{
			name: "no endpoints",
			config: webhook.Config{
				Endpoints: []webhook.Endpoint{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsConfigured()
			if result != tt.expected {
				t.Errorf("Expected IsConfigured() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEndpoint_IsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		endpoint webhook.Endpoint
		expected bool
	}{
		{
			name: "enabled endpoint",
			endpoint: webhook.Endpoint{
				Name: "test",
				URL:  "https://example.com/webhook",
			},
			expected: true,
		},
		{
			name: "disabled endpoint",
			endpoint: webhook.Endpoint{
				Name:     "test",
				URL:      "https://example.com/webhook",
				Disabled: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.endpoint.IsEnabled()
			if result != tt.expected {
				t.Errorf("Expected IsEnabled() to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEndpoint_GetName(t *testing.T) {
	endpoint := webhook.Endpoint{
		Name: "test-endpoint",
		URL:  "https://example.com/webhook",
	}

	result := endpoint.GetName()
	if result != "test-endpoint" {
		t.Errorf("Expected GetName() to return 'test-endpoint', got '%s'", result)
	}
}

func TestEndpoint_GetWeight(t *testing.T) {
	tests := []struct {
		name     string
		endpoint webhook.Endpoint
		expected int
	}{
		{
			name: "positive weight",
			endpoint: webhook.Endpoint{
				Name:   "test",
				Weight: 5,
			},
			expected: 5,
		},
		{
			name: "zero weight",
			endpoint: webhook.Endpoint{
				Name:   "test",
				Weight: 0,
			},
			expected: 1, // Default weight
		},
		{
			name: "negative weight",
			endpoint: webhook.Endpoint{
				Name:   "test",
				Weight: -1,
			},
			expected: 1, // Default weight
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.endpoint.GetWeight()
			if result != tt.expected {
				t.Errorf("Expected GetWeight() to return %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestEndpoint_GetType(t *testing.T) {
	endpoint := webhook.Endpoint{
		Name: "test",
		URL:  "https://example.com/webhook",
	}

	result := endpoint.GetType()
	if result != "" {
		t.Errorf("Expected GetType() to return empty string, got '%s'", result)
	}
}

func TestResponseConfig_Validation(t *testing.T) {
	// Test that ResponseConfig can be created without issues
	config := &webhook.ResponseConfig{
		SuccessStatusCodes: []int{200, 201},
		ValidateResponse:   true,
	}

	// Verify the config was created with expected values
	if config.SuccessStatusCodes[0] != 200 {
		t.Error("Expected first success status code to be 200")
	}
	if !config.ValidateResponse {
		t.Error("Expected ValidateResponse to be true")
	}
}

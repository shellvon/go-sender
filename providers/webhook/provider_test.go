package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/webhook"
)

// invalidMessage is a test message that doesn't implement webhook.Message.
type invalidMessage struct {
	core.DefaultMessage
}

func (m *invalidMessage) ProviderType() core.ProviderType {
	return "invalid"
}

func (m *invalidMessage) Validate() error {
	return nil
}

func TestNewProvider(t *testing.T) {
	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    "https://example.com/webhook",
				Method: "POST",
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	if provider.Name() != "webhook" {
		t.Errorf("Expected name 'webhook', got '%s'", provider.Name())
	}
}

func TestNewProvider_NotConfigured(t *testing.T) {
	config := webhook.Config{}

	_, err := webhook.New(config)
	if err == nil {
		t.Error("Expected error for not configured, got nil")
	}
	if !strings.Contains(err.Error(), "webhook provider is not configured or is disabled") {
		t.Errorf("Expected error to contain 'webhook provider is not configured or is disabled', got '%s'", err.Error())
	}
}

func TestNewProvider_Disabled(t *testing.T) {
	config := webhook.Config{
		BaseConfig: core.BaseConfig{Disabled: true},
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    "https://example.com/webhook",
				Method: "POST",
			},
		},
	}

	_, err := webhook.New(config)
	if err == nil {
		t.Error("Expected error for disabled config, got nil")
	}
	if !strings.Contains(err.Error(), "webhook provider is not configured or is disabled") {
		t.Errorf("Expected error to contain 'webhook provider is not configured or is disabled', got '%s'", err.Error())
	}
}

func TestNewProvider_NoEnabledEndpoints(t *testing.T) {
	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:     "test",
				URL:      "https://example.com/webhook",
				Method:   "POST",
				Disabled: true,
			},
		},
	}

	_, err := webhook.New(config)
	if err == nil {
		t.Error("Expected error for no enabled endpoints, got nil")
	}
	if !strings.Contains(err.Error(), "no enabled webhook endpoints found") {
		t.Errorf("Expected error to contain 'no enabled webhook endpoints found', got '%s'", err.Error())
	}
}

func TestProvider_Send_Success(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL,
				Method: "POST",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_WithPathParams(t *testing.T) {
	// Create mock server that checks path
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhook/user/123" {
			t.Errorf("Expected path '/webhook/user/123', got '%s'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL + "/webhook/{type}/{id}",
				Method: "POST",
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage(
		[]byte(`{"test": "data"}`),
		webhook.WithPathParams(map[string]string{
			"type": "user",
			"id":   "123",
		}),
	)

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_WithQueryParams(t *testing.T) {
	// Create mock server that checks query parameters
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("version") != "v1" {
			t.Errorf("Expected query param 'version=v1', got '%s'", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL,
				Method: "POST",
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage(
		[]byte(`{"test": "data"}`),
		webhook.WithQueryParams(map[string]string{
			"version": "v1",
		}),
	)

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_WithMethodOverride(t *testing.T) {
	// Create mock server that checks method
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method 'PUT', got '%s'", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL,
				Method: "POST", // This should be overridden
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage(
		[]byte(`{"test": "data"}`),
		webhook.WithMethod("PUT"),
	)

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_HTTPFailure(t *testing.T) {
	// Create mock server that returns error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL,
				Method: "POST",
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	err = provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for HTTP failure, got nil")
	}
}

func TestProvider_Send_JSONResponseValidation(t *testing.T) {
	// Create mock server that returns JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": "true", "message": "ok"}`))
	}))
	defer ts.Close()

	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    ts.URL,
				Method: "POST",
				ResponseConfig: &webhook.ResponseConfig{
					ValidateResponse: true,
					ResponseType:     "json",
					SuccessField:     "success",
					SuccessValue:     "true",
					ErrorField:       "error",
				},
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_InvalidMessageType(t *testing.T) {
	config := webhook.Config{
		Endpoints: []webhook.Endpoint{
			{
				Name:   "test",
				URL:    "https://example.com/webhook",
				Method: "POST",
			},
		},
	}

	provider, err := webhook.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := &invalidMessage{}

	err = provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for invalid message type, got nil")
	}
}

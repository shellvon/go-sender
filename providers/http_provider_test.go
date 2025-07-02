package providers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

// mockConfig implements core.Selectable for testing.
type mockConfig struct {
	name    string
	enabled bool
	typ     string
}

func (m *mockConfig) GetName() string   { return m.name }
func (m *mockConfig) IsEnabled() bool   { return m.enabled }
func (m *mockConfig) GetType() string   { return m.typ }
func (m *mockConfig) GetWeight() int    { return 1 }
func (m *mockConfig) IsAvailable() bool { return m.enabled }

// mockMessage implements core.Message for testing.
type mockMessage struct {
	subProvider string
}

func (m *mockMessage) Validate() error                 { return nil }
func (m *mockMessage) ProviderType() core.ProviderType { return core.ProviderTypeSMS }
func (m *mockMessage) MsgID() string                   { return "test-msg-id" }
func (m *mockMessage) GetSubProvider() string          { return m.subProvider }

// mockTransformer implements core.HTTPTransformer for testing.
type mockTransformer struct {
	shouldFail bool
	reqSpec    *core.HTTPRequestSpec
	handler    core.ResponseHandler
}

func (m *mockTransformer) CanTransform(_ core.Message) bool { return true }

func (m *mockTransformer) Transform(
	_ context.Context,
	_ core.Message,
	_ *mockConfig,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if m.shouldFail {
		return nil, nil, errors.New("transform failed")
	}
	return m.reqSpec, m.handler, nil
}

func TestNewHTTPProvider(t *testing.T) {
	configs := []*mockConfig{
		{name: "test1", enabled: true, typ: "type1"},
		{name: "test2", enabled: true, typ: "type2"},
	}
	transformer := &mockTransformer{}
	strategy := &core.RoundRobinStrategy{}

	provider := providers.NewHTTPProvider("test-provider", configs, transformer, strategy)

	if provider.Name() != "test-provider" {
		t.Errorf("Expected name 'test-provider', got '%s'", provider.Name())
	}

	configsResult := provider.GetConfigs()
	if len(configsResult) != 2 {
		t.Errorf("Expected 2 configs, got %d", len(configsResult))
	}
}

func TestHTTPProvider_Send_SingleConfig(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	configs := []*mockConfig{
		{name: "test", enabled: true, typ: "type1"},
	}

	reqSpec := &core.HTTPRequestSpec{
		URL:      ts.URL,
		Method:   "POST",
		BodyType: "json",
		Body:     []byte(`{"test": "data"}`),
	}

	transformer := &mockTransformer{
		reqSpec: reqSpec,
		handler: func(statusCode int, _ []byte) error {
			if statusCode != http.StatusOK {
				return errors.New("unexpected status code")
			}
			return nil
		},
	}

	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHTTPProvider_Send_NoConfigs(t *testing.T) {
	configs := []*mockConfig{}
	transformer := &mockTransformer{}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for no configs, got nil")
	}
	if err.Error() != "no available config" {
		t.Errorf("Expected 'no available config' error, got '%s'", err.Error())
	}
}

func TestHTTPProvider_Send_DisabledConfig(t *testing.T) {
	configs := []*mockConfig{
		{name: "test", enabled: false, typ: "type1"},
	}
	transformer := &mockTransformer{}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for disabled config, got nil")
	}
	if err.Error() != "the selected account is disabled" {
		t.Errorf("Expected 'the selected account is disabled' error, got '%s'", err.Error())
	}
}

func TestHTTPProvider_Send_TransformFailure(t *testing.T) {
	configs := []*mockConfig{
		{name: "test", enabled: true, typ: "type1"},
	}
	transformer := &mockTransformer{shouldFail: true}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for transform failure, got nil")
	}
	// HTTP provider wraps the error, so check if the error message contains the original error
	if !strings.Contains(err.Error(), "transform failed") {
		t.Errorf("Expected error to contain 'transform failed', got '%s'", err.Error())
	}
}

func TestHTTPProvider_ExecuteHTTPRequest_WithQueryParams(t *testing.T) {
	// Create mock server that checks query parameters
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("param1") != "value1" {
			t.Errorf("Expected query param 'param1=value1', got '%s'", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer ts.Close()

	reqSpec := &core.HTTPRequestSpec{
		URL: ts.URL,
		QueryParams: map[string]string{
			"param1": "value1",
		},
		Method:   "GET",
		BodyType: "json",
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	configs := []*mockConfig{{name: "test", enabled: true, typ: "type1"}}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHTTPProvider_ExecuteHTTPRequest_InvalidURL(t *testing.T) {
	reqSpec := &core.HTTPRequestSpec{
		URL:         "://invalid-url",
		Method:      "GET",
		BodyType:    "json",
		QueryParams: map[string]string{"param": "value"},
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	configs := []*mockConfig{{name: "test", enabled: true, typ: "type1"}}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestHTTPProvider_ExecuteHTTPRequest_HTTPFailure(t *testing.T) {
	// Create mock server that returns error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer ts.Close()

	reqSpec := &core.HTTPRequestSpec{
		URL:      ts.URL,
		Method:   "POST",
		BodyType: "json",
		Body:     []byte(`{"test": "data"}`),
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	configs := []*mockConfig{{name: "test", enabled: true, typ: "type1"}}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for HTTP failure, got nil")
	}
}

func TestHTTPProvider_ExecuteHTTPRequest_CustomHandler(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"custom": "response"}`))
	}))
	defer ts.Close()

	customHandlerCalled := false
	customHandler := func(statusCode int, body []byte) error {
		customHandlerCalled = true
		if statusCode != http.StatusOK {
			return errors.New("unexpected status code")
		}
		if string(body) != `{"custom": "response"}` {
			return errors.New("unexpected response body")
		}
		return nil
	}

	reqSpec := &core.HTTPRequestSpec{
		URL:      ts.URL,
		Method:   "POST",
		BodyType: "json",
		Body:     []byte(`{"test": "data"}`),
	}

	transformer := &mockTransformer{
		reqSpec: reqSpec,
		handler: customHandler,
	}
	configs := []*mockConfig{{name: "test", enabled: true, typ: "type1"}}
	provider := providers.NewHTTPProvider("test-provider", configs, transformer, &core.RoundRobinStrategy{})
	msg := &mockMessage{}

	err := provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !customHandlerCalled {
		t.Error("Expected custom handler to be called")
	}
}

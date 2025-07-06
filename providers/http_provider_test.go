package providers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
)

type mockSelectable struct {
	name    string
	enabled bool
	typ     string
}

// --- core.Selectable ---.
func (m *mockSelectable) GetName() string   { return m.name }
func (m *mockSelectable) IsEnabled() bool   { return m.enabled }
func (m *mockSelectable) GetType() string   { return m.typ }
func (m *mockSelectable) GetWeight() int    { return 1 }
func (m *mockSelectable) IsAvailable() bool { return m.enabled }

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

// Implements core.HTTPTransformer[*mockSelectable].
func (m *mockTransformer) CanTransform(_ core.Message) bool { return true }

func (m *mockTransformer) Transform(
	_ context.Context,
	_ core.Message,
	_ *mockSelectable,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if m.shouldFail {
		return nil, nil, errors.New("transform failed")
	}
	return m.reqSpec, m.handler, nil
}

func TestNewHTTPProvider(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test1", enabled: true, typ: "type1"},
			{name: "test2", enabled: true, typ: "type2"},
		},
	}
	transformer := &mockTransformer{}
	provider, err := providers.NewHTTPProvider("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}

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

	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test", enabled: true, typ: "type1"},
		},
	}

	reqSpec := &core.HTTPRequestSpec{
		URL:      ts.URL,
		Method:   "POST",
		BodyType: core.BodyTypeJSON,
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

	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHTTPProvider_Send_NoConfigs(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{}}
	transformer := &mockTransformer{}

	_, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err == nil {
		t.Error("Expected error for no configs, got nil")
	}
}

func TestHTTPProvider_Send_DisabledConfig(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: false, typ: "type1"}}}
	transformer := &mockTransformer{}

	_, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err == nil {
		t.Error("Expected error for disabled config, got nil")
	}
}

func TestHTTPProvider_Send_TransformFailure(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: true, typ: "type1"}}}
	transformer := &mockTransformer{shouldFail: true}

	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
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
		URL:         ts.URL,
		QueryParams: url.Values{"param1": {"value1"}},
		Method:      "GET",
		BodyType:    core.BodyTypeJSON,
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: true, typ: "type1"}}}
	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHTTPProvider_ExecuteHTTPRequest_InvalidURL(t *testing.T) {
	reqSpec := &core.HTTPRequestSpec{
		URL:         "://invalid-url",
		Method:      "GET",
		BodyType:    core.BodyTypeJSON,
		QueryParams: url.Values{"param": {"value"}},
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: true, typ: "type1"}}}
	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
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
		BodyType: core.BodyTypeJSON,
		Body:     []byte(`{"test": "data"}`),
	}

	transformer := &mockTransformer{reqSpec: reqSpec}
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: true, typ: "type1"}}}
	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
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
		BodyType: core.BodyTypeJSON,
		Body:     []byte(`{"test": "data"}`),
	}

	transformer := &mockTransformer{
		reqSpec: reqSpec,
		handler: customHandler,
	}
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: true, typ: "type1"}}}
	provider, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err != nil {
		t.Fatalf("unexpected error creating provider: %v", err)
	}
	msg := &mockMessage{}

	err = provider.Send(context.Background(), msg, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !customHandlerCalled {
		t.Error("Expected custom handler to be called")
	}
}

func TestHTTPProvider_New_AllDisabled(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{Items: []*mockSelectable{{name: "test", enabled: false, typ: "type1"}}}
	transformer := &mockTransformer{}

	_, err := providers.NewHTTPProvider[*mockSelectable]("test-provider", transformer, config)
	if err == nil {
		t.Error("Expected error when all configs are disabled, got nil")
	}
}

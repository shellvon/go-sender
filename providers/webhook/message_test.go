package webhook_test

import (
	"net/http"
	"testing"

	"github.com/shellvon/go-sender/providers/webhook"
)

func TestNewMessage(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	msg := webhook.NewMessage(body)

	if msg == nil {
		t.Fatal("NewMessage returned nil")
	}

	if string(msg.Body) != string(body) {
		t.Errorf("Expected Body to be %s, got %s", string(body), string(msg.Body))
	}

	if msg.Method != "" {
		t.Error("Expected Method to be empty by default")
	}

	if msg.Headers != nil {
		t.Error("Expected Headers to be nil by default")
	}

	if msg.PathParams != nil {
		t.Error("Expected PathParams to be nil by default")
	}

	if msg.QueryParams != nil {
		t.Error("Expected QueryParams to be nil by default")
	}
}

func TestMessageOptions(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	headers := map[string]string{"Authorization": "Bearer token"}
	pathParams := map[string]string{"id": "123"}
	queryParams := map[string]string{"version": "v1"}

	msg := webhook.NewMessage(
		body,
		webhook.WithMethod("PUT"),
		webhook.WithHeaders(headers),
		webhook.WithPathParams(pathParams),
		webhook.WithQueryParams(queryParams),
	)

	if msg.Method != http.MethodPut {
		t.Errorf("Expected Method to be 'PUT', got '%s'", msg.Method)
	}

	if len(msg.Headers) != 1 || msg.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Headers to contain Authorization, got %v", msg.Headers)
	}

	if len(msg.PathParams) != 1 || msg.PathParams["id"] != "123" {
		t.Errorf("Expected PathParams to contain id=123, got %v", msg.PathParams)
	}

	if len(msg.QueryParams) != 1 || msg.QueryParams["version"] != "v1" {
		t.Errorf("Expected QueryParams to contain version=v1, got %v", msg.QueryParams)
	}
}

func TestMessageOptions_MultipleHeaders(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	headers1 := map[string]string{"Authorization": "Bearer token"}
	headers2 := map[string]string{"Content-Type": "application/json", "X-Custom": "value"}

	msg := webhook.NewMessage(
		body,
		webhook.WithHeaders(headers1),
		webhook.WithHeaders(headers2),
	)

	// Later headers should override earlier ones
	expectedHeaders := map[string]string{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
		"X-Custom":      "value",
	}

	for k, v := range expectedHeaders {
		if msg.Headers[k] != v {
			t.Errorf("Expected header %s to be '%s', got '%s'", k, v, msg.Headers[k])
		}
	}
}

func TestMessageOptions_MultiplePathParams(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	params1 := map[string]string{"id": "123"}
	params2 := map[string]string{"type": "user", "action": "update"}

	msg := webhook.NewMessage(
		body,
		webhook.WithPathParams(params1),
		webhook.WithPathParams(params2),
	)

	expectedParams := map[string]string{
		"id":     "123",
		"type":   "user",
		"action": "update",
	}

	for k, v := range expectedParams {
		if msg.PathParams[k] != v {
			t.Errorf("Expected path param %s to be '%s', got '%s'", k, v, msg.PathParams[k])
		}
	}
}

func TestMessageOptions_MultipleQueryParams(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	params1 := map[string]string{"version": "v1"}
	params2 := map[string]string{"format": "json", "pretty": "true"}

	msg := webhook.NewMessage(
		body,
		webhook.WithQueryParams(params1),
		webhook.WithQueryParams(params2),
	)

	expectedParams := map[string]string{
		"version": "v1",
		"format":  "json",
		"pretty":  "true",
	}

	for k, v := range expectedParams {
		if msg.QueryParams[k] != v {
			t.Errorf("Expected query param %s to be '%s', got '%s'", k, v, msg.QueryParams[k])
		}
	}
}

func TestMessage_Validate(t *testing.T) {
	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	err := msg.Validate()
	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestMessage_ProviderType(t *testing.T) {
	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	providerType := msg.ProviderType()
	if providerType != "webhook" {
		t.Errorf("Expected ProviderType to be 'webhook', got '%s'", providerType)
	}
}

func TestMessage_MsgID(t *testing.T) {
	msg := webhook.NewMessage([]byte(`{"test": "data"}`))

	msgID := msg.MsgID()
	if msgID == "" {
		t.Error("Expected MsgID to be non-empty")
	}

	// MsgID should be consistent for the same message
	msgID2 := msg.MsgID()
	if msgID != msgID2 {
		t.Error("Expected MsgID to be consistent")
	}
}

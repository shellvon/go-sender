package webhook_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/shellvon/go-sender/providers/webhook"
)

func TestWebhookBuilder_Minimal(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	msg := webhook.Webhook().Body(body).Build()

	if msg == nil {
		t.Fatal("Webhook().Build() returned nil")
	}

	if string(msg.Body) != string(body) {
		t.Errorf("Expected Body to be %s, got %s", string(body), string(msg.Body))
	}

	if msg.Method != "" {
		t.Error("Expected Method to be empty by default")
	}

	if msg.Headers == nil || len(msg.Headers) != 0 {
		t.Error("Expected Headers to be empty by default")
	}

	if msg.PathParams == nil || len(msg.PathParams) != 0 {
		t.Error("Expected PathParams to be empty by default")
	}

	if msg.QueryParams == nil || len(msg.QueryParams) != 0 {
		t.Error("Expected QueryParams to be empty by default")
	}
}

func TestWebhookBuilder_AllOptions(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	headers := map[string]string{"Authorization": "Bearer token"}
	pathParams := map[string]string{"id": "123"}
	queryParams := map[string]string{"version": "v1"}

	msg := webhook.Webhook().
		Body(body).
		Method(http.MethodPut).
		Headers(headers).
		PathParams(pathParams).
		Queries(queryParams).
		Build()

	if msg.Method != http.MethodPut {
		t.Errorf("Expected Method to be '%s', got '%s'", http.MethodPut, msg.Method)
	}

	if len(msg.Headers) != 1 || msg.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Headers to contain Authorization, got %v", msg.Headers)
	}

	if len(msg.PathParams) != 1 || msg.PathParams["id"] != "123" {
		t.Errorf("Expected PathParams to contain id=123, got %v", msg.PathParams)
	}

	if len(msg.QueryParams) != 1 || msg.QueryParams.Get("version") != "v1" {
		t.Errorf("Expected QueryParams to contain version=v1, got %v", msg.QueryParams)
	}
}

func TestWebhookBuilder_MultipleHeaders(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	headers1 := map[string]string{"Authorization": "Bearer token"}
	headers2 := map[string]string{"Content-Type": "application/json", "X-Custom": "value"}

	msg := webhook.Webhook().
		Body(body).
		Headers(headers1).
		Headers(headers2).
		Build()

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

func TestWebhookBuilder_MultiplePathParams(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	params1 := map[string]string{"id": "123"}
	params2 := map[string]string{"type": "user", "action": "update"}

	msg := webhook.Webhook().
		Body(body).
		PathParams(params1).
		PathParams(params2).
		Build()

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

func TestWebhookBuilder_MultipleQueryParams(t *testing.T) {
	body := []byte(`{"test": "data"}`)
	params1 := map[string]string{"version": "v1"}
	params2 := map[string]string{"format": "json", "pretty": "true"}

	msg := webhook.Webhook().
		Body(body).
		Queries(params1).
		Queries(params2).
		Build()

	expectedParams := url.Values{
		"version": {"v1"},
		"format":  {"json"},
		"pretty":  {"true"},
	}

	for k, vals := range expectedParams {
		if msg.QueryParams.Get(k) != vals[0] {
			t.Errorf("Expected query param %s to be '%s', got '%s'", k, vals[0], msg.QueryParams.Get(k))
		}
	}
}

func TestWebhookBuilder_Validate(t *testing.T) {
	msg := webhook.Webhook().Body([]byte(`{"test": "data"}`)).Build()

	err := msg.Validate()
	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestWebhookBuilder_ProviderType(t *testing.T) {
	msg := webhook.Webhook().Body([]byte(`{"test": "data"}`)).Build()

	providerType := msg.ProviderType()
	if providerType != "webhook" {
		t.Errorf("Expected ProviderType to be 'webhook', got '%s'", providerType)
	}
}

func TestWebhookBuilder_MsgID(t *testing.T) {
	msg := webhook.Webhook().Body([]byte(`{"test": "data"}`)).Build()

	msgID := msg.MsgID()
	if msgID == "" {
		t.Error("Expected MsgID to be non-empty")
	}

	msgID2 := msg.MsgID()
	if msgID != msgID2 {
		t.Error("Expected MsgID to be consistent")
	}
}

//go:build legacy
// +build legacy

package telegram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/telegram"
)

// invalidMessage is a test message that doesn't implement telegram.Message.
type invalidMessage struct {
	core.DefaultMessage
}

func (m *invalidMessage) ProviderType() core.ProviderType {
	return "invalid"
}

func (m *invalidMessage) Validate() error {
	return nil
}

// rewriteRoundTripper 用于在测试中把原始 telegram API URL 重写到 httptest server.
type rewriteRoundTripper struct {
	base         http.RoundTripper
	targetHost   string
	targetScheme string
}

func (rt rewriteRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Host = rt.targetHost
	req.URL.Scheme = rt.targetScheme
	return rt.base.RoundTrip(req)
}

func TestNewProvider(t *testing.T) {
	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	if provider.Name() != "telegram" {
		t.Errorf("Expected name 'telegram', got '%s'", provider.Name())
	}
}

func TestNewProvider_NotConfigured(t *testing.T) {
	config := telegram.Config{}

	_, err := telegram.New(config)
	if err == nil {
		t.Error("Expected error for not configured, got nil")
	}
	if !strings.Contains(err.Error(), "telegram provider is not configured or is disabled") {
		t.Errorf(
			"Expected error to contain 'telegram provider is not configured or is disabled', got '%s'",
			err.Error(),
		)
	}
}

func TestNewProvider_Disabled(t *testing.T) {
	config := telegram.Config{
		BaseConfig: core.BaseConfig{Disabled: true},
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	_, err := telegram.New(config)
	if err == nil {
		t.Error("Expected error for disabled config, got nil")
	}
	if !strings.Contains(err.Error(), "telegram provider is not configured or is disabled") {
		t.Errorf(
			"Expected error to contain 'telegram provider is not configured or is disabled', got '%s'",
			err.Error(),
		)
	}
}

func TestNewProvider_NoEnabledAccounts(t *testing.T) {
	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:     "test",
				APIKey:   "bot123:token",
				Disabled: true,
			},
		},
	}

	_, err := telegram.New(config)
	if err == nil {
		t.Error("Expected error for no enabled accounts, got nil")
	}
	if !strings.Contains(err.Error(), "no enabled telegram accounts found") {
		t.Errorf("Expected error to contain 'no enabled telegram accounts found', got '%s'", err.Error())
	}
}

func TestProvider_Send_Success(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true, "result": {"message_id": 123}}`))
	}))
	defer ts.Close()

	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := telegram.NewTextMessage("123456789", "Hello, World!")

	// 覆盖 HTTPClient，将请求导向 mock server。
	tsURL := ts.URL // e.g., http://127.0.0.1:XXXXX
	parts := strings.SplitN(tsURL, "://", 2)
	scheme := parts[0]
	host := parts[1]
	client := &http.Client{
		Transport: rewriteRoundTripper{base: http.DefaultTransport, targetHost: host, targetScheme: scheme},
	}

	err = provider.Send(context.Background(), msg, &core.ProviderSendOptions{HTTPClient: client})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_WithOptions(t *testing.T) {
	// Create mock server that checks request body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request contains the expected fields
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true, "result": {"message_id": 123}}`))
	}))
	defer ts.Close()

	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := telegram.NewTextMessage(
		"123456789",
		"Hello, World!",
		telegram.WithSilent(true),
		telegram.WithProtectContent(true),
		telegram.WithParseMode("HTML"),
	)

	tsURL := ts.URL
	parts := strings.SplitN(tsURL, "://", 2)
	scheme, host := parts[0], parts[1]
	client := &http.Client{
		Transport: rewriteRoundTripper{base: http.DefaultTransport, targetHost: host, targetScheme: scheme},
	}

	err = provider.Send(context.Background(), msg, &core.ProviderSendOptions{HTTPClient: client})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestProvider_Send_HTTPFailure(t *testing.T) {
	// Create mock server that returns error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"ok": false, "description": "server error"}`))
	}))
	defer ts.Close()

	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := telegram.NewTextMessage("123456789", "Hello, World!")

	tsURL := ts.URL
	parts := strings.SplitN(tsURL, "://", 2)
	scheme, host := parts[0], parts[1]
	client := &http.Client{
		Transport: rewriteRoundTripper{base: http.DefaultTransport, targetHost: host, targetScheme: scheme},
	}

	err = provider.Send(context.Background(), msg, &core.ProviderSendOptions{HTTPClient: client})
	if err == nil {
		t.Error("Expected error for HTTP failure, got nil")
	}
}

func TestProvider_Send_InvalidMessageType(t *testing.T) {
	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	msg := &invalidMessage{}

	err = provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for invalid message type, got nil")
	}
}

func TestProvider_Send_InvalidMessage(t *testing.T) {
	config := telegram.Config{
		Accounts: []core.Account{
			{
				Name:   "test",
				APIKey: "bot123:token",
			},
		},
	}

	provider, err := telegram.New(config)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Create an invalid message (empty text)
	msg := telegram.NewTextMessage("123456789", "")

	err = provider.Send(context.Background(), msg, nil)
	if err == nil {
		t.Error("Expected error for invalid message, got nil")
	}
}

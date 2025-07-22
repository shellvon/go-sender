package dryrun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/cmd/gosender/internal/mock"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/webhook"
)

// Handler handles dry-run operations
type Handler struct {
	validator     *config.Validator
	mockClient    *http.Client
	mockTransport *mock.MockTransport
}

// NewHandler creates a new dry-run handler
func NewHandler() *Handler {
	mockClient := mock.NewHTTPClient()
	mockTransport := mock.GetMockTransport(mockClient)

	return &Handler{
		validator:     config.NewValidator(),
		mockClient:    mockClient,
		mockTransport: mockTransport,
	}
}

// HandleDryRun processes a dry-run request
func (h *Handler) HandleDryRun(flags *cli.CLIFlags, config *cli.RootConfig) (*cli.DryRunResult, error) {
	// Validate configuration first
	validationErrors := h.validator.ValidateConfiguration(config)

	// Create dry run result
	result := &cli.DryRunResult{
		Provider:         flags.Provider,
		Account:          flags.Account,
		MessageType:      flags.MessageType,
		ValidationErrors: validationErrors,
	}

	// Set default provider if not specified
	if result.Provider == "" {
		result.Provider = string(core.ProviderTypeSMS)
	}

	// Set default message type based on provider
	if result.MessageType == "" {
		result.MessageType = h.getDefaultMessageType(result.Provider, flags)
	}

	// Find matching account
	if result.Account == "" {
		result.Account = h.findMatchingAccount(config, result.Provider)
	}

	// If validation passed, try to build and "send" the message with mock client
	if len(validationErrors) == 0 {
		if err := h.performMockSend(flags, config, result); err != nil {
			result.ValidationErrors = append(result.ValidationErrors, fmt.Sprintf("Mock send failed: %v", err))
		}
	}

	return result, nil
}

// getDefaultMessageType returns the default message type for a provider
func (h *Handler) getDefaultMessageType(provider string, flags *cli.CLIFlags) string {
	switch provider {
	case string(core.ProviderTypeSMS):
		return "text"
	case string(core.ProviderTypeEmail):
		if flags.HTML {
			return "html"
		}
		return "text"
	case string(core.ProviderTypeDingtalk):
		return "text"
	case string(core.ProviderTypeWebhook):
		return "json"
	case string(core.ProviderTypeTelegram):
		return "text"
	case string(core.ProviderTypeLark):
		return "text"
	case string(core.ProviderTypeWecombot):
		return "text"
	case string(core.ProviderTypeServerChan):
		return "text"
	default:
		return "text"
	}
}

// findMatchingAccount finds the first enabled account for the specified provider
func (h *Handler) findMatchingAccount(config *cli.RootConfig, provider string) string {
	for _, account := range config.Accounts {
		if providerType, ok := account["provider"].(string); ok && providerType == provider {
			if enabled, ok := account["enabled"].(bool); !ok || enabled {
				if name, ok := account["name"].(string); ok {
					return name
				}
				return "default"
			}
		}
	}
	return "not-found"
}

// performMockSend attempts to build and send a message using the mock HTTP client
func (h *Handler) performMockSend(flags *cli.CLIFlags, config *cli.RootConfig, result *cli.DryRunResult) error {
	// Clear previous captures
	h.mockTransport.Clear()

	// Create sender with mock HTTP client
	sender := gosender.NewSender()

	// Set the mock HTTP client as the default HTTP client for all providers
	sender.SetDefaultHTTPClient(h.mockClient)

	// Set up providers with mock HTTP client
	if err := h.setupProvidersWithMockClient(sender, config); err != nil {
		return fmt.Errorf("failed to setup providers: %w", err)
	}

	// Build message based on provider type
	msg, err := h.buildMessage(flags, result.Provider)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	result.ValidatedMessage = msg

	// Attempt to send with mock client
	_, err = sender.SendWithResult(context.Background(), msg)
	if err != nil {
		return fmt.Errorf("mock send failed: %w", err)
	}

	// Capture the HTTP request details
	if capturedReq := h.mockTransport.GetLastCapturedRequest(); capturedReq != nil {
		result.HTTPRequest = capturedReq
	}

	return nil
}

// setupProvidersWithMockClient sets up providers to use the mock HTTP client
func (h *Handler) setupProvidersWithMockClient(sender *gosender.Sender, conf *cli.RootConfig) error {
	// Use the new parser instead of mapstructure
	parser := config.NewAccountParser()
	smsAccounts, emailAccounts, dingtalkAccounts, webhookEndpoints, err := parser.ParseAccounts(conf)
	if err != nil {
		return fmt.Errorf("failed to parse accounts: %w", err)
	}

	// Create and register providers
	if len(smsAccounts) > 0 {
		smsCfg := &sms.Config{Items: smsAccounts}
		smsProvider, err := sms.New(smsCfg)
		if err != nil {
			return fmt.Errorf("create sms provider: %w", err)
		}
		sender.RegisterProvider(core.ProviderTypeSMS, smsProvider, nil)
	}

	if len(emailAccounts) > 0 {
		emailCfg := &email.Config{Items: emailAccounts}
		emailProvider, err := email.New(emailCfg)
		if err != nil {
			return fmt.Errorf("create email provider: %w", err)
		}
		sender.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)
	}

	if len(webhookEndpoints) > 0 {
		webhookCfg := &webhook.Config{Items: webhookEndpoints}
		webhookProvider, err := webhook.New(webhookCfg)
		if err != nil {
			return fmt.Errorf("create webhook provider: %w", err)
		}
		sender.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)
	}

	if len(dingtalkAccounts) > 0 {
		dtCfg := &dingtalk.Config{Items: dingtalkAccounts}
		dtProvider, err := dingtalk.New(dtCfg)
		if err != nil {
			return fmt.Errorf("create dingtalk provider: %w", err)
		}
		sender.RegisterProvider(core.ProviderTypeDingtalk, dtProvider, nil)
	}

	return nil
}

// buildMessage builds a message based on the provider type and flags
func (h *Handler) buildMessage(flags *cli.CLIFlags, provider string) (core.Message, error) {
	switch provider {
	case string(core.ProviderTypeSMS):
		if len(flags.To) == 0 {
			return nil, fmt.Errorf("SMS requires at least one recipient")
		}
		return sms.Aliyun().To(flags.To[0]).Content(flags.Content).Build(), nil

	case string(core.ProviderTypeEmail):
		if len(flags.To) == 0 {
			return nil, fmt.Errorf("Email requires at least one recipient")
		}
		builder := email.Email().To(flags.To[0]).Body(flags.Content)
		if flags.Subject != "" {
			builder = builder.Subject(flags.Subject)
		}
		return builder.Build(), nil

	case string(core.ProviderTypeDingtalk):
		return dingtalk.Text().Content(flags.Content).Build(), nil

	case string(core.ProviderTypeWebhook):
		payload := map[string]interface{}{
			"text":    flags.Content,
			"subject": flags.Subject,
		}
		payloadBytes, _ := json.Marshal(payload)
		return webhook.Webhook().Body(payloadBytes).Build(), nil

	default:
		return nil, fmt.Errorf("unsupported provider type for message building: %s", provider)
	}
}

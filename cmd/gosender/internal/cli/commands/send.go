package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/output"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/providers"
	"github.com/shellvon/go-sender/cmd/gosender/internal/mock"
	"github.com/shellvon/go-sender/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewSendCommand creates the send command
func NewSendCommand(configLoader *config.ConfigLoader) *cobra.Command {
	var flags cli.CLIFlags

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send notifications via configured providers",
		Long: `Send notifications through various providers including SMS, Email, DingTalk, 
Telegram, Lark, WeComBot, Webhook, and ServerChan. Supports dry-run mode for testing
configurations without actually sending messages.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Bind command flags to viper for automatic configuration binding
			return configLoader.BindFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration with viper integration
			config, configSource, err := configLoader.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Get flag values from viper (which merges CLI flags, config file, and env vars)
			flags.To = viper.GetStringSlice("to")
			flags.Content = viper.GetString("content")
			flags.Subject = viper.GetString("subject")
			flags.Provider = viper.GetString("provider")
			flags.Account = viper.GetString("account")
			flags.Strategy = viper.GetString("strategy")
			flags.TemplateID = viper.GetString("template-id")                  // Use hyphen as in CLI flag
			flags.TemplateParams = viper.GetStringMapString("template-params") // Use hyphen as in CLI flag
			flags.Files = viper.GetStringSlice("file")
			flags.MessageType = viper.GetString("message-type") // Use hyphen as in CLI flag
			flags.HTML = viper.GetBool("html")
			flags.DryRun = viper.GetBool("dry-run") // Use hyphen as in CLI flag
			flags.Priority = viper.GetInt("priority")
			flags.Timeout = viper.GetDuration("timeout")
			flags.Metadata = viper.GetStringMapString("metadata")

			// 输出配置文件信息 - 根据 log_level 决定是否显示
			logLevel := strings.ToUpper(viper.GetString("log-level"))
			if logLevel == "DEBUG" {
				fmt.Printf("Using configuration from: %s\n", configSource)
				fmt.Printf("Log level: %s\n", config.LogLevel)
			}

			// Basic validation
			// WeComBot doesn't need --to parameter as it sends to groups
			// ServerChan doesn't need --to parameter as it sends to groups
			if flags.Provider != "wecombot" && flags.Provider != "serverchan" && len(flags.To) == 0 {
				return fmt.Errorf("--to is required: specify at least one recipient")
			}

			if flags.Content == "" && flags.TemplateID == "" {
				return fmt.Errorf("either --content or --template-id is required")
			}

			// Handle dry-run mode early to avoid creating real providers
			if flags.DryRun {
				return handleDryRun(&flags, config)
			}

			// Handle real send
			return handleRealSend(&flags, config)
		},
	}

	// Add comprehensive flags for the send command
	cmd.Flags().StringSliceVarP(&flags.To, "to", "t", []string{}, "recipient(s) - comma separated for multiple")
	cmd.Flags().StringVarP(&flags.Content, "content", "m", "", "message content")
	cmd.Flags().StringVarP(&flags.Subject, "subject", "s", "", "message subject (for email and applicable providers)")
	cmd.Flags().StringVarP(&flags.Provider, "provider", "p", "", "specific provider to use (sms, email, dingtalk, webhook, etc.)")
	cmd.Flags().StringVarP(&flags.Account, "account", "a", "", "specific account name to use")
	cmd.Flags().StringVar(&flags.Strategy, "strategy", "", "provider selection strategy (round_robin, weighted, failover)")
	cmd.Flags().StringVar(&flags.TemplateID, "template-id", "", "template ID for SMS and applicable providers")
	cmd.Flags().StringToStringVar(&flags.TemplateParams, "template-params", map[string]string{}, "template parameters as key=value pairs")
	cmd.Flags().StringSliceVar(&flags.Files, "file", []string{}, "file(s) to attach or send")
	cmd.Flags().StringVar(&flags.MessageType, "message-type", "", "specific message type for the provider")
	cmd.Flags().BoolVar(&flags.HTML, "html", false, "send as HTML (for email)")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "validate and preview without sending")
	cmd.Flags().IntVar(&flags.Priority, "priority", 0, "message priority")
	cmd.Flags().DurationVar(&flags.Timeout, "timeout", 0, "request timeout")
	cmd.Flags().StringToStringVar(&flags.Metadata, "metadata", map[string]string{}, "additional metadata as key=value pairs")

	// Don't mark flags as required here since we want to support env vars and config files
	// Validation will be done in RunE after viper processes all sources

	return cmd
}

// handleDryRun processes dry-run requests using the same flow as real send
func handleDryRun(flags *cli.CLIFlags, conf *cli.RootConfig) error {
	// 1. Get provider registry with all registered builders
	registry := providers.DefaultProviderRegistry()

	// 2. Create sender with mock HTTP client for dry-run
	sender := gosender.NewSender()
	mockClient := mock.NewHTTPClient()
	sender.SetDefaultHTTPClient(mockClient)

	// 3. Build all providers from configuration (same as real send)
	if err := registry.BuildProviders(sender, conf); err != nil {
		return fmt.Errorf("failed to build providers: %w", err)
	}

	// 4. Determine provider type
	providerType, err := providers.GetProviderType(flags.Provider)
	if err != nil {
		return err
	}

	// 5. Build message using the appropriate message builder (same as real send)
	msg, err := registry.BuildMessage(providerType, flags)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// 6. Handle different provider types for dry-run
	var httpRequest *cli.HTTPRequestCapture

	if providerType == core.ProviderTypeEmail {
		// For email provider, simulate SMTP connection details instead of actual sending
		httpRequest = simulateEmailRequest(flags, conf)
	} else {
		// For HTTP-based providers, "send" message with mock client
		_, err = sender.SendWithResult(context.Background(), msg)
		if err != nil {
			return fmt.Errorf("dry-run send failed: %w", err)
		}

		// Capture the HTTP request details
		mockTransport := mock.GetMockTransport(mockClient)
		if capturedReq := mockTransport.GetLastCapturedRequest(); capturedReq != nil {
			httpRequest = capturedReq
		}
	}

	// 8. Create dry-run result
	dryRunResult := &cli.DryRunResult{
		Provider:         string(providerType),
		Account:          flags.Account,
		MessageType:      flags.MessageType,
		ValidatedMessage: msg,
		HTTPRequest:      httpRequest,
		ValidationErrors: []string{}, // No validation errors if we got this far
	}

	// Set default message type if not specified
	if dryRunResult.MessageType == "" {
		dryRunResult.MessageType = getDefaultMessageType(string(providerType), flags)
	}

	// Find matching account name if not specified
	if dryRunResult.Account == "" {
		dryRunResult.Account = findMatchingAccount(conf, string(providerType))
	}

	// 9. Format and output result
	outputFormat := viper.GetString("output")
	formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	result := &cli.FormattedResult{
		Success: true,
		DryRun:  dryRunResult,
	}

	output, err := formatter.Format(result)
	if err != nil {
		return fmt.Errorf("format output: %w", err)
	}

	fmt.Print(output)
	return nil
}

// getDefaultMessageType returns the default message type for a provider
func getDefaultMessageType(provider string, flags *cli.CLIFlags) string {
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
func findMatchingAccount(config *cli.RootConfig, provider string) string {
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

// simulateEmailRequest creates a simulated HTTP request for email (SMTP) operations
func simulateEmailRequest(flags *cli.CLIFlags, config *cli.RootConfig) *cli.HTTPRequestCapture {
	// Find the email account
	var emailAccount map[string]interface{}
	for _, account := range config.Accounts {
		if provider, ok := account["provider"].(string); ok && provider == "email" {
			if enabled, ok := account["enabled"].(bool); !ok || enabled {
				emailAccount = account
				break
			}
		}
	}

	if emailAccount == nil {
		return nil
	}

	// Extract SMTP details
	host, _ := emailAccount["host"].(string)
	port, _ := emailAccount["port"].(int)
	if port == 0 {
		port = 587 // Default SMTP port
	}

	// Create a simulated SMTP request representation
	return &cli.HTTPRequestCapture{
		Method: "SMTP",
		URL:    fmt.Sprintf("smtp://%s:%d", host, port),
		Headers: map[string]string{
			"Protocol":    "SMTP",
			"Auth-Method": "PLAIN",
			"From":        fmt.Sprintf("%v", emailAccount["from"]),
			"To":          fmt.Sprintf("%v", flags.To),
			"Subject":     flags.Subject,
		},
		Body: map[string]interface{}{
			"content":      flags.Content,
			"html":         flags.HTML,
			"attachments":  flags.Files,
			"message_type": "email",
		},
		Timestamp: time.Now(),
		Duration:  50 * time.Millisecond, // Simulated duration
	}
}

// handleRealSend processes real send requests using the provider registry
func handleRealSend(flags *cli.CLIFlags, conf *cli.RootConfig) error {
	// 1. Get provider registry with all registered builders
	registry := providers.DefaultProviderRegistry()

	// 2. Create sender and build all providers from configuration
	sender := gosender.NewSender()
	if err := registry.BuildProviders(sender, conf); err != nil {
		return fmt.Errorf("failed to build providers: %w", err)
	}

	// 3. Determine provider type
	providerType, err := providers.GetProviderType(flags.Provider)
	if err != nil {
		return err
	}

	// 4. Build message using the appropriate message builder
	msg, err := registry.BuildMessage(providerType, flags)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// 5. Send message
	result, err := sender.SendWithResult(context.Background(), msg)
	if err != nil {
		return err
	}

	// 6. Format and output result
	outputFormat := viper.GetString("output")
	formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	formattedResult := &cli.FormattedResult{
		Success:    true,
		Provider:   string(providerType),
		StatusCode: result.StatusCode,
		Metadata: map[string]interface{}{
			"response_body": string(result.Body),
		},
	}

	output, err := formatter.Format(formattedResult)
	if err != nil {
		return fmt.Errorf("format output: %w", err)
	}

	fmt.Print(output)
	return nil
}

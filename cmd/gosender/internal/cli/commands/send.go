package commands

import (
	"context"
	"fmt"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/dryrun"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/output"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/dingtalk"
	"github.com/shellvon/go-sender/providers/email"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/providers/webhook"
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
			config, err := configLoader.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Get flag values from viper (which merges CLI flags, config file, and env vars)
			flags.To = viper.GetStringSlice("to")
			flags.Content = viper.GetString("content")
			flags.Subject = viper.GetString("subject")
			flags.Provider = viper.GetString("provider")
			flags.Account = viper.GetString("account")
			flags.TemplateID = viper.GetString("template-id")                  // Use hyphen as in CLI flag
			flags.TemplateParams = viper.GetStringMapString("template-params") // Use hyphen as in CLI flag
			flags.Files = viper.GetStringSlice("file")
			flags.MessageType = viper.GetString("message-type") // Use hyphen as in CLI flag
			flags.HTML = viper.GetBool("html")
			flags.DryRun = viper.GetBool("dry-run") // Use hyphen as in CLI flag
			flags.Priority = viper.GetInt("priority")
			flags.Timeout = viper.GetDuration("timeout")

			// Optional: Debug output for verbose mode
			if viper.GetBool("verbose") {
				fmt.Printf("Configuration loaded successfully\n")
			}

			if len(flags.To) == 0 || flags.Content == "" {
				return fmt.Errorf("--to and --content are required")
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
	cmd.Flags().StringVar(&flags.TemplateID, "template-id", "", "template ID for SMS and applicable providers")
	cmd.Flags().StringToStringVar(&flags.TemplateParams, "template-params", map[string]string{}, "template parameters as key=value pairs")
	cmd.Flags().StringSliceVar(&flags.Files, "file", []string{}, "file(s) to attach or send")
	cmd.Flags().StringVar(&flags.MessageType, "message-type", "", "specific message type for the provider")
	cmd.Flags().BoolVar(&flags.HTML, "html", false, "send as HTML (for email)")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "validate and preview without sending")
	cmd.Flags().IntVar(&flags.Priority, "priority", 0, "message priority")
	cmd.Flags().DurationVar(&flags.Timeout, "timeout", 0, "request timeout")

	// Don't mark flags as required here since we want to support env vars and config files
	// Validation will be done in RunE after viper processes all sources

	return cmd
}

// handleDryRun processes dry-run requests
func handleDryRun(flags *cli.CLIFlags, config *cli.RootConfig) error {
	handler := dryrun.NewHandler()
	dryRunResult, err := handler.HandleDryRun(flags, config)
	if err != nil {
		return fmt.Errorf("dry-run failed: %w", err)
	}

	// Format and output result
	outputFormat := viper.GetString("output")
	formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	result := &cli.FormattedResult{
		Success: len(dryRunResult.ValidationErrors) == 0,
		DryRun:  dryRunResult,
	}

	if len(dryRunResult.ValidationErrors) > 0 {
		result.Error = "Validation errors found"
	}

	output, err := formatter.Format(result)
	if err != nil {
		return fmt.Errorf("format output: %w", err)
	}

	fmt.Print(output)
	return nil
}

// handleRealSend processes real send requests
func handleRealSend(flags *cli.CLIFlags, conf *cli.RootConfig) error {
	// 1. Parse accounts using the new parser instead of mapstructure
	parser := config.NewAccountParser()
	smsAccounts, emailAccounts, dingtalkAccounts, webhookEndpoints, err := parser.ParseAccounts(conf)
	if err != nil {
		return fmt.Errorf("failed to parse accounts: %w", err)
	}

	sender := gosender.NewSender()

	// 2. Create and register providers
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

	// 3. Build and send message based on provider type
	var msg core.Message

	if flags.Provider == "" || flags.Provider == string(core.ProviderTypeSMS) {
		// Default to SMS for backward compatibility
		msg = sms.Aliyun().To(flags.To[0]).Content(flags.Content).Build()
	} else {
		return fmt.Errorf("provider-specific message building not yet implemented for: %s", flags.Provider)
	}

	result, err := sender.SendWithResult(context.Background(), msg)
	if err != nil {
		return err
	}

	// Format and output result
	outputFormat := viper.GetString("output")
	formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	formattedResult := &cli.FormattedResult{
		Success:    true,
		Provider:   string(core.ProviderTypeSMS),
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

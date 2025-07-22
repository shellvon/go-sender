package commands

import (
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/output"
	"github.com/shellvon/go-sender/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListProvidersCommand creates the list-providers command
func NewListProvidersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-providers",
		Short: "List all available providers and their capabilities",
		Long: `List all available notification providers supported by go-sender,
including their message types, required configuration fields, and capabilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			providers := getAvailableProviders()

			outputFormat := viper.GetString("output")
			formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
			if err != nil {
				return fmt.Errorf("invalid output format: %w", err)
			}

			result := &cli.FormattedResult{
				Success:  true,
				Metadata: map[string]interface{}{"providers": providers},
			}

			output, err := formatter.Format(result)
			if err != nil {
				return fmt.Errorf("format output: %w", err)
			}

			fmt.Print(output)
			return nil
		},
	}

	return cmd
}

// getAvailableProviders returns information about all available providers
func getAvailableProviders() []cli.ProviderInfo {
	return []cli.ProviderInfo{
		{
			Type:              core.ProviderTypeSMS,
			Name:              "SMS",
			Description:       "Send SMS messages via various providers (Aliyun, Tencent, Huawei, etc.)",
			SupportedTypes:    []string{"text", "template", "voice"},
			RequiredFields:    []string{"api_key", "api_secret", "sign_name"},
			OptionalFields:    []string{"subtype", "name", "enabled", "weight"},
			SupportsFiles:     false,
			SupportsHTML:      false,
			SupportsTemplates: true,
		},
		{
			Type:              core.ProviderTypeEmail,
			Name:              "Email",
			Description:       "Send emails via SMTP",
			SupportedTypes:    []string{"text", "html"},
			RequiredFields:    []string{"host", "port", "api_key", "api_secret", "from"},
			OptionalFields:    []string{"name", "enabled", "weight"},
			SupportsFiles:     true,
			SupportsHTML:      true,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeDingtalk,
			Name:              "DingTalk",
			Description:       "Send messages to DingTalk groups via webhook",
			SupportedTypes:    []string{"text", "markdown", "link", "action_card", "feed_card"},
			RequiredFields:    []string{"api_key"},
			OptionalFields:    []string{"api_secret", "name", "enabled", "weight"},
			SupportsFiles:     false,
			SupportsHTML:      false,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeWebhook,
			Name:              "Webhook",
			Description:       "Send custom JSON payloads to webhook endpoints",
			SupportedTypes:    []string{"json"},
			RequiredFields:    []string{"url", "method"},
			OptionalFields:    []string{"headers", "name", "enabled", "weight"},
			SupportsFiles:     false,
			SupportsHTML:      false,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeTelegram,
			Name:              "Telegram",
			Description:       "Send messages via Telegram Bot API",
			SupportedTypes:    []string{"text", "photo", "document", "audio", "video", "voice", "location", "contact"},
			RequiredFields:    []string{"api_key"},
			OptionalFields:    []string{"name", "enabled", "weight"},
			SupportsFiles:     true,
			SupportsHTML:      true,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeLark,
			Name:              "Lark/Feishu",
			Description:       "Send messages to Lark/Feishu groups via webhook",
			SupportedTypes:    []string{"text", "post", "image", "interactive", "share_chat"},
			RequiredFields:    []string{"api_key"},
			OptionalFields:    []string{"api_secret", "name", "enabled", "weight"},
			SupportsFiles:     false,
			SupportsHTML:      false,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeWecombot,
			Name:              "WeCom Bot",
			Description:       "Send messages to WeCom groups via webhook",
			SupportedTypes:    []string{"text", "markdown", "image", "news", "file", "template_card"},
			RequiredFields:    []string{"api_key"},
			OptionalFields:    []string{"name", "enabled", "weight"},
			SupportsFiles:     true,
			SupportsHTML:      false,
			SupportsTemplates: true,
		},
		{
			Type:              core.ProviderTypeServerChan,
			Name:              "ServerChan",
			Description:       "Send notifications via ServerChan service",
			SupportedTypes:    []string{"text", "markdown"},
			RequiredFields:    []string{"api_key"},
			OptionalFields:    []string{"name", "enabled", "weight"},
			SupportsFiles:     false,
			SupportsHTML:      false,
			SupportsTemplates: false,
		},
		{
			Type:              core.ProviderTypeEmailAPI,
			Name:              "Email API",
			Description:       "Send emails via API services (Resend, EmailJS, etc.)",
			SupportedTypes:    []string{"text", "html", "template"},
			RequiredFields:    []string{"api_key", "subtype"},
			OptionalFields:    []string{"api_secret", "from", "service_id", "template_id", "name", "enabled", "weight"},
			SupportsFiles:     true,
			SupportsHTML:      true,
			SupportsTemplates: true,
		},
	}
}

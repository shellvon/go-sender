package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"gopkg.in/yaml.v3"
)

// ConsoleFormatter formats output for console display.
type ConsoleFormatter struct{}

func (f *ConsoleFormatter) Format(result *cli.FormattedResult) (string, error) {
	var output strings.Builder

	if result.Success {
		output.WriteString("✅ Success\n")
	} else {
		output.WriteString("❌ Failed\n")
	}

	if result.Provider != "" {
		output.WriteString(fmt.Sprintf("Provider: %s\n", result.Provider))
	}

	if result.Account != "" {
		output.WriteString(fmt.Sprintf("Account: %s\n", result.Account))
	}

	if result.MessageID != "" {
		output.WriteString(fmt.Sprintf("Message ID: %s\n", result.MessageID))
	}

	if result.Duration > 0 {
		output.WriteString(fmt.Sprintf("Duration: %s\n", result.Duration))
	}

	if result.Error != "" {
		output.WriteString(fmt.Sprintf("Error: %s\n", result.Error))
	}

	// Handle metadata for special cases
	if result.Metadata != nil {
		if providers, ok := result.Metadata["providers"].([]cli.ProviderInfo); ok {
			output.WriteString("\nAvailable Providers:\n")
			for _, provider := range providers {
				output.WriteString(fmt.Sprintf("\n📦 %s (%s)\n", provider.Name, provider.Type))
				output.WriteString(fmt.Sprintf("   Description: %s\n", provider.Description))
				output.WriteString(fmt.Sprintf("   Message Types: %s\n", strings.Join(provider.SupportedTypes, ", ")))
				output.WriteString(fmt.Sprintf("   Required Fields: %s\n", strings.Join(provider.RequiredFields, ", ")))
				if len(provider.OptionalFields) > 0 {
					output.WriteString(
						fmt.Sprintf("   Optional Fields: %s\n", strings.Join(provider.OptionalFields, ", ")),
					)
				}

				var capabilities []string
				if provider.SupportsFiles {
					capabilities = append(capabilities, "Files")
				}
				if provider.SupportsHTML {
					capabilities = append(capabilities, "HTML")
				}
				if provider.SupportsTemplates {
					capabilities = append(capabilities, "Templates")
				}
				if len(capabilities) > 0 {
					output.WriteString(fmt.Sprintf("   Capabilities: %s\n", strings.Join(capabilities, ", ")))
				}
			}
		}

		if validationErrors, ok := result.Metadata["validation_errors"].([]string); ok {
			if len(validationErrors) > 0 {
				output.WriteString("\nValidation Errors:\n")
				for _, err := range validationErrors {
					output.WriteString(fmt.Sprintf("  ❌ %s\n", err))
				}
			} else {
				output.WriteString("\n✅ Configuration is valid\n")
			}

			if accountsCount, ok := result.Metadata["accounts_count"].(int); ok {
				output.WriteString(fmt.Sprintf("Accounts configured: %d\n", accountsCount))
			}
		}
	}

	if result.DryRun != nil {
		output.WriteString("\n🔍 Dry Run Results:\n")
		output.WriteString(fmt.Sprintf("Provider: %s\n", result.DryRun.Provider))
		output.WriteString(fmt.Sprintf("Account: %s\n", result.DryRun.Account))
		output.WriteString(fmt.Sprintf("Message Type: %s\n", result.DryRun.MessageType))

		if len(result.DryRun.ValidationErrors) > 0 {
			output.WriteString("Validation Errors:\n")
			for _, err := range result.DryRun.ValidationErrors {
				output.WriteString(fmt.Sprintf("  ❌ %s\n", err))
			}
		}

		if result.DryRun.HTTPRequest != nil {
			output.WriteString("\nHTTP Request Captured:\n")
			output.WriteString(fmt.Sprintf("Method: %s\n", result.DryRun.HTTPRequest.Method))
			output.WriteString(fmt.Sprintf("URL: %s\n", result.DryRun.HTTPRequest.URL))
			output.WriteString(
				fmt.Sprintf("Timestamp: %s\n", result.DryRun.HTTPRequest.Timestamp.Format("2006-01-02 15:04:05")),
			)
			if result.DryRun.HTTPRequest.Duration > 0 {
				output.WriteString(fmt.Sprintf("Duration: %s\n", result.DryRun.HTTPRequest.Duration))
			}
			if len(result.DryRun.HTTPRequest.Headers) > 0 {
				output.WriteString("Headers:\n")
				for k, v := range result.DryRun.HTTPRequest.Headers {
					output.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
				}
			}
			if result.DryRun.HTTPRequest.Body != nil {
				output.WriteString("Body:\n")
				if bodyBytes, err := json.MarshalIndent(result.DryRun.HTTPRequest.Body, "  ", "  "); err == nil {
					output.WriteString(fmt.Sprintf("  %s\n", string(bodyBytes)))
				}
			}
		}
	}

	return output.String(), nil
}

// JSONFormatter formats output as JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(result *cli.FormattedResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal JSON: %w", err)
	}
	return string(data) + "\n", nil
}

// YAMLFormatter formats output as YAML.
type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(result *cli.FormattedResult) (string, error) {
	data, err := yaml.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("marshal YAML: %w", err)
	}
	return string(data), nil
}

// GetFormatter returns the appropriate formatter for the given format.
func GetFormatter(format cli.OutputFormat) (cli.OutputFormatter, error) {
	switch format {
	case cli.OutputConsole:
		return &ConsoleFormatter{}, nil
	case cli.OutputJSON:
		return &JSONFormatter{}, nil
	case cli.OutputYAML:
		return &YAMLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

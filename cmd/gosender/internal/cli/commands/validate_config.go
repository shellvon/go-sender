package commands

import (
	"errors"
	"fmt"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewValidateConfigCommand creates the validate-config command.
func NewValidateConfigCommand(configLoader *config.ConfigLoader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-config",
		Short: "Validate configuration file without sending messages",
		Long: `Validate the configuration file to ensure all provider configurations
are correct and can be loaded successfully. This helps identify configuration
issues before attempting to send messages.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			conf, configSource, err := configLoader.LoadConfig()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			// 输出配置文件信息
			logLevel := conf.LogLevel
			if logLevel == "DEBUG" {
				fmt.Printf("Using configuration from: %s\n", configSource)
				fmt.Printf("Log level: %s\n", logLevel)
			}

			// Validate configuration
			validator := config.NewValidator()
			validationErrors := validator.ValidateConfiguration(conf)

			// Format and output result
			outputFormat := viper.GetString("output")
			formatter, err := output.GetFormatter(cli.OutputFormat(outputFormat))
			if err != nil {
				return fmt.Errorf("invalid output format: %w", err)
			}

			result := &cli.FormattedResult{
				Success: len(validationErrors) == 0,
				Metadata: map[string]interface{}{
					"validation_errors": validationErrors,
					"accounts_count":    len(conf.Accounts),
				},
			}

			if len(validationErrors) > 0 {
				result.Error = fmt.Sprintf("Configuration validation failed with %d errors", len(validationErrors))
			}

			output, err := formatter.Format(result)
			if err != nil {
				return fmt.Errorf("format output: %w", err)
			}

			fmt.Print(output)

			if len(validationErrors) > 0 {
				return errors.New("configuration validation failed")
			}

			return nil
		},
	}

	return cmd
}

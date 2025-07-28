package main

import (
	"errors"
	"log"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/commands"
	"github.com/shellvon/go-sender/cmd/gosender/internal/cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func main() {
	root := &cobra.Command{
		Use:   "gosender",
		Short: "Enhanced go-sender CLI - Send notifications via multiple providers",
		Long: `Enhanced go-sender CLI provides a comprehensive interface to send notifications
through various providers including SMS, Email, DingTalk, Telegram, Lark, WeComBot, 
Webhook, and ServerChan. Supports configuration files, environment variables, 
dry-run mode, and structured output formats.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
	}

	// Global flags
	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default: auto-discover)")
	root.PersistentFlags().String("log-level", "INFO", "log level (DEBUG, INFO, WARN, ERROR)")
	root.PersistentFlags().String("output", "console", "output format (console, json, yaml)")

	// Bind global flags to viper
	viper.BindPFlag("config", root.PersistentFlags().Lookup("config"))
	viper.BindPFlag("log-level", root.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("output", root.PersistentFlags().Lookup("output"))

	// Create config loader
	configLoader := config.NewConfigLoader(cfgFile)

	// Add commands
	root.AddCommand(commands.NewSendCommand(configLoader))
	root.AddCommand(commands.NewListProvidersCommand())
	root.AddCommand(commands.NewValidateConfigCommand(configLoader))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

// initializeConfig initializes viper configuration with CLI flag binding.
func initializeConfig(cmd *cobra.Command) error {
	// Set up viper with GO_SENDER prefix for environment variables
	viper.SetEnvPrefix("GO_SENDER")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory and home directory
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("$HOME/.gosender")
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if no config file is found, we can work with env vars and flags
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	return nil
}

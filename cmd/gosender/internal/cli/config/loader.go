package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigLoader handles loading configuration from various sources with viper.
type ConfigLoader struct {
	configFile string
	viper      *viper.Viper
}

// NewConfigLoader creates a new configuration loader with viper setup.
func NewConfigLoader(configFile string) *ConfigLoader {
	v := viper.New()

	// Set up environment variable support with GO_SENDER prefix
	v.SetEnvPrefix("GO_SENDER")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	return &ConfigLoader{
		configFile: configFile,
		viper:      v,
	}
}

// BindFlags binds CLI flags to viper for automatic configuration binding.
func (c *ConfigLoader) BindFlags(cmd *cobra.Command) error {
	// Use global viper instance for consistency
	globalViper := viper.GetViper()

	// Bind all flags to global viper for automatic configuration merging
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		globalViper.BindPFlag(flag.Name, flag)
	})

	return nil
}

// LoadConfig loads configuration from file, supporting multiple formats with priority.
func (c *ConfigLoader) LoadConfig() (*cli.RootConfig, string, error) {
	// Use global viper instance for better integration with CLI flags
	globalViper := viper.GetViper()

	// Get config file from global viper (which gets it from command line flags)
	configFile := globalViper.GetString("config")
	if configFile == "" && c.configFile != "" {
		configFile = c.configFile
	}

	var err error
	var actualConfigFile string

	if configFile != "" {
		// Use specified config file with global viper
		err = c.loadSpecificConfigFileWithViper(configFile, globalViper)
		actualConfigFile = configFile
	} else {
		// Implement configuration file discovery with priority (YAML > JSON)
		actualConfigFile, err = c.discoverAndLoadConfigFileWithViper(globalViper)
	}

	if err != nil {
		// If no config file found, try to load from environment variables only
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			config, envErr := c.loadFromEnvWithViper(globalViper)
			return config, "environment variables", envErr
		}
		return nil, "", fmt.Errorf("failed to read config file: %w", err)
	}

	var config cli.RootConfig
	if err := globalViper.Unmarshal(&config); err != nil {
		return nil, actualConfigFile, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, actualConfigFile, nil
}

// loadSpecificConfigFileWithViper loads a specific configuration file using provided viper instance.
func (c *ConfigLoader) loadSpecificConfigFileWithViper(configFile string, v *viper.Viper) error {
	// Determine config type from file extension
	ext := strings.ToLower(filepath.Ext(configFile))
	switch ext {
	case ".json":
		v.SetConfigType("json")
	case ".yaml", ".yml":
		v.SetConfigType("yaml")
	case ".toml":
		v.SetConfigType("toml")
	case ".env":
		return c.loadFromEnvFileToViperInstance(configFile, v)
	default:
		// Try to auto-detect format
		v.SetConfigType("yaml") // Default to YAML
	}

	v.SetConfigFile(configFile)
	return v.ReadInConfig()
}

// discoverAndLoadConfigFileWithViper discovers and loads config files with priority (YAML > JSON).
func (c *ConfigLoader) discoverAndLoadConfigFileWithViper(v *viper.Viper) (string, error) {
	searchPaths := []string{
		".",
		"./config",
		"$HOME/.gosender",
	}

	// Priority order: YAML > JSON
	configNames := []struct {
		name string
		ext  string
		typ  string
	}{
		{"conf", ".yml", "yaml"},
		{"conf", ".yaml", "yaml"},
		{"config", ".yml", "yaml"},
		{"config", ".yaml", "yaml"},
		{"gosender", ".yml", "yaml"},
		{"gosender", ".yaml", "yaml"},
		{"conf", ".json", "json"},
		{"config", ".json", "json"},
		{"gosender", ".json", "json"},
	}

	for _, path := range searchPaths {
		for _, cfg := range configNames {
			configFile := filepath.Join(path, cfg.name+cfg.ext)

			// Expand home directory if needed
			if strings.HasPrefix(configFile, "$HOME") {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					continue
				}
				configFile = strings.Replace(configFile, "$HOME", homeDir, 1)
			}

			if _, err := os.Stat(configFile); err == nil {
				v.SetConfigType(cfg.typ)
				v.SetConfigFile(configFile)
				err := v.ReadInConfig()
				return configFile, err
			}
		}
	}

	return "", viper.ConfigFileNotFoundError{}
}

// loadFromEnvWithViper loads configuration from environment variables only using provided viper instance.
func (c *ConfigLoader) loadFromEnvWithViper(v *viper.Viper) (*cli.RootConfig, error) {
	config := &cli.RootConfig{
		LogLevel: v.GetString("log_level"),
		Accounts: []map[string]interface{}{},
	}

	// Load accounts from environment variables with enhanced parsing
	c.loadSMSAccountsFromEnvWithViper(config, v)
	c.loadEmailAccountsFromEnvWithViper(config, v)
	c.loadDingTalkAccountsFromEnvWithViper(config, v)
	c.loadTelegramAccountsFromEnvWithViper(config, v)
	c.loadWebhookAccountsFromEnvWithViper(config, v)

	return config, nil
}

// loadFromEnvFileToViperInstance loads configuration from a .env file into provided viper instance.
func (c *ConfigLoader) loadFromEnvFileToViperInstance(filename string, v *viper.Viper) error {
	// Read .env file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	// Parse .env file and set environment variables
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}

	return nil
}

// loadSMSAccountsFromEnvWithViper loads SMS accounts from environment variables using provided viper instance.
func (c *ConfigLoader) loadSMSAccountsFromEnvWithViper(config *cli.RootConfig, v *viper.Viper) {
	// Support multiple SMS providers through indexed environment variables
	providers := []string{
		"aliyun",
		"tencent",
		"huawei",
		"cl253",
		"juhe",
		"luosimao",
		"smsbao",
		"submail",
		"ucp",
		"volc",
		"yunpian",
		"yuntongxun",
	}

	for _, provider := range providers {
		keyPrefix := fmt.Sprintf("sms_%s", provider)
		if apiKey := v.GetString(fmt.Sprintf("%s_api_key", keyPrefix)); apiKey != "" {
			account := map[string]interface{}{
				"provider":   "sms",
				"subtype":    provider,
				"name":       v.GetString(fmt.Sprintf("%s_name", keyPrefix)),
				"api_key":    apiKey,
				"api_secret": v.GetString(fmt.Sprintf("%s_api_secret", keyPrefix)),
				"sign_name":  v.GetString(fmt.Sprintf("%s_sign_name", keyPrefix)),
				"region":     v.GetString(fmt.Sprintf("%s_region", keyPrefix)),
				"app_id":     v.GetString(fmt.Sprintf("%s_app_id", keyPrefix)),
				"enabled":    v.GetBool(fmt.Sprintf("%s_enabled", keyPrefix)),
				"weight":     v.GetInt(fmt.Sprintf("%s_weight", keyPrefix)),
			}

			// Set defaults
			if account["name"] == "" {
				account["name"] = fmt.Sprintf("%s-default", provider)
			}
			if account["enabled"] == false && !v.IsSet(fmt.Sprintf("%s_enabled", keyPrefix)) {
				account["enabled"] = true
			}
			if account["weight"] == 0 && !v.IsSet(fmt.Sprintf("%s_weight", keyPrefix)) {
				account["weight"] = 10
			}

			config.Accounts = append(config.Accounts, account)
		}
	}
}

// loadEmailAccountsFromEnvWithViper loads Email accounts from environment variables using provided viper instance.
func (c *ConfigLoader) loadEmailAccountsFromEnvWithViper(config *cli.RootConfig, v *viper.Viper) {
	if host := v.GetString("email_host"); host != "" {
		account := map[string]interface{}{
			"provider":   "email",
			"name":       v.GetString("email_name"),
			"host":       host,
			"port":       v.GetInt("email_port"),
			"api_key":    v.GetString("email_api_key"),
			"api_secret": v.GetString("email_api_secret"),
			"from":       v.GetString("email_from"),
			"enabled":    v.GetBool("email_enabled"),
			"weight":     v.GetInt("email_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "email-default"
		}
		if account["port"] == 0 {
			account["port"] = 587
		}
		if account["enabled"] == false && !v.IsSet("email_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !v.IsSet("email_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadDingTalkAccountsFromEnvWithViper loads DingTalk accounts from environment variables using provided viper instance.
func (c *ConfigLoader) loadDingTalkAccountsFromEnvWithViper(config *cli.RootConfig, v *viper.Viper) {
	if apiKey := v.GetString("dingtalk_api_key"); apiKey != "" {
		account := map[string]interface{}{
			"provider":   "dingtalk",
			"name":       v.GetString("dingtalk_name"),
			"api_key":    apiKey,
			"api_secret": v.GetString("dingtalk_api_secret"),
			"enabled":    v.GetBool("dingtalk_enabled"),
			"weight":     v.GetInt("dingtalk_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "dingtalk-default"
		}
		if account["enabled"] == false && !v.IsSet("dingtalk_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !v.IsSet("dingtalk_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadTelegramAccountsFromEnvWithViper loads Telegram accounts from environment variables using provided viper instance.
func (c *ConfigLoader) loadTelegramAccountsFromEnvWithViper(config *cli.RootConfig, v *viper.Viper) {
	if apiKey := v.GetString("telegram_api_key"); apiKey != "" {
		account := map[string]interface{}{
			"provider": "telegram",
			"name":     v.GetString("telegram_name"),
			"api_key":  apiKey,
			"enabled":  v.GetBool("telegram_enabled"),
			"weight":   v.GetInt("telegram_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "telegram-default"
		}
		if account["enabled"] == false && !v.IsSet("telegram_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !v.IsSet("telegram_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadWebhookAccountsFromEnvWithViper loads Webhook accounts from environment variables using provided viper instance.
func (c *ConfigLoader) loadWebhookAccountsFromEnvWithViper(config *cli.RootConfig, v *viper.Viper) {
	if url := v.GetString("webhook_url"); url != "" {
		account := map[string]interface{}{
			"provider": "webhook",
			"name":     v.GetString("webhook_name"),
			"url":      url,
			"method":   v.GetString("webhook_method"),
			"enabled":  v.GetBool("webhook_enabled"),
			"weight":   v.GetInt("webhook_weight"),
		}

		// Parse headers from environment variables
		if headersStr := v.GetString("webhook_headers"); headersStr != "" {
			headers := make(map[string]string)
			pairs := strings.Split(headersStr, ",")
			for _, pair := range pairs {
				kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
				if len(kv) == 2 {
					headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			account["headers"] = headers
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "webhook-default"
		}
		if account["method"] == "" {
			account["method"] = "POST"
		}
		if account["enabled"] == false && !v.IsSet("webhook_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !v.IsSet("webhook_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadSpecificConfigFile loads a specific configuration file.
func (c *ConfigLoader) loadSpecificConfigFile(configFile string) error {
	// Determine config type from file extension
	ext := strings.ToLower(filepath.Ext(configFile))
	switch ext {
	case ".json":
		c.viper.SetConfigType("json")
	case ".yaml", ".yml":
		c.viper.SetConfigType("yaml")
	case ".toml":
		c.viper.SetConfigType("toml")
	case ".env":
		return c.loadFromEnvFileToViper(configFile)
	default:
		// Try to auto-detect format
		c.viper.SetConfigType("yaml") // Default to YAML
	}

	c.viper.SetConfigFile(configFile)
	return c.viper.ReadInConfig()
}

// discoverAndLoadConfigFile discovers and loads config files with priority (YAML > JSON).
func (c *ConfigLoader) discoverAndLoadConfigFile() error {
	searchPaths := []string{
		".",
		"./config",
		"$HOME/.gosender",
	}

	// Priority order: YAML > JSON
	configNames := []struct {
		name string
		ext  string
		typ  string
	}{
		{"conf", ".yml", "yaml"},
		{"conf", ".yaml", "yaml"},
		{"config", ".yml", "yaml"},
		{"config", ".yaml", "yaml"},
		{"gosender", ".yml", "yaml"},
		{"gosender", ".yaml", "yaml"},
		{"conf", ".json", "json"},
		{"config", ".json", "json"},
		{"gosender", ".json", "json"},
	}

	for _, path := range searchPaths {
		for _, cfg := range configNames {
			configFile := filepath.Join(path, cfg.name+cfg.ext)

			// Expand home directory if needed
			if strings.HasPrefix(configFile, "$HOME") {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					continue
				}
				configFile = strings.Replace(configFile, "$HOME", homeDir, 1)
			}

			if _, err := os.Stat(configFile); err == nil {
				c.viper.SetConfigType(cfg.typ)
				c.viper.SetConfigFile(configFile)
				return c.viper.ReadInConfig()
			}
		}
	}

	return viper.ConfigFileNotFoundError{}
}

// loadFromEnv loads configuration from environment variables only.
func (c *ConfigLoader) loadFromEnv() (*cli.RootConfig, error) {
	config := &cli.RootConfig{
		LogLevel: c.viper.GetString("log_level"),
		Accounts: []map[string]interface{}{},
	}

	// Load accounts from environment variables with enhanced parsing
	c.loadSMSAccountsFromEnv(config)
	c.loadEmailAccountsFromEnv(config)
	c.loadDingTalkAccountsFromEnv(config)
	c.loadTelegramAccountsFromEnv(config)
	c.loadWebhookAccountsFromEnv(config)

	return config, nil
}

// loadSMSAccountsFromEnv loads SMS accounts from environment variables.
func (c *ConfigLoader) loadSMSAccountsFromEnv(config *cli.RootConfig) {
	// Support multiple SMS providers through indexed environment variables
	providers := []string{
		"aliyun",
		"tencent",
		"huawei",
		"cl253",
		"juhe",
		"luosimao",
		"smsbao",
		"submail",
		"ucp",
		"volc",
		"yunpian",
		"yuntongxun",
	}

	for _, provider := range providers {
		keyPrefix := fmt.Sprintf("sms_%s", provider)
		if apiKey := c.viper.GetString(fmt.Sprintf("%s_api_key", keyPrefix)); apiKey != "" {
			account := map[string]interface{}{
				"provider":   "sms",
				"subtype":    provider,
				"name":       c.viper.GetString(fmt.Sprintf("%s_name", keyPrefix)),
				"api_key":    apiKey,
				"api_secret": c.viper.GetString(fmt.Sprintf("%s_api_secret", keyPrefix)),
				"sign_name":  c.viper.GetString(fmt.Sprintf("%s_sign_name", keyPrefix)),
				"region":     c.viper.GetString(fmt.Sprintf("%s_region", keyPrefix)),
				"app_id":     c.viper.GetString(fmt.Sprintf("%s_app_id", keyPrefix)),
				"enabled":    c.viper.GetBool(fmt.Sprintf("%s_enabled", keyPrefix)),
				"weight":     c.viper.GetInt(fmt.Sprintf("%s_weight", keyPrefix)),
			}

			// Set defaults
			if account["name"] == "" {
				account["name"] = fmt.Sprintf("%s-default", provider)
			}
			if account["enabled"] == false && !c.viper.IsSet(fmt.Sprintf("%s_enabled", keyPrefix)) {
				account["enabled"] = true
			}
			if account["weight"] == 0 && !c.viper.IsSet(fmt.Sprintf("%s_weight", keyPrefix)) {
				account["weight"] = 10
			}

			config.Accounts = append(config.Accounts, account)
		}
	}
}

// loadEmailAccountsFromEnv loads Email accounts from environment variables.
func (c *ConfigLoader) loadEmailAccountsFromEnv(config *cli.RootConfig) {
	if host := c.viper.GetString("email_host"); host != "" {
		account := map[string]interface{}{
			"provider":   "email",
			"name":       c.viper.GetString("email_name"),
			"host":       host,
			"port":       c.viper.GetInt("email_port"),
			"api_key":    c.viper.GetString("email_api_key"),
			"api_secret": c.viper.GetString("email_api_secret"),
			"from":       c.viper.GetString("email_from"),
			"enabled":    c.viper.GetBool("email_enabled"),
			"weight":     c.viper.GetInt("email_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "email-default"
		}
		if account["port"] == 0 {
			account["port"] = 587
		}
		if account["enabled"] == false && !c.viper.IsSet("email_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !c.viper.IsSet("email_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadDingTalkAccountsFromEnv loads DingTalk accounts from environment variables.
func (c *ConfigLoader) loadDingTalkAccountsFromEnv(config *cli.RootConfig) {
	if apiKey := c.viper.GetString("dingtalk_api_key"); apiKey != "" {
		account := map[string]interface{}{
			"provider":   "dingtalk",
			"name":       c.viper.GetString("dingtalk_name"),
			"api_key":    apiKey,
			"api_secret": c.viper.GetString("dingtalk_api_secret"),
			"enabled":    c.viper.GetBool("dingtalk_enabled"),
			"weight":     c.viper.GetInt("dingtalk_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "dingtalk-default"
		}
		if account["enabled"] == false && !c.viper.IsSet("dingtalk_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !c.viper.IsSet("dingtalk_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadTelegramAccountsFromEnv loads Telegram accounts from environment variables.
func (c *ConfigLoader) loadTelegramAccountsFromEnv(config *cli.RootConfig) {
	if apiKey := c.viper.GetString("telegram_api_key"); apiKey != "" {
		account := map[string]interface{}{
			"provider": "telegram",
			"name":     c.viper.GetString("telegram_name"),
			"api_key":  apiKey,
			"enabled":  c.viper.GetBool("telegram_enabled"),
			"weight":   c.viper.GetInt("telegram_weight"),
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "telegram-default"
		}
		if account["enabled"] == false && !c.viper.IsSet("telegram_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !c.viper.IsSet("telegram_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadWebhookAccountsFromEnv loads Webhook accounts from environment variables.
func (c *ConfigLoader) loadWebhookAccountsFromEnv(config *cli.RootConfig) {
	if url := c.viper.GetString("webhook_url"); url != "" {
		account := map[string]interface{}{
			"provider": "webhook",
			"name":     c.viper.GetString("webhook_name"),
			"url":      url,
			"method":   c.viper.GetString("webhook_method"),
			"enabled":  c.viper.GetBool("webhook_enabled"),
			"weight":   c.viper.GetInt("webhook_weight"),
		}

		// Parse headers from environment variables
		if headersStr := c.viper.GetString("webhook_headers"); headersStr != "" {
			headers := make(map[string]string)
			pairs := strings.Split(headersStr, ",")
			for _, pair := range pairs {
				kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
				if len(kv) == 2 {
					headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			account["headers"] = headers
		}

		// Set defaults
		if account["name"] == "" {
			account["name"] = "webhook-default"
		}
		if account["method"] == "" {
			account["method"] = "POST"
		}
		if account["enabled"] == false && !c.viper.IsSet("webhook_enabled") {
			account["enabled"] = true
		}
		if account["weight"] == 0 && !c.viper.IsSet("webhook_weight") {
			account["weight"] = 10
		}

		config.Accounts = append(config.Accounts, account)
	}
}

// loadFromEnvFileToViper loads configuration from a .env file into viper.
func (c *ConfigLoader) loadFromEnvFileToViper(filename string) error {
	// Read .env file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	// Parse .env file and set environment variables
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}

	return nil
}

// loadFromEnvFile loads configuration from a .env file.
func (c *ConfigLoader) loadFromEnvFile(filename string) (*cli.RootConfig, error) {
	// Read .env file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	// Parse .env file and set environment variables
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}

	// Now load from environment variables
	return c.loadFromEnv()
}

// NeedsConfig determines if a command requires configuration.
func NeedsConfig(commandName string) bool {
	configRequiredCommands := []string{
		"send",
		"validate-config",
	}

	for _, cmd := range configRequiredCommands {
		if cmd == commandName {
			return true
		}
	}
	return false
}

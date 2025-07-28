package cli

import (
	"time"

	"github.com/shellvon/go-sender/core"
)

// CLIFlags defines all command-line flags for the enhanced CLI
type CLIFlags struct {
	// Configuration
	ConfigFile string
	DryRun     bool
	LogLevel   string
	Output     string // json, yaml, console

	// Provider Selection
	Provider    string
	SubProvider string // sub-provider type for SMS, EmailAPI, etc.
	Account     string
	Strategy    string // provider selection strategy

	// Message Content
	To             []string
	Content        string
	Subject        string
	TemplateID     string
	TemplateParams map[string]string
	Files          []string

	// Message Type Specific
	MessageType string // for providers supporting multiple types
	HTML        bool   // for email

	// Advanced Options
	Priority int
	Timeout  time.Duration
	Metadata map[string]string // additional metadata
}

// OutputFormat represents the output format type
type OutputFormat string

const (
	OutputConsole OutputFormat = "console"
	OutputJSON    OutputFormat = "json"
	OutputYAML    OutputFormat = "yaml"
)

// FormattedResult represents the structured result output
type FormattedResult struct {
	Success    bool                   `json:"success" yaml:"success"`
	Provider   string                 `json:"provider" yaml:"provider"`
	Account    string                 `json:"account" yaml:"account"`
	MessageID  string                 `json:"message_id" yaml:"message_id"`
	RequestID  string                 `json:"request_id,omitempty" yaml:"request_id,omitempty"`
	Duration   time.Duration          `json:"duration" yaml:"duration"`
	StatusCode int                    `json:"status_code,omitempty" yaml:"status_code,omitempty"`
	Error      string                 `json:"error,omitempty" yaml:"error,omitempty"`
	DryRun     *DryRunResult          `json:"dry_run,omitempty" yaml:"dry_run,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// DryRunResult represents the result of a dry-run operation
type DryRunResult struct {
	Provider         string              `json:"provider"`
	Account          string              `json:"account"`
	MessageType      string              `json:"message_type"`
	ValidatedMessage core.Message        `json:"validated_message"`
	HTTPRequest      *HTTPRequestCapture `json:"http_request,omitempty"`
	ValidationErrors []string            `json:"validation_errors,omitempty"`
	EstimatedCost    *CostEstimate       `json:"estimated_cost,omitempty"`
}

// HTTPRequestCapture captures the actual HTTP request details
type HTTPRequestCapture struct {
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      interface{}       `json:"body"`
	Timestamp time.Time         `json:"timestamp"`
	Duration  time.Duration     `json:"duration,omitempty"`
}

// CostEstimate provides cost estimation for the message
type CostEstimate struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
}

// ProviderMessageBuilder defines the interface for building provider-specific messages
type ProviderMessageBuilder interface {
	BuildMessage(flags *CLIFlags, config *RootConfig) (core.Message, error)
	ValidateFlags(flags *CLIFlags) error
	GetSupportedMessageTypes() []string
	GetRequiredFlags() []string
	GetOptionalFlags() []string
}

// OutputFormatter handles different output formats
type OutputFormatter interface {
	Format(result *FormattedResult) (string, error)
}

// ProviderInfo contains information about a provider's capabilities
type ProviderInfo struct {
	Type              core.ProviderType `json:"type"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	SupportedTypes    []string          `json:"supported_message_types"`
	RequiredFields    []string          `json:"required_config_fields"`
	OptionalFields    []string          `json:"optional_config_fields"`
	SupportsFiles     bool              `json:"supports_files"`
	SupportsHTML      bool              `json:"supports_html"`
	SupportsTemplates bool              `json:"supports_templates"`
}

// RootConfig defines the top-level configuration structure
type RootConfig struct {
	LogLevel string                   `mapstructure:"log_level"`
	Accounts []map[string]interface{} `mapstructure:"accounts"`
}

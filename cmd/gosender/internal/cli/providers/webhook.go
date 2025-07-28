package providers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shellvon/go-sender/cmd/gosender/internal/cli"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/webhook"
)

// createWebhookProvider creates a Webhook Provider from endpoints list
func createWebhookProvider(endpoints []*webhook.Endpoint) (core.Provider, error) {
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no valid webhook endpoints found")
	}

	cfg := &webhook.Config{Items: endpoints}
	return webhook.New(cfg)
}

// buildWebhookMessage builds a webhook message from CLI flags
func buildWebhookMessage(flags *cli.CLIFlags) (core.Message, error) {
	builder := webhook.Webhook()

	// Set body content
	// First check if content is JSON
	if flags.Content != "" {
		var jsonObj interface{}
		err := json.Unmarshal([]byte(flags.Content), &jsonObj)
		if err == nil {
			// Content is valid JSON, use it directly
			builder.Body([]byte(flags.Content))
		} else {
			// Content is not JSON, treat as plain text
			builder.Body([]byte(flags.Content))
			// Set default content-type to text/plain
			builder.Header("Content-Type", "text/plain")
		}
	}
	
	// Support advanced features through metadata for flexibility
	// These are optional and advanced users can use them if needed
	if flags.Metadata != nil {
		// Support for path parameters via metadata (for URL templates with {param})
		if pathParamsStr, ok := flags.Metadata["path_params"]; ok && pathParamsStr != "" {
			pathParams := parseKeyValuePairs(pathParamsStr, "=", ",")
			if len(pathParams) > 0 {
				builder.PathParams(pathParams)
			}
		}
		
		// Support for query parameters via metadata
		if queryParamsStr, ok := flags.Metadata["query_params"]; ok && queryParamsStr != "" {
			queryParams := parseKeyValuePairs(queryParamsStr, "=", "&")
			if len(queryParams) > 0 {
				builder.Queries(queryParams)
			}
		}
		
		// Allow method override via metadata
		if method, ok := flags.Metadata["method"]; ok && method != "" {
			builder.Method(method)
		}
		
		// Allow adding custom headers via metadata
		if headersStr, ok := flags.Metadata["headers"]; ok && headersStr != "" {
			headers := parseKeyValuePairs(headersStr, ":", ",")
			if len(headers) > 0 {
				builder.Headers(headers)
			}
		}
	}

	return builder.Build(), nil
}

// parseKeyValuePairs 解析字符串中的键值对
// separator: 键值对分隔符 (如 "=" 或 ":")
// delimiter: 多个键值对之间的分隔符 (如 "," 或 "&")
func parseKeyValuePairs(s, separator, delimiter string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(s, delimiter)
	for _, pair := range pairs {
		kv := strings.SplitN(pair, separator, 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}

// validateWebhookFlags validates CLI flags for webhook messages
func validateWebhookFlags(flags *cli.CLIFlags) error {
	if flags.Content == "" {
		return fmt.Errorf("webhook requires content")
	}
	return nil
}

// NewWebhookBuilder creates a new Webhook GenericBuilder
func NewWebhookBuilder() *GenericBuilder[*webhook.Endpoint, core.Message] {
	return NewGenericBuilder(
		core.ProviderTypeWebhook,
		createWebhookProvider,
		buildWebhookMessage,
		validateWebhookFlags,
	)
} 
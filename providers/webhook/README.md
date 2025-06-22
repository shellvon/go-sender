# Webhook Provider

This provider supports sending messages via HTTP webhooks to any endpoint that accepts HTTP requests.

## Features

- **Universal HTTP Support**: Send messages to any HTTP endpoint
- **Multiple Methods**: Support for GET, POST, PUT, PATCH, DELETE methods
- **Custom Headers**: Add custom headers for authentication and content type
- **Flexible Body Format**: Support for JSON, form data, and raw text
- **Timeout Control**: Configurable request timeout
- **Retry Support**: Built-in retry mechanism with exponential backoff
- **Multiple Endpoints**: Support multiple webhook endpoints with load balancing

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

// Create webhook configuration
config := webhook.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "primary-webhook",
            Key:      "https://api.example.com/webhook",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup-webhook",
            Key:      "https://backup.example.com/webhook",
            Weight:   50,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := webhook.New(config)
if err != nil {
    log.Fatalf("Failed to create webhook provider: %v", err)
}
```

## Message Types

### 1. JSON Message

```go
// Send JSON data
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer your-token",
    }),
    webhook.WithJSONBody(map[string]interface{}{
        "event": "user.created",
        "data": map[string]interface{}{
            "user_id": "12345",
            "email": "user@example.com",
            "timestamp": time.Now().Unix(),
        },
    }),
)
```

### 2. Form Data Message

```go
// Send form data
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    }),
    webhook.WithFormData(map[string]string{
        "action": "notify",
        "message": "Hello from webhook",
        "priority": "high",
    }),
)
```

### 3. Raw Text Message

```go
// Send raw text
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "text/plain",
    }),
    webhook.WithRawBody("Simple text message"),
)
```

### 4. GET Request with Query Parameters

```go
// Send GET request with query parameters
msg := webhook.NewMessage(
    webhook.WithMethod("GET"),
    webhook.WithQueryParams(map[string]string{
        "action": "ping",
        "timestamp": fmt.Sprintf("%d", time.Now().Unix()),
    }),
)
```

## Advanced Configuration

### Custom Timeout and Retry

```go
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithTimeout(30*time.Second),
    webhook.WithRetryAttempts(3),
    webhook.WithRetryDelay(2*time.Second),
    webhook.WithJSONBody(map[string]interface{}{
        "message": "Important notification",
    }),
)
```

### Complex Headers and Authentication

```go
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer your-api-token",
        "X-Custom-Header": "custom-value",
        "User-Agent": "Go-Sender/1.0",
    }),
    webhook.WithJSONBody(map[string]interface{}{
        "event": "system.alert",
        "level": "critical",
        "message": "System is down",
    }),
)
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

// Create sender
s := gosender.NewSender(nil)

// Register webhook provider
webhookProvider, err := webhook.New(config)
if err != nil {
    log.Fatalf("Failed to create webhook provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

// Send webhook message
ctx := context.Background()
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithJSONBody(map[string]interface{}{
        "message": "Hello from Go-Sender",
        "timestamp": time.Now().Unix(),
    }),
)
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send webhook: %v", err)
}
```

## Message Options

### HTTP Method Options

- `WithMethod(method string)`: Set HTTP method (GET, POST, PUT, PATCH, DELETE)
- `WithTimeout(timeout time.Duration)`: Set request timeout
- `WithRetryAttempts(attempts int)`: Set retry attempts
- `WithRetryDelay(delay time.Duration)`: Set retry delay

### Header and Authentication Options

- `WithHeaders(headers map[string]string)`: Set custom headers
- `WithBasicAuth(username, password string)`: Set basic authentication
- `WithBearerToken(token string)`: Set bearer token authentication

### Body Options

- `WithJSONBody(data interface{})`: Set JSON body
- `WithFormData(data map[string]string)`: Set form data body
- `WithRawBody(body string)`: Set raw text body
- `WithQueryParams(params map[string]string)`: Set query parameters (for GET requests)

### Advanced Options

- `WithCustomClient(client *http.Client)`: Use custom HTTP client
- `WithFollowRedirects(follow bool)`: Control redirect following
- `WithSkipSSLVerification(skip bool)`: Skip SSL certificate verification

## Configuration Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of webhook endpoint configurations

### Account (core.Account)

- `Name`: Account name for identification
- `Key`: Webhook URL (endpoint)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled
- `Webhook`: Optional webhook URL (alternative to Key)

### Message

- `Method`: HTTP method (default: POST)
- `URL`: Target URL (from account)
- `Headers`: HTTP headers
- `Body`: Request body
- `Timeout`: Request timeout
- `RetryAttempts`: Number of retry attempts
- `RetryDelay`: Delay between retries

## Error Handling

The provider handles:

- Network timeouts and connection errors
- HTTP error status codes (4xx, 5xx)
- Automatic retries with exponential backoff
- Provider selection based on strategy
- Fallback to alternative endpoints on failure

## Rate Limits and Security

- **Rate Limiting**: Respect target endpoint rate limits
- **Authentication**: Support for various authentication methods
- **SSL/TLS**: Full SSL/TLS support with certificate verification
- **Headers**: Custom headers for API keys, tokens, etc.

## Best Practices

### 1. Use Appropriate HTTP Methods

```go
// For notifications
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithJSONBody(data),
)

// For status checks
msg := webhook.NewMessage(
    webhook.WithMethod("GET"),
    webhook.WithQueryParams(params),
)
```

### 2. Handle Authentication Properly

```go
// Use headers for API keys
msg := webhook.NewMessage(
    webhook.WithHeaders(map[string]string{
        "X-API-Key": "your-api-key",
    }),
)

// Or use bearer tokens
msg := webhook.NewMessage(
    webhook.WithBearerToken("your-bearer-token"),
)
```

### 3. Set Appropriate Timeouts

```go
msg := webhook.NewMessage(
    webhook.WithTimeout(10*time.Second), // Short timeout for critical notifications
    webhook.WithRetryAttempts(3),
)
```

### 4. Use Multiple Endpoints for Reliability

```go
config := webhook.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name: "primary",
            Key:  "https://primary.example.com/webhook",
            Weight: 100,
        },
        {
            Name: "backup",
            Key:  "https://backup.example.com/webhook",
            Weight: 50,
        },
    },
}
```

## Common Use Cases

### 1. Slack Webhook Integration

```go
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/json",
    }),
    webhook.WithJSONBody(map[string]interface{}{
        "text": "Hello from Go-Sender!",
        "channel": "#general",
        "username": "Go-Sender Bot",
    }),
)
```

### 2. Discord Webhook Integration

```go
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/json",
    }),
    webhook.WithJSONBody(map[string]interface{}{
        "content": "Hello from Go-Sender!",
        "embeds": []map[string]interface{}{
            {
                "title": "Notification",
                "description": "This is a test message",
                "color": 0x00ff00,
            },
        },
    }),
)
```

### 3. Custom API Integration

```go
msg := webhook.NewMessage(
    webhook.WithMethod("POST"),
    webhook.WithHeaders(map[string]string{
        "Content-Type": "application/json",
        "Authorization": "Bearer your-token",
        "X-Event-Type": "user.created",
    }),
    webhook.WithJSONBody(map[string]interface{}{
        "user_id": "12345",
        "email": "user@example.com",
        "created_at": time.Now().Format(time.RFC3339),
    }),
)
```

## API Reference

### Constructor Functions

- `New(config Config) (*Provider, error)`: Create new webhook provider
- `NewMessage(options ...MessageOption) Message`: Create new webhook message

### Message Options

- `WithMethod(method string)`: Set HTTP method
- `WithHeaders(headers map[string]string)`: Set headers
- `WithJSONBody(data interface{})`: Set JSON body
- `WithFormData(data map[string]string)`: Set form data
- `WithRawBody(body string)`: Set raw body
- `WithQueryParams(params map[string]string)`: Set query parameters
- `WithTimeout(timeout time.Duration)`: Set timeout
- `WithRetryAttempts(attempts int)`: Set retry attempts
- `WithRetryDelay(delay time.Duration)`: Set retry delay
- `WithBasicAuth(username, password string)`: Set basic auth
- `WithBearerToken(token string)`: Set bearer token
- `WithCustomClient(client *http.Client)`: Set custom client
- `WithFollowRedirects(follow bool)`: Control redirects
- `WithSkipSSLVerification(skip bool)`: Skip SSL verification

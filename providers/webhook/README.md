# Webhook Provider

This provider supports sending messages via HTTP webhooks to any endpoint that accepts HTTP requests.

## Features

- **Universal HTTP Support**: Send messages to any HTTP endpoint
- **Multiple Methods**: Support for GET, POST, PUT, PATCH, DELETE methods (configured in endpoint)
- **Custom Headers**: Add custom headers for authentication and content type
- **Flexible Body Format**: Support for any raw body content
- **Response Validation**: Configurable response validation for different webhook formats
- **Multiple Response Types**: Support for JSON, text, and XML response validation
- **Custom Status Codes**: Configure custom success status codes
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
    Endpoints: []webhook.Endpoint{
        {
            Name:     "primary-webhook",
            URL:      "https://api.example.com/webhook",
            Method:   "POST",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup-webhook",
            URL:      "https://backup.example.com/webhook",
            Method:   "POST",
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

## Response Validation

The webhook provider supports configurable response validation to handle different webhook response formats.

### 1. Simple Status Code Validation (Default)

```go
// Default behavior - only checks HTTP status codes (2xx = success)
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "simple-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
        },
    },
}
```

### 2. JSON Response Validation

```go
// Validate JSON responses with success/error fields
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "json-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "json",
                SuccessField:     "success",    // Field name indicating success
                SuccessValue:     "true",       // Expected value for success
                ErrorField:       "error",      // Field name containing error message
                MessageField:     "message",    // Field name containing response message
            },
        },
    },
}
```

**Example JSON responses:**

```json
// Success response
{
    "success": "true",
    "message": "Message sent successfully",
    "id": "12345"
}

// Error response
{
    "success": "false",
    "error": "Invalid API key",
    "code": 401
}
```

### 3. Custom Status Codes

```go
// Accept only specific status codes as success
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "custom-status-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                SuccessStatusCodes: []int{200, 201, 202}, // Only these codes = success
            },
        },
    },
}
```

### 4. Text Response Validation

```go
// Validate text responses using regex patterns
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "text-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "text",
                SuccessPattern:   "^OK$",           // Regex for success response
                ErrorPattern:     "^ERROR:",        // Regex for error response
            },
        },
    },
}
```

**Example text responses:**

```
// Success response
OK

// Error response
ERROR: Invalid request
```

### 5. No Response Validation

```go
// Skip response body validation (only check status code)
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "no-validation-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: false, // or omit ResponseConfig entirely
            },
        },
    },
}
```

## Message Types

### 1. JSON Message

```go
// Send JSON data
jsonData := map[string]interface{}{
    "event": "user.created",
    "data": map[string]interface{}{
        "user_id": "12345",
        "email": "user@example.com",
        "timestamp": time.Now().Unix(),
    },
}
body, _ := json.Marshal(jsonData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token",
}))
```

### 2. Form Data Message

```go
// Send form data
formData := "action=notify&message=Hello from webhook&priority=high"
body := []byte(formData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/x-www-form-urlencoded",
}))
```

### 3. Raw Text Message

```go
// Send raw text
body := []byte("Simple text message")

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "text/plain",
}))
```

### 4. GET Request with Query Parameters

```go
// For GET requests, use endpoint QueryParams configuration
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "get-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "GET",
            QueryParams: map[string]string{
                "action": "ping",
                "timestamp": fmt.Sprintf("%d", time.Now().Unix()),
            },
        },
    },
}

// Empty body for GET requests
msg := webhook.NewMessage([]byte{})
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
jsonData := map[string]interface{}{
    "message": "Hello from Go-Sender",
    "timestamp": time.Now().Unix(),
}
body, _ := json.Marshal(jsonData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}))

err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send webhook: %v", err)
}
```

## Message Options

### Header Options

- `WithHeaders(headers map[string]string)`: Set custom headers

## Configuration Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Endpoints`: Array of webhook endpoint configurations

### Endpoint

- `Name`: Endpoint name for identification
- `URL`: Webhook URL (endpoint)
- `Method`: HTTP method (default: POST)
- `Headers`: Fixed request headers
- `QueryParams`: Fixed query parameters
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this endpoint is disabled
- `ResponseConfig`: Response validation configuration

### ResponseConfig

- `SuccessStatusCodes`: Custom success status codes (default: 2xx range)
- `ValidateResponse`: Whether to validate response body (default: false)
- `ResponseType`: Response type for validation ("json", "text", "xml", "none")
- `SuccessField`: JSON field name indicating success
- `SuccessValue`: Expected value for success field
- `ErrorField`: JSON field name containing error message
- `MessageField`: JSON field name containing response message
- `SuccessPattern`: Regex pattern for success text response
- `ErrorPattern`: Regex pattern for error text response

### Message

- `Body`: Request body (raw bytes)
- `Headers`: HTTP headers

## Error Handling

The provider handles:

- Network timeouts and connection errors
- HTTP error status codes (4xx, 5xx)
- Custom response validation failures
- Provider selection based on strategy
- Fallback to alternative endpoints on failure

## Best Practices

### 1. Configure Response Validation

```go
// For APIs that return structured responses
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "api-webhook",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "json",
                SuccessField:     "status",
                SuccessValue:     "ok",
                ErrorField:       "error",
            },
        },
    },
}
```

### 2. Use Appropriate HTTP Methods

```go
// Configure method in endpoint
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "notifications",
            URL:    "https://api.example.com/webhook",
            Method: "POST", // For notifications
        },
        {
            Name:   "status-check",
            URL:    "https://api.example.com/status",
            Method: "GET", // For status checks
        },
    },
}
```

### 3. Handle Authentication Properly

```go
// Use headers for API keys
msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "X-API-Key": "your-api-key",
    "Content-Type": "application/json",
}))

// Or use bearer tokens
msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Authorization": "Bearer your-bearer-token",
    "Content-Type": "application/json",
}))
```

### 4. Use Multiple Endpoints for Reliability

```go
config := webhook.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Endpoints: []webhook.Endpoint{
        {
            Name: "primary",
            URL:  "https://primary.example.com/webhook",
            Weight: 100,
        },
        {
            Name: "backup",
            URL:  "https://backup.example.com/webhook",
            Weight: 50,
        },
    },
}
```

## Common Use Cases

### 1. Slack Webhook Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "slack-webhook",
            URL:    "https://hooks.slack.com/services/YOUR/WEBHOOK/URL",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "text",
                SuccessPattern:   "^ok$",
            },
        },
    },
}

slackData := map[string]interface{}{
    "text": "Hello from Go-Sender!",
    "channel": "#general",
    "username": "Go-Sender Bot",
}
body, _ := json.Marshal(slackData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}))
```

### 2. Discord Webhook Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "discord-webhook",
            URL:    "https://discord.com/api/webhooks/YOUR/WEBHOOK/URL",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "json",
                SuccessField:     "id", // Discord returns message ID on success
            },
        },
    },
}

discordData := map[string]interface{}{
    "content": "Hello from Go-Sender!",
    "embeds": []map[string]interface{}{
        {
            "title": "Notification",
            "description": "This is a test message",
            "color": 0x00ff00,
        },
    },
}
body, _ := json.Marshal(discordData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}))
```

### 3. Custom API Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "custom-api",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
            ResponseConfig: &webhook.ResponseConfig{
                ValidateResponse: true,
                ResponseType:     "json",
                SuccessField:     "success",
                SuccessValue:     "true",
                ErrorField:       "error",
                MessageField:     "message",
            },
        },
    },
}

apiData := map[string]interface{}{
    "user_id": "12345",
    "email": "user@example.com",
    "created_at": time.Now().Format(time.RFC3339),
}
body, _ := json.Marshal(apiData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token",
    "X-Event-Type": "user.created",
}))
```

## API Reference

### Constructor Functions

- `New(config Config) (*Provider, error)`: Create new webhook provider
- `NewMessage(body []byte, opts ...MessageOption) *Message`: Create new webhook message

### Message Options

- `WithHeaders(headers map[string]string)`: Set custom headers

## Supported Push Services

The webhook provider can be used to send notifications to many popular push services, including but not limited to:

- **Pushover** ([API](https://pushover.net/api))
- **SimplePush** ([API](https://simplepush.io/api))
- **Bark** ([API](https://github.com/Finb/Bark))
- **PushDeer** ([API](https://github.com/easychen/pushdeer))
- ...and any service that accepts HTTP POST/GET requests

### Example Configurations

#### Pushover

```json
{
  "endpoints": [
    {
      "name": "pushover",
      "url": "https://api.pushover.net/1/messages.json",
      "method": "POST",
      "headers": { "Content-Type": "application/x-www-form-urlencoded" },
      "body": "token=YOUR_APP_TOKEN&user=USER_KEY&message=Hello+from+go-sender"
    }
  ]
}
```

#### SimplePush

```json
{
  "endpoints": [
    {
      "name": "simplepush",
      "url": "https://simplepu.sh",
      "method": "POST",
      "headers": { "Content-Type": "application/json" },
      "body": "{\"key\":\"YOUR_KEY\",\"msg\":\"Hello from go-sender!\"}"
    }
  ]
}
```

#### Bark

```json
{
  "endpoints": [
    {
      "name": "bark",
      "url": "https://api.day.app/YOUR_DEVICE_KEY/Hello%20from%20go-sender",
      "method": "GET"
    }
  ]
}
```

#### PushDeer

```json
{
  "endpoints": [
    {
      "name": "pushdeer",
      "url": "https://api2.pushdeer.com/message/push?pushkey=YOUR_KEY&text=Hello+from+go-sender",
      "method": "GET"
    }
  ]
}
```

> You can use the webhook provider to integrate with any service that supports HTTP requests. Just configure the URL, method, headers, and body as needed.

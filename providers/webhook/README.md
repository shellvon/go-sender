# Webhook Provider

> Unified HTTP webhook messaging for any endpoint.

[⬅️ Back to project README](../../README.md)

---

## Supported Providers

| Provider    | Website           |
| ----------- | ----------------- |
| **Webhook** | Any HTTP endpoint |

---

## Capabilities

| Provider | GET | POST | PUT | PATCH | DELETE | Headers | Query Params | Body | Notes                       |
| -------- | --- | ---- | --- | ----- | ------ | ------- | ------------ | ---- | --------------------------- |
| Webhook  | ✅  | ✅   | ✅  | ✅    | ✅     | ✅      | ✅           | ✅   | Fully customizable requests |

---

## Features

- Universal HTTP support for any endpoint.
- Multiple HTTP methods (GET, POST, PUT, PATCH, DELETE).
- Custom headers for authentication and content type.
- Flexible raw body content.
- Response validation for JSON, text, and XML.
- Custom success status codes.
- Multiple endpoints with load-balancing strategies.
- Chainable builder API for type-safe message construction.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

config := webhook.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Items: []*webhook.Endpoint{
        {
            Name:     "primary-webhook",
            URL:      "https://api.example.com/webhook",
            Method:   "POST",
            Weight:   100,
            Disabled: false,
        },
    },
}

provider, err := webhook.New(config)
if err != nil {
    log.Fatalf("Failed to create webhook provider: %v", err)
}
```

---

## Quick Builder

```go
msg := webhook.Webhook().
    Method("POST").
    Body([]byte(`{"foo": "bar"}`)).
    Header("Content-Type", "application/json").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := webhook.New(&config)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := webhook.New(&config)
sender.RegisterProvider(core.ProviderTypeWebhook, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## SendVia Helper

`SendVia(accountName, msg)` lets you choose a specific endpoint by name at runtime:

```go
msg := webhook.Webhook().
    Body([]byte(`{"message": "Hello"}`)).
    Header("Content-Type", "application/json").
    Build()

if err := sender.SendVia("primary-webhook", msg); err != nil {
    log.Printf("Primary failed, trying backup: %v", err)
    _ = sender.SendVia("backup-webhook", msg)
}
```

SendVia only switches between endpoints **inside the Webhook provider**; it does not allow cross-provider reuse of one message instance.

---

## Message Types

### 1. JSON Message

```go
jsonData := map[string]interface{}{
    "event": "user.created",
    "data": map[string]interface{}{
        "user_id": "12345",
        "email": "user@example.com",
        "timestamp": time.Now().Unix(),
    },
}
body, _ := json.Marshal(jsonData)

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "application/json").
    Build()
```

### 2. Form Data Message

```go
formData := "action=notify&message=Hello from webhook&priority=high"
body := []byte(formData)

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "application/x-www-form-urlencoded").
    Build()
```

### 3. Raw Text Message

```go
body := []byte("Simple text message")

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "text/plain").
    Build()
```

### 4. GET Request with Query Parameters

```go
msg := webhook.Webhook().
    Method("GET").
    Query("action", "ping").
    Query("timestamp", fmt.Sprintf("%d", time.Now().Unix())).
    Build()
```

---

## Configuration Reference

### Config

- `BaseConfig`: Common configuration fields.
- `Endpoints`: Array of webhook endpoint configurations.

### Endpoint

- `Name`: Endpoint identifier.
- `URL`: Webhook endpoint URL.
- `Method`: HTTP method (default: POST).
- `Headers`: Fixed request headers.
- `QueryParams`: Fixed query parameters.
- `Weight`: Weight for weighted strategy (default: 1).
- `Disabled`: Whether the endpoint is disabled.
- `ResponseConfig`: Response validation configuration.

### ResponseConfig

- `SuccessStatusCodes`: Custom success status codes (default: 2xx).
- `ValidateResponse`: Whether to validate response body (default: false).
- `ResponseType`: Response type ("json", "text", "xml", "none").
- `SuccessField`: JSON field for success indication.
- `SuccessValue`: Expected value for success field.
- `ErrorField`: JSON field for error message.
- `MessageField`: JSON field for response message.
- `SuccessPattern`: Regex pattern for text response success.
- `ErrorPattern`: Regex pattern for text response error.

### Message

- `Body`: Raw request body.
- `Headers`: HTTP headers.

---

## Error Handling

The provider handles:

- Network timeouts and connection errors.
- HTTP error status codes (4xx, 5xx).
- Custom response validation failures.
- Endpoint selection based on strategy.
- Fallback to alternative endpoints on failure.

---

## Best Practices

### 1. Configure Response Validation

```go
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
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "notifications",
            URL:    "https://api.example.com/webhook",
            Method: "POST",
        },
        {
            Name:   "status-check",
            URL:    "https://api.example.com/status",
            Method: "GET",
        },
    },
}
```

### 3. Handle Authentication

```go
msg := webhook.Webhook().
    Body(body).
    Header("X-API-Key", "your-api-key").
    Header("Content-Type", "application/json").
    Build()
```

### 4. Use Multiple Endpoints for Reliability

- Configure multiple endpoints for failover and load balancing.
- Always use `SendVia` for precise endpoint selection.
- Specify endpoint names clearly for maintainability.

---

## Integration Examples

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
    "text":    "Hello from Go-Sender!",
    "channel": "#general",
    "username": "Go-Sender Bot",
}
body, _ := json.Marshal(slackData)

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "application/json").
    Build()
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
                SuccessField:     "id",
            },
        },
    },
}

discordData := map[string]interface{}{
    "content": "Hello from Go-Sender!",
    "embeds": []map[string]interface{}{
        {
            "title":       "Notification",
            "description": "This is a test message",
            "color":       0x00ff00,
        },
    },
}
body, _ := json.Marshal(discordData)

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "application/json").
    Build()
```

### 3. Bark (iOS Push) Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "bark",
            URL:    "https://api.day.app/YOUR_DEVICE_KEY/Hello%20from%20go-sender",
            Method: "GET",
        },
    },
}

msg := webhook.Webhook().Build()
```

### 4. PushDeer Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "pushdeer",
            URL:    "https://api2.pushdeer.com/message/push?pushkey=YOUR_KEY&text=Hello+from+go-sender",
            Method: "GET",
        },
    },
}

msg := webhook.Webhook().Build()
```

### 5. Pushover Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "pushover",
            URL:    "https://api.pushover.net/1/messages.json",
            Method: "POST",
            Headers: map[string]string{
                "Content-Type": "application/x-www-form-urlencoded",
            },
        },
    },
}

form := "token=YOUR_APP_TOKEN&user=USER_KEY&message=Hello+from+go-sender"
msg := webhook.Webhook().
    Body([]byte(form)).
    Build()
```

### 6. SimplePush Integration

```go
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "simplepush",
            URL:    "https://simplepu.sh",
            Method: "POST",
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        },
    },
}

body := []byte(`{"key":"YOUR_KEY","msg":"Hello from go-sender!"}`)
msg := webhook.Webhook().
    Body(body).
    Build()
```

### 7. Custom API Integration (with JSON Response Validation)

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
    "user_id":   "12345",
    "email":     "user@example.com",
    "created_at": time.Now().Format(time.RFC3339),
}
body, _ := json.Marshal(apiData)

msg := webhook.Webhook().
    Body(body).
    Header("Content-Type", "application/json").
    Header("Authorization", "Bearer your-token").
    Header("X-Event-Type", "user.created").
    Build()
```

---

## API Documentation

- [Webhook Provider Guide](https://github.com/shellvon/go-sender)

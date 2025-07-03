[⬅️ Back to Main README](../../README.md)

# Webhook Provider | 通用 Webhook 推送组件

The Webhook Provider for go-sender allows you to send HTTP requests to any endpoint with full control over method, headers, path/query parameters, and body.

本组件支持通过 HTTP Webhook 向任意支持 HTTP 请求的接口推送消息，支持自定义请求方法、Header、路径参数、查询参数和请求体。

---

## Usage | 用法

### Constructing Webhook Messages | 构造 Webhook 消息

Use the builder API for all options. 所有参数均通过 builder 链式设置。

#### API

- `Webhook()` - create a builder instance | 创建 builder 实例
- `Body(body []byte)` - set request body | 设置请求体
- `Method(method string)` - set HTTP method (e.g., http.MethodPost) | 设置 HTTP 方法（如 http.MethodPost）
- `Header(key, value string)` - set a single header | 设置单个 Header
- `Headers(headers map[string]string)` - set multiple headers | 批量设置 Header
- `PathParam(key, value string)` - set a single path parameter | 设置单个路径参数
- `PathParams(params map[string]string)` - set multiple path parameters | 批量设置路径参数
- `Query(key, value string)` - set a single query parameter | 设置单个查询参数
- `Queries(params map[string]string)` - set multiple query parameters | 批量设置查询参数
- `Build()` - build the Message instance | 生成消息实例

**Example | 示例**

```go
msg := webhook.Webhook().
    Method(http.MethodPost). // English: Set HTTP method | 中文：设置 HTTP 方法
    Body([]byte(`{"foo": "bar"}`)). // English: Set request body | 中文：设置请求体
    Header("Authorization", "Bearer token"). // English: Set single header | 中文：设置单个 Header
    Headers(map[string]string{"X-Custom": "value"}). // English: Set multiple headers | 中文：批量设置 Header
    PathParam("id", "123"). // English: Set single path param | 中文：设置单个路径参数
    Query("version", "v1"). // English: Set single query param | 中文：设置单个查询参数
    Build()
```

---

## Features | 功能特性

- Supports all HTTP methods (GET, POST, PUT, DELETE, PATCH, etc.) | 支持所有 HTTP 方法（GET、POST、PUT、DELETE、PATCH 等）
- Flexible headers, path and query parameters | 灵活设置 Header、路径参数、查询参数
- Arbitrary request body (JSON, XML, form, etc.) | 支持任意请求体（JSON、XML、表单等）
- Chainable builder API for type safety and clarity | 链式 builder API，类型安全、易于 IDE 补全

---

## Configuration Example | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "default", // 英文：端点名称
            URL:    "https://example.com/webhook/{id}", // 英文：URL，支持路径参数
            Method: http.MethodPost, // 英文：HTTP 方法
            Headers: map[string]string{
                "Authorization": "Bearer token", // 英文：Header
            },
        },
    },
}

provider, err := webhook.New(config)
if err != nil {
    panic(err)
}
```

---

## Sending a Webhook | 发送 Webhook

```go
msg := webhook.Webhook().
    Method(http.MethodPost).
    Body([]byte(`{"foo": "bar"}`)).
    PathParam("id", "123").
    Query("version", "v1").
    Build()

err := provider.Send(context.Background(), msg, nil)
if err != nil {
    log.Printf("Failed to send webhook: %v", err) // 英文：发送失败 | 中文：发送失败
}
```

---

## Notes | 注意事项

- All fields are optional; set only what you need. | 所有字段均为可选，按需设置。
- If Method is not set, the provider or endpoint config may determine the default (usually POST or GET). | 如未设置 Method，将由 provider 或 endpoint 配置决定（通常为 POST 或 GET）。
- Path and query parameters are merged into the endpoint URL. | 路径参数和查询参数会自动合并到 URL。
- Headers can be set individually or in bulk. | Header 可单独或批量设置。
- Body can be any []byte (JSON, XML, form, etc.). | 请求体可为任意 []byte（JSON、XML、表单等）。

---

## Features | 功能特性

- **Universal HTTP Support 通用 HTTP 支持**: Send messages to any HTTP endpoint | 支持任意 HTTP 接口
- **Multiple Methods 多种请求方法**: Support for GET, POST, PUT, PATCH, DELETE methods | 支持 GET、POST、PUT、PATCH、DELETE 等方法
- **Custom Headers 自定义请求头**: Add custom headers for authentication and content type | 支持自定义请求头
- **Flexible Body Format 灵活请求体**: Support for any raw body content | 支持任意原始请求体
- **Response Validation 响应校验**: Configurable response validation for different webhook formats | 支持多种响应校验方式
- **Multiple Response Types 多种响应类型**: Support for JSON, text, and XML response validation | 支持 JSON、文本、XML 响应校验
- **Custom Status Codes 自定义状态码**: Configure custom success status codes | 支持自定义成功状态码
- **Multiple Endpoints 多端点支持**: Support multiple webhook endpoints with load balancing | 支持多端点负载均衡

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

// English: Create webhook configuration
// 中文：创建 Webhook 配置
config := webhook.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机、加权等
    },
    Endpoints: []webhook.Endpoint{
        {
            Name:     "primary-webhook",
            URL:      "https://api.example.com/webhook",
            Method:   "POST",
            Weight:   100,
            Disabled: false,
        },
        // ... more endpoints
    },
}

provider, err := webhook.New(config)
if err != nil {
    log.Fatalf("Failed to create webhook provider: %v", err) // 创建失败
}
```

---

## Message Types | 消息类型

### 1. JSON Message | JSON 消息

```go
// English: Send JSON data
// 中文：发送 JSON 数据
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
}))
```

### 2. Form Data Message | 表单消息

```go
// English: Send form data
// 中文：发送表单数据
formData := "action=notify&message=Hello from webhook&priority=high"
body := []byte(formData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/x-www-form-urlencoded",
}))
```

### 3. Raw Text Message | 纯文本消息

```go
// English: Send raw text
// 中文：发送纯文本
body := []byte("Simple text message")

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "text/plain",
}))
```

### 4. GET Request with Query Parameters | GET 请求带查询参数

```go
// English: GET request with query params
// 中文：GET 请求带查询参数
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

// Empty body for GET requests | GET 请求 body 为空
msg := webhook.NewMessage([]byte{})
```

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/webhook"
)

s := gosender.NewSender(nil)
webhookProvider, err := webhook.New(config)
if err != nil {
    log.Fatalf("Failed to create webhook provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWebhook, webhookProvider, nil)

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

---

## Message Options | 消息选项

### Header Options | 请求头选项

- `WithHeaders(headers map[string]string)`: Set custom headers | 设置自定义请求头

---

## Configuration Reference | 配置参考

### Config | 配置

- `BaseConfig`: Common configuration fields | 通用配置字段
- `Endpoints`: Array of webhook endpoint configurations | 端点配置数组

### Endpoint | 端点

- `Name`: Endpoint name for identification | 端点名称
- `URL`: Webhook URL (endpoint) | Webhook 地址
- `Method`: HTTP method (default: POST) | HTTP 方法（默认 POST）
- `Headers`: Fixed request headers | 固定请求头
- `QueryParams`: Fixed query parameters | 固定查询参数
- `Weight`: Weight for weighted strategy (default: 1) | 权重（加权策略）
- `Disabled`: Whether this endpoint is disabled | 是否禁用
- `ResponseConfig`: Response validation configuration | 响应校验配置

### ResponseConfig | 响应校验配置

- `SuccessStatusCodes`: Custom success status codes (default: 2xx range) | 自定义成功状态码
- `ValidateResponse`: Whether to validate response body (default: false) | 是否校验响应体
- `ResponseType`: Response type for validation ("json", "text", "xml", "none") | 响应类型
- `SuccessField`: JSON field name indicating success | JSON 成功字段
- `SuccessValue`: Expected value for success field | 成功字段值
- `ErrorField`: JSON field name containing error message | 错误字段
- `MessageField`: JSON field name containing response message | 消息字段
- `SuccessPattern`: Regex pattern for success text response | 文本成功正则
- `ErrorPattern`: Regex pattern for error text response | 文本错误正则

### Message | 消息

- `Body`: Request body (raw bytes) | 原始请求体
- `Headers`: HTTP headers | HTTP 请求头

---

## Error Handling | 错误处理

The provider handles:

本组件处理如下错误：

- Network timeouts and connection errors | 网络超时与连接错误
- HTTP error status codes (4xx, 5xx) | HTTP 错误状态码（4xx, 5xx）
- Custom response validation failures | 自定义响应校验失败
- Provider selection based on strategy | 策略选择账号失败
- Fallback to alternative endpoints on failure | 失败自动切换备用端点

---

## Best Practices | 最佳实践

### 1. Configure Response Validation | 配置响应校验

```go
// English: For APIs that return structured responses
// 中文：对于返回结构化响应的 API
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

### 2. Use Appropriate HTTP Methods | 合理选择 HTTP 方法

```go
// English: Configure method in endpoint
// 中文：在端点配置中指定方法
config := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:   "notifications",
            URL:    "https://api.example.com/webhook",
            Method: "POST", // For notifications | 通知用 POST
        },
        {
            Name:   "status-check",
            URL:    "https://api.example.com/status",
            Method: "GET", // For status checks | 状态查询用 GET
        },
    },
}
```

### 3. Handle Authentication Properly | 正确处理认证

```go
// English: Use headers for API keys
// 中文：通过请求头传递 API Key
msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "X-API-Key": "your-api-key",
    "Content-Type": "application/json",
}))

// Or use bearer tokens | 或使用 Bearer Token
msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Authorization": "Bearer your-bearer-token",
    "Content-Type": "application/json",
}))
```

### 4. Use Multiple Endpoints for Reliability | 多端点提升可靠性

- Configure multiple endpoints for failover and load balancing | 配置多个端点以实现故障切换和负载均衡

---

## API Documentation | 官方文档

- [Webhook Provider Guide | Webhook 组件文档](https://github.com/shellvon/go-sender)

## Integration Examples (English Only)

Below are practical integration examples for popular third-party push services. Each example includes configuration, message construction, and sending with go-sender. You can adapt these patterns for any HTTP-based push service.

### 1. Slack Webhook Integration

```go
import (
    "github.com/shellvon/go-sender/providers/webhook"
    "encoding/json"
)

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

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}))
```

### 2. Discord Webhook Integration

```go
import (
    "github.com/shellvon/go-sender/providers/webhook"
    "encoding/json"
)

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
            "title":       "Notification",
            "description": "This is a test message",
            "color":       0x00ff00,
        },
    },
}
body, _ := json.Marshal(discordData)

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}))
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

msg := webhook.NewMessage([]byte{})
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

msg := webhook.NewMessage([]byte{})
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
msg := webhook.NewMessage([]byte(form))
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
msg := webhook.NewMessage(body)
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

msg := webhook.NewMessage(body, webhook.WithHeaders(map[string]string{
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token",
    "X-Event-Type":  "user.created",
}))
```

### 8. go-sender Unified Send Example

```go
import (
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/webhook"
    "context"
)

sender := gosender.New()
webhookProvider, _ := webhook.New(config)
sender.AddProvider(webhookProvider)

err := sender.Send(context.Background(), msg)
if err != nil {
    log.Printf("Failed to send webhook: %v", err)
}
```

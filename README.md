# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)
[![Go Report Card](https://goreportcard.com/badge/github.com/shellvon/go-sender)](https://goreportcard.com/report/github.com/shellvon/go-sender)
[![GoDoc](https://godoc.org/github.com/shellvon/go-sender?status.svg)](https://godoc.org/github.com/shellvon/go-sender)

> 🚀 The easiest way to send SMS, Email, IM, and Webhook notifications in Go.

[中文文档](./README_CN.md) | [Docs](./docs/getting-started.md)

---

## Why go-sender?

- 🪶 **Lightweight**: Pure Go, zero bloat, minimal dependencies.
- 🧩 **Flexible**: Plug-and-play for SMS, Email, IM, Webhook, and more.
- 🚀 **Simple**: Send a message in just a few lines.
- 🔌 **Extensible**: Add new channels, middleware **and Before/After Hooks** easily.

---

## ⚡ Quick Start

### Method 1: Direct Provider (Simplest)

Send a message **without any setup** - just use the provider directly:

```go
import (
    "context"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // Create provider and send in one go
    account := wecombot.NewAccount("your-webhook-key")
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    
    msg := wecombot.Text().Content("Hello from go-sender!").Build()
    provider.Send(context.Background(), msg, nil)
}
```

### Method 2: Using Sender (For Middleware)

Need retry, rate limiting, or other advanced features? Use the 4-step pattern:

```go
package main

import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // 1️⃣ Create sender
    sender := gosender.NewSender()
    
    // 2️⃣ Create account  
    account := wecombot.NewAccount("your-webhook-key")
    
    // 3️⃣ Register provider
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    sender.RegisterProvider(core.ProviderTypeWecombot, provider, nil)
    
    // 4️⃣ Send message
    msg := wecombot.Text().Content("Hello from go-sender!").Build()
    sender.Send(context.Background(), msg)
}
```

**That's it!** 🎉 This same 4-step pattern works for **any channel**.

---

## 🧠 How It Works

### The Universal Pattern

**Most modern notification services are HTTP APIs** (SMS, IM, Webhooks), while others use specific protocols (SMTP for email). go-sender abstracts this complexity:

1. **Message → Provider Auto-routing**: Any message implementing `core.Message` gets automatically routed to the correct provider based on `ProviderType()`

2. **Protocol Abstraction**: HTTP-based providers use `BaseHTTPTransformer`, while SMTP uses specialized handling. Just tell us the endpoint, parameters, and auth method.

3. **Decorator Pattern**: Want retry? Rate limiting? Queuing? Add middleware without changing your business logic.

4. **Multi-account Strategy**: Round-robin, weighted, or manually specify which account to use.

### Real Examples

```go
// SMS (automatically routed to SMS provider)
msg := sms.Aliyun().To("13800138000").Content("Hello").Build()

// Email (automatically routed to Email provider)
msg := email.NewMessage().To("user@example.com").Subject("Hi").Body("Hello").Build()

// Telegram (automatically routed to Telegram provider)
msg := telegram.Text().Chat("@channel").Text("Hello").Build()

// All use the same sender.Send(ctx, msg) - zero coupling!
```

---

## 🎯 More Examples

### Advanced Features Made Simple

#### Want Middleware? Use Decorators
```go
// Add retry + rate limiting without changing business logic
sender.SetRetryPolicy(core.NewRetryPolicy(core.WithRetryMaxAttempts(3)))
sender.SetRateLimiter(ratelimiter.NewTokenBucketRateLimiter(100, 100))

// Business logic stays the same
sender.Send(ctx, msg) // Now has retry + rate limiting!
```

#### Multiple Accounts? Use Strategies
```go
// Load balance across accounts
accounts := []*sms.Account{
    sms.NewAccount("aliyun", "key1", "secret1"),
    sms.NewAccount("tencent", "key2", "secret2"),
}
provider, _ := sms.NewProvider(accounts,
    sms.Strategy(core.StrategyRoundRobin))

// Or manually pick account
sender.Send(ctx, msg, core.WithSendAccount("backup-account"))
```

#### Complex Authentication? Use Transformers
```go
// WeChat Work App needs access_token? We handle it:
cache := core.NewMemoryCache[*wecomapp.AccessToken]()
provider, _ := wecomapp.New(&config, cache)

// Token fetching, caching, renewal - all automatic!
```

#### Multiple SMS Providers? Use SubProviders
```go
// Same SMS interface, different providers
msg1 := sms.Aliyun().To("phone").Content("Hello").Build()    // → Aliyun API
msg2 := sms.Tencent().To("phone").Content("Hello").Build()   // → Tencent API
msg3 := sms.Huawei().To("phone").Content("Hello").Build()    // → Huawei API

// Same sender handles all
sender.Send(ctx, msg1)
sender.Send(ctx, msg2)
sender.Send(ctx, msg3)
```

## 📦 Installation

```bash
go get github.com/shellvon/go-sender
```

---

## ✨ Supported Providers

| Provider Type | Implementations | Status |
|---------------|-----------------|--------|
| **SMS** | Aliyun, Tencent, Huawei, Yunpian, CL253, Volc, etc. | ✅ Production Ready |
| **Email** | SMTP, EmailJS, Resend | ✅ Production Ready |
| **IM/Bot** | WeChat Work, DingTalk, Lark, Telegram, ServerChan | ✅ Production Ready |
| **Webhook** | Generic HTTP, Custom APIs | ✅ Production Ready |

See [docs/providers.md](docs/providers.md) for complete provider list and configurations.

---

## 🛠 Don't See Your Provider?

**No problem!** go-sender is designed for extensibility:

### 1. Use Generic Webhook
```go
// Step 1: Configure the webhook endpoint
endpoint := &webhook.Endpoint{
    Name:    "my-api",
    URL:     "https://api.example.com/send",
    Method:  "POST",
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
        "Content-Type":  "application/json",
    },
}

provider, _ := webhook.New(&webhook.Config{
    Items: []*webhook.Endpoint{endpoint},
})

// Step 2: Create and send message
msg := webhook.Webhook().
    Body([]byte(`{"message": "Hello World", "recipient": "user123"}`)).
    Build()

provider.Send(context.Background(), msg, nil)
```

### 2. Create Custom Provider
Building custom providers is simple - just implement the `core.Message` interface using `core.BaseMessage`:

```go
// Define your message type
type CustomMessage struct {
    core.BaseMessage  // Handles routing automatically
    Content   string `json:"content"`
    Recipient string `json:"recipient"`
}

func (m *CustomMessage) ProviderType() core.ProviderType {
    return "custom_api"  // This enables auto-routing
}

// Create transformer for HTTP protocol conversion
// See existing providers like wecombot/, sms/, email/ for patterns
```

**Want to dive deeper?** Study these provider implementations:
- **Simple**: [`providers/wecombot/`](./providers/wecombot/) - Basic HTTP webhook
- **Authentication**: [`providers/wecomapp/`](./providers/wecomapp/) - OAuth with caching
- **Multi-vendor**: [`providers/sms/`](./providers/sms/) - SubProvider pattern

See [docs/advanced.md](./docs/advanced.md) for the complete custom provider guide.

---

## 📚 Documentation

- 🚀 [Quick Start Guide](./docs/getting-started.md) - Get running in 5 minutes
- 💡 [Core Concepts](./docs/concepts.md) - Understand the architecture
- 📖 [All Examples](./docs/examples.md) - Real-world usage patterns
- 🔧 [Advanced Usage](./docs/advanced.md) - Custom providers, middleware, deployment
- 🔌 [Provider Reference](./docs/providers.md) - All supported providers

---

**go-sender** — Send anything, anywhere, with Go. One API, All Providers 🚀

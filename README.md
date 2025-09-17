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
- 🔌 **Extensible**: Add new providers, middleware **and Before/After Hooks** easily.

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

**That's it!** 🎉 This same 4-step pattern works for **any provider**.

> 📚 **Want to learn more?** Check out our [comprehensive guides](./docs/getting-started.md)

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

**No problem!** go-sender is designed for extensibility. You have **two approaches**:

### Option 1: Use Generic Webhook (Recommended for HTTP APIs)
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

### Option 2: Create Custom Provider (For Complex Requirements)

For complex authentication, custom protocols, or special requirements:

```go
// 1. Define your message type
type CustomMessage struct {
    core.BaseMessage
    Content string `json:"content"`
}

// 2. Implement provider interface
// See docs/advanced.md for complete guide
```

**Want the full tutorial?** See [Advanced: Custom Providers](./docs/advanced.md#custom-providers)

---

## 📚 Documentation

| **Getting Started** | **Advanced** | **Reference** |
|---------------------|--------------|---------------|
| [📖 Getting Started](./docs/getting-started.md) | [🔧 Advanced Usage](./docs/advanced.md) | [🔌 Providers](./docs/providers.md) |
| [💡 Core Concepts](./docs/concepts.md) | [🧪 Examples](./docs/examples.md) | [❓ FAQ](./docs/faq.md) |
| [🏗️ Architecture](./docs/architecture.md) | [🚦 Middleware](./docs/middleware.md) | [🔧 Troubleshooting](./docs/troubleshooting.md) |

**Quick Navigation:**
- 🆕 **New user?** Start with [Getting Started](./docs/getting-started.md)
- 🔍 **Need a specific provider?** Check [Providers](./docs/providers.md)  
- 🛠 **Want to build custom provider?** See [Advanced Usage](./docs/advanced.md)

---

**go-sender** — Send anything, anywhere, with Go. One API, All Providers 🚀

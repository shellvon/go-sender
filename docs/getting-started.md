# Getting Started

## 🎯 Why go-sender?

**The Problem:** Modern applications need to send notifications through multiple channels - SMS, Email, IM, Webhooks. Each service has different APIs, authentication methods, and failure modes. This leads to:

- 🔄 **Repeated Integration Work**: Every new notification channel means new HTTP clients, error handling, retry logic
- 🚨 **Operational Complexity**: Different monitoring, rate limiting, and failure recovery for each service  
- 🧩 **Tight Coupling**: Business logic scattered across multiple service-specific implementations

**The Solution:** go-sender provides a **unified interface** with **progressive complexity** - start simple, add production features only when needed.

## 🏗️ Design Philosophy

go-sender follows these core principles:

```
┌─────────────┐    ┌──────────────┐    ┌─────────────────┐
│   Message   │───▶│   go-sender  │───▶│   Target API    │
│             │    │   (Router)   │    │                 │
│ - Content   │    │              │    │ - Aliyun SMS    │
│ - To        │    │ ┌──────────┐ │    │ - Gmail SMTP    │
│ - Type()────┼────┼▶│ Provider │ │    │ - WeChat Bot    │
└─────────────┘    │ │ Selection│ │    │ - Telegram API  │
                   │ └──────────┘ │    │ - Webhook...    │
                   │              │    └─────────────────┘
                   │ ┌──────────┐ │
                   │ │Middleware│ │    
                   │ │- Retry   │ │    
                   │ │- Limit   │ │    
                   │ │- Circuit │ │    
                   │ └──────────┘ │    
                   └──────────────┘    
```

**How It Works:**

1. **📝 Create Message**: Build a message using provider-specific builders (e.g., `sms.Aliyun()`, `email.NewMessage()`)
2. **🎯 Auto-Route**: Message's `ProviderType()` determines which provider handles it - no manual routing needed
3. **⚖️ Account Selection**: If multiple accounts exist, strategies (round-robin, weighted, health-based) pick the best one
4. **🔄 Middleware Processing**: Optional retry, rate limiting, circuit breaker, metrics collection
5. **🚀 API Call**: Provider transforms message to target API format and sends

**Key Benefits:**

- **🎯 Auto-Routing**: Messages define their destination - no manual switching
- **🔄 Decorator Pattern**: Add middleware without changing business logic  
- **🏗️ Progressive Complexity**: Start simple, add features when needed
- **⚖️ Multi-Account Strategy**: Built-in load balancing and health checks

## 📈 Progressive Learning Path

> **Already read the README?** Jump to [Level 3: Production Ready](#level-3-production-ready) or [Level 4: Multi-Provider](#level-4-multi-provider)

Choose your starting level based on your needs:

| **Use Case** | **Level** | **Why** |
|--------------|-----------|---------|
| Quick script / Testing | [Level 1](#level-1-simplest-start) | Zero setup, immediate results |
| Small application | [Level 2](#level-2-add-structure) | Better structure, error handling |
| Production service | [Level 3](#level-3-production-ready) | Reliability, monitoring, load balancing |
| Multi-provider platform | [Level 4](#level-4-multi-provider) | Unified API for all notification types |

### Level 1: Simplest Start (30 seconds)

**Already know this?** This is identical to README examples.

**Send your first message with zero configuration:**

```go
// Same as README - minimal setup
account := wecombot.NewAccount("your-webhook-key")
provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
msg := wecombot.Text().Content("Hello go-sender!").Build()
provider.Send(context.Background(), msg, nil)
```

### Level 2: Add Structure (2 minutes)

**Use the Sender for better organization:**

```go
package main

import (
    "context"
    "log"
    
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // 1️⃣ Initialize sender
    sender := gosender.NewSender()
    
    // 2️⃣ Register provider
    account := wecombot.NewAccount("your-webhook-key")
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    sender.RegisterProvider(core.ProviderTypeWecombot, provider, nil)
    
    // 3️⃣ Send with detailed results
    msg := wecombot.Text().Content("Hello go-sender!").Build()
    result, err := sender.SendWithResult(context.Background(), msg)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    log.Printf("Success! Status: %d", result.StatusCode)
}
```

**What you gain:** Structured error handling, detailed results, and preparation for advanced features.

### Level 3: Production Ready (5 minutes)

**Add retry, rate limiting, and multi-account support:**

```go
import (
    "time"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/ratelimiter"
)

func main() {
    sender := gosender.NewSender()
    
    // Multiple accounts for high availability
    accounts := []*wecombot.Account{
        wecombot.NewAccount("primary-webhook-key"),
        wecombot.NewAccount("backup-webhook-key"),
    }
    
    provider, _ := wecombot.NewProvider(accounts)
    
    // Production middleware
    middleware := &core.SenderMiddleware{
        RateLimiter: ratelimiter.NewTokenBucketRateLimiter(10, 5), // 10 QPS, burst 5
        Retry: &core.RetryPolicy{
            MaxAttempts: 3,
            InitialDelay: time.Second,
            MaxDelay: 10 * time.Second,
        },
    }
    
    sender.RegisterProvider(core.ProviderTypeWecombot, provider, middleware)
    
    // Your code stays the same!
    msg := wecombot.Text().Content("Production ready!").Build()
    sender.Send(context.Background(), msg)
}
```

**Enterprise features unlocked:** Multi-account load balancing, rate limiting, exponential backoff retry, circuit breaker.

### Level 4: Multi-Provider (10 minutes)

**Add SMS, Email, and other providers:**

```go
import (
    "github.com/shellvon/go-sender/providers/sms"
    "github.com/shellvon/go-sender/providers/email"
)

func main() {
    sender := gosender.NewSender()
    
    // Register multiple providers
    registerWeCom(sender)
    registerSMS(sender)
    registerEmail(sender)
    
    // Auto-routing: message type determines the provider
    wecomMsg := wecombot.Text().Content("WeChat notification").Build()
    smsMsg := sms.Aliyun().To("13800138000").Content("SMS alert").Build()
    emailMsg := email.NewMessage("Alert", "Email notification", "admin@company.com")
    
    // All use the same API - go-sender routes automatically
    sender.Send(context.Background(), wecomMsg)  // → WeCom provider
    sender.Send(context.Background(), smsMsg)    // → SMS provider
    sender.Send(context.Background(), emailMsg)  // → Email provider
}

func registerSMS(sender *gosender.Sender) {
    account := sms.NewAccount("aliyun", "key", "secret", sms.WithSignName("MyApp"))
    provider, _ := sms.NewProvider([]*sms.Account{account})
    sender.RegisterProvider(core.ProviderTypeSMS, provider, nil)
}
```

**The power:** One API for all providers. Message types automatically route to the right provider.

## 🚀 What's Next?

Congratulations! You've learned the progressive path from simple scripts to enterprise-ready notification systems.

### Choose Your Path:

| **If you want to...** | **Go to** |
|-----------------------|-----------|
| 🔍 **Find specific providers** | [Providers Reference](./providers.md) |
| 🛠 **Build custom providers** | [Advanced Usage](./advanced.md) |
| 🧪 **See real examples** | [Examples](./examples.md) |
| 🚦 **Add middleware** | [Middleware Guide](./middleware.md) |
| 💡 **Understand architecture** | [Core Concepts](./concepts.md) |
| 🔧 **Troubleshoot issues** | [Troubleshooting](./troubleshooting.md) |

---

**Ready to send your first message?** Pick your level above and start coding! 🎉

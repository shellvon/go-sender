# Getting Started

Welcome to **go-sender** – the unified notification library that grows with your needs, from simple scripts to enterprise applications.

## 🚀 Quick Installation

```bash
go get github.com/shellvon/go-sender
```

## 📈 Progressive Learning Path

### Level 1: Simplest Start (30 seconds)

**Send your first message with zero configuration:**

```go
import (
    "context"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // Create and send in one go - no setup needed!
    account := wecombot.NewAccount("your-webhook-key")
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    
    msg := wecombot.Text().Content("Hello go-sender!").Build()
    provider.Send(context.Background(), msg, nil)
}
```

**Why this works:** Most providers need only an API key or webhook URL. Perfect for testing or simple automation scripts.

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

**Enterprise features unlocked:** Automatic failover, rate limiting, exponential backoff retry, circuit breaker.

### Level 4: Multi-Channel (10 minutes)

**Add SMS, Email, and other channels:**

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

**The power:** One API for all channels. Message types automatically route to the right provider.

## 🎯 What's Your Use Case?

Choose your starting level based on your needs:

| **Scenario** | **Start At** | **Why** |
|--------------|-------------|---------|
| Quick script / Testing | Level 1 | Zero setup, immediate results |
| Small application | Level 2 | Better structure, error handling |
| Production service | Level 3 | Reliability, monitoring, failover |
| Multi-channel platform | Level 4 | Unified API for all notification types |

## 📋 Supported Providers

- **SMS**: Aliyun, Tencent, Huawei, Volc, Yunpian, CL253, and more
- **Email**: SMTP, EmailJS, Resend
- **IM/Bot**: WeCom, DingTalk, Lark, Telegram
- **Webhook**: Universal HTTP integration for any API

See [providers.md](./providers.md) for the complete list.

## 💡 Key Benefits

### 🔄 **Auto-Routing**
Message types automatically route to the right provider - no manual switching needed.

### 🏗️ **Progressive Architecture**
Start simple, add complexity only when you need it. Your code doesn't break as you scale.

### 🛡️ **Production Ready**
Built-in retry, rate limiting, circuit breaker, and multi-account failover.

### 🧩 **Extensible**
Can't find a provider? Create custom ones in ~50 lines of code.

## 🔧 Advanced Features

### Custom HTTP Client
```go
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
        TLSClientConfig: &tls.Config{/* custom TLS */},
    },
}

sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

### Hooks for Custom Logic
```go
middleware := &core.SenderMiddleware{}
middleware.UseBeforeHook(func(ctx context.Context, msg core.Message, opts *core.SendOptions) error {
    log.Printf("Sending message: %s", msg.MsgID())
    return nil
})

sender.RegisterProvider(providerType, provider, middleware)
```

## 🚀 Ready to Level Up?

### Next Steps:
- **Level 1-2 Users**: Check out [providers.md](./providers.md) for more channels
- **Level 3-4 Users**: Explore [middleware.md](./middleware.md) for advanced features
- **Custom Needs**: See [advanced.md](./advanced.md) for custom providers and deep customization
- **Real Examples**: Browse [examples.md](./examples.md) for production scenarios

### Quick References:
- [Core Concepts](./concepts.md) - Understanding the architecture
- [Troubleshooting](./troubleshooting.md) - Common issues and solutions

---

**Ready to send your first message?** Pick your level above and start coding! 🎉

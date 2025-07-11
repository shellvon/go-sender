# Getting Started

Welcome to **go-sender** – the simple, flexible, and extensible notification library for Go developers.

## 🚀 Quick Installation

```bash
go get github.com/shellvon/go-sender
```

## 🏁 Your First Message

### Method 1: Direct Provider Usage (No Middleware)

Send SMS directly using the provider:

```go
package main

import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

func main() {
    // Create SMS provider
    config := sms.Config{
        ProviderMeta: core.ProviderMeta{
            Strategy: core.StrategyRoundRobin,
        },
        Items: []*sms.Account{{
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "default",
                    SubType: "aliyun",
                },
                Credentials: core.Credentials{
                    APIKey:    "your-access-key",
                    APISecret: "your-secret-key",
                },
            },
            // 可选: Region, Callback 等 sms.Account 字段
        }},
    }

    provider, err := sms.New(config)
    if err != nil {
        panic(err)
    }

    // Create and send message
    msg := sms.Aliyun().
        To("13800138000").
        Content("Hello from go-sender!").
        TemplateID("SMS_xxx").
        Build()

    err = provider.Send(context.Background(), msg)
    if err != nil {
        panic(err)
    }
}
```

### Method 2: Using Sender with Provider Registration

Register provider with sender for middleware support:

```go
package main

import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

func main() {
    // Create sender
    sender := gosender.NewSender()

    // Create and register SMS provider
    config := sms.Config{
        ProviderMeta: core.ProviderMeta{
            Strategy: core.StrategyRoundRobin,
        },
        Items: []*sms.Account{{
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "default",
                    SubType: "aliyun",
                },
                Credentials: core.Credentials{
                    APIKey:    "your-access-key",
                    APISecret: "your-secret-key",
                },
            },
            // 可选: Region, Callback 等 sms.Account 字段
        }},
    }

    smsProvider, err := sms.New(config)
    if err != nil {
        panic(err)
    }
    sender.RegisterProvider(core.ProviderTypeSMS, smsProvider, nil)

    // Create and send message
    msg := sms.Aliyun().
        To("13800138000").
        Content("Hello from go-sender!").
        TemplateID("SMS_xxx").
        Build()

    err = sender.Send(context.Background(), msg)
    if err != nil {
        panic(err)
    }
}
```

## ✉️ Supported Channels

- SMS: Aliyun, Tencent, Huawei, Yunpian, etc.
- Email: SMTP, EmailJS, Resend
- IM/Bot: WeCom, DingTalk, Lark, Telegram
- Webhook: Universal HTTP integration

See [providers.md](./providers.md) for the full list.

## 🧑‍💻 FAQ

**Q: Is go-sender production ready?**  
A: Yes, but always test with your own provider credentials and templates.

**Q: How do I add a new provider?**  
A: See [advanced.md](./advanced.md) for custom provider instructions.

**Q: Can I use go-sender in microservices?**  
A: Absolutely! It is designed for both monoliths and microservices.

**Q: When should I use Method 1 vs Method 2?**  
A: Use Method 1 for simple cases without middleware. Use Method 2 when you need rate limiting, retry, circuit breaker, or other middleware features.

**Q: Can I use a custom HTTP client?**  
A: Yes! You can pass a custom `*http.Client` to `sender.Send()` for advanced features:

```go
import (
    "crypto/tls"
    "net/http"
    "net/url"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Custom HTTP client with timeout and proxy
client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        Proxy: http.ProxyURL(&url.URL{
            Scheme: "http",
            Host:   "proxy.example.com:8080",
        }),
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: false,
        },
    },
}

// Use custom client for all requests
err = sender.Send(context.Background(), msg, core.WithSendHTTPClient(client))
```

**Benefits of custom HTTP client:**

- **Timeout Control**: Prevent hanging requests
- **Proxy Support**: Route through corporate proxies
- **TLS Configuration**: Custom certificates and security settings
- **Connection Pooling**: Optimize performance for high-volume sending
- **Retry Logic**: Built-in retry with exponential backoff
- **Load Balancing**: Distribute requests across multiple endpoints
- **Authentication**: Custom auth headers or certificates
- **Monitoring**: Add request/response logging and metrics
- **Caching**: Implement response caching for repeated requests

## 🪝 Using Hooks (Before / After)

Need to run custom logic before or after each send? go-sender provides **Hooks**:

```go
mw := &core.SenderMiddleware{}
mw.UseBeforeHook(func(_ context.Context, m core.Message, _ *core.SendOptions) error {
    fmt.Println("GLOBAL BEFORE", m.MsgID())
    return nil
})

sender.RegisterProvider(core.ProviderTypeSMS, smsProvider, mw)

// Per-request hooks:
sender.Send(ctx, msg,
    core.WithSendAfterHooks(func(_ context.Context, _ core.Message, _ *core.SendOptions, _ *core.SendResult, err error) {
        fmt.Println("PER-REQ AFTER, err:", err)
    }),
)
```

Execution order: `global before → per-request before → send → global after → per-request after`.

## 📚 Next Steps

- [Core Concepts](./concepts.md)
- [Provider Usage](./providers.md)
- [Middleware & Advanced Features](./middleware.md)
- [Examples](./examples.md)

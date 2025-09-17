# Core Concepts

go-sender is designed to be simple, decoupled, and extensible. This guide explains the key architectural concepts and how they work together.

> **Quick Reference**: Looking for specific usage patterns? See [Getting Started](./getting-started.md) for hands-on examples.

## Sender

The **central orchestrator** that manages providers and middleware chains. Created via `gosender.NewSender()`.

**Key Responsibilities:**
- **Provider Management**: Register and route messages to appropriate providers
- **Middleware Execution**: Apply rate limiting, retry, circuit breaker transparently  
- **Message Routing**: Auto-route based on message's `ProviderType()`
- **Runtime Configuration**: Override account/strategy per request with `core.WithSendAccount()`

**API Methods:**
- `Send(ctx, msg)` - Fire-and-forget sending
- `SendWithResult(ctx, msg)` - Get detailed response and metrics

## Provider System

Providers are **pluggable components** that implement the `core.Provider` interface for specific communication channels.

**Provider Ecosystem:**

| **Category** | **Providers** | **Key Features** |
|-------------|---------------|------------------|
| **SMS** | Aliyun, Tencent, Huawei, CL253, Volc | Template support, signature management, multi-region |
| **Email** | SMTP, EmailJS, Resend, Mailgun | Direct SMTP or API-based delivery |
| **IM/Bot** | WeCom, DingTalk, Lark, Telegram | Rich media, markdown, interactive cards |
| **Webhook** | Generic HTTP, Custom APIs | Universal integration for any HTTP API |

**Provider Features:**
- **Account Management**: Multiple accounts per provider with load balancing
- **Message Transformation**: Convert messages to provider-specific API formats
- **Error Handling**: Provider-specific error mapping and retry logic

## Message Construction

Messages are built using **fluent builder patterns** specific to each provider type.

**Builder Examples:**
```go
// SMS Message
smsMsg := sms.Aliyun().To("13800138000").Content("Hello").Build()

// Email Message  
emailMsg := email.NewMessage().To("user@example.com").Subject("Hi").Body("Hello").Build()

// IM Message
imMsg := wecombot.Text().Content("Notification").Build()
```

**Key Features:**
- **Type Safety**: Compile-time validation of required fields
- **Provider-Specific**: Each provider exposes relevant options (templates, attachments, etc.)
- **Auto-Routing**: Message's `ProviderType()` determines which provider handles it

## Integration Approaches

go-sender supports two primary integration patterns:

### Pattern 1: Direct Provider Usage
**Best for:** Simple scenarios without middleware requirements

```go
// Direct provider usage - minimal setup
account := wecombot.NewAccount("webhook-key")
provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
msg := wecombot.Text().Content("Hello").Build()
provider.Send(ctx, msg, nil)
```

### Pattern 2: Sender with Provider Registration  
**Best for:** Production systems needing middleware (retry, rate limiting, circuit breakers)

```go
// Sender orchestrator with middleware support
sender := gosender.NewSender()
provider, _ := wecombot.NewProvider(accounts)
middleware := &core.SenderMiddleware{
    RateLimiter: ratelimiter.NewTokenBucketRateLimiter(10, 5),
    Retry: &core.RetryPolicy{MaxAttempts: 3},
}
sender.RegisterProvider(core.ProviderTypeWecombot, provider, middleware)
sender.Send(ctx, msg)
```

## Middleware Architecture

The library implements a **decorator pattern** where `SenderMiddleware` wraps providers with cross-cutting concerns:

| **Component** | **Purpose** | **Example** |
|---------------|-------------|-------------|
| **Rate Limiter** | Prevent API rate limit violations | `10 QPS, burst 5` |
| **Retry Policy** | Handle transient failures | `3 attempts, exponential backoff` |
| **Circuit Breaker** | Prevent cascading failures | `Open after 5 failures` |
| **Queue** | Async processing and batching | `Redis, memory queue` |
| **Metrics** | Observability and monitoring | `Success/failure rates` |

## Extensibility Model

The library is designed for extensibility through well-defined interfaces:

| **Extension Point** | **Interface** | **Use Case** |
|-------------------|---------------|--------------|
| **Custom Providers** | `core.Provider` | New communication channels |
| **Custom Middleware** | Middleware interfaces | Cross-cutting concerns |
| **HTTP Transformers** | `HTTPRequestTransformer` | Custom request/response handling |
| **Selection Strategies** | `SelectionStrategy` | Account selection algorithms |

**Extension Examples:**
```go
// Custom Provider
type MyProvider struct { /* implementation */ }
func (p *MyProvider) Send(ctx, msg, opts) (*SendResult, error) { /* ... */ }

// Custom Strategy  
type GeoStrategy struct { /* implementation */ }
func (s *GeoStrategy) Select(accounts []Selectable) Selectable { /* ... */ }
```

---

## What's Next?

| **To Learn About** | **Go To** |
|-------------------|-----------|
| 🚀 **Hands-on Examples** | [Getting Started](./getting-started.md) |
| 🔌 **Provider Details** | [Providers](./providers.md) |
| 🚦 **Middleware Setup** | [Middleware](./middleware.md) |
| 🛠 **Custom Extensions** | [Advanced Usage](./advanced.md) |

# Advanced Usage

Unlock the full power of go-sender with advanced customization and extension points.

## Custom Providers

go-sender provides two approaches for creating custom providers:

### Quick Start: Simple Provider

For basic use cases, implement the `core.Provider` interface directly:

```go
import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
)

const ProviderTypeCustom core.ProviderType = "custom"

type CustomProvider struct{}

func (p *CustomProvider) Send(ctx context.Context, msg core.Message, _ *core.ProviderSendOptions) (*core.SendResult, error) {
    // 1) Validate / cast the message to your concrete type
    // 2) Talk to remote service
    // 3) Build & return *core.SendResult
    return &core.SendResult{StatusCode: 200}, nil
}

func (p *CustomProvider) Name() string { return string(ProviderTypeCustom) }

sender := gosender.NewSender()
sender.RegisterProvider(ProviderTypeCustom, &CustomProvider{}, nil)
```

### Production Ready: HTTP Provider with BaseMessage

For HTTP-based services (most modern APIs), use the structured approach with `BaseMessage` and `HTTPProvider`:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers"
)

// Step 1: Define provider type constant
const CustomProviderType core.ProviderType = "custom"

// Step 2: Create your message type using BaseMessage
type CustomMessage struct {
    core.BaseMessage  // Provides default implementations
    Content   string `json:"content"`
    Recipient string `json:"recipient"`
}

// Only need to implement ProviderType() - BaseMessage handles the rest
func (m *CustomMessage) ProviderType() core.ProviderType {
    return CustomProviderType
}

// Optional: Add validation logic
func (m *CustomMessage) Validate() error {
    if m.Content == "" || m.Recipient == "" {
        return fmt.Errorf("content and recipient are required")
    }
    return nil
}

// Step 3: Create account type
type CustomAccount struct {
    core.BaseAccount
    APIEndpoint string `json:"api_endpoint"`
}

// Step 4: Create the transformer (handles HTTP protocol conversion)
type customTransformer struct{}

func (t *customTransformer) CanTransform(msg core.Message) bool {
    _, ok := msg.(*CustomMessage)
    return ok
}

func (t *customTransformer) Transform(ctx context.Context, msg core.Message, account *CustomAccount) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
    customMsg := msg.(*CustomMessage)
    
    // Build request body
    body, err := json.Marshal(map[string]interface{}{
        "message": customMsg.Content,
        "to":      customMsg.Recipient,
    })
    if err != nil {
        return nil, nil, err
    }
    
    return &core.HTTPRequestSpec{
        Method: http.MethodPost,
        URL:    account.APIEndpoint,
        Headers: map[string]string{
            "Content-Type":  "application/json",
            "Authorization": "Bearer " + account.APIKey,
        },
        Body:     body,
        BodyType: core.BodyTypeJSON,
    }, nil, nil
}

// Step 5: Create the provider factory
func NewCustomProvider(accounts []*CustomAccount) (*providers.HTTPProvider[*CustomAccount], error) {
    config := &core.BaseConfig[*CustomAccount]{
        Items: accounts,
    }
    
    return providers.NewHTTPProvider(
        string(CustomProviderType),
        &customTransformer{},
        config,
    )
}

// Usage example:
func main() {
    // Create account
    account := &CustomAccount{
        BaseAccount: core.BaseAccount{
            AccountMeta: core.AccountMeta{Name: "main"},
            Credentials: core.Credentials{APIKey: "your-api-key"},
        },
        APIEndpoint: "https://api.example.com/send",
    }
    
    // Create provider
    provider, _ := NewCustomProvider([]*CustomAccount{account})
    
    // Create message with BaseMessage convenience
    msg := &CustomMessage{
        BaseMessage: *core.NewBaseMessage(CustomProviderType),
        Content:     "Hello World",
        Recipient:   "user123",
    }
    
    // Send message
    result, err := provider.Send(context.Background(), msg, nil)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Sent! Status: %d\n", result.StatusCode)
    }
}
```

### Understanding the Architecture

The structured approach above follows go-sender's core principles:

1. **Message Routing**: `ProviderType()` enables automatic routing by the sender
2. **Protocol Abstraction**: `Transformer` converts your message to HTTP requests
3. **Account Management**: Multiple accounts with selection strategies and load balancing
4. **Middleware Integration**: Automatic retry, rate limiting, circuit breaker support

**Study Existing Providers**: The best way to understand the patterns is to examine existing providers:
- **Simple HTTP**: [`providers/wecombot/`](../providers/wecombot/) - Basic webhook
- **Authentication**: [`providers/wecomapp/`](../providers/wecomapp/) - OAuth with token caching
- **Multi-Vendor**: [`providers/sms/`](../providers/sms/) - SubProvider pattern for multiple backends
- **Complex Protocols**: [`providers/email/`](../providers/email/) - SMTP integration

## Custom Middleware

Every cross-cutting component is just an interface. Drop-in your own implementation and wire it into the sender:

```go
// 1. Rate limiter
sender.SetRateLimiter(myRateLimiter)              // implements core.RateLimiter

// 2. Retry policy (struct)
sender.SetRetryPolicy(&core.RetryPolicy{ /* … */ })

// 3. Circuit breaker
sender.SetCircuitBreaker(myCircuitBreaker)        // implements core.CircuitBreaker

// 4. Queue implementation (e.g. Redis, Kafka, RabbitMQ)
sender.SetQueue(myQueue)                          // implements core.Queue

// 5. Metrics collector
sender.SetMetrics(myCollector)                    // implements core.MetricsCollector

// Or provide everything at once when registering a provider
mw := &core.SenderMiddleware{Queue: myQueue, Retry: customRetry}
sender.RegisterProvider(core.ProviderTypeSMS, smsProvider, mw)
```

### Example: Simple Token Bucket Rate-Limiter

```go
type SimpleLimiter struct{
    limiter *rate.Limiter
}

func NewSimpleLimiter(qps float64, burst int) *SimpleLimiter {
    return &SimpleLimiter{limiter: rate.NewLimiter(rate.Limit(qps), burst)}
}

func (l *SimpleLimiter) Allow() bool                  { return l.limiter.Allow() }
func (l *SimpleLimiter) Wait(ctx context.Context) error { return l.limiter.Wait(ctx) }
func (l *SimpleLimiter) Close() error                 { return nil }

sender.SetRateLimiter(NewSimpleLimiter(20, 5))
```

## Custom Selection Strategy

Need a bespoke account-selection rule? Implement `core.SelectionStrategy`.

```go
// GeoStrategy picks the first account whose Region matches the phone prefix.
type GeoStrategy struct{}

func (g *GeoStrategy) Name() core.StrategyType { return "geo_based" }

func (g *GeoStrategy) Select(items []core.Selectable) core.Selectable {
    for _, it := range items {
        if !it.IsEnabled() { continue }
        if strings.HasPrefix(userPhone, "+1") && it.GetName() == "us" {
            return it
        }
        // … add more rules
    }
    return items[0] // fallback
}

// Register to the global registry so every provider can see it.
core.GlobalStrategyRegistry.Register(g.Name(), g)

cfg.ProviderMeta.Strategy = g.Name() // use custom strategy
```

Now your provider will pick accounts according to `GeoStrategy` whenever `Strategy` is set to `geo_based`.

## Custom HTTP Client

You can use a custom `http.Client` for any HTTP-based provider:

```go
import (
    "crypto/tls"
    "net/http"
    "net/url"
    "time"
    "github.com/shellvon/go-sender/core"
)

// Create custom HTTP client with advanced configuration
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
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// Use custom client for all requests
sender.Send(ctx, msg, core.WithSendHTTPClient(client))
```

**See [examples.md](./examples.md) for detailed examples of:**

- Timeout control and proxy configuration
- TLS security and certificate management
- Connection pooling for high-performance scenarios
- Custom retry logic and load balancing
- Authentication, monitoring, and caching implementations

## Strategies

- Weighted round-robin, failover, and more.
- See [core/strategy.go](../core/strategy.go) for details.

## Extending Message Types

You can define your own message types for new channels or advanced scenarios.

---

**See [examples.md](./examples.md) for real-world scenarios.**

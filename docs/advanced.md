# Advanced Usage

Unlock the full power of go-sender with advanced customization and extension points.

## Custom Providers

You can add your own provider by implementing the Provider interface and registering it with the sender.

```go
import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
)

const ProviderTypeCustom core.ProviderType = "custom"

// CustomProvider demonstrates how to plug in your own channel.
// Implement the core.Provider interface.
type CustomProvider struct{}

// Send implements your actual delivery logic. Options let you access
// a per-send *http.Client or other future extensions.
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

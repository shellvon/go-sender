# Advanced Usage

Unlock the full power of go-sender with advanced customization and extension points.

## Custom Providers

You can add your own provider by implementing the Provider interface and registering it with the sender.

```go
import (
    "context"
    "github.com/shellvon/go-sender/core"
)

const ProviderTypeCustom core.ProviderType = "custom"

type CustomProvider struct{}

func (p *CustomProvider) Send(ctx context.Context, msg core.Message) error {
    // Your sending logic here
    return nil
}

func (p *CustomProvider) Name() string {
    return string(ProviderTypeCustom)
}

sender := gosender.NewSender()
sender.RegisterProvider(ProviderTypeCustom, &CustomProvider{}, nil)
```

## Custom Middleware

支持设置指定的中间件，实现 Middleware 接口即可：

```go
// 设置 rate limiter
sender.SetRateLimiter(myRateLimiter)

// 设置 retry policy
sender.SetRetryPolicy(&core.RetryPolicy{...})

// 设置 circuit breaker
sender.SetCircuitBreaker(myCircuitBreaker)

// 设置 queue
sender.SetQueue(myQueue)
```

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

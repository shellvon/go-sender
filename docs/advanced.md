# Advanced Usage

Unlock the full power of go-sender with advanced customization and extension points.

## Custom Providers

You can add your own provider by implementing the Provider interface and registering it with the sender.

```go
// Example: CustomProvider implements core.Provider
import "github.com/shellvon/go-sender/core"

type CustomProvider struct{}

func (p *CustomProvider) Send(ctx context.Context, msg core.Message) error {
    // Your sending logic here
    return nil
}

sender := sender.NewSender()
sender.RegisterProvider("custom", &CustomProvider{})
```

## Custom Middleware

支持设置指定的中间件，实现 Middleware 接口即可：

```go
type MyMiddleware struct{}

func (m *MyMiddleware) Handle(next core.SendFunc) core.SendFunc {
    return func(ctx context.Context, msg core.Message) error {
        // Pre-processing
        err := next(ctx, msg)
        // Post-processing
        return err
    }
}

sender.SetRateLimiter(&MyMiddleware{}) // Or use SetMiddlewareChain
```

## Custom HTTP Client

You can use a custom `http.Client` for any HTTP-based provider:

```go
sender.Send(ctx, msg, core.WithSendHTTPClient(myClient))
```

## Strategies

- Weighted round-robin, failover, and more.
- See [core/strategy.go](../core/strategy.go) for details.

## Extending Message Types

You can define your own message types for new channels or advanced scenarios.

---

**See [examples.md](./examples.md) for real-world scenarios.**

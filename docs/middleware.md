# Middleware

go-sender supports powerful middleware for cross-cutting concerns.

## Built-in Middleware

- **Rate Limiter**: Limit the rate of sending messages.
- **Retry Policy**: Automatic retry on failure.
- **Circuit Breaker**: Prevent repeated failures from overwhelming providers.
- **Queue**: Asynchronous message queue.
- **Metrics**: Collect and export metrics.

## Usage Example

```go
import (
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/ratelimiter"
)

sender := sender.NewSender()
limiter := ratelimiter.NewTokenBucketRateLimiter(10, 1) // 10 reqs/sec
sender.SetRateLimiter(limiter)
```

## Custom Middleware

- You can replace any built-in component with your own implementation.
- Implement the corresponding interface (e.g. `RateLimiter`, `Queue`, `CircuitBreaker`, `MetricsCollector`).
- Attach the component via the dedicated setter on `Sender`, for example:

```go
myLimiter := mypkg.NewAwesomeLimiter()
sender.SetRateLimiter(myLimiter) // any struct that satisfies core.RateLimiter
```

If you need to bundle several components at once, build a `core.SenderMiddleware` struct and pass it when registering the provider:

```go
mw := &core.SenderMiddleware{
    Queue:       myQueueImpl,      // implements core.Queue
    Retry:       myRetryPolicy,    // *core.RetryPolicy
    CircuitBreaker: myCB,          // implements core.CircuitBreaker
}

sender.RegisterProvider(core.ProviderTypeSMS, smsProvider, mw)
```

---

**See [advanced.md](./advanced.md) for more customization.**

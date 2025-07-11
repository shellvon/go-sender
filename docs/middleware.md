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
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/ratelimiter"
)

sender := gosender.NewSender()
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

## Hooks vs Middleware

| Aspect          | Middleware (RateLimiter / Retry …)                                | Hooks (Before / After)                                     |
| --------------- | ----------------------------------------------------------------- | ---------------------------------------------------------- |
| Purpose         | Infrastructure & reliability (rate-limit, retry, circuit-breaker) | Lightweight business extensions: logging, tracing, masking |
| Lifecycle       | Long-lived, initialised once                                      | Executed on every send                                     |
| Can abort flow? | Depends on implementation                                         | **BeforeHook can abort**                                   |
| Alters result?  | May change result (e.g. Retry)                                    | AfterHook never alters result                              |

Choose **Hooks** for logging, tracing, data-masking, etc.
Choose **Middleware** for flow-control features like rate-limit, retry or circuit-breaker.

---

**See [advanced.md](./advanced.md) for more customization.**

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

You can implement your own middleware by following the Middleware interface.

---

**See [advanced.md](./advanced.md) for more customization.**

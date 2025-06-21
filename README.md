# Go-Sender

English | [中文](./README_CN.md)

A high-performance, extensible Go message sending framework supporting multiple notification channels and rich middleware capabilities.

## 🎯 Design Philosophy

Go-Sender is designed around the **Decorator Pattern** and **Plugin Architecture**, making it easy to add new notification channels or cross-cutting concerns without changing your business logic.

### Core Design Principles

- **🔄 Decoupling**: Business code only cares about sending messages, not how they're delivered
- **🔌 Pluggable**: Easy to add new providers or middleware through interfaces
- **🛡️ Reliability**: Built-in retry, circuit breaker, and rate limiting
- **📊 Observable**: Comprehensive metrics and health checks
- **⚡ Flexible**: Support for multiple instances, strategies, and configurations

### Architecture Overview

```
Business Logic → Sender → ProviderDecorator → Provider
                      ↓
                Middleware Chain:
                - Rate Limiter
                - Circuit Breaker
                - Retry Policy
                - Queue
                - Metrics
```

## ✨ Features

### 🚀 Multi-Channel Support

- **Email**: SMTP with multi-account support
- **WeCom Bot**: Enterprise WeChat bot messages
- **Webhook**: Generic HTTP webhook calls
- **Extensible**: Easy to add Telegram, Slack, Discord, etc.

### 🛡️ Advanced Reliability Features

- **Smart Retry**: Configurable retry policies with exponential backoff
- **Circuit Breaker**: Prevents cascading failures
- **Rate Limiting**: Token bucket and sliding window algorithms
- **Queue Support**: In-memory and distributed queues
- **Health Checks**: Comprehensive health monitoring

### 🎛️ Multi-Instance & Strategy Support

- **Multiple Accounts**: Support multiple email accounts, WeCom bots, webhook endpoints
- **Load Balancing**: Round-robin, random, weighted, and health-based strategies
- **Context-Aware**: Override strategies per request via context

### 📊 Observability

- **Metrics Collection**: Performance and outcome metrics
- **Health Monitoring**: Provider and system health checks
- **Structured Logging**: Pluggable logger interface

## 🚀 Quick Start

### Installation

```bash
go get github.com/shellvon/go-sender
```

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/email"
    "github.com/shellvon/go-sender/circuitbreaker"
)

func main() {
    // Create sender instance
    sender := gosender.NewSender(nil)

    // Configure email provider with multiple accounts
    emailConfig := email.Config{
        Accounts: []email.Account{
            {
                Name:     "primary",
                Host:     "smtp.gmail.com",
                Port:     587,
                Username: "primary@gmail.com",
                Password: "password",
                From:     "primary@gmail.com",
                Weight:   2, // Higher weight for primary account
            },
            {
                Name:     "backup",
                Host:     "smtp.outlook.com",
                Port:     587,
                Username: "backup@outlook.com",
                Password: "password",
                From:     "backup@outlook.com",
                Weight:   1, // Lower weight for backup account
            },
        },
        Strategy: "weighted", // Use weighted strategy
    }

    emailProvider, err := email.New(emailConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Register provider
    sender.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

    // Send message
    ctx := context.Background()
    emailMsg := &email.Message{
        To:      []string{"recipient@example.com"},
        Subject: "Hello from Go-Sender",
        Body:    "This is a test message",
    }

    err = sender.Send(ctx, emailMsg)
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }

    // Circuit breaker
    circuitBreaker := circuitbreaker.NewMemoryCircuitBreaker(
        "email-provider",
        5,                    // maxFailures
        30*time.Second,       // resetTimeout
    )
    sender.SetCircuitBreaker(circuitBreaker)

    defer sender.Close()
}
```

## 🔧 Advanced Features

### 1. Custom Retry Policies

```go
// Disable retry for specific message (method 1: set MaxAttempts to 0)
noRetryPolicy := core.NewRetryPolicy(core.WithRetryMaxAttempts(0))
err := sender.Send(ctx, message, core.WithSendRetryPolicy(noRetryPolicy))

// Disable retry for specific message (method 2: no retry policy)
err := sender.Send(ctx, message) // No retry if no global policy is set

// Custom retry policy
retryPolicy := core.NewRetryPolicy(
    core.WithRetryMaxAttempts(5),
    core.WithRetryInitialDelay(time.Second),
    core.WithRetryBackoffFactor(2.0),
    core.WithRetryFilter(func(attempt int, err error) bool {
        // Only retry on network errors
        return strings.Contains(err.Error(), "connection")
    }),
)

// Set global retry policy
sender.SetRetryPolicy(retryPolicy)

// Or use per-message retry policy (overrides global)
err := sender.Send(ctx, message, core.WithSendRetryPolicy(retryPolicy))
```

### 2. Multi-Instance with Load Balancing

```go
// WeCom Bot with multiple instances
wecomConfig := wecombot.Config{
    Bots: []wecombot.Bot{
        {
            Name:     "bot1",
            WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=key1",
            Weight:   2,
        },
        {
            Name:     "bot2",
            WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=key2",
            Weight:   1,
        },
    },
    Strategy: "weighted", // or "round_robin", "random"
    // Note: "health_based" strategy requires custom HealthChecker setup
}

// Webhook with multiple endpoints
webhookConfig := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:     "primary",
            URL:      "https://api1.example.com/webhook",
            Weight:   3,
        },
        {
            Name:     "backup",
            URL:      "https://api2.example.com/webhook",
            Weight:   1,
        },
    },
    Strategy: "weighted",
}
```

### 3. Queue with Callbacks

```go
// Set up in-memory queue
queue := queue.NewMemoryQueue[*core.QueueItem](1000)
sender.SetQueue(queue)

// Send with callback
err := sender.Send(ctx, message,
    core.WithSendAsync(),
    core.WithSendCallback(func(err error) {
        if err != nil {
            log.Printf("Message failed: %v", err)
        } else {
            log.Printf("Message sent successfully")
        }
    }),
)
```

### 4. Circuit Breaker & Rate Limiting

```go
// Circuit breaker
circuitBreaker := circuitbreaker.NewMemoryCircuitBreaker(
    "email-provider",
    5,                    // maxFailures
    30*time.Second,       // resetTimeout
)
sender.SetCircuitBreaker(circuitBreaker)

// Rate limiter
rateLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 20) // 10 QPS, burst 20
sender.SetRateLimiter(rateLimiter)
```

### 5. Health Monitoring

```go
// Check system health
health := sender.HealthCheck(ctx)
if health.Status != core.HealthStatusHealthy {
    log.Printf("System unhealthy: %+v", health)

    // Check specific provider
    if providerHealth, exists := health.Providers[core.ProviderTypeEmail]; exists {
        log.Printf("Email provider: %s", providerHealth.Status)
    }
}
```

## 🔌 Extending Go-Sender

### Adding a New Provider

```go
type MyProvider struct{}

func (p *MyProvider) Send(ctx context.Context, msg core.Message) error {
    // Your implementation
    return nil
}

func (p *MyProvider) Name() string {
    return "my-provider"
}

// Register your provider
sender.RegisterProvider("my-provider", &MyProvider{}, nil)
```

### Adding Custom Middleware

```go
type MyMiddleware struct{}

func (m *MyMiddleware) Execute(ctx context.Context, fn func() error) error {
    // Pre-processing
    log.Println("Before sending")

    err := fn()

    // Post-processing
    log.Println("After sending")

    return err
}
```

## 📊 Supported Strategies

| Strategy       | Description                 | Use Case               |
| -------------- | --------------------------- | ---------------------- |
| `round_robin`  | Distributes requests evenly | Load balancing         |
| `random`       | Random selection            | Simple distribution    |
| `weighted`     | Weight-based selection      | Priority-based routing |
| `health_based` | Health-aware selection      | Custom health checks   |

## 🏗️ Architecture Benefits

### Decorator Pattern

- **ProviderDecorator** wraps basic providers with middleware
- Each middleware can be enabled/disabled independently
- Easy to add new middleware without changing existing code

### Plugin Architecture

- All providers implement the same interface
- Easy to swap implementations
- No vendor lock-in

### Strategy Pattern

- Multiple selection strategies for load balancing
- Context-aware strategy overrides
- Easy to add new strategies

## 📈 Performance & Reliability

- **High Throughput**: Efficient queue processing
- **Low Latency**: Optimized middleware chain
- **Fault Tolerance**: Circuit breaker prevents cascading failures
- **Observability**: Built-in metrics and health checks
- **Scalability**: Support for distributed queues

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

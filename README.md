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

#### Currently Supported Providers

- **📧 Email**: SMTP with multi-account support
- **📱 SMS**: Multi-platform SMS support

  - **Tencent Cloud SMS**: [Official Documentation](https://cloud.tencent.com/document/product/382) | [Official Site](https://cloud.tencent.com/product/sms)
  - **Alibaba Cloud SMS**: [Official Documentation](https://help.aliyun.com/document_detail/101300.html) | [Official Site](https://www.aliyun.com/product/sms)
  - **Huawei Cloud SMS**: [Official Documentation](https://support.huaweicloud.com/sms/index.html) | [Official Site](https://www.huaweicloud.com/product/sms.html)
  - **NetEase Cloud SMS**: [Official Documentation](https://dev.yunxin.163.com/docs/product/短信服务) | [Official Site](https://www.163yun.com/product/sms)
  - **Yunpian SMS**: [Official Documentation](https://www.yunpian.com/doc/zh_CN/api/single_send.html) | [Official Site](https://www.yunpian.com/)
  - **UCP SMS**: [Official Documentation](https://www.ucpaas.com/doc/) | [Official Site](https://www.ucpaas.com/)
  - **CL253 SMS**: [Official Documentation](http://www.253.com/) | [Official Site](http://www.253.com/)
  - **SMSBao**: [Official Documentation](https://www.smsbao.com/openapi/) | [Official Site](https://www.smsbao.com/)
  - **Juhe SMS**: [Official Documentation](https://www.juhe.cn/docs/api/sms) | [Official Site](https://www.juhe.cn/)
  - **Luosimao SMS**: [Official Documentation](https://luosimao.com/docs/api/) | [Official Site](https://luosimao.com/)

  > **Note**: SMS provider implementations are based on code from the [smsBomb](https://github.com/shellvon/smsBomb) project, translated to Go using AI. Not all platforms have been individually tested.

- **🤖 WeCom Bot**: Enterprise WeChat bot messages | [Official Documentation](https://developer.work.weixin.qq.com/document/path/91770)
- **🔔 DingTalk Bot**: DingTalk group bot messages | [Official Documentation](https://open.dingtalk.com/document/robots/custom-robot-access)
- **📢 Lark/Feishu**: Lark (International) and Feishu (China) bot messages | [Official Documentation](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- **💬 Slack**: Slack bot messages | [Official Documentation](https://api.slack.com/messaging/webhooks)
- **📨 ServerChan**: ServerChan push service | [Official Site](https://sct.ftqq.com/)
- **📱 Telegram**: Telegram Bot messages | [Official Documentation](https://core.telegram.org/bots/api)
- **🔗 Webhook**: Generic HTTP webhook calls

### 🛡️ Advanced Reliability Features

- **Smart Retry**: Configurable retry policies with exponential backoff
- **Circuit Breaker**: Prevents cascading failures
- **Rate Limiting**: Token bucket and sliding window algorithms
- **Queue Support**: In-memory and distributed queues
- **Health Checks**: Comprehensive health monitoring

### 🎛️ Multi-Instance & Strategy Support

- **Multiple Accounts**: Support multiple email accounts, bots, webhook endpoints
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

    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

func main() {
    // Create sender instance
    sender := gosender.NewSender(nil)

    // Configure email provider
    emailConfig := email.Config{
        BaseConfig: core.BaseConfig{
            Strategy: core.StrategyRoundRobin,
        },
        Accounts: []email.Account{
            {
                Name:     "primary",
                Host:     "smtp.gmail.com",
                Port:     587,
                Username: "your-email@gmail.com",
                Password: "your-password",
                From:     "your-email@gmail.com",
                Weight:   1,
            },
        },
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

    defer sender.Close()
}
```

## 🔧 Advanced Features

### 1. Custom Retry Policies

```go
// Set global retry policy
retryPolicy := core.NewRetryPolicy(
    core.WithRetryMaxAttempts(5),
    core.WithRetryInitialDelay(time.Second),
    core.WithRetryBackoffFactor(2.0),
)
sender.SetRetryPolicy(retryPolicy)

// Or use per-message retry policy (overrides global)
err := sender.Send(ctx, message, core.WithSendRetryPolicy(retryPolicy))
```

### 2. Multi-Instance with Load Balancing

```go
// WeCom Bot with multiple instances
wecomConfig := wecombot.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyWeighted,
    },
    Accounts: []core.Account{
        {
            Name:     "bot1",
            Key:      "YOUR_KEY_1",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "bot2",
            Key:      "YOUR_KEY_2",
            Weight:   80,
            Disabled: false,
        },
    },
}
```

### 3. Queue and Async Sending

```go
// Set memory queue
queue := queue.NewMemoryQueue[*core.QueueItem](1000)
sender.SetQueue(queue)

// Send message with callback
err := sender.Send(ctx, message,
    core.WithSendAsync(),
    core.WithSendCallback(func(err error) {
        if err != nil {
            log.Printf("Message send failed: %v", err)
        } else {
            log.Printf("Message sent successfully")
        }
    }),
)
```

### 4. Circuit Breaker and Rate Limiting

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
        log.Printf("Email provider status: %s", providerHealth.Status)
    }
}
```

## 🎯 Extending Go-Sender

### Adding New Providers

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

## 📊 Supported Strategies

| Strategy       | Description                | Use Case               |
| -------------- | -------------------------- | ---------------------- |
| `round_robin`  | Distribute requests evenly | Load balancing         |
| `random`       | Random selection           | Simple distribution    |
| `weighted`     | Weight-based selection     | Priority-based routing |
| `health_based` | Health-based selection     | Custom health checks   |

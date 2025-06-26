# Go-Sender

> ⚠️ **Project Status: In Active Development**
>
> This project is under heavy development. APIs may be unstable and subject to change. Please use with caution in production environments.

English | [中文](./README_CN.md)

A high-performance, extensible Go message sending framework supporting multiple notification channels and rich middleware capabilities.

---

## Design Philosophy

Go-Sender is designed around the **Decorator Pattern** and **Plugin Architecture**, making it easy to add new notification channels or cross-cutting concerns without changing your business logic.

### Core Design Principles

- **🔄 Decoupling**: Business code only cares about sending messages, not how they're delivered
- **🔌 Pluggable**: Easy to add new providers or middleware through interfaces
- **🛡️ Reliability**: Built-in retry, circuit breaker, and rate limiting
- **📊 Observable**: Comprehensive metrics and health checks
- **🧩 Flexible**: Support for multiple instances, strategies, and configurations

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

### 🚦 Supported Providers (Grouped by Type)

### 📱 SMS & Voice

| Provider                  | Website                                        | API Docs                                                                                                                             | Provider Doc                            |
| ------------------------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------- |
| Aliyun (阿里云)           | [aliyun.com](https://www.aliyun.com)           | [API](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms)                                            | [SMS README](./providers/sms/README.md) |
| Tencent Cloud (腾讯云)    | [cloud.tencent.com](https://cloud.tencent.com) | [SMS API](https://cloud.tencent.com/document/product/382/55981) / [Voice API](https://cloud.tencent.com/document/product/1128/51559) | [SMS README](./providers/sms/README.md) |
| Huawei Cloud (华为云)     | [huaweicloud.com](https://www.huaweicloud.com) | [API](https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html)                                                        | [SMS README](./providers/sms/README.md) |
| Volcano Engine (火山引擎) | [volcengine.com](https://www.volcengine.com)   | [API](https://www.volcengine.com/docs/63933)                                                                                         | [SMS README](./providers/sms/README.md) |
| Yunpian (云片)            | [yunpian.com](https://www.yunpian.com)         | [API](https://www.yunpian.com/official/document/sms/zh_CN/domestic_list)                                                             | [SMS README](./providers/sms/README.md) |
| CL253 (创蓝 253)          | [253.com](https://www.253.com)                 | [API](https://www.253.com/api)                                                                                                       | [SMS README](./providers/sms/README.md) |
| Submail (赛邮)            | [mysubmail.com](https://www.mysubmail.com/)    | [API](https://www.mysubmail.com/documents)                                                                                           | [SMS README](./providers/sms/README.md) |
| UCP (云之讯)              | [ucpaas.com](https://www.ucpaas.com)           | [API](http://docs.ucpaas.com)                                                                                                        | [SMS README](./providers/sms/README.md) |
| Juhe (聚合数据)           | [juhe.cn](https://www.juhe.cn)                 | [API](https://www.juhe.cn/docs)                                                                                                      | [SMS README](./providers/sms/README.md) |
| SMSBao (短信宝)           | [smsbao.com](https://www.smsbao.com)           | [API](https://www.smsbao.com/openapi)                                                                                                | [SMS README](./providers/sms/README.md) |
| Yuntongxun (云讯通)       | [yuntongxun.com](https://www.yuntongxun.com)   | [API](https://www.yuntongxun.com/developer-center)                                                                                   | [SMS README](./providers/sms/README.md) |

### 📧 Email

| Provider           | Website                                        | API Docs                                                              | Provider Doc                                |
| ------------------ | ---------------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------- |
| go-mail (SMTP)     | [go-mail](https://github.com/wneessen/go-mail) | [Docs](https://pkg.go.dev/github.com/wneessen/go-mail)                | [Email README](./providers/email/README.md) |
| (Planned) Mailgun  | [mailgun.com](https://www.mailgun.com/)        | [API](https://documentation.mailgun.com/en/latest/api_reference.html) | N/A                                         |
| (Planned) Mailjet  | [mailjet.com](https://www.mailjet.com/)        | [API](https://dev.mailjet.com/email/guides/send-api-v31/)             | N/A                                         |
| (Planned) Mailtrap | [mailtrap.io](https://mailtrap.io/)            | [API](https://api-docs.mailtrap.io/docs)                              | N/A                                         |
| (Planned) Brevo    | [brevo.com](https://www.brevo.com/)            | [API](https://developers.brevo.com/docs)                              | N/A                                         |
| (Planned) Braze    | [braze.com](https://www.braze.com/)            | [API](https://www.braze.com/docs/api/)                                | N/A                                         |

### 🤖 IM/Bot/Enterprise Notification

- [WeCom Bot (企业微信机器人)](https://developer.work.weixin.qq.com/document/path/91770) ([Provider Doc](./providers/wecombot/README.md))
- [DingTalk Bot (钉钉机器人)](https://open.dingtalk.com/document/robots/custom-robot-access) ([Provider Doc](./providers/dingtalk/README.md))
- [Lark/Feishu (飞书/国际版)](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) ([Provider Doc](./providers/lark/README.md))
- [Telegram](https://core.telegram.org/bots/api) ([Provider Doc](./providers/telegram/README.md))
- (Planned) Slack ([API](https://api.slack.com/messaging/webhooks))
- [ServerChan](https://sct.ftqq.com/) ([Provider Doc](./providers/serverchan/README.md))

### 🌐 Universal Push / Webhook

All the following are supported via the [Webhook Provider](./providers/webhook/README.md) (generic HTTP integration):

- [ntfy](https://ntfy.sh/)
- [IFTTT](https://ifttt.com/)
- [Bark](https://github.com/Finb/Bark)
- [PushDeer](https://github.com/easychen/pushdeer)
- [PushPlus](https://pushplus.hxtrip.com/)
- [PushAll](https://pushall.ru/)
- [PushBack](https://pushback.io/)
- [Pushy](https://pushy.me/)
- [Pushbullet](https://www.pushbullet.com/)
- [Gotify](https://gotify.net/)
- [OneBot](https://github.com/botuniverse/onebot)
- [Push](https://push.techulus.com/)
- [Pushjet](https://pushjet.io/)
- [Pushsafer](https://www.pushsafer.com/)
- [Pushover](https://pushover.net/)
- [Simplepush](https://simplepush.io/)
- [Zulip](https://zulip.com/)
- [Mattermost](https://mattermost.com/)
- [Discord](https://discord.com/) (message push supported via webhook; for advanced/interaction features, a dedicated provider is needed)

> See [Webhook Provider documentation](./providers/webhook/README.md) for details and examples of supported push platforms.

### 🚀 Push Providers

| Provider                                 | Website                                                                     | API Docs                                                           | Provider Doc |
| ---------------------------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------ | ------------ |
| (Planned) FCM (Firebase Cloud Messaging) | [firebase.google.com](https://firebase.google.com/products/cloud-messaging) | [API](https://firebase.google.com/docs/cloud-messaging)            | N/A          |
| (Planned) JPush (极光推送)               | [jiguang.cn](https://www.jiguang.cn/)                                       | [API](https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push/) | N/A          |

---

### 🛡️ Advanced Reliability

- Built-in retry, circuit breaker, and rate limiting
- Token bucket and sliding window algorithms
- Health checks and observability

### 🎛️ Multi-Instance & Strategy Support

- Multiple accounts/providers per channel
- Load balancing: round-robin, random, weighted, health-based
- Context-aware strategy override

### 🧩 Middleware & Plugin Architecture

- Rate limiter, circuit breaker, retry, queue, metrics, etc.

### 📊 Observability

- Metrics, tracing, health checks

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

// Send message asynchronously
err := sender.Send(ctx, message, core.WithSendAsync())
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

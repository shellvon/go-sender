# Go-Sender

[English](./README.md) | 中文

一个高性能、可扩展的 Go 消息发送框架，支持多种通知渠道和丰富的中间件功能。

一个灵活的 Go 通知发送库，支持 Webhook、企业微信、邮件等多种渠道。专注于可扩展性和可靠性设计。

## 🎯 设计理念

Go-Sender 基于**装饰器模式**和**插件架构**设计，让你可以轻松添加新的通知渠道或横切关注点，而无需改变业务逻辑。

### 核心设计原则

- **🔄 解耦**: 业务代码只关心发送消息，不关心如何传递
- **🔌 可插拔**: 通过接口轻松添加新的提供者或中间件
- **🛡️ 可靠性**: 内置重试、熔断器和限流机制
- **📊 可观测**: 全面的指标和健康检查
- **⚡ 灵活性**: 支持多实例、策略和配置

### 架构概览

```
业务逻辑 → Sender → ProviderDecorator → Provider
                ↓
          中间件链:
          - 限流器
          - 熔断器
          - 重试策略
          - 队列
          - 指标收集
```

## ✨ 功能特性

### 🚀 多渠道支持

- **邮件**: SMTP 多账号支持
- **企业微信机器人**: 企业微信机器人消息
- **Webhook**: 通用 HTTP webhook 调用
- **可扩展**: 轻松添加 Telegram、Slack、Discord 等

### 🛡️ 高级可靠性功能

- **智能重试**: 可配置的重试策略，支持指数退避
- **熔断器**: 防止级联故障
- **限流**: 令牌桶和滑动窗口算法
- **队列支持**: 内存队列和分布式队列
- **健康检查**: 全面的健康监控

### 🎛️ 多实例和策略支持

- **多账号**: 支持多个邮件账号、企业微信机器人、webhook 端点
- **负载均衡**: 轮询、随机、权重和基于健康状态的策略
- **上下文感知**: 通过上下文覆盖每个请求的策略

### 📊 可观测性

- **指标收集**: 性能和结果指标
- **健康监控**: 提供者和系统健康检查
- **结构化日志**: 可插拔的日志接口

## 🚀 快速开始

### 安装

```bash
go get github.com/shellvon/go-sender
```

### 基本使用

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
    // 创建sender实例
    sender := gosender.NewSender(nil)

    // 配置邮件提供者，支持多账号
    emailConfig := email.Config{
        Accounts: []email.Account{
            {
                Name:     "primary",
                Host:     "smtp.gmail.com",
                Port:     587,
                Username: "primary@gmail.com",
                Password: "password",
                From:     "primary@gmail.com",
                Weight:   2, // 主账号权重更高
            },
            {
                Name:     "backup",
                Host:     "smtp.outlook.com",
                Port:     587,
                Username: "backup@outlook.com",
                Password: "password",
                From:     "backup@outlook.com",
                Weight:   1, // 备用账号权重较低
            },
        },
        Strategy: "weighted", // 使用权重策略
    }

    emailProvider, err := email.New(emailConfig)
    if err != nil {
        log.Fatal(err)
    }

    // 注册提供者
    sender.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

    // 设置熔断器
    circuitBreaker := circuitbreaker.NewMemoryCircuitBreaker(
        "email-provider",
        5,                    // maxFailures
        30*time.Second,       // resetTimeout
    )
    sender.SetCircuitBreaker(circuitBreaker)

    // 发送消息
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

## 🔧 高级功能

### 1. 自定义重试策略

```go
// 禁用特定消息的重试（方法1：设置MaxAttempts为0）
noRetryPolicy := core.NewRetryPolicy(core.WithRetryMaxAttempts(0))
err := sender.Send(ctx, message, core.WithSendRetryPolicy(noRetryPolicy))

// 禁用特定消息的重试（方法2：不设置重试策略）
err := sender.Send(ctx, message) // 如果没有全局重试策略，就不会重试

// 自定义重试策略
retryPolicy := core.NewRetryPolicy(
    core.WithRetryMaxAttempts(5),
    core.WithRetryInitialDelay(time.Second),
    core.WithRetryBackoffFactor(2.0),
    core.WithRetryFilter(func(attempt int, err error) bool {
        // 只对网络错误重试
        return strings.Contains(err.Error(), "connection")
    }),
)

// 设置全局重试策略
sender.SetRetryPolicy(retryPolicy)

// 或使用每条消息的重试策略（覆盖全局策略）
err := sender.Send(ctx, message, core.WithSendRetryPolicy(retryPolicy))
```

### 2. 多实例负载均衡

```go
// 企业微信机器人多实例
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
    Strategy: "weighted", // 或 "round_robin", "random"
}

// Webhook 多端点
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

### 3. 队列回调

```go
// 设置内存队列
queue := queue.NewMemoryQueue[*core.QueueItem](1000)
sender.SetQueue(queue)

// 发送带回调的消息
err := sender.Send(ctx, message,
    core.WithSendAsync(),
    core.WithSendCallback(func(err error) {
        if err != nil {
            log.Printf("消息发送失败: %v", err)
        } else {
            log.Printf("消息发送成功")
        }
    }),
)
```

### 4. 熔断器和限流

```go
// 熔断器
circuitBreaker := circuitbreaker.NewMemoryCircuitBreaker(
    "email-provider",
    5,                    // maxFailures
    30*time.Second,       // resetTimeout
)
sender.SetCircuitBreaker(circuitBreaker)

// 限流器
rateLimiter := ratelimiter.NewTokenBucketRateLimiter(10, 20) // 10 QPS, 突发 20
sender.SetRateLimiter(rateLimiter)
```

### 5. 健康监控

```go
// 检查系统健康状态
health := sender.HealthCheck(ctx)
if health.Status != core.HealthStatusHealthy {
    log.Printf("系统不健康: %+v", health)

    // 检查特定提供者
    if providerHealth, exists := health.Providers[core.ProviderTypeEmail]; exists {
        log.Printf("邮件提供者状态: %s", providerHealth.Status)
    }
}
```

## 🎯 扩展 Go-Sender

### 通过 Webhook 实现其他渠道

虽然当前版本没有直接支持 Telegram、飞书等渠道，但你可以通过 webhook 提供者轻松实现。由于 webhook 消息的 ProviderType 固定为 "webhook"，需要直接使用对应的 webhook 提供者实例：

#### Telegram Bot 示例

```go
// 创建 Telegram webhook 配置
telegramConfig := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:    "telegram-bot",
            URL:     "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage",
            Method:  "POST", // 请求方法来自 Endpoint
            Headers: map[string]string{
                "Content-Type": "application/json", // 固定头部
            },
        },
    },
    Strategy: "round_robin",
}

// 创建 Telegram webhook 提供者
telegramProvider, _ := webhook.New(telegramConfig)

// 创建 Telegram 消息（只支持 JSON 格式）
telegramMsg := &webhook.Message{
    EndpointName: "telegram-bot",
    Body: map[string]interface{}{
        "chat_id":    "@your_channel",
        "text":       "Hello from Go-Sender!",
        "parse_mode": "Markdown",
    },
    // Headers 字段可选，会与 Endpoint 的 Headers 合并
}

// 直接使用提供者发送消息
err := telegramProvider.Send(ctx, telegramMsg)
```

#### 飞书 Webhook 示例

```go
// 创建飞书 webhook 配置
feishuConfig := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:    "feishu-webhook",
            URL:     "https://open.feishu.cn/open-apis/bot/v2/hook/<YOUR_WEBHOOK_TOKEN>",
            Method:  "POST",
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        },
    },
    Strategy: "round_robin",
}

// 创建飞书 webhook 提供者
feishuProvider, _ := webhook.New(feishuConfig)

// 创建飞书消息（只支持 JSON 格式）
feishuMsg := &webhook.Message{
    EndpointName: "feishu-webhook",
    Body: map[string]interface{}{
        "msg_type": "text",
        "content": map[string]interface{}{
            "text": "Hello from Go-Sender!",
        },
    },
}

// 直接使用提供者发送消息
err := feishuProvider.Send(ctx, feishuMsg)
```

#### 钉钉 Webhook 示例

```go
// 创建钉钉 webhook 配置
dingtalkConfig := webhook.Config{
    Endpoints: []webhook.Endpoint{
        {
            Name:    "dingtalk-webhook",
            URL:     "https://oapi.dingtalk.com/robot/send?access_token=<YOUR_ACCESS_TOKEN>",
            Method:  "POST",
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
        },
    },
    Strategy: "round_robin",
}

// 创建钉钉 webhook 提供者
dingtalkProvider, _ := webhook.New(dingtalkConfig)

// 创建钉钉消息（只支持 JSON 格式）
dingtalkMsg := &webhook.Message{
    EndpointName: "dingtalk-webhook",
    Body: map[string]interface{}{
        "msgtype": "text",
        "text": map[string]interface{}{
            "content": "Hello from Go-Sender!",
        },
    },
}

// 直接使用提供者发送消息
err := dingtalkProvider.Send(ctx, dingtalkMsg)
```

#### 使用 Sender 的统一接口（需要自定义消息类型）

如果你想使用 `sender.Send()` 的统一接口，需要创建自定义的消息类型：

```go
// 自定义钉钉消息类型
type DingTalkMessage struct {
    webhook.Message
}

func (m *DingTalkMessage) ProviderType() core.ProviderType {
    return "dingtalk" // 返回自定义的提供者类型
}

// 注册钉钉提供者
sender.RegisterProvider("dingtalk", dingtalkProvider, nil)

// 创建自定义钉钉消息
dingtalkMsg := &DingTalkMessage{
    Message: webhook.Message{
        EndpointName: "dingtalk-webhook",
        Body: map[string]interface{}{
            "msgtype": "text",
            "text": map[string]interface{}{
                "content": "Hello from Go-Sender!",
            },
        },
    },
}

// 使用统一接口发送
err := sender.Send(ctx, dingtalkMsg)
```

#### 重要说明

1. **请求方法**：来自 `webhook.Endpoint.Method`，不是 `webhook.Message`
2. **内容类型**：只支持 `application/json`，Body 会被自动序列化为 JSON
3. **头部合并**：Message 的 Headers 会覆盖 Endpoint 的 Headers
4. **查询参数**：支持在 Endpoint 和 Message 中配置，会自动合并

### 添加新的提供者

```go
type MyProvider struct{}

func (p *MyProvider) Send(ctx context.Context, msg core.Message) error {
    // 你的实现
    return nil
}

func (p *MyProvider) Name() string {
    return "my-provider"
}

// 注册你的提供者
sender.RegisterProvider("my-provider", &MyProvider{}, nil)
```

### 添加自定义中间件

```go
type MyMiddleware struct{}

func (m *MyMiddleware) Execute(ctx context.Context, fn func() error) error {
    // 预处理
    log.Println("发送前")

    err := fn()

    // 后处理
    log.Println("发送后")

    return err
}
```

## 📊 支持的策略

| 策略           | 描述         | 使用场景         |
| -------------- | ------------ | ---------------- |
| `round_robin`  | 均匀分配请求 | 负载均衡         |
| `random`       | 随机选择     | 简单分发         |
| `weighted`     | 基于权重选择 | 基于优先级的路由 |
| `health_based` | 基于健康状态 | 自定义健康检查   |

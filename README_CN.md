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

- **📧 邮件**: SMTP 多账号支持
- **📱 短信**: 多平台短信支持

  - **腾讯云短信**: [官方文档](https://cloud.tencent.com/document/product/382) | [官方网站](https://cloud.tencent.com/product/sms)
  - **阿里云短信**: [官方文档](https://help.aliyun.com/document_detail/101300.html) | [官方网站](https://www.aliyun.com/product/sms)
  - **华为云短信**: [官方文档](https://support.huaweicloud.com/sms/index.html) | [官方网站](https://www.huaweicloud.com/product/sms.html)
  - **网易云短信**: [官方文档](https://dev.yunxin.163.com/docs/product/短信服务) | [官方网站](https://www.163yun.com/product/sms)
  - **云片网**: [官方文档](https://www.yunpian.com/doc/zh_CN/api/single_send.html) | [官方网站](https://www.yunpian.com/)
  - **云之讯**: [官方文档](https://www.ucpaas.com/doc/) | [官方网站](https://www.ucpaas.com/)
  - **蓝创 253**: [官方文档](http://www.253.com/) | [官方网站](http://www.253.com/)
  - **短信宝**: [官方文档](https://www.smsbao.com/openapi/) | [官方网站](https://www.smsbao.com/)
  - **聚合服务**: [官方文档](https://www.juhe.cn/docs/api/sms) | [官方网站](https://www.juhe.cn/)
  - **螺丝帽**: [官方文档](https://luosimao.com/docs/api/) | [官方网站](https://luosimao.com/)

  > **注意**: 短信提供者实现基于 [smsBomb](https://github.com/shellvon/smsBomb) 项目代码，通过 AI 翻译到 Go 语言。并非所有平台都经过单独测试。

- **🤖 企业微信机器人**: 企业微信机器人消息 | [官方文档](https://developer.work.weixin.qq.com/document/path/91770)
- **🔔 钉钉机器人**: 钉钉群机器人消息 | [官方文档](https://open.dingtalk.com/document/robots/custom-robot-access)
- **📢 飞书/国际版**: Lark/Feishu 机器人消息 | [官方文档](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- **💬 Slack**: Slack 机器人消息 | [官方文档](https://api.slack.com/messaging/webhooks)
- **📨 ServerChan**: ServerChan 推送服务 | [官方网站](https://sct.ftqq.com/)
- **📱 Telegram**: Telegram Bot 消息 | [官方文档](https://core.telegram.org/bots/api)
- **🔗 Webhook**: 通用 HTTP webhook 调用

### 🛡️ 高级可靠性功能

- **智能重试**: 可配置的重试策略，支持指数退避
- **熔断器**: 防止级联故障
- **限流**: 令牌桶和滑动窗口算法
- **队列支持**: 内存队列和分布式队列
- **健康检查**: 全面的健康监控

### 🎛️ 多实例和策略支持

- **多账号**: 支持多个邮件账号、机器人、webhook 端点
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

    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

func main() {
    // 创建sender实例
    sender := gosender.NewSender(nil)

    // 配置邮件提供者
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

    // 注册提供者
    sender.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

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
// 设置全局重试策略
retryPolicy := core.NewRetryPolicy(
    core.WithRetryMaxAttempts(5),
    core.WithRetryInitialDelay(time.Second),
    core.WithRetryBackoffFactor(2.0),
)
sender.SetRetryPolicy(retryPolicy)

// 或使用每条消息的重试策略（覆盖全局策略）
err := sender.Send(ctx, message, core.WithSendRetryPolicy(retryPolicy))
```

### 2. 多实例负载均衡

```go
// 企业微信机器人多实例
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

### 3. 队列和异步发送

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

## 📊 支持的策略

| 策略           | 描述         | 使用场景         |
| -------------- | ------------ | ---------------- |
| `round_robin`  | 均匀分配请求 | 负载均衡         |
| `random`       | 随机选择     | 简单分发         |
| `weighted`     | 基于权重选择 | 基于优先级的路由 |
| `health_based` | 基于健康状态 | 自定义健康检查   |

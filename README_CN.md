# Go-Sender

> ⚠️ **项目状态：开发中**
>
> 本项目正在积极开发中，API 可能不稳定且随时变动。请在生产环境中谨慎使用。

[English](./README.md) | 中文

一个高性能、可扩展的 Go 消息发送框架，支持多种通知渠道和丰富的中间件功能。

一个灵活的 Go 通知发送库，支持 Webhook、企业微信、邮件等多种渠道。专注于可扩展性和可靠性设计。

## 🚦 支持的通道（按类型分组）

### 📱 短信/语音

| 提供商   | 官网                                           | API 文档                                                                                                                             | Provider 文档                           |
| -------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------- |
| 阿里云   | [aliyun.com](https://www.aliyun.com)           | [API](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms)                                            | [SMS README](./providers/sms/README.md) |
| 腾讯云   | [cloud.tencent.com](https://cloud.tencent.com) | [短信 API](https://cloud.tencent.com/document/product/382/55981) / [语音 API](https://cloud.tencent.com/document/product/1128/51559) | [SMS README](./providers/sms/README.md) |
| 华为云   | [huaweicloud.com](https://www.huaweicloud.com) | [API](https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html)                                                        | [SMS README](./providers/sms/README.md) |
| 火山引擎 | [volcengine.com](https://www.volcengine.com)   | [API](https://www.volcengine.com/docs/63933)                                                                                         | [SMS README](./providers/sms/README.md) |
| 云片     | [yunpian.com](https://www.yunpian.com)         | [API](https://www.yunpian.com/official/document/sms/zh_CN/domestic_list)                                                             | [SMS README](./providers/sms/README.md) |
| 创蓝 253 | [253.com](https://www.253.com)                 | [API](https://www.253.com/api)                                                                                                       | [SMS README](./providers/sms/README.md) |
| 赛邮     | [mysubmail.com](https://www.mysubmail.com/)    | [API](https://www.mysubmail.com/documents)                                                                                           | [SMS README](./providers/sms/README.md) |
| 云之讯   | [ucpaas.com](https://www.ucpaas.com)           | [API](http://docs.ucpaas.com)                                                                                                        | [SMS README](./providers/sms/README.md) |
| 聚合数据 | [juhe.cn](https://www.juhe.cn)                 | [API](https://www.juhe.cn/docs)                                                                                                      | [SMS README](./providers/sms/README.md) |
| 短信宝   | [smsbao.com](https://www.smsbao.com)           | [API](https://www.smsbao.com/openapi)                                                                                                | [SMS README](./providers/sms/README.md) |
| 云讯通   | [yuntongxun.com](https://www.yuntongxun.com)   | [API](https://www.yuntongxun.com/developer-center)                                                                                   | [SMS README](./providers/sms/README.md) |

> 详细短信/彩信/语音能力请见 [能力矩阵](./providers/sms/capabilities.md)。

### 📧 邮件

| 提供方             | 官网                                           | API 文档                                                              | Provider 文档                               | 状态   |
| ------------------ | ---------------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------- | ------ |
| go-mail (SMTP)     | [go-mail](https://github.com/wneessen/go-mail) | [Docs](https://pkg.go.dev/github.com/wneessen/go-mail)                | [Email README](./providers/email/README.md) | 已实现 |
| EmailJS (API)      | [emailjs.com](https://www.emailjs.com/)        | [API](https://www.emailjs.com/docs/rest-api/send/)                    | [emailapi](./providers/emailapi/README.md)  | 已实现 |
| Resend (API)       | [resend.com](https://resend.com/)              | [API](https://resend.com/docs/api-reference/emails/send-batch-emails) | [emailapi](./providers/emailapi/README.md)  | 已实现 |
| （计划）Mailgun    | [mailgun.com](https://www.mailgun.com/)        | [API](https://documentation.mailgun.com/en/latest/api_reference.html) | N/A                                         | 计划中 |
| （计划）Mailjet    | [mailjet.com](https://www.mailjet.com/)        | [API](https://dev.mailjet.com/email/guides/send-api-v31/)             | N/A                                         | 计划中 |
| （计划）Brevo      | [brevo.com](https://www.brevo.com/)            | [API](https://developers.brevo.com/docs)                              | N/A                                         | 计划中 |
| （计划）Mailersend | [mailersend.com](https://www.mailersend.com/)  | [API](https://developers.mailersend.com/)                             | N/A                                         | 计划中 |
| （计划）Mailtrap   | [mailtrap.io](https://mailtrap.io/)            | [API](https://api-docs.mailtrap.io/docs)                              | N/A                                         | 计划中 |

> **注意：** `emailapi` 类型为实验性特性，API 可能随时变更。

### 🤖 IM/Bot/企业通知

- [企业微信机器人](https://developer.work.weixin.qq.com/document/path/91770) ([Provider 文档](./providers/wecombot/README.md))
- [钉钉机器人](https://open.dingtalk.com/document/robots/custom-robot-access) ([Provider 文档](./providers/dingtalk/README.md))
- [飞书/Lark](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN) ([Provider 文档](./providers/lark/README.md))
- [Telegram](https://core.telegram.org/bots/api) ([Provider 文档](./providers/telegram/README.md))
- （计划）Slack（[API](https://api.slack.com/messaging/webhooks)）
- [Server 酱](https://sct.ftqq.com/) ([Provider 文档](./providers/serverchan/README.md))

### 🌐 通用推送 / Webhook

以下所有平台均通过 [Webhook Provider](./providers/webhook/README.md)（通用 HTTP 集成）支持：

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
- [Discord](https://discord.com/)（仅支持消息推送，交互/事件需专用 provider）

> 详见 [Webhook Provider 文档](./providers/webhook/README.md)，了解已支持的推送平台和用法示例。

### 🚀 推送服务

| 推送服务                                | 官网                                                                        | API 文档                                                           | Provider 文档 |
| --------------------------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------ | ------------- |
| （计划）FCM（Firebase Cloud Messaging） | [firebase.google.com](https://firebase.google.com/products/cloud-messaging) | [API](https://firebase.google.com/docs/cloud-messaging)            | N/A           |
| （计划）极光推送（JPush）               | [jiguang.cn](https://www.jiguang.cn/)                                       | [API](https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push/) | N/A           |

## 🎯 设计理念

Go-Sender 基于**装饰器模式**和**插件架构**设计，让你可以轻松添加新的通知渠道或横切关注点，而无需改变业务逻辑。

### 核心设计原则

- **🔄 解耦**: 业务代码只关心发送消息，不关心如何传递
- **🔌 可插拔**: 通过接口轻松添加新的提供者或中间件
- **🛡️ 可靠性**: 内置重试、熔断器和限流机制
- **🧩 灵活性**: 支持多实例、策略和配置

### 架构概览

```
业务逻辑 → Sender → ProviderDecorator → Provider
                ↓
          中间件链:
          - 限流器
          - 熔断器
          - 重试策略
          - 队列
```

## ✨ 功能特性

### 🚀 多渠道支持

#### 当前支持的提供者

- **📧 邮件**: 使用 [wneessen/go-mail](https://github.com/wneessen/go-mail) 的 SMTP 多账号支持（[通道文档](./providers/email/README.md)）
- **📱 短信**: 多平台短信支持（[通道文档](./providers/sms/README.md)）

  - **Aliyun SMS (阿里云, 中国大陆)**: [官方文档](https://help.aliyun.com/document_detail/419273.html)（[通道文档](./providers/sms/README.md)）
  - **Aliyun Intl SMS (阿里云国际)**: [官方文档](https://help.aliyun.com/document_detail/108146.html)（[通道文档](./providers/sms/README.md)）
  - **Huawei Cloud SMS (华为云)**: [官方文档](https://support.huaweicloud.com/sms/index.html)（[通道文档](./providers/sms/README.md)）
  - **Luosimao (螺丝帽)**: [官方文档](https://luosimao.com/docs/api/)（[通道文档](./providers/sms/README.md)）
  - **CL253 (创蓝 253)**: [官方文档](http://www.253.com/)（[通道文档](./providers/sms/README.md)）
  - **Juhe (聚合数据)**: [官方文档](https://www.juhe.cn/docs/api/id/54)（[通道文档](./providers/sms/README.md)）
  - **SMSBao (短信宝)**: [官方文档](https://www.smsbao.com/openapi/213.html)（[通道文档](./providers/sms/README.md)）
  - **UCP (云之讯)**: [官方文档](https://doc.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sms:index)（[通道文档](./providers/sms/README.md)）
  - **Tencent Cloud SMS (腾讯云)** (开发中)（[通道文档](./providers/sms/README.md)）
  - **Yunpian (云片)** (开发中)（[通道文档](./providers/sms/README.md)）
  - **Submail (赛邮)** (开发中)（[通道文档](./providers/sms/README.md)）
  - **Volcano Engine (火山引擎)** (开发中)（[通道文档](./providers/sms/README.md)）

- **🤖 企业微信机器人**: 企业微信机器人消息（[通道文档](./providers/wecombot/README.md)） | [官方文档](https://developer.work.weixin.qq.com/document/path/91770)
- **🔔 钉钉机器人**: 钉钉群机器人消息（[通道文档](./providers/dingtalk/README.md)） | [官方文档](https://open.dingtalk.com/document/robots/custom-robot-access)
- **📢 飞书/国际版**: Lark/Feishu 机器人消息（[通道文档](./providers/lark/README.md)） | [官方文档](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- **💬 Slack**: Slack 机器人消息 | [官方文档](https://api.slack.com/messaging/webhooks)
- **📨 ServerChan**: ServerChan 推送服务（[通道文档](./providers/serverchan/README.md)） | [官方网站](https://sct.ftqq.com/)
- **📱 Telegram**: Telegram Bot 消息（[通道文档](./providers/telegram/README.md)） | [官方文档](https://core.telegram.org/bots/api)
- **🔗 Webhook**: 通用 HTTP webhook 调用（[通道文档](./providers/webhook/README.md)）

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

// 异步发送消息
err := sender.Send(ctx, message, core.WithSendAsync())
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

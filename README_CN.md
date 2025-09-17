# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)
[![Go Report Card](https://goreportcard.com/badge/github.com/shellvon/go-sender)](https://goreportcard.com/report/github.com/shellvon/go-sender)
[![GoDoc](https://godoc.org/github.com/shellvon/go-sender?status.svg)](https://pkg.go.dev/github.com/shellvon/go-sender)

> Send anything, anywhere, with Go. One API, All Providers 🚀

[English](./README.md) | **简体中文**

**一个统一的多渠道通知系统**，支持短信、邮件、IM、Webhook 等，具备重试、限流、熔断等企业级特性。

---

## 🌟 为什么选择 go-sender？

### Go 的优势：
- **🪶 轻量级**：纯 Go 实现，零臃肿，依赖极少
- **⚡ 高性能**：协程并发，内存占用小
- **🔒 类型安全**：编译时检查，运行时稳定
- **📦 简单部署**：单一二进制文件，容器友好

### 架构优势：
- **🎯 自动路由**：消息类型自动选择对应 Provider
- **🏗️ 渐进式**：从简单脚本到企业级应用，API 不变
- **🛡️ 生产就绪**：内置重试、限流、熔断、多账号故障转移
- **🧩 高扩展性**：约 50 行代码即可实现自定义 Provider

---

## ⚡ 快速开始

### 方法 1：直接使用 Provider（最简单）

无需任何配置即可发送消息：

```go
import (
    "context"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // 创建 Provider 并发送消息
    account := wecombot.NewAccount("your-webhook-key")
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    
    msg := wecombot.Text().Content("Hello from go-sender!").Build()
    provider.Send(context.Background(), msg, nil)
}
```

### 方法 2：使用 Sender（支持中间件）

需要重试、限流等高级功能？使用 4 步模式：

```go
package main

import (
    "context"
    "log"
    
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // 1️⃣ 初始化 Sender 实例（可稍后添加中间件）
    sender := gosender.NewSender()

    // 2️⃣ 创建企业微信机器人账号和 Provider
    account := wecombot.NewAccount("your-webhook-key")
    wecomProvider, err := wecombot.NewProvider([]*wecombot.Account{account})
    if err != nil {
        log.Fatalf("创建 Provider 失败: %v", err)
    }
    // 向 Sender 注册（nil = 使用全局中间件设置）
    sender.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

    // 3️⃣ 构造要发送的消息
    msg := wecombot.Text().Content("Hello from go-sender!").Build()

    // 4️⃣ 发送消息并获取详细结果
    _, err = sender.SendWithResult(context.Background(), msg)
    if err != nil {
        log.Fatalf("发送失败: %v", err)
    }
    log.Println("消息发送成功！")
}
```

---

## 🔧 工作原理

go-sender 采用现代化的设计模式：

1. **🎯 自动路由**：任何消息只要实现了 `ProviderType()`，系统会自动分发给对应的 Provider 处理
2. **🔄 装饰器模式**：通过中间件为您增加重试、限流、熔断等策略
3. **⚖️ 多账号策略**：内置轮询、权重、故障转移等账号选择策略
4. **🌐 HTTP 抽象**：现代通知服务大多是 HTTP APIs，而其他使用特定协议（如邮件的 SMTP）

想要重试？队列？限流？我们通过装饰器模式实现切面编程，无需装饰器时可直接使用 Provider 发送。

---

## 📦 更多示例

### 高级特性（中间件）

```go
import (
    "time"
    "github.com/shellvon/go-sender/ratelimiter"
)

// 带重试、限流的生产级配置
middleware := &core.SenderMiddleware{
    RateLimiter: ratelimiter.NewTokenBucketRateLimiter(10, 5), // 10 QPS，突发 5
    Retry: &core.RetryPolicy{
        MaxAttempts: 3,
        InitialDelay: time.Second,
        MaxDelay: 10 * time.Second,
    },
}

sender.RegisterProvider(core.ProviderTypeWecombot, provider, middleware)
```

### 多账号与策略

```go
// 多个账号实现高可用
accounts := []*wecombot.Account{
    wecombot.NewAccount("primary-webhook"),
    wecombot.NewAccount("backup-webhook"),
}

config := &wecombot.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyFailover, // 故障转移策略
    },
    Items: accounts,
}
```

### 复杂认证（企业微信应用）

```go
import "github.com/shellvon/go-sender/providers/wecomapp"

// 自动 OAuth 令牌管理
account := wecomapp.NewAccount("corp-id", "agent-id", "app-secret")
provider, _ := wecomapp.New(&wecomapp.Config{Items: []*wecomapp.Account{account}}, nil)

msg := wecomapp.Text().Content("来自企业应用的消息").Build()
provider.Send(context.Background(), msg, nil)
```

### 子 Provider（短信多厂商）

```go
import "github.com/shellvon/go-sender/providers/sms"

// 同一个短信 Provider 支持多个厂商
aliyunMsg := sms.Aliyun().To("13800138000").Content("阿里云短信").Build()
tencentMsg := sms.Tencent().To("13800138000").Content("腾讯云短信").Build()

// 自动路由到对应的厂商 API
sender.Send(context.Background(), aliyunMsg)  // → 阿里云 API
sender.Send(context.Background(), tencentMsg) // → 腾讯云 API
```

---

## 🛠 支持的 Provider

| Provider | 状态 | 说明 |
|----------|------|------|
| **短信** |
| Aliyun SMS | ✅ | 阿里云短信服务 |
| Tencent SMS | ✅ | 腾讯云短信服务 |
| Huawei SMS | ✅ | 华为云短信服务 |
| Volc SMS | ✅ | 火山引擎短信服务 |
| Yunpian SMS | ✅ | 云片短信服务 |
| **邮件** |
| SMTP | ✅ | 标准 SMTP 协议 |
| EmailJS | ✅ | EmailJS API 服务 |
| Resend | ✅ | Resend API 服务 |
| **IM/机器人** |
| 企业微信机器人 | ✅ | WeCom Bot Webhook |
| 企业微信应用 | ✅ | WeCom App API |
| 钉钉机器人 | ✅ | DingTalk Bot |
| 飞书/Lark | ✅ | Lark/Feishu API |
| Telegram | ✅ | Telegram Bot API |
| **Webhook** |
| 通用 Webhook | ✅ | 支持任意 HTTP API |

[查看完整 Provider 列表 →](./docs/providers.md)

---

## 🛠 找不到您的 Provider？

**没问题！** go-sender 专为扩展性而设计：

### 1. 使用通用 Webhook

```go
// 步骤 1：配置 webhook 端点
endpoint := &webhook.Endpoint{
    Name:    "my-api",
    URL:     "https://api.example.com/send",
    Method:  "POST",
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
        "Content-Type":  "application/json",
    },
}

provider, _ := webhook.New(&webhook.Config{
    Items: []*webhook.Endpoint{endpoint},
})

// 步骤 2：创建并发送消息
msg := webhook.Webhook().
    Body([]byte(`{"message": "Hello World", "recipient": "user123"}`)).
    Build()

provider.Send(context.Background(), msg, nil)
```

### 2. 创建自定义 Provider

构建自定义 Provider 很简单 - 只需使用 `core.BaseMessage` 实现 `core.Message` 接口：

```go
// 定义消息类型
type CustomMessage struct {
    core.BaseMessage  // 自动处理路由
    Content   string `json:"content"`
    Recipient string `json:"recipient"`
}

func (m *CustomMessage) ProviderType() core.ProviderType {
    return "custom_api"  // 这将启用自动路由
}

// 创建 transformer 进行 HTTP 协议转换
// 参考现有 Provider 如 wecombot/、sms/、email/ 的模式
```

**想深入了解？** 研究这些 Provider 实现：
- **简单**：[`providers/wecombot/`](./providers/wecombot/) - 基础 HTTP webhook
- **认证**：[`providers/wecomapp/`](./providers/wecomapp/) - OAuth 与缓存
- **多厂商**：[`providers/sms/`](./providers/sms/) - SubProvider 模式

查看 [docs/advanced.md](./docs/advanced.md) 获取完整的自定义 Provider 指南。

---

## 📚 文档

| 文档 | 说明 |
|------|------|
| [快速入门](./docs/getting-started.md) | 从简单脚本到企业级应用的渐进式指南 |
| [核心概念](./docs/concepts.md) | 理解 go-sender 的架构设计 |
| [Provider 文档](./docs/providers.md) | 所有支持的 Provider 详细说明 |
| [中间件](./docs/middleware.md) | 重试、限流、熔断等高级特性 |
| [高级用法](./docs/advanced.md) | 自定义 Provider、中间件、策略 |
| [示例](./docs/examples.md) | 生产环境使用案例 |
| [故障排除](./docs/troubleshooting.md) | 常见问题与解决方案 |

---

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

## 📄 许可证

[MIT License](./LICENSE)

---

**go-sender** - 让 Go 语言的通知发送变得简单而强大 🚀

# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)
[![Go Report Card](https://goreportcard.com/badge/github.com/shellvon/go-sender)](https://goreportcard.com/report/github.com/shellvon/go-sender)
[![GoDoc](https://godoc.org/github.com/shellvon/go-sender?status.svg)](https://godoc.org/github.com/shellvon/go-sender)

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
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

func main() {
    // 1️⃣ 创建 sender
    sender := gosender.NewSender()
    
    // 2️⃣ 创建账号  
    account := wecombot.NewAccount("your-webhook-key")
    
    // 3️⃣ 注册 provider
    provider, _ := wecombot.NewProvider([]*wecombot.Account{account})
    sender.RegisterProvider(core.ProviderTypeWecombot, provider, nil)
    
    // 4️⃣ 发送消息
    msg := wecombot.Text().Content("Hello from go-sender!").Build()
    sender.Send(context.Background(), msg)
}
```

**就是这样！** 🎉 这个相同的 4 步模式适用于**任何 Provider**。

> 📚 **想了解更多？** 查看我们的[详细指南](./docs/getting-started.md)


## 📦 安装

```bash
go get github.com/shellvon/go-sender
```


---

## ✨ 支持的 Providers

| Provider 类型 | 实现 | 状态 |
|---------------|------|--------|
| **短信** | 阿里云、腾讯云、华为云、云片、CL253、火山引擎等 | ✅ 生产就绪 |
| **邮件** | SMTP、EmailJS、Resend | ✅ 生产就绪 |
| **IM/机器人** | 企业微信、钉钉、飞书、Telegram、ServerChan | ✅ 生产就绪 |
| **Webhook** | 通用 HTTP、自定义 APIs | ✅ 生产就绪 |

查看 [docs/providers.md](docs/providers.md) 获取完整的 provider 列表和配置。

---

## 🛠 找不到您的 Provider？

**没问题！** go-sender 专为扩展性而设计。您有**两种选择**：

### 选择 1：使用通用 Webhook（推荐用于 HTTP APIs）

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

### 选择 2：创建自定义 Provider（用于复杂需求）

对于复杂认证、自定义协议或特殊需求：

```go
// 1. 定义消息类型
type CustomMessage struct {
    *core.BaseMessage
    // 可选：如果需要额外字段支持
    *core.WithExtraFields
    Content string `json:"content"`
}

// 2. 实现 provider 接口
// 查看 docs/advanced.md 获取完整指南
```

**想要完整教程？** 查看 [高级用法：自定义 Providers](./docs/advanced.md#custom-providers)

---

## 📚 文档

| **入门指南** | **高级使用** | **参考文档** |
|-------------|-------------|-------------|
| [📖 快速入门](./docs/getting-started.md) | [🔧 高级用法](./docs/advanced.md) | [🔌 Providers](./docs/providers.md) |
| [💡 核心概念](./docs/concepts.md) | [🧪 示例](./docs/examples.md) | [❓ FAQ](./docs/faq.md) |
| [🏗️ 架构概览](./docs/architecture.md) | [🚦 中间件](./docs/middleware.md) | [🔧 故障排除](./docs/troubleshooting.md) |

**快速导航：**
- 🆕 **新用户？** 从[快速入门](./docs/getting-started.md)开始
- 🔍 **需要特定 Provider？** 查看[Providers](./docs/providers.md)  
- 🛠 **想构建自定义 Provider？** 参考[高级用法](./docs/advanced.md)

---

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

## 📄 许可证

[MIT License](./LICENSE)

---

**go-sender** - 让 Go 语言的通知发送变得简单而强大 🚀

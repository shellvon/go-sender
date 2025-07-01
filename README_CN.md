**⚠️ 本项目仍在开发中，API 可能变动，请谨慎用于生产环境。**

# go-sender

> 🚀 Go 语言下最简单、灵活的多渠道消息推送库。

[English](./README.md) | [文档](./docs/getting-started.md)

---

## 🚀 项目路线图

请参阅我们的 [项目路线图与任务追踪](https://github.com/your-org/go-sender/issues/1)，了解当前优先级、计划特性和进度。

## 为什么选择 go-sender？

- 🪶 **轻量**：纯 Go 实现，零臃肿，极少依赖。
- 🧩 **灵活**：即插即用，支持短信、邮件、IM、Webhook 等。
- 🚀 **简单**：几行代码即可发消息。
- 🔌 **可扩展**：轻松添加新渠道或自定义功能。

---

## 🚀 快速上手

```go
import (
    "context"
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/sms"
)

func main() {
    sender := sender.NewSender()
    msg := sms.Aliyun().
        To([]string{"***REMOVED***"}).
        Content("Hello from go-sender!").
        TemplateCode("SMS_xxx").
        Build()
    if err := sender.Send(context.Background(), msg); err != nil {
        panic(err)
    }
}
```

安装：

```bash
go get github.com/shellvon/go-sender
```

---

## ✨ 支持的渠道

- **短信**：阿里云、腾讯云、华为、云片、创蓝 253 等
- **邮件**：SMTP、EmailJS、Resend
- **IM/机器人**：企业微信、钉钉、飞书、Telegram、ServerChan
- **Webhook/推送**：ntfy、Bark、PushDeer、PushPlus、Discord 等

完整支持列表和详细用法见 [docs/providers.md](docs/providers.md)。

---

## 🧑‍💻 进阶文档

- [快速入门](./docs/getting-started.md)
- [核心概念](./docs/concepts.md)
- [示例](./docs/examples.md)
- [高级用法](./docs/advanced.md)

---

**go-sender** —— Go 语言的万能消息推送利器。

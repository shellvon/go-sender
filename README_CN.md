**⚠️ 本项目仍在开发中，API 可能变动，请谨慎用于生产环境。**

# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)

> 🚀 Go 语言下最简单、灵活的多渠道消息推送库。

[English](./README.md) | [文档](./docs/getting-started.md)

---

## 🚀 项目路线图

请参阅我们的 [项目路线图与任务追踪](https://github.com/shellvon/go-sender/issues/1)，了解当前优先级、计划特性和进度。

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
    "log"

    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

func main() {
    // 1️⃣ 创建 Sender 实例（可按需再设置中间件）
    sender := gosender.NewSender()

    // 2️⃣ 初始化并注册短信 Provider（以阿里云为例）
    cfg := sms.Config{
        ProviderMeta: core.ProviderMeta{
            Strategy: core.StrategyRoundRobin,
        },
        Items: []*sms.Account{{
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "aliyun-default",
                    SubType: "aliyun",
                },
                Credentials: core.Credentials{
                    APIKey:    "your-access-key",
                    APISecret: "your-secret-key",
                },
            },
        }},
    }
    aliyunProvider, err := sms.New(cfg)
    if err != nil {
        log.Fatalf("创建 Provider 失败: %v", err)
    }
    sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, nil)

    // 3️⃣ 构造要发送的消息
    msg := sms.Aliyun().
        To("***REMOVED***").
        Content("Hello from go-sender!").
        TemplateID("SMS_xxx").
        Build()

    // 4️⃣ 发送并获取详细结果
    res, err := sender.SendWithResult(context.Background(), msg)
    if err != nil {
        log.Fatalf("发送失败: %v", err)
    }
    log.Printf("RequestID: %s, Provider: %s, 耗时: %v", res.RequestID, res.ProviderName, res.Elapsed)
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

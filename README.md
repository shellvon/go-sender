**⚠️ This project is under active development. APIs may change. Use with caution in production.**

# go-sender

> 🚀 The easiest way to send SMS, Email, IM, and Webhook notifications in Go.

[中文文档](./README_CN.md) | [Docs](./docs/getting-started.md)

---

## Why go-sender?

- 🪶 **Lightweight**: Pure Go, zero bloat, minimal dependencies.
- 🧩 **Flexible**: Plug-and-play for SMS, Email, IM, Webhook, and more.
- 🚀 **Simple**: Send a message in just a few lines.
- 🔌 **Extensible**: Add new channels or features easily.

---

## 🚀 Quick Start

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

Install:

```bash
go get github.com/shellvon/go-sender
```

---

## ✨ Supported Channels

- **SMS**: Aliyun, Tencent, Huawei, Yunpian, CL253, etc.
- **Email**: SMTP, EmailJS, Resend
- **IM/Bot**: WeCom, DingTalk, Lark, Telegram, ServerChan
- **Webhook/Push**: ntfy, Bark, PushDeer, PushPlus, Discord, etc.

See [docs/providers.md](docs/providers.md) for the full list and details.

---

## 🧑‍💻 Next Steps

- [Getting Started](./docs/getting-started.md)
- [Core Concepts](./docs/concepts.md)
- [Examples](./docs/examples.md)
- [Advanced Usage](./docs/advanced.md)

---

**go-sender** — Send anything, anywhere, with Go.

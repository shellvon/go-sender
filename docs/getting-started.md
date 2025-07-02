# Getting Started

Welcome to **go-sender** – the simple, flexible, and extensible notification library for Go developers.

## 🚀 Quick Installation

```bash
go get github.com/shellvon/go-sender
```

## 🏁 Your First Message

Send your first SMS with Aliyun in just a few lines:

```go
package main

import (
    "context"
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/sms"
)

func main() {
    sender := gosender.NewSender()
	msg := sms.Aliyun().
		To("13800138000").
		Content("Hello from go-sender!").
		TemplateID("SMS_xxx").
		Build()
    err := sender.Send(context.Background(), msg)
    if err != nil {
        panic(err)
    }
}
```

## ✉️ Supported Channels

- SMS: Aliyun, Tencent, Huawei, Yunpian, etc.
- Email: SMTP, EmailJS, Resend
- IM/Bot: WeCom, DingTalk, Lark, Telegram
- Webhook: Universal HTTP integration

See [providers.md](./providers.md) for the full list.

## 🧑‍💻 FAQ

**Q: Is go-sender production ready?**  
A: Yes, but always test with your own provider credentials and templates.

**Q: How do I add a new provider?**  
A: See [advanced.md](./advanced.md) for custom provider instructions.

**Q: Can I use go-sender in microservices?**  
A: Absolutely! It is designed for both monoliths and microservices.

## 📚 Next Steps

- [Core Concepts](./concepts.md)
- [Provider Usage](./providers.md)
- [Middleware & Advanced Features](./middleware.md)
- [Examples](./examples.md)

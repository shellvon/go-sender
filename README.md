**⚠️ This project is under active development. APIs may change. Use with caution in production.**

# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)

> 🚀 The easiest way to send SMS, Email, IM, and Webhook notifications in Go.

[中文文档](./README_CN.md) | [Docs](./docs/getting-started.md)

---

## 🚀 Project Roadmap

See our [Project Roadmap & Task Tracking](https://github.com/shellvon/go-sender/issues/1) for current priorities, planned features, and progress tracking.

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
	"log"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

func main() {
	// 1️⃣ Initialize a new sender instance (middleware can be added later)
	sender := gosender.NewSender()

	// 2️⃣ Prepare and register an SMS provider (Aliyun as an example)
	config := sms.Config{
		ProviderMeta: core.ProviderMeta{   // provider-level config
			Strategy: core.StrategyRoundRobin, // account selection strategy
		},
		Items: []*sms.Account{             // one or more sub-accounts (AK/SK)
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{
						Name:    "aliyun-default", // custom account name
						SubType: "aliyun",        // sms sub-provider
					},
					Credentials: core.Credentials{
						APIKey:    "your-access-key",
						APISecret: "your-secret-key",
					},
				},
				// Optional: Region, Callback, SignName ...
			},
		},
	}
	aliyunProvider, err := sms.New(config)
	if err != nil {
		log.Fatalf("failed to create provider: %v", err)
	}
	// Register with sender (nil = use global middleware settings)
	sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, nil)

	// 3️⃣ Build the message to send
	msg := sms.Aliyun().
		To("***REMOVED***").
		Content("Hello from go-sender!").
		TemplateID("SMS_xxx").
		Build()

	// 4️⃣ Send the message and receive detailed result
	res, err := sender.SendWithResult(context.Background(), msg)
	if err != nil {
		log.Fatalf("send failed: %v", err)
	}
	log.Printf("request id: %s, provider: %s, cost: %v", res.RequestID, res.ProviderName, res.Elapsed)
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

**⚠️ This project is under active development. APIs may change. Use with caution in production.**

# go-sender

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shellvon/go-sender)

> 🚀 The easiest way to send SMS, Email, IM, and Webhook notifications in Go.

[中文文档](./README_CN.md) | [Docs](./docs/getting-started.md)

---

## Why go-sender?

- 🪶 **Lightweight**: Pure Go, zero bloat, minimal dependencies.
- 🧩 **Flexible**: Plug-and-play for SMS, Email, IM, Webhook, and more.
- 🚀 **Simple**: Send a message in just a few lines.
- 🔌 **Extensible**: Add new channels, middleware **and Before/After Hooks** easily.

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

	// 2️⃣ Create an SMS account and provider (Aliyun as an example)
	account := sms.NewAccount("aliyun", "your-access-key", "your-secret-key",
		sms.Name("aliyun-default"),        // Custom account name
		sms.WithSignName("MyApp"),         // Optional SMS-specific settings
		sms.WithRegion("cn-hangzhou"))

	aliyunProvider, err := sms.NewProvider([]*sms.Account{account},
		sms.Strategy(core.StrategyRoundRobin)) // Round-robin strategy
	if err != nil {
		log.Fatalf("failed to create provider: %v", err)
	}
	// Register with sender (nil = use global middleware settings)
	sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, nil)

	// 3️⃣ Build the message to send
	msg := sms.Aliyun().
		To("13800138000").
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

// --- Mini Hook Demo ---------------------------------------------------

// Add a global before-hook: run for every message
senderMiddleware := &core.SenderMiddleware{}
senderMiddleware.UseBeforeHook(func(_ context.Context, m core.Message, _ *core.SendOptions) error {
   log.Printf("about to send %s", m.MsgID())
   return nil
})

// Register provider with custom middleware containing the hook
sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, senderMiddleware)

// Or add a per-request hook only for this message:
_, _ = sender.SendWithResult(
    context.Background(), msg,
    core.WithSendAfterHooks(func(_ context.Context, _ core.Message, _ *core.SendOptions, _ *core.SendResult, err error) {
        log.Printf("done, err=%v", err)
    }),
)

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

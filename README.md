**⚠️ This project is under active development. APIs may change. Use with caution in production.**

# go-sender

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
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/sms"
)

package main

import (
	"context"
	"github.com/yourpackage/gosender"
	"github.com/yourpackage/sms"
)

func main() {
	// Initialize a new sender instance
	sender := gosender.NewSender()

	// Create an SMS message using Aliyun provider
	msg := sms.Aliyun().
		To("13800138000").
		Content("Hello from go-sender!").
		TemplateID("SMS_xxx").
		Build()

	// Option 1: Send directly using provider
	provider, _ := sms.New(&cfg)
	provider.Send(context.Background(), msg, nil)

	// Option 2: Send using GoSender with registered provider
	sender.RegisterProvider(core.ProviderTypeSMS, provider, nil)
	// Automatically applies middleware chain
	sender.Send(context.Background(), msg)
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

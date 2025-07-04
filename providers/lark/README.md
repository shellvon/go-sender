# Lark (Feishu) Provider

> Send rich messages to Lark / Feishu group bots via WebHook.

[⬅️ Back to project README](../../README.md)

---

## Features

- Multiple bot tokens with round-robin / random / weighted load-balancing.
- Builder API for Text, Post (rich-text), Image, **Interactive Card (schema 2.0)**, Share Chat messages.
- Internationalisation support for Post & Card builders.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/lark"
)

cfg := lark.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin,
    },
    Items: []*lark.Account{
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name: "team-bot",
                },
                Credentials: core.Credentials{
                    APIKey: "YOUR_WEBHOOK_KEY", // part after /hook/
                },
            },
        },
    },
}
```

---

## Quick Builders

```go
// Text
msg := lark.Text().
    Content("Hello from go-sender!").
    Build()

// Interactive Card (schema 2.0)
msg := lark.Interactive().
    HeaderTitle("plain_text", "System Report").
    HeaderTemplate("blue").
    BodyDirection("vertical").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := lark.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := lark.New(&cfg)
sender.RegisterProvider(core.ProviderTypeLark, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## References

- Lark Bot API: <https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN>
- Card JSON 2.0: <https://open.feishu.cn/document/feishu-cards/card-json-v2-structure>

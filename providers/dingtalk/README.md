# DingTalk Provider

> Send rich messages to DingTalk (钉钉) group robots via WebHook.

[⬅️ Back to project README](../../README.md)

---

## Features

- Multiple accounts with round-robin / random / weighted load-balancing.
- Rich message builders: Text, Markdown, Link, Action Card, Feed Card.
- Chainable API for clean message composition.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

cfg := dingtalk.Config{
    ProviderMeta: core.ProviderMeta{                // global flags
        Strategy: core.StrategyRoundRobin,         // load-balancing strategy
    },
    Items: []*dingtalk.Account{                    // one or more bot tokens
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{     // unique name inside provider
                    Name: "main",
                },
                Credentials: core.Credentials{    // access_token from WebHook URL
                    APIKey: "YOUR_ACCESS_TOKEN",
                },
            },
        },
    },
}
```

---

## Quick Builders

```go
// Plain text
msg := dingtalk.Text().
    Content("Hello from go-sender!").
    Build()

// Markdown
md := "**CPU**: 45%  \n**Memory**: 60%"
msg := dingtalk.Markdown().
    Title("System Report").
    Text(md).
    Build()
```

See the GoDoc for all builders (`Text`, `Markdown`, `Link`, `ActionCard`, `FeedCard`).

---

## Usage

### 1. Direct Provider

```go
provider, err := dingtalk.New(&cfg)
if err != nil {
    log.Fatalf("init provider: %v", err)
}

ctx := context.Background()
err = provider.Send(ctx, msg, nil)   // msg built above
```

### 2. Using GoSender (with middleware, queue, etc.)

```go
import (
    gosender "github.com/shellvon/go-sender"
)

sender := gosender.NewSender()                       // global sender with middleware support
provider, _ := dingtalk.New(&cfg)
sender.RegisterProvider(core.ProviderTypeDingtalk, provider, nil)

_ = sender.Send(context.Background(), msg)           // middleware chain applied automatically
```

---

## References

- DingTalk Bot API: <https://open.dingtalk.com/document/robots/custom-robot-access>
- Message type docs: <https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type>

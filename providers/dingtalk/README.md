# DingTalk Provider

> Push notifications via [DingTalk](https://open.dingtalk.com/) custom robot.

[⬅️ Back to project README](../../README.md)

---

## Features

- Multiple webhook tokens with round-robin / random / weighted load-balancing
- Optional security signature support
- Rich message types:
  - Text (with @mentions)
  - Markdown
  - Link Card
  - ActionCard (single/multi buttons)
  - FeedCard

---

## Security Settings

The DingTalk custom robot supports optional security settings. When enabled, you need to provide the `APISecret` in addition to the `APIKey`:

```go
cfg := dingtalk.Config{
    Items: []*dingtalk.Account{{
        BaseAccount: core.BaseAccount{
            AccountMeta: core.AccountMeta{
                Name: "default",
            },
            Credentials: core.Credentials{
                APIKey:    "YOUR_ACCESS_TOKEN",    // Required
                APISecret: "YOUR_SIGN_SECRET",     // Optional, for signature
            },
        },
    }},
}
```

When `APISecret` is provided, the provider will automatically:

1. Generate timestamp
2. Calculate signature using HMAC-SHA256
3. Append signature parameters to webhook URL

For more details about security settings, see [DingTalk Documentation](https://open.dingtalk.com/document/orgapp/customize-robot-security-settings#title-7fs-kgs-36x).

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

cfg := dingtalk.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin, // load-balancing strategy
    },
    Items: []*dingtalk.Account{
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name: "primary",
                },
                Credentials: core.Credentials{
                    APIKey: "YOUR_ACCESS_TOKEN",
                },
            },
        },
    },
}
```

---

## Quick Builder

### Text Message

```go
msg := dingtalk.Text().
    Content("System Alert").
    AtMobiles([]string{"***REMOVED***"}).
    AtAll().
    Build()
```

### Markdown Message

```go
msg := dingtalk.Markdown().
    Title("Release Notes").
    Text("# Version 1.0.0\n## Features\n...").
    Build()
```

### Link Card

```go
msg := dingtalk.Link().
    Title("New Feature").
    Text("Check out our latest updates").
    MessageURL("https://example.com/news").
    PicURL("https://example.com/image.png").
    Build()
```

### ActionCard (Single Button)

```go
msg := dingtalk.ActionCard().
    Title("System Update").
    Text("# Important Update\nPlease review...").
    SingleButton("View Details", "https://example.com").
    BtnOrientation("0"). // 0=vertical, 1=horizontal
    Build()
```

### ActionCard (Multiple Buttons)

```go
msg := dingtalk.ActionCard().
    Title("Choose Action").
    Text("# Options\nPlease select...").
    AddButton("Accept", "https://example.com/accept").
    AddButton("Reject", "https://example.com/reject").
    BtnOrientation("1").
    Build()
```

### FeedCard

```go
msg := dingtalk.FeedCard().
    AddLink(
        "News Title",
        "https://example.com/news",
        "https://example.com/image.png",
    ).
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := dingtalk.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := dingtalk.New(&cfg)
sender.RegisterProvider(core.ProviderTypeDingtalk, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## References

- DingTalk Custom Robot: <https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type>
- Security Settings: <https://open.dingtalk.com/document/orgapp/customize-robot-security-settings#title-7fs-kgs-36x>

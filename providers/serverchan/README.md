# ServerChan Provider

> Push notifications via [Server 酱](https://sct.ftqq.com/) WebHook (SendKey).

[⬅️ Back to project README](../../README.md)

---

## Features

- Multiple SendKeys with round-robin / random / weighted load-balancing.
- Builder API for text messages with rich Markdown.
- Channel routing (Android, WeCom, DingTalk, Bark, …).
- Optional short summary, hide-IP flag, OpenID copy.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/serverchan"
)

cfg := serverchan.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin, // load-balancing strategy
    },
    Items: []*serverchan.Account{
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name: "primary",
                },
                Credentials: core.Credentials{
                    APIKey: "YOUR_SENDKEY_HERE", // SendKey from ServerChan
                },
            },
        },
    },
}
```

---

## Quick Builder

```go
msg := serverchan.Text().
    Title("System Alert").
    Content("CPU 90%, please check.").
    Short("CPU High").
    Channel("wecom|dingtalk").
    Build()
```

### Supported Channels

| Channel   | Code | Description         |
| --------- | ---- | ------------------- |
| android   | 98   | Android client push |
| wecom     | 66   | WeCom (企业微信)    |
| wecom_bot | 1    | WeCom Bot           |
| dingtalk  | 2    | DingTalk Bot        |
| feishu    | 3    | Lark / Feishu Bot   |
| bark      | 8    | Bark iOS            |
| test      | 0    | Test only           |
| custom    | 88   | Custom WebHook      |
| pushdeer  | 18   | PushDeer            |
| service   | 9    | Service notice      |

---

## Usage

### 1. Direct Provider

```go
provider, _ := serverchan.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := serverchan.New(&cfg)
sender.RegisterProvider(core.ProviderTypeServerChan, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## References

- ServerChan: <https://sct.ftqq.com/>
- API Docs: <https://sct.ftqq.com/sendkey>

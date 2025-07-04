# SMS Providers

> Unified SMS sending across Aliyun, Tencent Cloud, Huawei, CL253, Yunpian, and more.

[⬅️ Back to project README](../../README.md)

---

## Supported Providers

| Provider (brand)  | Website                     |
| ----------------- | --------------------------- |
| **Aliyun**        | https://www.aliyun.com      |
| **Tencent Cloud** | https://cloud.tencent.com   |
| **CL253**         | http://www.cl253.com        |
| **Huawei Cloud**  | https://www.huaweicloud.com |
| **Juhe**          | https://www.juhe.cn         |
| **Luosimao**      | https://luosimao.com        |
| **Smsbao**        | https://www.smsbao.com      |
| **Submail**       | https://www.submail.cn      |
| **UCP**           | https://www.ucpaas.com      |
| **Volcengine**    | https://www.volcengine.com  |
| **Yuntongxun**    | https://www.yuntongxun.com  |
| **Yunpian**       | https://www.yunpian.com     |

---

## Capabilities

| Provider      | CN SMS | Intl SMS | Voice | MMS | Notes                        |
| ------------- | ------ | -------- | ----- | --- | ---------------------------- |
| Aliyun        | ✅     | ✅       | ✅    | ✅  | Single, batch & MMS APIs     |
| Tencent Cloud | ✅     | ✅       | ✅    | ❌  | Voice via TTS / IVR          |
| Huawei Cloud  | ✅     | ✅       | ✅    | ❌  | Voice call API               |
| CL253         | ✅     | ❌       | ❌    | ❌  | Mainland only                |
| Yunpian       | ✅     | ✅       | ❌    | ❌  | Separate intl endpoint       |
| Volcengine    | ✅     | ✅       | ❌    | ❌  | Mainland & intl              |
| Submail       | ✅     | ✅       | ✅    | ❌  | XSend voice supported        |
| Juhe          | ✅     | ✅       | ❌    | ✅  | Mainland & intl; MMS support |
| Luosimao      | ✅     | ❌       | ❌    | ❌  | Mainland only                |
| Smsbao        | ✅     | ❌       | ❌    | ❌  | Mainland only                |
| UCP           | ✅     | ✅       | ✅    | ❌  | Voice & intl                 |
| Yuntongxun    | ✅     | ✅       | ✅    | ✅  | Rich APIs                    |

---

## Features

- Multiple accounts per provider with load-balancing strategies.
- Builder API for Text SMS; provider-specific helpers (Aliyun, Tencent, …).
- Template ID, sign, extra params per vendor.
- Supports batch & international SMS, voice SMS (where available).

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/sms"
)

cfg := sms.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin,
    },
    Items: []*sms.Account{
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "aliyun-main",
                    SubType: "aliyun",
                },
                Credentials: core.Credentials{
                    APIKey:    "ALIYUN_ACCESS_KEY_ID",
                    APISecret: "ALIYUN_ACCESS_KEY_SECRET",
                },
            },
            Region: "cn-hangzhou", // optional
        },
    },
}
```

---

## Quick Builder

```go
msg := sms.Aliyun().
    To("***REMOVED***").
    Content("Hello from go-sender!").
    TemplateID("SMS_1234567").
    Sign("YourSign").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := sms.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := sms.New(&cfg)
sender.RegisterProvider(core.ProviderTypeSMS, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## SendVia Helper

`SendVia(accountName, msg)` lets you choose a specific account (e.g., Aliyun-cn vs Aliyun-intl) at runtime:

```go
msg := sms.Aliyun().
    To("+86***REMOVED***").
    Content("Verification code: 5678").
    TemplateID("SMS_1234567").
    Build()

// try primary aliyun account first
if err := sender.SendVia("aliyun-main", msg); err != nil {
    // fallback to backup account (maybe in another region)
    _ = sender.SendVia("aliyun-backup", msg)
}
```

SendVia only switches between accounts **inside the SMS provider**; it does not allow cross-provider reuse of one message instance.

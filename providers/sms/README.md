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

| Provider      | CN SMS | Intl SMS | Voice | MMS | Notes                         |
| ------------- | ------ | -------- | ----- | --- | ----------------------------- |
| Aliyun        | ✅     | ✅       | ✅    | ✅  | Card SMS; voice domestic only |
| Tencent Cloud | ✅     | ✅       | ✅    | ❌  | Voice (TTS / IVR) domestic    |
| Huawei Cloud  | ✅     | ✅       | ❌    | ❌  | SMS only                      |
| CL253         | ✅     | ✅       | ❌    | ❌  |                               |
| Yunpian       | ✅     | ✅       | ✅    | ✅  | Supports voice & MMS          |
| Volcengine    | ✅     | ❌       | ❌    | ❌  | Mainland only                 |
| Submail       | ✅     | ✅       | ✅    | ✅  | XSend APIs                    |
| Juhe          | ✅     | ❌       | ❌    | ✅  | MMS only                      |
| Luosimao      | ✅     | ❌       | ✅    | ❌  | Voice API                     |
| Smsbao        | ✅     | ❌       | ✅    | ❌  |                               |
| UCP           | ✅     | ✅       | ❌    | ❌  |                               |
| Yuntongxun    | ✅     | ✅       | ✅    | ❌  | Voice / video SMS             |

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
_, _ = provider.Send(context.Background(), msg, nil) // provider.Send already returns (*SendResult, error)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := sms.New(&cfg)
sender.RegisterProvider(core.ProviderTypeSMS, provider, nil)
_, _ = sender.SendWithResult(context.Background(), msg)
```

---

## Per-Account Send (WithSendAccount)

Use `core.WithSendAccount("accountName")` to pick a specific account (e.g., Aliyun-cn vs Aliyun-intl) at runtime:

```go
msg := sms.Aliyun().
    To("+86***REMOVED***").
    Content("Verification code: 5678").
    TemplateID("SMS_1234567").
    Build()

ctx := context.Background()

// 先尝试主账号发送
if _, err := sender.SendWithResult(ctx, msg, core.WithSendAccount("aliyun-main")); err != nil {
    // 失败后回退到备用账号
    _, _ = sender.SendWithResult(ctx, msg, core.WithSendAccount("aliyun-backup"))
}
```

`core.WithSendAccount()` only switches between accounts **inside the SMS provider**; it does not allow cross-provider reuse of one message instance.

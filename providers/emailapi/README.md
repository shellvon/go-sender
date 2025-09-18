# Email API Provider (Experimental)

> Unified integration for hosted email-API services such as Resend, EmailJS, Mailgun, etc.

[⬅️ Back to project README](../../README.md)

---

- This package is **experimental** – APIs may change.


## Supported Providers

| Provider   | API Docs                                                       |
| ---------- | -------------------------------------------------------------- |
| Resend     | https://resend.com/docs/api-reference                          |
| EmailJS    | https://www.emailjs.com/docs/rest-api/send/                    |
| Mailgun    | https://documentation.mailgun.com/docs/mailgun                 |
| Mailjet    | https://dev.mailjet.com/email/guides/send-api-v31/             |
| Brevo      | https://developers.brevo.com/reference/sendtransacemail        |
| Mailersend | https://developers.mailersend.com/                             |
| Mailtrap   | https://api-docs.mailtrap.io/docs/mailtrap-api-docs            |

---

## Features

- Consistent builder API for every SaaS email service.
- JSON / HTTP based – no SMTP required.
- Multiple accounts with load-balancing strategies.
- Supports text & HTML body, attachments, custom headers, templates.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/emailapi"
)

cfg := emailapi.Config{
    ProviderMeta: core.ProviderMeta{                 // global flags
        Strategy: core.StrategyRoundRobin,          // round-robin by default
    },
    Items: []*emailapi.Account{                     // Account == API credentials per vendor
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "resend-us",          // unique name
                    SubType: "resend",            // vendor identifier
                },
                Credentials: core.Credentials{
                    APIKey: "YOUR_RESEND_API_KEY",
                },
            },
            Region: "us",                          // optional extra fields
        },
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name:   "emailjs-eu",
                    SubType: "emailjs",
                },
                Credentials: core.Credentials{
                    APIKey:    "EMAILJS_PUBLIC_KEY",
                    APISecret: "EMAILJS_SERVICE_ID",   // example
                },
            },
            Region: "eu",
        },
    },
}
```

---

## Quick Builder

```go
msg := emailapi.Email().
    To("alice@example.com").
    Subject("Welcome").
    HTML().
    Body(`<h1>Hello</h1><p>Thanks for joining!</p>`).
    From("no-reply@example.com").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := emailapi.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := emailapi.New(&cfg)
sender.RegisterProvider(core.ProviderTypeEmailAPI, provider, nil)
_ = sender.Send(context.Background(), msg)
```

If you want send email via STMP, see [Email Provider](../email/README.md)

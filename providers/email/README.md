# Email Provider

> SMTP email sending via [go-mail](https://github.com/wneessen/go-mail).

[⬅️ Back to project README](../../README.md)

---

## Features

- Multiple SMTP accounts with round-robin / random / weighted load-balancing.
- Plain‐text & HTML email, attachments, CC/BCC, Reply-To.
- TLS / SSL and username + password authentication.
- Builder API with full RFC 5322 address support.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

cfg := email.Config{
    ProviderMeta: core.ProviderMeta{                   // global flags
        Strategy: core.StrategyRoundRobin,            // load-balancing strategy
    },
    Items: []*email.Account{                          // one or more SMTP accounts
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name: "primary",                 // unique inside provider
                },
                Credentials: core.Credentials{
                    APIKey:    "your-smtp-username", // username
                    APISecret: "your-app-password",  // password / app token
                },
            },
            Host: "smtp.gmail.com",
            Port: 587,
            From: "no-reply@example.com",
        },
    },
}
```

---

## Quick Builders

```go
// Plain-text email
msg := email.Email().
    To("alice@example.com").
    Body("Hello from go-sender!").
    Subject("Greetings").
    From("no-reply@example.com").
    Build()

// HTML email with attachment
html := `<h1>Monthly Report</h1><p>See attachment.</p>`
msg := email.Email().
    To("team@example.com").
    HTML().
    Body(html).
    Subject("Report").
    From("reports@example.com").
    Attach("./report.pdf").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, err := email.New(&cfg)
if err != nil {
    log.Fatalf("init provider: %v", err)
}

ctx := context.Background()
_ = provider.Send(ctx, msg, nil)   // msg built above
```

### 2. Using GoSender (with middleware, queue, etc.)

```go
import (
    gosender "github.com/shellvon/go-sender"
)

sender := gosender.NewSender()                       // global sender
provider, _ := email.New(&cfg)
sender.RegisterProvider(core.ProviderTypeEmail, provider, nil)

_ = sender.Send(context.Background(), msg)           // middleware chain applied automatically
```

---

## Common SMTP Provider Settings

| Provider             | SMTP Server              | Port(s)      | Username / Description  | Auth Method            | Documentation                                                                                                   |
| -------------------- | ------------------------ | ------------ | ----------------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------- |
| Apple iCloud         | smtp.mail.me.com         | 587 / 465    | iCloud email            | App-specific password  | https://support.apple.com/en-us/HT202304                                                                        |
| Gmail (Google)       | smtp.gmail.com           | 587 / 465    | Gmail address           | App password / XOAUTH2 | https://support.google.com/mail/answer/7126229                                                                  |
| Outlook / Office 365 | smtp.office365.com       | 587          | Email address           | Password / XOAUTH2     | https://learn.microsoft.com/exchange/clients-and-mobile-in-exchange-online/authenticated-client-smtp-submission |
| QQ Mail              | smtp.qq.com              | 465 / 587    | QQ email                | Auth code              | https://service.mail.qq.com/detail/0/428                                                                        |
| 163 Mail             | smtp.163.com             | 465 / 994    | 163 email               | Auth code              | https://help.mail.163.com/faqDetail.do?code=d7a5dc8471cd0c0e8b4b8f4f8e49998b374173cfe9171305fa1ce630d7f67ac25ef2e192b234ae4d                                                             |
| SendGrid             | smtp.sendgrid.net        | 587 / 465    | Literal string "apikey" | API Key                | https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp                                     |
| Mailgun              | smtp.mailgun.org         | 587 / 465    | postmaster@<domain>     | SMTP password          | https://help.mailgun.com/hc/en-us/articles/203380100                                                            |
| Mailtrap             | sandbox.smtp.mailtrap.io | 2525/465/587 | Mailtrap username       | Mailtrap password      | https://help.mailtrap.io/article/122-mailtrap-email-sending-smtp-integration                                    |
| Zoho Mail            | smtp.zoho.com            | 465 / 587    | Zoho email              | App password           | https://www.zoho.com/mail/help/zoho-smtp.html                                                                   |
| Yandex               | smtp.yandex.com          | 465 / 587    | Yandex email            | Password               | https://yandex.com/support/mail/mail-clients/others.html                                                        |
| Alibaba Cloud        | smtp.mxhichina.com       | 465 / 25     | Email address           | Password               | https://help.aliyun.com/document_detail/36576.html                                                              |
| Tencent Exmail       | smtp.exmail.qq.com       | 465 / 587    | Email address           | Password               | https://open.work.weixin.qq.com/help2/pc/1988744458                                                           |

> Tip: for any other SMTP provider, consult its official documentation for host, port and authentication details.

---

## References

- SMTP RFC 5321 / RFC 5322
- go-mail library: <https://github.com/wneessen/go-mail>

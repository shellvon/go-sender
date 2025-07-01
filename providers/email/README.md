[⬅️ Back to Main README](../../README.md)

# Email Provider

The Email Provider for go-sender supports sending emails through SMTP servers using the [go-mail](https://github.com/wneessen/go-mail) library.

## Features

- **Multiple Account Support**: Configure multiple email accounts with load balancing strategies
- **SMTP Authentication**: Support for username/password authentication
- **TLS/SSL Support**: Secure email transmission
- **HTML and Text Emails**: Support for both HTML and plain text email formats
- **Attachments**: File attachment support
- **CC and BCC**: Carbon copy and blind carbon copy support
- **Reply-To Support**: Set custom reply-to address for email responses
- **[RFC 5322](https://tools.ietf.org/html/rfc5322) Address Format**: Full support for RFC 5322 email address format with display names

## Email Address Format

All email addresses (From, To, Cc, Bcc, ReplyTo) support [RFC 5322](https://tools.ietf.org/html/rfc5322) format, which allows you to include display names along with email addresses:

- **Simple format**: `user@example.com`
- **With display name**: `"John Doe" <john@example.com>` or `John Doe <john@example.com>`

When you need to specify a display name for an email address, you can use the RFC 5322 format directly in the string:

```go
// Example with display names
msg := email.NewMessage(
    []string{"John Doe <john@example.com>", "Jane Smith <jane@example.com>"},
    "Hello team!",
    email.WithSubject("Team Update"),
    email.WithFrom("Manager <manager@company.com>"),
    email.WithCc("HR Department <hr@company.com>"),
    email.WithReplyTo("Support Team <support@company.com>"),
)
```

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // or StrategyWeighted
    },
    Accounts: []email.Account{
        {
            Name:     "primary",
            Host:     "smtp.gmail.com",
            Port:     587,
            Username: "your-email@gmail.com",
            Password: "your-app-password",
            From:     "your-email@gmail.com",
            Weight:   1,
        },
        {
            Name:     "backup",
            Host:     "smtp.outlook.com",
            Port:     587,
            Username: "your-email@outlook.com",
            Password: "your-password",
            From:     "your-email@outlook.com",
            Weight:   1,
        },
    },
}
```

## Usage

### Basic Usage

```go
// Create provider
provider, err := email.New(config)
if err != nil {
    log.Fatal(err)
}

// Create message with functional options
msg := email.NewMessage(
    []string{"recipient@example.com"},
    "Hello, this is a test email!",
    email.WithSubject("Test Email"),
    email.WithFrom("sender@example.com"),
)

// Send email
err = provider.Send(context.Background(), msg)
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
```

### HTML Email

```go
htmlBody := `
<html>
    <body>
        <h1>Welcome!</h1>
        <p>This is an <strong>HTML</strong> email.</p>
    </body>
</html>`

msg := email.NewMessage(
    []string{"recipient@example.com"},
    htmlBody,
    email.WithSubject("HTML Email"),
    email.WithFrom("sender@example.com"),
    email.WithHTML(),
)
```

### Email with CC and BCC

```go
msg := email.NewMessage(
    []string{"recipient@example.com"},
    "Meeting reminder",
    email.WithSubject("Team Meeting"),
    email.WithFrom("organizer@company.com"),
    email.WithCc("manager@company.com", "team@company.com"),
    email.WithBcc("hr@company.com"),
)
```

### Email with Reply-To

```go
msg := email.NewMessage(
    []string{"customer@example.com"},
    "Thank you for your inquiry",
    email.WithSubject("Customer Support"),
    email.WithFrom("noreply@company.com"),
    email.WithReplyTo("support@company.com"), // Replies will go to support team
)
```

### Email with Attachments

```go
msg := email.NewMessage(
    []string{"recipient@example.com"},
    "Please find the attached report.",
    email.WithSubject("Monthly Report"),
    email.WithFrom("reports@company.com"),
    email.WithAttachments("/path/to/report.pdf", "/path/to/data.xlsx"),
)
```

### Email with Display Names (RFC 5322 Format)

```go
msg := email.NewMessage(
    []string{"John Doe <john@example.com>", "Jane Smith <jane@example.com>"},
    "Hello team!",
    email.WithSubject("Team Update"),
    email.WithFrom("Manager <manager@company.com>"),
    email.WithCc("HR Department <hr@company.com>"),
)
```

### Using with go-sender

```go
import (
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/email"
)

// Create sender with email provider
sender := gosender.New()
sender.AddProvider(email.New(config))

// Send email
msg := email.NewMessage(
    []string{"user@example.com"},
    "Welcome to our service!",
    email.WithSubject("Welcome"),
    email.WithFrom("noreply@company.com"),
)

err := sender.Send(context.Background(), msg)
```

### Send via Specific Account

```go
// Get provider instance
emailProvider := sender.GetProvider(core.ProviderTypeEmail).(*email.Provider)

// Send directly through the provider
msg := email.NewMessage(
    []string{"recipient@example.com"},
    "Direct message",
    email.WithSubject("Direct"),
    email.WithFrom("direct@company.com"),
)

err := emailProvider.Send(context.Background(), msg)
```

## Message Options

The email provider uses the functional options pattern for message construction:

- `WithFrom(from string)`: Set the sender email address (supports RFC 5322 format)
- `WithSubject(subject string)`: Set the email subject
- `WithCc(cc ...string)`: Add CC recipients (supports RFC 5322 format)
- `WithBcc(bcc ...string)`: Add BCC recipients (supports RFC 5322 format)
- `WithReplyTo(replyTo string)`: Set the Reply-To address (supports RFC 5322 format)
- `WithHTML()`: Mark the email as HTML content
- `WithAttachments(attachments ...string)`: Add file attachments

## Configuration Options

### Account Configuration

- `Name`: Unique identifier for the account
- `Host`: SMTP server hostname
- `Port`: SMTP server port (typically 25, 465, or 587)
- `Username`: SMTP username
- `Password`: SMTP password or app password
- `From`: Default sender email address (supports RFC 5322 format)
- `Weight`: Weight for weighted selection strategy
- `Disabled`: Whether this account is disabled

### Strategy Options

- `StrategyRoundRobin`: Rotate through accounts in order
- `StrategyWeighted`: Select accounts based on their weights

## Error Handling

The provider returns descriptive errors for common issues:

- Invalid email addresses
- Missing required fields (recipients, body)
- SMTP connection failures
- Authentication errors
- File attachment errors

## Testing

Run the integration tests with proper environment variables:

```bash
export EMAIL_HOST="smtp.gmail.com"
export EMAIL_PORT="587"
export EMAIL_USERNAME="your-email@gmail.com"
export EMAIL_PASSWORD="your-app-password"
export EMAIL_FROM="your-email@gmail.com"
export EMAIL_TO="test@example.com"
export EMAIL_CC="cc@example.com"
export EMAIL_BCC="bcc@example.com"

go test -v ./providers/email/
```

## Dependencies

- [go-mail](https://github.com/wneessen/go-mail): Modern, actively maintained email library
- [go-sender/core](https://github.com/shellvon/go-sender): Core framework interfaces and utilities

## Common SMTP Provider Settings & Official Documentation

| Provider          | SMTP Server Address      | Port(s)      | Username/Description            | Password/Auth Method           | Official Documentation                                                                                                    |
| ----------------- | ------------------------ | ------------ | ------------------------------- | ------------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
| Apple iCloud      | smtp.mail.me.com         | 587/465      | Apple ID email                  | Apple ID password/app password | [Apple Docs](https://support.apple.com/en-us/HT202304)                                                                    |
| Gmail (Google)    | smtp.gmail.com           | 587/465      | Gmail address                   | App password/XOAUTH2           | [Gmail Docs](https://support.google.com/a/answer/176600?hl=en)                                                              |
| Outlook/Office365 | smtp.office365.com       | 587          | Email address                   | Login password/XOAUTH2         | [Microsoft Docs](https://support.microsoft.com/zh-hk/office/outlook-com-%E7%9A%84-pop-imap-%E5%92%8C-smtp-%E8%A8%AD%E5%AE%9A-d088b986-291d-42b8-9564-9c414e2aa040) |
| Alibaba Cloud     | smtp.mxhichina.com       | 465/25       | Email address                   | Login password                 | [Alibaba Docs](https://help.aliyun.com/document_detail/36576.html)                              |
| Tencent Exmail    | smtp.exmail.qq.com       | 465/587      | Email address                   | Login password                 | [Tencent Docs](https://www.tencentcloud.com/document/product/1084/44458)                                     |
| QQ Mail           | smtp.qq.com              | 465/587      | QQ number or full email address | Auth code                      | [QQ Mail Docs](https://service.mail.qq.com/detail/0/428)                                     |
| 163 Mail          | smtp.163.com             | 465/994      | 163 email address               | Auth code                      | [163 Mail Docs](https://help.mail.163.com/faq.do?m=OTUw&id=MjQ5Nw==)                                                      |
| SendGrid          | smtp.sendgrid.net        | 587/465      | apikey (literal string)         | API Key                        | [SendGrid Docs](https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp)                              |
| Mailgun           | smtp.mailgun.org         | 587/465      | postmaster@yourdomain           | SMTP password                  | [Mailgun Docs](https://help.mailgun.com/hc/en-us/articles/203380100-How-Do-I-Use-Mailgun-SMTP-)                           |
| Mailtrap          | sandbox.smtp.mailtrap.io | 2525/465/587 | Mailtrap username               | Mailtrap password              | [Mailtrap Docs](https://help.mailtrap.io/article/122-mailtrap-email-sending-smtp-integration)                                                                           |
| Zoho Mail         | smtp.zoho.com            | 465/587      | Zoho email address              | Login password/app password    | [Zoho Docs](https://www.zoho.com/mail/help/zoho-smtp.html)                                                                |
| Yandex            | smtp.yandex.com          | 465/587      | Yandex email address            | Login password                 | [Yandex Docs](https://yandex.com/support/mail/mail-clients/others.html)                                                   |

### Example Configuration (SendGrid)

```go
config := email.Config{
    Accounts: []email.Account{
        {
            Name:     "sendgrid",
            Host:     "smtp.sendgrid.net",
            Port:     587,
            Username: "apikey", // literal string
            Password: "<your_sendgrid_api_key>",
            From:     "your@email.com",
        },
    },
}
```

### Example Configuration (QQ Mail)

```go
config := email.Config{
    Accounts: []email.Account{
        {
            Name:     "qq",
            Host:     "smtp.qq.com",
            Port:     465,
            Username: "123456@qq.com",
            Password: "<your_qq_email_auth_code>", // Get the auth code from QQ Mail settings
            From:     "123456@qq.com",
        },
    },
}
```

> For more provider settings, please refer to the table above and each provider's official documentation.

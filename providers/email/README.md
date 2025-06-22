# Email Provider

This provider supports sending emails via SMTP with support for multiple email accounts, templates, and advanced features.

## Features

- **Multi-Account Support**: Configure multiple SMTP accounts with load balancing
- **SMTP Authentication**: Support for various authentication methods (PLAIN, LOGIN, CRAM-MD5)
- **TLS/SSL Support**: Full support for STARTTLS and SSL/TLS connections
- **Template Support**: HTML and text email templates
- **Attachments**: Support for file attachments
- **CC/BCC**: Support for carbon copy and blind carbon copy
- **Priority Levels**: Set email priority (High, Normal, Low)
- **Custom Headers**: Add custom email headers
- **Retry Mechanism**: Built-in retry with exponential backoff

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

// Create email configuration
config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []email.Account{
        {
            Name:     "primary",
            Host:     "smtp.gmail.com",
            Port:     587,
            Username: "your-email@gmail.com",
            Password: "your-app-password",
            From:     "your-email@gmail.com",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup",
            Host:     "smtp.office365.com",
            Port:     587,
            Username: "your-email@outlook.com",
            Password: "your-password",
            From:     "your-email@outlook.com",
            Weight:   50,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := email.New(config)
if err != nil {
    log.Fatalf("Failed to create email provider: %v", err)
}
```

## Message Types

### 1. Simple Text Email

```go
// Basic text email
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Hello from Go-Sender",
    Body:    "This is a simple text email.",
}
```

### 2. HTML Email

```go
// HTML email with both text and HTML content
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Welcome to Go-Sender",
    Body:    "This is the text version for email clients that don't support HTML.",
    HTMLBody: `
        <html>
            <body>
                <h1>Welcome to Go-Sender!</h1>
                <p>This is an <strong>HTML email</strong> with styling.</p>
                <p>Sent at: <em>%s</em></p>
            </body>
        </html>
    `,
}
```

### 3. Email with CC and BCC

```go
// Email with carbon copy and blind carbon copy
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Cc:      []string{"cc@example.com"},
    Bcc:     []string{"bcc@example.com"},
    Subject: "Meeting Invitation",
    Body:    "You are invited to attend the meeting.",
}
```

### 4. Email with Attachments

```go
// Email with file attachments
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Report Attached",
    Body:    "Please find the attached report.",
    Attachments: []email.Attachment{
        {
            Filename: "report.pdf",
            Data:     pdfData, // []byte
            ContentType: "application/pdf",
        },
        {
            Filename: "data.csv",
            Data:     csvData, // []byte
            ContentType: "text/csv",
        },
    },
}
```

### 5. Email with Custom Headers

```go
// Email with custom headers
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Custom Header Email",
    Body:    "This email has custom headers.",
    Headers: map[string]string{
        "X-Custom-Header": "custom-value",
        "X-Priority": "1",
        "X-MSMail-Priority": "High",
    },
}
```

## Advanced Configuration

### SMTP with TLS/SSL

```go
config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []email.Account{
        {
            Name:     "gmail-ssl",
            Host:     "smtp.gmail.com",
            Port:     465, // SSL port
            Username: "your-email@gmail.com",
            Password: "your-app-password",
            From:     "your-email@gmail.com",
            UseSSL:   true, // Enable SSL
            Weight:   100,
        },
        {
            Name:     "gmail-starttls",
            Host:     "smtp.gmail.com",
            Port:     587, // STARTTLS port
            Username: "your-email@gmail.com",
            Password: "your-app-password",
            From:     "your-email@gmail.com",
            UseTLS:   true, // Enable STARTTLS
            Weight:   50,
        },
    },
}
```

### Custom SMTP Settings

```go
config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []email.Account{
        {
            Name:         "custom-smtp",
            Host:         "smtp.example.com",
            Port:         587,
            Username:     "your-username",
            Password:     "your-password",
            From:         "noreply@example.com",
            UseTLS:       true,
            AuthType:     "PLAIN", // PLAIN, LOGIN, CRAM-MD5
            Timeout:      30 * time.Second,
            Weight:       100,
        },
    },
}
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

// Create sender
s := gosender.NewSender(nil)

// Register email provider
emailProvider, err := email.New(config)
if err != nil {
    log.Fatalf("Failed to create email provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeEmail, emailProvider, nil)

// Send email
ctx := context.Background()
emailMsg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Hello from Go-Sender",
    Body:    "This is a test email sent via Go-Sender.",
}
err = s.Send(ctx, emailMsg)
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
```

## Message Options

### Basic Options

- `To`: Primary recipients (required)
- `Cc`: Carbon copy recipients
- `Bcc`: Blind carbon copy recipients
- `Subject`: Email subject (required)
- `Body`: Plain text body
- `HTMLBody`: HTML body

### Advanced Options

- `Attachments`: File attachments
- `Headers`: Custom email headers
- `Priority`: Email priority (High, Normal, Low)
- `ReplyTo`: Reply-to address
- `From`: Override sender address

## Configuration Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of email account configurations

### Account

- `Name`: Account name for identification
- `Host`: SMTP server hostname
- `Port`: SMTP server port
- `Username`: SMTP username
- `Password`: SMTP password
- `From`: Default sender address
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled
- `UseSSL`: Use SSL connection (default: false)
- `UseTLS`: Use STARTTLS (default: false)
- `AuthType`: Authentication type (PLAIN, LOGIN, CRAM-MD5)
- `Timeout`: Connection timeout

### Message

- `To`: Recipient email addresses
- `Cc`: Carbon copy addresses
- `Bcc`: Blind carbon copy addresses
- `Subject`: Email subject
- `Body`: Plain text body
- `HTMLBody`: HTML body
- `Attachments`: File attachments
- `Headers`: Custom headers
- `Priority`: Email priority
- `ReplyTo`: Reply-to address

### Attachment

- `Filename`: Attachment filename
- `Data`: Attachment data ([]byte)
- `ContentType`: MIME content type
- `Inline`: Whether attachment is inline

## Error Handling

The provider handles:

- SMTP connection errors
- Authentication failures
- Network timeouts
- Invalid email addresses
- Provider selection based on strategy
- Fallback to alternative accounts on failure

## Security Best Practices

### 1. Use App Passwords for Gmail

```go
// For Gmail, use App Passwords instead of regular passwords
{
    Name:     "gmail",
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-16-char-app-password", // App password, not regular password
    From:     "your-email@gmail.com",
    UseTLS:   true,
}
```

### 2. Use Environment Variables

```go
// Store sensitive data in environment variables
{
    Name:     "gmail",
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: os.Getenv("GMAIL_USERNAME"),
    Password: os.Getenv("GMAIL_APP_PASSWORD"),
    From:     os.Getenv("GMAIL_FROM"),
    UseTLS:   true,
}
```

### 3. Enable TLS/SSL

```go
// Always use TLS/SSL for security
{
    Name:     "secure-smtp",
    Host:     "smtp.example.com",
    Port:     587,
    Username: "username",
    Password: "password",
    From:     "noreply@example.com",
    UseTLS:   true, // Enable STARTTLS
}
```

## Common SMTP Providers

### Gmail

```go
{
    Name:     "gmail",
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your-email@gmail.com",
    Password: "your-app-password",
    From:     "your-email@gmail.com",
    UseTLS:   true,
}
```

### Outlook/Office 365

```go
{
    Name:     "outlook",
    Host:     "smtp.office365.com",
    Port:     587,
    Username: "your-email@outlook.com",
    Password: "your-password",
    From:     "your-email@outlook.com",
    UseTLS:   true,
}
```

### Yahoo Mail

```go
{
    Name:     "yahoo",
    Host:     "smtp.mail.yahoo.com",
    Port:     587,
    Username: "your-email@yahoo.com",
    Password: "your-app-password",
    From:     "your-email@yahoo.com",
    UseTLS:   true,
}
```

### Custom SMTP Server

```go
{
    Name:     "custom",
    Host:     "smtp.yourcompany.com",
    Port:     587,
    Username: "noreply@yourcompany.com",
    Password: "your-password",
    From:     "noreply@yourcompany.com",
    UseTLS:   true,
    AuthType: "PLAIN",
}
```

## Best Practices

### 1. Use Multiple Accounts for Reliability

```go
config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []email.Account{
        {
            Name: "primary",
            Host: "smtp.gmail.com",
            Port: 587,
            Username: "primary@gmail.com",
            Password: "app-password",
            From: "primary@gmail.com",
            UseTLS: true,
            Weight: 100,
        },
        {
            Name: "backup",
            Host: "smtp.office365.com",
            Port: 587,
            Username: "backup@outlook.com",
            Password: "password",
            From: "backup@outlook.com",
            UseTLS: true,
            Weight: 50,
        },
    },
}
```

### 2. Set Appropriate Timeouts

```go
{
    Name:     "timeout-configured",
    Host:     "smtp.example.com",
    Port:     587,
    Username: "username",
    Password: "password",
    From:     "noreply@example.com",
    UseTLS:   true,
    Timeout:  30 * time.Second, // 30 second timeout
}
```

### 3. Use HTML and Text Versions

```go
msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Welcome Email",
    Body:    "Welcome to our service! Click here: https://example.com",
    HTMLBody: `
        <html>
            <body>
                <h1>Welcome to our service!</h1>
                <p>Click <a href="https://example.com">here</a> to get started.</p>
            </body>
        </html>
    `,
}
```

### 4. Handle Attachments Properly

```go
// Read file and create attachment
fileData, err := os.ReadFile("report.pdf")
if err != nil {
    log.Fatal(err)
}

msg := &email.Message{
    To:      []string{"recipient@example.com"},
    Subject: "Monthly Report",
    Body:    "Please find the attached monthly report.",
    Attachments: []email.Attachment{
        {
            Filename:    "monthly-report.pdf",
            Data:        fileData,
            ContentType: "application/pdf",
        },
    },
}
```

## API Reference

### Constructor Functions

- `New(config Config) (*Provider, error)`: Create new email provider

### Message Types

- `Message`: Email message structure
- `Attachment`: File attachment structure

### Configuration Types

- `Config`: Email provider configuration
- `Account`: Individual email account configuration

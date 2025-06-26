# emailapi Provider (Experimental)

> ⚠️ **This package is under active development. APIs are unstable and may change at any time. Use with caution in production!**

The `emailapi` package provides a unified interface for sending emails via popular third-party email APIs. It is designed for extensibility and future support of multiple providers, similar to the SMS provider architecture.

## 🚧 Status

- **This package is experimental and not stable.**
- Only a few providers are currently implemented. More are planned.
- The API surface and configuration may change frequently.

## ✨ Supported & Planned Providers

| Provider   | Website                                       | API Docs                                                              | Status      |
| ---------- | --------------------------------------------- | --------------------------------------------------------------------- | ----------- |
| EmailJS    | [emailjs.com](https://www.emailjs.com/)       | [API](https://www.emailjs.com/docs/rest-api/send/)                    | Implemented |
| Resend     | [resend.com](https://resend.com/)             | [API](https://resend.com/docs/api-reference/emails/send-batch-emails) | Implemented |
| Mailgun    | [mailgun.com](https://www.mailgun.com/)       | [API](https://documentation.mailgun.com/en/latest/api_reference.html) | Planned     |
| Mailjet    | [mailjet.com](https://www.mailjet.com/)       | [API](https://dev.mailjet.com/email/guides/send-api-v31/)             | Planned     |
| Brevo      | [brevo.com](https://www.brevo.com/)           | [API](https://developers.brevo.com/docs)                              | Planned     |
| Mailersend | [mailersend.com](https://www.mailersend.com/) | [API](https://developers.mailersend.com/)                             | Planned     |
| Mailtrap   | [mailtrap.io](https://mailtrap.io/)           | [API](https://api-docs.mailtrap.io/docs)                              | Planned     |

## Features

- Unified message structure for all API providers
- Support for To, Cc, Bcc, Subject, Body (text/html), Attachments, Headers, Templates, etc.
- Extensible: add new providers with minimal effort
- Provider routing via `provider` field in message (like SMS)

## Usage

```go
import (
    "github.com/shellvon/go-sender/providers/emailapi"
)

// See each provider's Go file for configuration and usage examples.
```

## Contributing

- Contributions and feedback are welcome!
- Please note that the API is not stable and may change frequently.

---

For more details, see the code and comments in this directory.

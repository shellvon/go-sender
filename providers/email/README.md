[⬅️ Back to Main README](../../README.md)

# Email Provider | 邮件发送组件

The Email Provider for go-sender supports sending emails through SMTP servers using the [go-mail](https://github.com/wneessen/go-mail) library.

go-sender 的邮件组件支持通过 SMTP 协议发送邮件，底层基于 [go-mail](https://github.com/wneessen/go-mail) 实现。

---

## Usage

### Constructing Email Messages

- Use the builder style for all advanced options.
- For the simplest use case, use NewMessage for required fields only.

#### API

- `Email()` - create a builder instance
- `To(recipients ...string)` - set recipients
- `Subject(subject string)` - set subject
- `Body(body string)` - set body
- `From(from string)` - set sender
- `Cc(cc ...string)` - set CC
- `Bcc(bcc ...string)` - set BCC
- `ReplyTo(replyTo string)` - set reply-to
- `HTML()` - mark as HTML content
- `Attach(files ...string)` - replace attachments
- `AddAttach(files ...string)` - append attachments
- `Build()` - build the Message instance

**Example**

```go
msg := email.Email().
    To("a@b.com").
    Subject("Subject").
    Body("Body").
    From("noreply@b.com").
    Cc("c@b.com").
    Bcc("d@b.com").
    ReplyTo("reply@b.com").
    HTML().
    Attach("a.txt").
    AddAttach("b.txt", "c.pdf").
    Build()
```

For the simplest use case:

```go
msg := email.NewMessage([]string{"user@example.com"}, "Hello!")
```

All advanced options (subject, from, cc, bcc, replyTo, html, attachments, etc.) must be set via the builder API.

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple email accounts with load balancing strategies | 支持多账号、负载均衡
- **SMTP Authentication SMTP 认证**: Support for username/password authentication | 支持用户名/密码认证
- **TLS/SSL Support 安全传输**: Secure email transmission | 支持 TLS/SSL 加密
- **HTML and Text Emails HTML/文本邮件**: Support for both HTML and plain text email formats | 支持 HTML 与纯文本格式
- **Attachments 附件支持**: File attachment support | 支持文件附件
- **CC and BCC 抄送/密送**: Carbon copy and blind carbon copy support | 支持抄送与密送
- **Reply-To Support 回复地址**: Set custom reply-to address for email responses | 支持自定义回复地址
- **[RFC 5322](https://tools.ietf.org/html/rfc5322) Address Format**: Full support for RFC 5322 email address format with display names | 完全支持 RFC 5322 邮箱格式（含显示名）

---

## Email Address Format | 邮箱地址格式

All email addresses (From, To, Cc, Bcc, ReplyTo) support [RFC 5322](https://tools.ietf.org/html/rfc5322) format, which allows you to include display names along with email addresses.

所有邮箱地址（发件人、收件人、抄送、密送、回复）均支持 RFC 5322 格式，可带显示名：

- **Simple format 简单格式**: `user@example.com`
- **With display name 含显示名**: `"John Doe" <john@example.com>` 或 `John Doe <john@example.com>`

```go
// English: Example with display names
// 中文：带显示名的邮箱示例
msg := email.NewMessage(
    []string{"John Doe <john@example.com>", "Jane Smith <jane@example.com>"},
    "Hello team!",
    email.WithSubject("Team Update"),
    email.WithFrom("Manager <manager@company.com>"),
    email.WithCc("HR Department <hr@company.com>"),
    email.WithReplyTo("Support Team <support@company.com>"),
)
```

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/email"
)

// English: Create email config
// 中文：创建邮件配置
config := email.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询或加权
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
        // ... more accounts
    },
}
```

---

## Usage | 用法

### Basic Usage | 基本用法

```go
// English: Create provider and send email
// 中文：创建 provider 并发送邮件
provider, err := email.New(config)
if err != nil {
    log.Fatal(err)
}

msg := email.NewMessage(
    []string{"recipient@example.com"},
    "Hello, this is a test email!",
    email.WithSubject("Test Email"),
    email.WithFrom("sender@example.com"),
)

err = provider.Send(context.Background(), msg)
if err != nil {
    log.Printf("Failed to send email: %v", err)
}
```

### HTML Email | 发送 HTML 邮件

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

### Email with CC and BCC | 抄送与密送

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

### Email with Reply-To | 自定义回复地址

```go
msg := email.Email().
    To("customer@example.com").
    Body("Thank you for your inquiry").
    Subject("Customer Support").
    From("noreply@company.com").
    ReplyTo("support@company.com"). // Replies will go to support team
    Build()
```

### Email with Attachments | 带附件邮件

```go
msg := email.Email().
    To("recipient@example.com").
    Body("Please find the attached report.").
    Subject("Monthly Report").
    From("reports@company.com").
    Attach("/path/to/report.pdf").
    AddAttach("/path/to/data.xlsx").
    Build()
```

### Using with go-sender | 与 go-sender 集成

```go
import (
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/email"
)

// English: Create sender with email provider
// 中文：创建 sender 并添加邮件 provider
sender := gosender.New()
sender.AddProvider(email.New(config))

msg := email.Email().
    To("user@example.com").
    Body("Welcome to our service!").
    Subject("Welcome").
    From("noreply@company.com").
    Build()

err := sender.Send(context.Background(), msg)
```

### Send via Specific Account | 指定账号发送

```go
// English: Get provider instance and send
// 中文：获取 provider 实例并发送
emailProvider := sender.GetProvider(core.ProviderTypeEmail).(*email.Provider)

msg := email.Email().
    To("recipient@example.com").
    Body("Direct message").
    Subject("Direct").
    From("direct@company.com").
    Build()

err := emailProvider.Send(context.Background(), msg)
```

---

## Message Construction

All advanced options (subject, from, cc, bcc, replyTo, html, attachments, etc.) must be set via the builder API.

---

## Configuration Options | 配置参数

### Account Configuration | 账号配置

- `Name`: Unique identifier for the account | 账号唯一标识
- `Host`: SMTP server hostname | SMTP 服务器地址
- `Port`: SMTP server port (typically 25, 465, or 587) | SMTP 端口（常用 25/465/587）
- `Username`: SMTP username | SMTP 用户名
- `Password`: SMTP password or app password | SMTP 密码或授权码
- `From`: Default sender email address (supports RFC 5322 format) | 默认发件人（支持 RFC 5322 格式）
- `Weight`: Weight for weighted selection strategy | 加权策略权重
- `Disabled`: Whether this account is disabled | 是否禁用

### Strategy Options | 策略选项

- `StrategyRoundRobin`: Rotate through accounts in order | 轮询
- `StrategyWeighted`: Select accounts based on their weights | 加权选择

---

## Error Handling | 错误处理

The provider returns descriptive errors for common issues:

组件会针对常见问题返回详细错误：

- Invalid email addresses | 邮箱地址无效
- Missing required fields (recipients, body) | 缺少必填项（收件人、正文）
- SMTP connection failures | SMTP 连接失败
- Authentication errors | 认证失败
- File attachment errors | 附件处理失败

---

## Testing | 测试

Run the integration tests with proper environment variables:

设置环境变量后运行集成测试：

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

---

## Dependencies | 依赖

- [go-mail](https://github.com/wneessen/go-mail): Modern, actively maintained email library | 现代邮件库
- [go-sender/core](https://github.com/shellvon/go-sender): Core framework interfaces and utilities | go-sender 核心框架

---

## Common SMTP Provider Settings & Official Documentation | 常见 SMTP 服务商配置与官方文档

| Provider          | SMTP Server Address      | Port(s)      | Username/Description            | Password/Auth Method           | Official Documentation                                                                                                                                             | 服务商        | SMTP 服务器地址          | 端口         | 用户名/说明          | 密码/认证方式          | 官方文档                                                                                                                                                     |
| ----------------- | ------------------------ | ------------ | ------------------------------- | ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------- | ------------------------ | ------------ | -------------------- | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Apple iCloud      | smtp.mail.me.com         | 587/465      | Apple ID email                  | Apple ID password/app password | [Apple Docs](https://support.apple.com/en-us/HT202304)                                                                                                             | 苹果 iCloud   | smtp.mail.me.com         | 587/465      | Apple ID 邮箱        | Apple ID 密码/专用密码 | [Apple 文档](https://support.apple.com/zh-cn/HT202304)                                                                                                       |
| Gmail (Google)    | smtp.gmail.com           | 587/465      | Gmail address                   | App password/XOAUTH2           | [Gmail Docs](https://support.google.com/a/answer/176600?hl=en)                                                                                                     | 谷歌 Gmail    | smtp.gmail.com           | 587/465      | Gmail 邮箱           | 应用专用密码/XOAUTH2   | [Gmail 文档](https://support.google.com/mail/answer/7126229?hl=zh-Hans)                                                                                      |
| Outlook/Office365 | smtp.office365.com       | 587          | Email address                   | Login password/XOAUTH2         | [Microsoft Docs](https://support.microsoft.com/zh-hk/office/outlook-com-%E7%9A%84-pop-imap-%E5%92%8C-smtp-%E8%A8%AD%E5%AE%9A-d088b986-291d-42b8-9564-9c414e2aa040) | 微软 Outlook  | smtp.office365.com       | 587          | 邮箱地址             | 登录密码/XOAUTH2       | [微软文档](https://support.microsoft.com/zh-cn/office/outlook-com-%E7%9A%84-pop-imap-%E5%92%8C-smtp-%E8%AE%BE%E7%BD%AE-d088b986-291d-42b8-9564-9c414e2aa040) |
| Alibaba Cloud     | smtp.mxhichina.com       | 465/25       | Email address                   | Login password                 | [Alibaba Docs](https://help.aliyun.com/document_detail/36576.html)                                                                                                 | 阿里云        | smtp.mxhichina.com       | 465/25       | 邮箱地址             | 登录密码               | [阿里云文档](https://help.aliyun.com/document_detail/36576.html)                                                                                             |
| Tencent Exmail    | smtp.exmail.qq.com       | 465/587      | Email address                   | Login password                 | [Tencent Docs](https://www.tencentcloud.com/document/product/1084/44458)                                                                                           | 腾讯企业邮    | smtp.exmail.qq.com       | 465/587      | 邮箱地址             | 登录密码               | [腾讯文档](https://cloud.tencent.com/document/product/1084/44458)                                                                                            |
| QQ Mail           | smtp.qq.com              | 465/587      | QQ number or full email address | Auth code                      | [QQ Mail Docs](https://service.mail.qq.com/detail/0/428)                                                                                                           | QQ 邮箱       | smtp.qq.com              | 465/587      | QQ 号或完整邮箱      | 授权码                 | [QQ 邮箱文档](https://service.mail.qq.com/detail/0/428)                                                                                                      |
| 163 Mail          | smtp.163.com             | 465/994      | 163 email address               | Auth code                      | [163 Mail Docs](https://help.mail.163.com/faq.do?m=OTUw&id=MjQ5Nw==)                                                                                               | 网易 163 邮箱 | smtp.163.com             | 465/994      | 163 邮箱地址         | 授权码                 | [网易文档](https://help.mail.163.com/faq.do?m=OTUw&id=MjQ5Nw==)                                                                                              |
| SendGrid          | smtp.sendgrid.net        | 587/465      | apikey (literal string)         | API Key                        | [SendGrid Docs](https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp)                                                                       | SendGrid      | smtp.sendgrid.net        | 587/465      | apikey（固定字符串） | API Key                | [SendGrid 文档](https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp)                                                                 |
| Mailgun           | smtp.mailgun.org         | 587/465      | postmaster@yourdomain           | SMTP password                  | [Mailgun Docs](https://help.mailgun.com/hc/en-us/articles/203380100-How-Do-I-Use-Mailgun-SMTP-)                                                                    | Mailgun       | smtp.mailgun.org         | 587/465      | postmaster@域名      | SMTP 密码              | [Mailgun 文档](https://help.mailgun.com/hc/zh-cn/articles/203380100)                                                                                         |
| Mailtrap          | sandbox.smtp.mailtrap.io | 2525/465/587 | Mailtrap username               | Mailtrap password              | [Mailtrap Docs](https://help.mailtrap.io/article/122-mailtrap-email-sending-smtp-integration)                                                                      | Mailtrap      | sandbox.smtp.mailtrap.io | 2525/465/587 | Mailtrap 用户名      | Mailtrap 密码          | [Mailtrap 文档](https://help.mailtrap.io/article/122-mailtrap-email-sending-smtp-integration)                                                                |
| Zoho Mail         | smtp.zoho.com            | 465/587      | Zoho email address              | Login password/app password    | [Zoho Docs](https://www.zoho.com/mail/help/zoho-smtp.html)                                                                                                         | Zoho 邮箱     | smtp.zoho.com            | 465/587      | Zoho 邮箱地址        | 登录密码/专用密码      | [Zoho 文档](https://www.zoho.com/mail/help/zoho-smtp.html)                                                                                                   |
| Yandex            | smtp.yandex.com          | 465/587      | Yandex email address            | Login password                 | [Yandex Docs](https://yandex.com/support/mail/mail-clients/others.html)                                                                                            | Yandex 邮箱   | smtp.yandex.com          | 465/587      | Yandex 邮箱地址      | 登录密码               | [Yandex 文档](https://yandex.com/support/mail/mail-clients/others.html)                                                                                      |

### Example Configuration (SendGrid) | SendGrid 配置示例

```go
config := email.Config{
    Accounts: []email.Account{
        {
            Name:     "sendgrid",
            Host:     "smtp.sendgrid.net",
            Port:     587,
            Username: "apikey", // literal string | 固定字符串
            Password: "<your_sendgrid_api_key>",
            From:     "your@email.com",
        },
    },
}
```

### Example Configuration (QQ Mail) | QQ 邮箱配置示例

```go
config := email.Config{
    Accounts: []email.Account{
        {
            Name:     "qq",
            Host:     "smtp.qq.com",
            Port:     465,
            Username: "123456@qq.com",
            Password: "<your_qq_email_auth_code>", // Get the auth code from QQ Mail settings | QQ 邮箱设置中获取授权码
            From:     "123456@qq.com",
        },
    },
}
```

> For more provider settings, please refer to the table above and each provider's official documentation. | 更多服务商配置请参考上表及各自官方文档。

---

## API Documentation | 官方文档

- [Email Provider Guide | 邮件组件文档](https://github.com/shellvon/go-sender)

### 构造邮件消息 Constructing Email Messages

- 推荐使用 builder 风格（推荐）Recommended: Builder style
- 也可直接用 NewMessage（仅支持必填参数）Or use NewMessage (required fields only)

#### Builder API

- `Email()`：创建 builder 实例 Create a builder instance
- `To(recipients ...string)`：设置收件人 Set recipients
- `Subject(subject string)`：设置主题 Set subject
- `Body(body string)`：设置正文 Set body
- `From(from string)`：设置发件人 Set sender
- `Cc(cc ...string)`：设置抄送 Set CC
- `Bcc(bcc ...string)`：设置密送 Set BCC
- `ReplyTo(replyTo string)`：设置回复地址 Set reply-to
- `HTML()`：标记为 HTML 内容 Mark as HTML content
- `Attach(files ...string)`：替换附件列表 Set (replace) attachments
- `AddAttach(files ...string)`：追加附件 Append attachments
- `Build()`：生成 Message 实例 Build the Message instance

**示例 Example**

```go
// builder 风格 Builder style
msg := email.Email().
    To("a@b.com").
    Subject("Subject").
    Body("Body").
    From("noreply@b.com").
    Cc("c@b.com").
    Bcc("d@b.com").
    ReplyTo("reply@b.com").
    HTML().
    Attach("a.txt").
    AddAttach("b.txt", "c.pdf").
    Build()

// 只保留最后一次 SetAttach 的附件 Only the last SetAttach takes effect
email.Email().To("a@b.com").Attach("a.txt").Attach("b.txt").Build() // 只包含 b.txt only b.txt attached

// 追加多个附件 Append multiple attachments
email.Email().To("a@b.com").AddAttach("a.txt").AddAttach("b.txt", "c.pdf").Build() // 包含 a.txt, b.txt, c.pdf
```

#### NewMessage API

- `NewMessage(to []string, body string)`：仅支持必填参数 Only required fields supported

**示例 Example**

```go
msg := email.NewMessage([]string{"a@b.com"}, "Body")
```

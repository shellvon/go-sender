# WeCom Bot Provider

This provider supports sending messages to WeCom (Enterprise WeChat) group robots via webhooks.

[⬅️ Back to project README](../../README.md)

---

## Features

- **Multiple Account Support**: Configure multiple accounts with different strategies (round-robin, random, weighted).
- **Message Types**: Supports Text, Markdown, Image, News, Template Card, File, and Voice messages.
- **Media Upload**: Supports uploading media files (images, files, voice).
- **Mention Support**: Supports mentioning users or @all.

---

## Configuration Example

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// Create WeCom Bot configuration
config := wecombot.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // Round-robin, random, or weighted
    },
    Items: []*wecombot.Account{
        {
            Name:     "main",
            Key:      "YOUR_WEBHOOK_KEY", // Part after /key= in webhook URL
            Weight:   100,
            Disabled: false,
        },
        // ... more accounts
    },
}

provider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
```

---

## Message Types (Builder Style)

### 1. Text Message

```go
// Simple text message
msg := wecombot.Text().
    Content("Hello from go-sender!").
    Build()

// Text message with mentions
msg := wecombot.Text().
    Content("Hello @all, this is a test message").
    MentionUsers([]string{"@all"}).
    MentionMobiles([]string{"***REMOVED***"}).
    Build()
```

### 2. Markdown Message

```go
markdownContent := `# System Notice\n\n- **CPU**: 45%\n- **Memory**: 60%\n- **Disk**: 75%\n\n> System OK\n\n[View Details](https://example.com)`

msg := wecombot.Markdown().
    Content(markdownContent).
    Build()
```

### 3. Image Message

```go
// Send image with base64 and MD5
msg := wecombot.Image().
    Base64(imgBase64).
    MD5(imgMD5).
    Build()
```

### 4. News Message

```go
msg := wecombot.News().
    AddArticle("Important", "Description", "https://example.com", "https://example.com/image.jpg").
    Build()
```

### 5. Template Card Message

```go
msg := wecombot.Card(wecombot.CardTypeTextNotice).
    MainTitle("Main Title", "Description").
    SubTitle("Click for details").
    JumpURL("https://example.com").
    Build()
```

### 6. File Message

```go
// Send a file by local path
msg := wecombot.File().
    LocalPath("/path/to/report.pdf").
    Build()
```

### 7. Voice Message

```go
// Send a voice message by local path (AMR format, ≤2MB, ≤60s)
msg := wecombot.Voice().
    LocalPath("/path/to/voice.amr").
    Build()
```

---

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

s := gosender.NewSender()
wecomProvider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

ctx := context.Background()
msg := wecombot.File().LocalPath("/path/to/report.pdf").Build()
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## Media Upload Notes

- **File/Voice Upload**: Use `.LocalPath("/path/to/file")` in the builder. The SDK will upload the file/voice and fill in the media_id.
- **Image Upload**: Use `provider.UploadMedia(ctx, "image.jpg", file)` to get media_id, then use it in the builder.
- **Constraints**:
  - File: ≤20MB
  - Voice: ≤2MB, ≤60s, AMR format only
  - All media: >5 bytes
  - Media files are valid for 3 days and usable only by the uploading account.

---

## API Reference

### Config

- `BaseConfig`: Common configuration fields.
- `Accounts`: Array of account configurations.

### Message Types

- `TextMessage`: Text
- `MarkdownMessage`: Markdown
- `ImageMessage`: Image
- `NewsMessage`: News
- `TemplateCardMessage`: Template Card
- `FileMessage`: File
- `VoiceMessage`: Voice

### Builder Options

- `.MentionUsers(users []string)`: Mention users.
- `.MentionMobiles(mobiles []string)`: Mention by mobile.
- `.LocalPath(path string)`: Auto upload file/voice.
- `.MediaID(id string)`: Use existing media_id.

---

## Notes

- **Webhook Key**: Obtain from WeCom group robot settings.
- **Media Upload**: Media files valid for 3 days, only usable by the uploading account.
- **Mention**: Use `@all` to mention all.
- **Voice**: Only AMR format, ≤2MB, ≤60s.
- **File**: ≤20MB.

---

## API Documentation

- [WeCom Bot API Documentation](https://developer.work.weixin.qq.com/document/path/91770)

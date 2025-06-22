# WeCom Bot Provider

This provider supports sending messages to WeCom (企业微信) group robots via webhooks.

## Features

- **Multiple Account Support**: Configure multiple accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for text, markdown, image, news, and template card messages
- **Media Upload**: Support for uploading media files (images, documents)
- **Mention Support**: Support for mentioning users or @all

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// Create WeCom Bot configuration
config := wecombot.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_WEBHOOK_KEY",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup",
            Key:      "YOUR_BACKUP_KEY",
            Weight:   80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
```

## Message Types

### 1. Text Message

```go
// Simple text message
textMsg := wecombot.NewTextMessage("Hello from go-sender!")

// Text message with mentions
textMsg := wecombot.NewTextMessage("Hello everyone!",
    wecombot.WithMentionedList([]string{"@all"}),
    wecombot.WithMentionedMobileList([]string{"***REMOVED***"}),
)
```

### 2. Markdown Message

```go
markdownContent := `# 系统通知

## 状态报告
- **CPU 使用率**: 45%
- **内存使用率**: 60%
- **磁盘空间**: 75%

> 系统运行正常

[查看详情](https://example.com)`

markdownMsg := wecombot.NewMarkdownMessage(markdownContent)
```

### 3. Image Message

```go
// Upload image first
file, err := os.Open("image.jpg")
if err != nil {
    log.Fatalf("Failed to open image: %v", err)
}
defer file.Close()

mediaId, account, err := provider.UploadMedia(ctx, "image.jpg", file)
if err != nil {
    log.Fatalf("Failed to upload image: %v", err)
}

// Send image message
imageMsg := wecombot.NewImageMessage(mediaId)
```

### 4. News Message

```go
newsMsg := wecombot.NewNewsMessage([]wecombot.NewsArticle{
    {
        Title:       "重要通知",
        Description: "这是一条重要通知的描述",
        URL:         "https://example.com",
        PicURL:      "https://example.com/image.jpg",
    },
    {
        Title:       "系统更新",
        Description: "系统已更新到最新版本",
        URL:         "https://example.com/update",
        PicURL:      "https://example.com/update.jpg",
    },
})
```

### 5. Template Card Message

```go
cardMsg := wecombot.NewTemplateCardMessage().
    SetCardType(wecombot.CardTypeTextNotice).
    SetSource(&wecombot.CardSource{
        IconURL: "https://example.com/icon.png",
        Desc:    "消息来源",
    }).
    SetMainTitle(&wecombot.CardMainTitle{
        Title: "模板卡片标题",
        Desc:  "模板卡片描述",
    }).
    SetHorizontalContentList([]wecombot.CardHorizontalContent{
        {
            KeyName: "项目名称",
            Value:   "go-sender",
        },
        {
            KeyName: "版本",
            Value:   "1.0.0",
        },
    }).
    SetJumpList([]wecombot.CardJump{
        {
            Type:  wecombot.JumpTypeUrl,
            Title: "查看详情",
            URL:   "https://example.com",
        },
    }).
    SetCardAction(&wecombot.CardAction{
        Type: wecombot.ActionTypeUrl,
        URL:  "https://example.com",
    })
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// Create sender
s := gosender.NewSender(nil)

// Register WeCom Bot provider
wecomProvider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

// Send message
ctx := context.Background()
textMsg := wecombot.NewTextMessage("Hello from go-sender!")
err = s.Send(ctx, textMsg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

## Media Upload

WeCom Bot supports uploading media files for use in image and file messages:

```go
// Upload an image
file, err := os.Open("image.jpg")
if err != nil {
    log.Fatalf("Failed to open file: %v", err)
}
defer file.Close()

mediaId, account, err := provider.UploadMedia(ctx, "image.jpg", file)
if err != nil {
    log.Fatalf("Failed to upload media: %v", err)
}

// Use the mediaId in image message
imageMsg := wecombot.NewImageMessage(mediaId)
```

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of account configurations

### Account

- `Name`: Account name for identification
- `Key`: WeCom webhook key (get from WeCom group robot settings)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled

### Message Types

- `TextMessage`: Simple text messages with mention support
- `MarkdownMessage`: Rich text messages with Markdown formatting
- `ImageMessage`: Image messages (requires media upload)
- `NewsMessage`: News article messages
- `TemplateCardMessage`: Interactive template card messages

### Message Options

- `WithMentionedList(users []string)`: Mention specific users
- `WithMentionedMobileList(mobiles []string)`: Mention users by mobile number
- `WithMentionedEmailList(emails []string)`: Mention users by email

## Notes

- **Webhook Key**: Get your webhook key from WeCom group robot settings
- **Media Upload**: Media files are valid for 3 days and can only be used by the account that uploaded them
- **Mention Support**: Use `@all` to mention all group members
- **Markdown**: Supports standard Markdown syntax
- **Template Cards**: Support various card types for rich interactions

## API Documentation

For detailed API documentation, visit:

- [WeCom Bot API Documentation](https://developer.work.weixin.qq.com/document/path/91770)
- [WeCom Group Robot](https://developer.work.weixin.qq.com/document/path/91770)

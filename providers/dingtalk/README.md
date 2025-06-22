# DingTalk Provider

This provider supports sending messages to DingTalk (钉钉) group robots via webhooks.

## Features

- **Multiple Account Support**: Configure multiple accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for text, markdown, link, action card, and feed card messages
- **Rich Formatting**: Support for Markdown formatting in text and markdown messages
- **Interactive Cards**: Support for action cards with buttons and feed cards with multiple articles

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

// Create DingTalk configuration
config := dingtalk.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_ACCESS_TOKEN",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup",
            Key:      "YOUR_BACKUP_TOKEN",
            Weight:   80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := dingtalk.New(config)
if err != nil {
    log.Fatalf("Failed to create DingTalk provider: %v", err)
}
```

## Message Types

### 1. Text Message

```go
// Simple text message
textMsg := dingtalk.NewTextMessage("Hello from go-sender!")

// Text message with @ mentions
textMsg := dingtalk.NewTextMessage("Hello @all, this is a test message",
    dingtalk.WithAtMobiles([]string{"13800138000"}),
    dingtalk.WithAtUserIds([]string{"user123"}),
    dingtalk.WithIsAtAll(true),
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

markdownMsg := dingtalk.NewMarkdownMessage("系统通知", markdownContent,
    dingtalk.WithAtMobiles([]string{"13800138000"}),
    dingtalk.WithAtUserIds([]string{"user123"}),
    dingtalk.WithIsAtAll(false),
)
```

### 3. Link Message

```go
linkMsg := dingtalk.NewLinkMessage(
    "重要通知",
    "这是一条重要通知的摘要内容",
    "https://example.com",
    "https://example.com/image.jpg",
)
```

### 4. Action Card Message

```go
// Single action card
singleCardMsg := dingtalk.NewActionCardMessage(
    "系统通知",
    "系统已更新到最新版本，请查看详情",
    "查看详情",
    "https://example.com",
    dingtalk.WithBtnOrientation("0"), // 0: vertical, 1: horizontal
)

// Multiple action card
multipleCardMsg := dingtalk.NewActionCardMessage(
    "系统通知",
    "系统已更新到最新版本，请选择操作",
    "", // No single button
    "",
    dingtalk.WithBtnOrientation("1"),
    dingtalk.WithActionButtons([]dingtalk.ActionButton{
        {
            Title:     "查看详情",
            ActionURL: "https://example.com/details",
        },
        {
            Title:     "下载更新",
            ActionURL: "https://example.com/download",
        },
    }),
)
```

### 5. Feed Card Message

```go
feedCardMsg := dingtalk.NewFeedCardMessage([]dingtalk.FeedCardLink{
    {
        Title:      "重要通知",
        MessageURL: "https://example.com/notice1",
        PicURL:     "https://example.com/image1.jpg",
    },
    {
        Title:      "系统更新",
        MessageURL: "https://example.com/notice2",
        PicURL:     "https://example.com/image2.jpg",
    },
})
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

// Create sender
s := gosender.NewSender(nil)

// Register DingTalk provider
dingtalkProvider, err := dingtalk.New(config)
if err != nil {
    log.Fatalf("Failed to create DingTalk provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeDingtalk, dingtalkProvider, nil)

// Send message
ctx := context.Background()
textMsg := dingtalk.NewTextMessage("Hello from go-sender!")
err = s.Send(ctx, textMsg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

## Message Options

### Text Message Options

- `WithAtMobiles(mobiles []string)`: @ mention users by mobile number
- `WithAtUserIds(userIds []string)`: @ mention users by user ID
- `WithIsAtAll(atAll bool)`: @ mention all users

### Markdown Message Options

- `WithAtMobiles(mobiles []string)`: @ mention users by mobile number
- `WithAtUserIds(userIds []string)`: @ mention users by user ID
- `WithIsAtAll(atAll bool)`: @ mention all users

### Action Card Message Options

- `WithBtnOrientation(orientation string)`: Set button orientation ("0": vertical, "1": horizontal)
- `WithActionButtons(buttons []ActionButton)`: Set multiple action buttons
- `WithHideAvatar(hide string)`: Hide avatar ("0": show, "1": hide)

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of account configurations

### Account

- `Name`: Account name for identification
- `Key`: DingTalk access token (get from DingTalk group robot settings)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled

### Message Types

- `TextMessage`: Simple text messages with @ mentions
- `MarkdownMessage`: Rich text messages with Markdown formatting
- `LinkMessage`: Link messages with title, text, and image
- `ActionCardMessage`: Interactive cards with buttons
- `FeedCardMessage`: Feed cards with multiple articles

### Message Validation

- Text content: Required, supports @ mentions
- Markdown content: Required, supports Markdown syntax
- Link content: Title, text, message URL, and picture URL required
- Action card: Title and text required, buttons optional
- Feed card: At least one link required

## Notes

- **Access Token**: Get your access token from DingTalk group robot settings
- **@ Mentions**: Support for mentioning users by mobile number or user ID
- **Markdown**: Supports standard Markdown syntax
- **Action Cards**: Can have single button or multiple buttons
- **Feed Cards**: Support multiple articles in one message
- **Rate Limits**: Respect DingTalk's rate limits

## API Documentation

For detailed API documentation, visit:

- [DingTalk Bot API Documentation](https://open.dingtalk.com/document/robots/custom-robot-access)
- [DingTalk Bot Message Types](https://open.dingtalk.com/document/robots/custom-robot-access)

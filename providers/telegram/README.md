# Telegram Provider

This provider supports sending messages via Telegram Bot API.

## Features

- **Multiple Account Support**: Configure multiple bot accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for text, photo, document, location, contact, and poll messages
- **Rich Formatting**: Support for Markdown and HTML formatting
- **Interactive Features**: Support for polls, location sharing, and contact sharing

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

// Create Telegram configuration
config := telegram.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "main-bot",
            Key:      "YOUR_BOT_TOKEN",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup-bot",
            Key:      "YOUR_BACKUP_BOT_TOKEN",
            Weight:   80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err)
}
```

## Message Types

### 1. Text Message

```go
// Simple text message
textMsg := telegram.NewTextMessage("@your_channel", "Hello from go-sender!")

// Text message with Markdown formatting
textMsg := telegram.NewTextMessage("@your_channel",
    "*Bold text*\n_Italic text_\n`Code text`\n[Link](https://example.com)",
    telegram.WithParseMode("MarkdownV2"),
)

// Text message with options
textMsg := telegram.NewTextMessage("@your_channel", "Important message",
    telegram.WithDisableWebPreview(true),
    telegram.WithSilent(true),
    telegram.WithProtectContent(true),
    telegram.WithReplyTo(12345),
)
```

### 2. Photo Message

```go
// Photo message with caption
photoMsg := telegram.NewPhotoMessage("@your_channel", "https://example.com/image.jpg",
    telegram.WithPhotoCaption("Beautiful image!"),
    telegram.WithPhotoParseMode("MarkdownV2"),
    telegram.WithPhotoSilent(true),
)
```

### 3. Document Message

```go
// Document message
docMsg := telegram.NewDocumentMessage("@your_channel", "https://example.com/document.pdf",
    telegram.WithDocumentCaption("Important document"),
    telegram.WithDocumentParseMode("MarkdownV2"),
)
```

### 4. Location Message

```go
// Location message
locationMsg := telegram.NewLocationMessage("@your_channel", 40.7128, -74.0060,
    telegram.WithLocationSilent(true),
)
```

### 5. Contact Message

```go
// Contact message
contactMsg := telegram.NewContactMessage("@your_channel", "+1234567890", "John Doe",
    telegram.WithContactLastName("Smith"),
    telegram.WithContactVCard("BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+1234567890\nEND:VCARD"),
)
```

### 6. Poll Message

```go
// Poll message
pollMsg := telegram.NewPollMessage("@your_channel", "What's your favorite programming language?",
    []string{"Go", "Python", "JavaScript", "Rust"},
    telegram.WithPollIsAnonymous(false),
    telegram.WithPollType("quiz"),
    telegram.WithPollAllowsMultipleAnswers(false),
)
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

// Create sender
s := gosender.NewSender(nil)

// Register Telegram provider
telegramProvider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeTelegram, telegramProvider, nil)

// Send message
ctx := context.Background()
textMsg := telegram.NewTextMessage("@your_channel", "Hello from go-sender!")
err = s.Send(ctx, textMsg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

## Message Options

### Text Message Options

- `WithParseMode(mode string)`: Set parse mode ("MarkdownV2", "HTML")
- `WithDisableWebPreview(disable bool)`: Disable web page preview
- `WithSilent(silent bool)`: Send message silently
- `WithProtectContent(protect bool)`: Protect message content
- `WithReplyTo(replyTo int)`: Reply to specific message

### Photo Message Options

- `WithPhotoCaption(caption string)`: Set photo caption
- `WithPhotoParseMode(mode string)`: Set caption parse mode
- `WithPhotoSilent(silent bool)`: Send photo silently
- `WithPhotoProtectContent(protect bool)`: Protect photo content
- `WithPhotoReplyTo(replyTo int)`: Reply to specific message

### Document Message Options

- `WithDocumentCaption(caption string)`: Set document caption
- `WithDocumentParseMode(mode string)`: Set caption parse mode
- `WithDocumentSilent(silent bool)`: Send document silently
- `WithDocumentProtectContent(protect bool)`: Protect document content
- `WithDocumentReplyTo(replyTo int)`: Reply to specific message

### Location Message Options

- `WithLocationSilent(silent bool)`: Send location silently
- `WithLocationProtectContent(protect bool)`: Protect location content
- `WithLocationReplyTo(replyTo int)`: Reply to specific message

### Contact Message Options

- `WithContactLastName(lastName string)`: Set contact last name
- `WithContactVCard(vcard string)`: Set contact vCard
- `WithContactSilent(silent bool)`: Send contact silently
- `WithContactProtectContent(protect bool)`: Protect contact content
- `WithContactReplyTo(replyTo int)`: Reply to specific message

### Poll Message Options

- `WithPollIsAnonymous(anonymous bool)`: Set poll as anonymous
- `WithPollType(pollType string)`: Set poll type ("quiz", "regular")
- `WithPollAllowsMultipleAnswers(allow bool)`: Allow multiple answers
- `WithPollSilent(silent bool)`: Send poll silently
- `WithPollProtectContent(protect bool)`: Protect poll content
- `WithPollReplyTo(replyTo int)`: Reply to specific message

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of account configurations

### Account

- `Name`: Account name for identification
- `Key`: Telegram Bot Token (get from @BotFather)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled

### Message Types

- `TextMessage`: Simple text messages with formatting
- `PhotoMessage`: Photo messages with captions
- `DocumentMessage`: Document/file messages
- `LocationMessage`: Location sharing messages
- `ContactMessage`: Contact sharing messages
- `PollMessage`: Poll/quiz messages

## Notes

- **Bot Token**: Get your bot token from [@BotFather](https://t.me/botfather) on Telegram
- **Chat ID**: Use channel username (e.g., "@your_channel") or chat ID
- **Parse Modes**: Support for MarkdownV2 and HTML formatting
- **File URLs**: For photos and documents, use direct URLs or file IDs
- **Poll Limits**: Polls can have 2-10 options
- **Rate Limits**: Respect Telegram's rate limits

## API Documentation

For detailed API documentation, visit:

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [Telegram Bot API Methods](https://core.telegram.org/bots/api#available-methods)

# Telegram Provider

> Unified messaging via Telegram Bot API for chats, groups, and channels.

[⬅️ Back to project README](../../README.md)

---

## Supported Providers

| Provider     | Website                            |
| ------------ | ---------------------------------- |
| **Telegram** | https://core.telegram.org/bots/api |

---

## Capabilities

| Provider | Text | Media | Location | Contact | Poll | Dice | Notes                              |
| -------- | ---- | ----- | -------- | ------- | ---- | ---- | ---------------------------------- |
| Telegram | ✅   | ✅    | ✅       | ✅      | ✅   | ✅   | Supports all Bot API message types |

---

## Features

- Multiple bot accounts with load-balancing strategies (round-robin, random, weighted).
- Builder API for message types (Text, Photo, Audio, Poll, etc.).
- Rich text formatting (HTML, Markdown).
- File support via file_id or public HTTP URLs.
- Interactive elements: polls, dice, custom keyboards.

---

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

cfg := telegram.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Items: []*telegram.Account{
        {
            Name:   "test",
            APIKey: "bot123:token",
        },
    },
}
```

---

## Quick Builder

```go
msg := telegram.Text().
    Chat("@channel").
    Text("Hello from go-sender!").
    Build()
```

---

## Usage

### 1. Direct Provider

```go
provider, _ := telegram.New(&cfg)
_ = provider.Send(context.Background(), msg, nil)
```

### 2. Using GoSender

```go
sender := gosender.NewSender()
provider, _ := telegram.New(&cfg)
sender.RegisterProvider(core.ProviderTypeTelegram, provider, nil)
_ = sender.Send(context.Background(), msg)
```

---

## SendVia Helper

`SendVia(accountName, msg)` lets you choose a specific bot account at runtime:

```go
msg := telegram.Text().
    Chat("@channel").
    Text("Hello from go-sender!").
    Build()

// try primary bot account first
if err := sender.SendVia("main-bot", msg); err != nil {
    // fallback to backup bot account
    _ = sender.SendVia("backup-bot", msg)
}
```

SendVia only switches between accounts **inside the Telegram provider**; it does not allow cross-provider reuse of one message instance.

---

## Message Types

### 1. Text Message

```go
msg := telegram.Text().
    Chat("@channel").
    Text("Hello from go-sender!").
    ParseMode("Markdown").
    Build()
```

### 2. Photo Message

```go
msg := telegram.Photo().
    Chat("@channel").
    File("https://example.com/image.jpg").
    Caption("Beautiful image").
    ParseMode("HTML").
    Build()
```

### 3. Audio Message

```go
msg := telegram.Audio().
    Chat("@channel").
    File("https://example.com/audio.mp3").
    Title("Song Title").
    Performer("Artist Name").
    Duration(180).
    Build()
```

### 4. Poll Message

```go
msg := telegram.Poll().
    Chat("@channel").
    Question("What's your favorite color?").
    Options(
        telegram.Option("Option 1"),
        telegram.Option("Option 2"),
        telegram.Option("Option 3"),
    ).
    IsAnonymous(false).
    AllowsMultipleAnswers(true).
    Build()
```

---

## Notes

- **Bot Token**: Obtain from [BotFather](https://core.telegram.org/bots#botfather).
- **File Upload**: Supports file_id or public HTTP URLs; local file upload not supported.
- **Formatting**: Use `ParseMode("Markdown")` or `ParseMode("HTML")` for rich text.
- **Interactive Elements**: Use `telegram.Poll()`, `telegram.Dice()`, etc., for polls, dice, and keyboards.
- **API Reference**: See [Telegram Bot API Documentation](https://core.telegram.org/bots/api) for advanced options.

---

## API Documentation

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Sending Files Guide](https://core.telegram.org/bots/api#sending-files)
- [Bot Creation Guide](https://core.telegram.org/bots#how-do-i-create-a-bot)

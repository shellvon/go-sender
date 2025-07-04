[⬅️ Back to Main README](../../README.md)

# Telegram Provider | Telegram 机器人推送组件

This provider supports sending messages to Telegram chats, groups, and channels via the [Telegram Bot API](https://core.telegram.org/bots/api).

本组件支持通过 [Telegram Bot API](https://core.telegram.org/bots/api) 向 Telegram 群组、频道、私聊发送消息。

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple bot accounts with different strategies (round-robin, random, weighted) | 支持多账号、灵活负载均衡（轮询、随机、加权）
- **Message Types 消息类型**: Support for all Telegram Bot API message types including text, media, location, contact, poll, dice, etc. | 支持所有官方消息类型：文本、媒体、位置、联系人、投票、骰子等
- **Rich Formatting 富文本格式**: Support for HTML and Markdown formatting in text messages | 文本消息支持 HTML/Markdown 富文本格式
- **File Support 文件支持**: Support for sending files via file_id or public HTTP URLs | 支持通过 file_id 或公网 URL 发送文件
- **Interactive Elements 交互元素**: Support for polls, dice, custom keyboards, etc. | 支持投票、骰子、键盘等交互元素

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

// English: Create Telegram configuration
// 中文：创建 Telegram 配置
config := telegram.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机、加权等
    },
    Accounts: []core.Account{
        {
            Name:     "main-bot",
            Key:      "YOUR_BOT_TOKEN", // 机器人 token
            Weight:   100,
            Disabled: false,
        },
        // ... more accounts
    },
}

provider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err) // 创建失败
}
```

---

## Message Types (Builder Style) | 消息类型（链式构建）

### 1. Text Message | 文本消息

```go
// English: Simple text message
// 中文：简单文本消息
msg := telegram.Text().
    Chat("@channel").
    Text("Hello from go-sender! 你好，世界！").
    Build()

// English: Text message with Markdown formatting
// 中文：带 Markdown 格式的文本消息
msg := telegram.Text().
    Chat("@channel").
    Text("**Bold text** and *italic text*").
    ParseMode("Markdown").
    Build()

// English: Text message with entities
// 中文：带实体的文本消息
msg := telegram.Text().
    Chat("@channel").
    Text("Hello @username").
    Entities([]telegram.MessageEntity{
        {Type: "mention", Offset: 6, Length: 9},
    }).
    Build()
```

### 2. Photo Message | 图片消息

```go
// English: Photo from URL
// 中文：通过 URL 发送图片
msg := telegram.Photo().
    Chat("@channel").
    File("https://example.com/image.jpg").
    Caption("Beautiful image").
    ParseMode("HTML").
    Build()

// English: Photo from file_id
// 中文：通过 file_id 发送图片
msg := telegram.Photo().
    Chat("@channel").
    File("AgACAgIAAxkBAAIB...").
    Caption("Reused image").
    HasSpoiler(true).
    Build()
```

### 3. Audio Message | 音频消息

```go
// English: Audio from URL
// 中文：通过 URL 发送音频
msg := telegram.Audio().
    Chat("@channel").
    File("https://example.com/audio.mp3").
    Title("Song Title").
    Performer("Artist Name").
    Duration(180).
    Build()

// English: Audio from file_id
// 中文：通过 file_id 发送音频
msg := telegram.Audio().
    Chat("@channel").
    File("CQACAgIAAxkBAAIB...").
    Caption("Listen to this!").
    Build()
```

### 4. Poll Message | 投票消息

```go
// English: Regular poll
// 中文：普通投票
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

// English: Quiz poll
// 中文：测验投票
msg := telegram.Poll().
    Chat("@channel").
    Question("What is 2+2?").
    Options(
        telegram.Option("3"),
        telegram.Option("4"),
        telegram.Option("5"),
    ).
    Type("quiz").
    CorrectOptionID(1).
    Explanation("The correct answer is 4").
    Build()
```

> **Note**: Other message types can be added similarly using their respective builders. See source code for all available builders.
>
> **注意**：其他类型消息可参考源码使用对应 builder。所有可用的 builder 请参考源码。

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

s := gosender.NewSender(nil)
telegramProvider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeTelegram, telegramProvider, nil)

ctx := context.Background()
msg := telegram.Text().Chat("@channel").Text("Hello from go-sender! 你好，世界！").Build()
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## API Reference | API 参考

- Each message type (Text, Photo, Audio, Poll, etc.) has a corresponding builder: `telegram.Text()`, `telegram.Photo()`, `telegram.Audio()`, `telegram.Poll()`, etc. | 每种消息类型都有对应 builder，如 `telegram.Text()`、`telegram.Photo()`、`telegram.Audio()`、`telegram.Poll()` 等
- All builders support chainable methods for setting fields, and end with `.Build()` to produce the message object. | 所有 builder 支持链式设置参数，最后 `.Build()` 生成消息对象
- For more advanced options, see the GoDoc or source code for each builder. | 更多高级用法请参考 GoDoc 或源码

---

## Notes | 说明

- **Bot Token**: Get your bot token from [BotFather](https://core.telegram.org/bots#botfather) | 机器人 token 请通过 [BotFather](https://core.telegram.org/bots#botfather) 获取
- **File Upload**: For files, you can use file_id or HTTP URL | 文件支持 file_id 或公网 URL，暂不支持本地文件直传
- **Markdown/HTML**: Use `ParseMode("Markdown")` or `ParseMode("HTML")` for rich formatting | 富文本格式请用 `ParseMode("Markdown")` 或 `ParseMode("HTML")`
- **Polls**: Use `telegram.Poll()` builder for regular and quiz polls | 投票请用 `telegram.Poll()` builder
- **Dice/Other Types**: See source code for additional builder types | 骰子等其他类型请参考源码

---

## API Documentation | 官方文档

- [Telegram Bot API Documentation | 官方文档](https://core.telegram.org/bots/api)
- [Sending Files Guide | 文件发送指南](https://core.telegram.org/bots/api#sending-files)
- [Bot Creation Guide | 机器人创建指南](https://core.telegram.org/bots#how-do-i-create-a-bot)

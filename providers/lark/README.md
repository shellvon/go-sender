[⬅️ Back to Main README](../../README.md)

# Lark / Feishu Provider | Lark / 飞书 消息推送组件

This provider supports sending messages to [Lark (飞书)](https://www.larksuite.com/) / [Feishu](https://www.feishu.cn/) group bots via webhooks.

本组件支持通过 webhook 向 [Lark (飞书)](https://www.larksuite.com/) / [Feishu](https://www.feishu.cn/) 群机器人发送消息。

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple bots/accounts with flexible load balancing | 支持多账号、灵活负载均衡
- **Message Types 消息类型**: Text, Post (rich text), Image, **Interactive Card (schema 2.0, 推荐)**, Share Chat | 文本、富文本、图片、卡片（推荐）、群分享
- **Internationalization 国际化**: Post and Card messages support both Chinese and English content | 富文本和卡片支持中英文内容
- **Builder API 构建器风格**: All messages use a modern, chainable builder pattern | 所有消息均为链式 builder 构建

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/lark"
)

config := lark.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机、加权等
    },
    Accounts: []core.Account{
        {
            Name:    "main-bot",
            Key:     "your-webhook-key", // webhook URL 末尾 /hook/ 后的部分
            Weight:  100,
            Disabled: false,
        },
        // 可添加更多账号
    },
}

provider, err := lark.New(config)
if err != nil {
    log.Fatalf("Failed to create Lark provider: %v", err) // 创建 Lark provider 失败
}
```

---

## Message Types (Builder Style) | 消息类型（链式构建）

### 1. Text Message | 文本消息

```go
msg := lark.Text().
    Content("Hello from go-sender! 你好，世界！").
    Build()
```

### 2. Post (Rich Text) Message | 富文本消息

```go
msg := lark.Post().
    ZhCN("测试标题", [][]lark.PostElement{
        {
            {Tag: "text", Text: "这是一条测试消息"},
        },
        {
            {Tag: "a", Text: "点击这里", Href: "https://www.feishu.cn"},
        },
    }).
    EnUS("Test Title", [][]lark.PostElement{
        {
            {Tag: "text", Text: "This is a test message"},
        },
        {
            {Tag: "a", Text: "Click here", Href: "https://www.larksuite.com"},
        },
    }).
    Build()
```

### 3. Image Message | 图片消息

```go
msg := lark.Image().
    ImageKey("img_1234567890abcdef").
    Build()
```

### 4. Interactive Card Message (schema 2.0, 推荐) | 交互卡片消息（schema 2.0，推荐）

```go
msg := lark.Interactive().
    HeaderTitle("plain_text", "主标题 Main Title").
    HeaderSubtitle("plain_text", "副标题 Subtitle").
    HeaderTemplate("blue").
    AddElement(divElement).
    AddElement(actionElement).
    BodyDirection("vertical").
    Build()
```

- You can use the builder to flexibly set config, card_link, header, body, elements, i18n, style, etc. | 你可以通过 builder 灵活设置 config、card_link、header、body、elements、i18n、style 等所有 schema 2.0 支持的字段。
- elements support any JSON structure, covering all official components. | elements 支持任意 JSON 结构，满足所有官方组件。
- See: [Card JSON 2.0 Structure (Official)](https://open.feishu.cn/document/feishu-cards/card-json-v2-structure) | 详见：[官方卡片 JSON 2.0 结构文档](https://open.feishu.cn/document/feishu-cards/card-json-v2-structure)

### 5. Share Chat Message | 群分享消息

```go
msg := lark.ShareChat().
    ChatID("oc_1234567890abcdef").
    Build()
```

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/lark"
)

s := gosender.NewSender(nil)
larkProvider, err := lark.New(config)
if err != nil {
    log.Fatalf("Failed to create Lark provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeLark, larkProvider, nil)

ctx := context.Background()
msg := lark.Text().Content("Hello from go-sender! 你好，世界！").Build()
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## API Reference | API 参考

### Config | 配置

- `BaseConfig`: Common configuration fields (strategy, disabled, etc.) | 通用配置字段（策略、禁用等）
- `Accounts`: List of bot accounts. Each account: | 账号列表，每个账号：
  - `Name`: Logical name for the bot | 逻辑名称
  - `Key`: The webhook key (from Lark/Feishu bot settings) | webhook key（飞书机器人设置中获取）
  - `Weight`: Used for weighted load balancing | 加权负载均衡
  - `Disabled`: Whether this account is disabled | 是否禁用

### Message Types | 消息类型

- `TextMessage`: Plain text | 文本消息
- `PostMessage`: Rich text (multi-language) | 富文本（多语言）
- `ImageMessage`: Image (requires image_key) | 图片（需 image_key）
- `InteractiveMessage`: **Card (schema 2.0, 推荐，支持所有官方组件)** | 卡片（schema 2.0，推荐，支持所有官方组件）
- `ShareChatMessage`: Share a group chat | 群分享

All messages implement `core.Message` and have `Validate()` and `ProviderType()` methods. | 所有消息实现 `core.Message` 接口，包含 `Validate()` 和 `ProviderType()` 方法。

---

## Notes | 说明

- **Webhook Key**: Get your key from Lark/Feishu bot settings (the part after `/hook/`) | webhook key 请在飞书机器人设置中获取（/hook/ 后的部分）
- **Image Key**: For image messages, upload the image to Lark/Feishu to get the `image_key` | 图片消息需先上传图片获取 image_key
- **Chat ID**: For share chat messages, use the correct IDs from Lark/Feishu | 分享消息请使用正确的 chat_id
- **Interactive Card**: 推荐使用 schema 2.0 卡片，支持所有新特性和组件，详见官方文档 | 推荐使用 schema 2.0 卡片，详见官方文档

---

## 最新官方文档 | Official Documentation

- [Lark/Feishu Bot API](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- [卡片消息 JSON 2.0 结构（推荐） | Card JSON 2.0 Structure](https://open.feishu.cn/document/feishu-cards/card-json-v2-structure)
- [卡片消息开发指引 | Card Message Dev Guide](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot#5a997364)

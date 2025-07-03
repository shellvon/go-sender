[⬅️ Back to Main README](../../README.md)

# DingTalk Provider | 钉钉机器人推送组件

This provider supports sending messages to DingTalk (钉钉) group robots via webhooks.

本组件支持通过 webhook 向钉钉（DingTalk）群机器人发送消息。

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple accounts with different strategies (round-robin, random, weighted) | 支持多账号、灵活负载均衡（轮询、随机、加权）
- **Message Types 消息类型**: Text, Markdown, Link, Action Card, Feed Card | 文本、Markdown、链接、动作卡片、Feed 卡片
- **Rich Formatting 富文本格式**: Support for Markdown formatting in text and markdown messages | 文本和 Markdown 消息支持 Markdown 格式
- **Interactive Cards 交互卡片**: Support for action cards with buttons and feed cards with multiple articles | 支持带按钮的动作卡片和多条内容的 Feed 卡片

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

// English: Create DingTalk configuration
// 中文：创建钉钉机器人配置
config := dingtalk.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机、加权等
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_ACCESS_TOKEN", // webhook URL 末尾 /access_token= 后的部分
            Weight:   100,
            Disabled: false,
        },
        // ... more accounts
    },
}

provider, err := dingtalk.New(config)
if err != nil {
    log.Fatalf("Failed to create DingTalk provider: %v", err) // 创建失败
}
```

---

## Message Types (Builder Style) | 消息类型（链式构建）

### 1. Text Message | 文本消息

```go
// English: Simple text message
// 中文：简单文本消息
msg := dingtalk.Text().
    Content("Hello from go-sender! 你好，世界！").
    Build()

// English: Text message with @ mentions
// 中文：带@的文本消息
msg := dingtalk.Text().
    Content("Hello @all, this is a test message").
    AtMobiles([]string{"***REMOVED***"}).
    AtUserIDs([]string{"user123"}).
    AtAll().
    Build()
```

### 2. Markdown Message | Markdown 消息

```go
markdownContent := `# 系统通知 System Notice\n\n- **CPU**: 45%\n- **Memory**: 60%\n- **Disk**: 75%\n\n> 系统运行正常 System OK\n\n[查看详情 View Details](https://example.com)`

msg := dingtalk.Markdown().
    Title("系统通知 System Notice").
    Text(markdownContent).
    AtMobiles([]string{"***REMOVED***"}).
    AtUserIDs([]string{"user123"}).
    Build()
```

### 3. Link Message | 链接消息

```go
// English: Link message
// 中文：链接消息
msg := dingtalk.Link().
    Title("重要通知 Important Notice").
    Text("这是一条重要通知的摘要内容 This is a summary of an important notice").
    MessageURL("https://example.com").
    PicURL("https://example.com/image.jpg").
    Build()
```

### 4. Action Card Message | 动作卡片消息

```go
// English: Single action card
// 中文：单按钮动作卡片
msg := dingtalk.ActionCard().
    Title("系统通知 System Notice").
    Text("系统已更新到最新版本，请查看详情 The system has been updated, please check details").
    SingleButton("查看详情 View Details", "https://example.com").
    BtnOrientation("0"). // 0: vertical, 1: horizontal
    Build()

// English: Multiple action card
// 中文：多按钮动作卡片
msg := dingtalk.ActionCard().
    Title("系统通知 System Notice").
    Text("系统已更新到最新版本，请选择操作 The system has been updated, please choose an action").
    BtnOrientation("1").
    AddButton("查看详情 View Details", "https://example.com/details").
    AddButton("下载更新 Download Update", "https://example.com/download").
    Build()
```

### 5. Feed Card Message | Feed 卡片消息

```go
// English: Feed card message
// 中文：Feed 卡片消息
msg := dingtalk.FeedCard().
    AddLink("重要通知 Important Notice", "https://example.com/notice1", "https://example.com/image1.jpg").
    AddLink("系统更新 System Update", "https://example.com/notice2", "https://example.com/image2.jpg").
    Build()
```

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/dingtalk"
)

s := gosender.NewSender(nil)
dingtalkProvider, err := dingtalk.New(config)
if err != nil {
    log.Fatalf("Failed to create DingTalk provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeDingtalk, dingtalkProvider, nil)

ctx := context.Background()
msg := dingtalk.Text().Content("Hello from go-sender! 你好，世界！").Build()
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## API Reference | API 参考

### Config | 配置

- `BaseConfig`: Common configuration fields | 通用配置字段
- `Accounts`: Array of account configurations | 账号配置数组

### Message Types | 消息类型

- `TextMessage`: Text | 文本
- `MarkdownMessage`: Markdown
- `LinkMessage`: Link | 链接
- `ActionCardMessage`: Action Card | 动作卡片
- `FeedCardMessage`: Feed Card | Feed 卡片

### Builder Options | 构建器选项

- `.AtMobiles([]string)`: @手机号 | Mention by mobile
- `.AtUserIDs([]string)`: @用户 ID | Mention by user ID
- `.AtAll()`: @所有人 | Mention all

---

## Notes | 说明

- **Access Token**: Get from DingTalk group robot settings | access_token 请在钉钉群机器人设置中获取
- **@ Mentions**: Support for mentioning users by mobile number or user ID | 支持通过手机号或用户 ID @成员
- **Markdown**: Supports standard Markdown syntax | 支持标准 Markdown 语法
- **Action Cards**: Can have single button or multiple buttons | 动作卡片支持单按钮或多按钮
- **Feed Cards**: Support multiple articles in one message | Feed 卡片支持多条内容
- **Rate Limits**: Respect DingTalk's rate limits | 请遵守钉钉官方限流规则

---

## API Documentation | 官方文档

- [DingTalk Bot API Documentation | 钉钉机器人官方文档](https://open.dingtalk.com/document/robots/custom-robot-access)
- [DingTalk Bot Message Types | 钉钉消息类型官方文档](https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type)

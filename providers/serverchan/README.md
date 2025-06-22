# ServerChan Provider

This provider supports sending messages via ServerChan (方糖/Server 酱) service.

## Features

- **Multiple Account Support**: Configure multiple accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for text messages with Markdown formatting
- **Channel Support**: Support for multiple push channels (WeCom, DingTalk, Lark, etc.)
- **Message Options**: Short content, IP hiding, channel selection, and more

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/serverchan"
)

// Create ServerChan configuration
config := serverchan.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_SENDKEY_HERE",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup",
            Key:      "YOUR_BACKUP_SENDKEY",
            Weight:   80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := serverchan.New(config)
if err != nil {
    log.Fatalf("Failed to create ServerChan provider: %v", err)
}
```

## Message Types

### Basic Text Message

```go
// Simple text message
msg := serverchan.NewMessage("测试标题", "这是消息内容")
```

### Message with Options

```go
// Message with short content and channel selection
msg := serverchan.NewMessage(
    "系统通知",
    "## 系统状态报告\n\n- CPU 使用率: 45%\n- 内存使用率: 60%\n- 磁盘空间: 75%\n\n**状态**: 正常",
    serverchan.WithShort("系统运行正常"),
    serverchan.WithChannel("wecom|dingtalk"), // 企业微信 + 钉钉
    serverchan.WithNoIP(), // 隐藏发送 IP
)
```

### Long Text Message

```go
longContent := `# 详细报告

## 项目状态
这是一个使用 go-sender 库发送的详细报告。

### 功能特性
- ✅ 支持多种消息类型
- ✅ 支持多渠道推送
- ✅ 支持 Markdown 格式
- ✅ 支持自定义配置

### 技术栈
- Go 语言
- ServerChan API
- HTTP 客户端

---
*此消息由 go-sender 自动发送*`

msg := serverchan.NewMessage(
    "详细报告",
    longContent,
    serverchan.WithShort("项目状态报告"),
)
```

## Supported Channels

ServerChan supports multiple push channels. You can specify channels using either names or codes:

```go
// Using channel names
msg := serverchan.NewMessage("标题", "内容",
    serverchan.WithChannel("wecom|dingtalk|feishu"))

// Using channel codes
msg := serverchan.NewMessage("标题", "内容",
    serverchan.WithChannel("66|2|3"))

// Available channels:
// - android (98): Android 推送
// - wecom (66): 企业微信应用消息
// - wecom_bot (1): 企业微信群机器人
// - dingtalk (2): 钉钉
// - feishu (3): 飞书
// - bark (8): Bark 推送
// - test (0): 测试
// - custom (88): 自定义
// - pushdeer (18): PushDeer
// - service (9): 方糖服务号
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/serverchan"
)

// Create sender
s := gosender.NewSender(nil)

// Register ServerChan provider
serverchanProvider, err := serverchan.New(config)
if err != nil {
    log.Fatalf("Failed to create ServerChan provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeServerChan, serverchanProvider, nil)

// Send message
ctx := context.Background()
msg := serverchan.NewMessage("测试消息", "这是一条测试消息")
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of account configurations

### Account

- `Name`: Account name for identification
- `Key`: ServerChan SendKey (get from https://sct.ftqq.com/)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled

### Message Options

- `WithShort(short string)`: Set short content for message card (max 64 chars)
- `WithNoIP()`: Hide the calling IP
- `WithChannel(channel string)`: Set message channel(s)
- `WithOpenID(openid string)`: Set message copy openid

### Message Validation

- Title: Required, max 32 characters
- Content: Optional, supports Markdown, max 32KB
- Short: Optional, max 64 characters

## Notes

- **SendKey**: Get your SendKey from [ServerChan website](https://sct.ftqq.com/)
- **Enterprise Version**: If using enterprise version, SendKey format is `sctp{num}t{key}`
- **Channel Selection**: You can specify multiple channels separated by `|`
- **IP Hiding**: Use `WithNoIP()` to hide the sender's IP address
- **Message History**: View message history on ServerChan website

## API Documentation

For detailed API documentation, visit:

- [ServerChan API Documentation](https://sct.ftqq.com/)
- [ServerChan Demo](https://github.com/easychen/serverchan-demo)

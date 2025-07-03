[⬅️ Back to Main README](../../README.md)

# ServerChan Provider | 方糖/Server 酱 推送组件

This provider supports sending messages via ServerChan (方糖/Server 酱) service.

本组件支持通过 ServerChan (方糖/Server 酱) 服务推送消息。

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple accounts with different strategies (round-robin, random, weighted) | 支持多账号、负载均衡
- **Message Types 消息类型**: Support for text messages with Markdown formatting | 支持文本消息、Markdown 格式
- **Channel Support 多渠道支持**: Support for multiple push channels (WeCom, DingTalk, Lark, etc.) | 支持多渠道推送（企业微信、钉钉、飞书等）
- **Message Options 消息选项**: Short content, IP hiding, channel selection, and more | 支持卡片摘要、隐藏 IP、渠道选择等

---

## Builder API (Recommended) | 推荐用法：Builder 构造

All new usage should use the builder pattern for type safety and clarity.

推荐使用 builder 链式构造，类型安全且易于 IDE 补全：

```go
import "github.com/shellvon/go-sender/providers/serverchan"

// English: Build a message with all options
// 中文：构造带所有选项的消息
msg := serverchan.Text().
    Title("System Alert"). // 必填/required
    Content("**CPU Usage:** 95%\n**Memory:** 80%\n**Status:** Warning"). // 支持 Markdown
    Short("High CPU"). // 卡片摘要
    Channel("wecom|dingtalk"). // 多渠道
    NoIP(). // 隐藏IP
    Build()
```

---

## Minimal Constructor | 最简用法

For the simplest use case, you can use:

最简单的用法：

```go
// English: Minimal message
// 中文：最简消息
msg := serverchan.NewMessage("Title", "Content")
```

All advanced options (short, channel, noip, openid, etc.) must be set via the builder API.
所有高级选项（摘要、渠道、隐藏 IP、openid 等）请用 builder 设置。

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/serverchan"
)

// English: Create ServerChan config
// 中文：创建 ServerChan 配置
config := serverchan.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、加权等
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_SENDKEY_HERE",
            Weight:   100,
            Disabled: false,
        },
        // ... more accounts
    },
}

provider, err := serverchan.New(config)
if err != nil {
    log.Fatalf("Failed to create ServerChan provider: %v", err) // 创建失败
}
```

---

## Message Types (Builder Examples) | 消息类型（Builder 示例）

### Basic Text Message | 基本文本消息

```go
// English: Basic text message
// 中文：基础文本消息
msg := serverchan.Text().
    Title("Test Title").
    Content("This is the message content").
    Build()
```

### Message with Options | 带选项的消息

```go
// English: Message with short, channel, noip
// 中文：带摘要、渠道、隐藏IP的消息
msg := serverchan.Text().
    Title("System Notification").
    Content("## System Status\n- CPU: 45%\n- Memory: 60%\n- Disk: 75%\n**Status:** OK").
    Short("System OK").
    Channel("wecom|dingtalk").
    NoIP().
    Build()
```

### Long Text Message | 长文本消息

```go
// English: Long content message
// 中文：长内容消息
longContent := `# Detailed Report\n\n## Project Status\nThis is a detailed report sent by go-sender.\n\n### Features\n- Multi-type\n- Multi-channel\n- Markdown\n- Custom config\n\n---\n*Sent by go-sender*`

msg := serverchan.Text().
    Title("Detailed Report").
    Content(longContent).
    Short("Project Status").
    Build()
```

---

## Supported Channels | 支持的渠道

You can specify channels using either names or codes.
可用渠道支持名称或数字代码：

```go
// English: Use channel names
// 中文：用渠道名
msg := serverchan.Text().
    Title("Title").
    Content("Content").
    Channel("wecom|dingtalk|feishu").
    Build()

// English: Use channel codes
// 中文：用渠道代码
msg2 := serverchan.Text().
    Title("Title").
    Content("Content").
    Channel("66|2|3").
    Build()
```

// Available channels | 可用渠道：
// - android (98): Android
// - wecom (66): WeCom
// - wecom_bot (1): WeCom Bot
// - dingtalk (2): DingTalk
// - feishu (3): Lark/Feishu
// - bark (8): Bark
// - test (0): Test
// - custom (88): Custom
// - pushdeer (18): PushDeer
// - service (9): Service

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/providers/serverchan"
    "context"
)

// English: Register and send
// 中文：注册 provider 并发送
sender := gosender.New()
serverchanProvider, _ := serverchan.New(config)
sender.AddProvider(serverchanProvider)

msg := serverchan.Text().
    Title("Test Message").
    Content("This is a test message").
    Build()

err := sender.Send(context.Background(), msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## API Reference | API 参考

### Config | 配置

- `BaseConfig`: Common configuration fields | 通用配置字段
- `Accounts`: Array of account configurations | 账号配置数组

### Account | 账号

- `Name`: Account name for identification | 账号名称
- `Key`: ServerChan SendKey (get from https://sct.ftqq.com/) | SendKey（在官网获取）
- `Weight`: Weight for weighted strategy (default: 1) | 权重
- `Disabled`: Whether this account is disabled | 是否禁用

### Message Builder Methods | 消息构造方法

- `Text()`: Entry point for builder | builder 入口
- `.Title(string)`: Set title (required) | 设置标题（必填）
- `.Content(string)`: Set content (supports Markdown) | 设置内容（支持 Markdown）
- `.Short(string)`: Set short content (max 64 chars) | 设置摘要
- `.Channel(string)`: Set channel(s) | 设置渠道
- `.NoIP()`: Hide sender IP | 隐藏 IP
- `.OpenID(string)`: Set openid for message copy | 设置 openid
- `.Build()`: Build the message | 构建消息

### Message Validation | 消息校验

- Title: Required, max 32 characters | 标题必填，最长 32 字符
- Content: Optional, supports Markdown, max 32KB | 内容可选，最长 32KB
- Short: Optional, max 64 characters | 摘要可选，最长 64 字符

---

## Notes | 注意事项

- **SendKey**: Get your SendKey from [ServerChan website](https://sct.ftqq.com/) | SendKey 请在官网获取
- **Enterprise Version**: If using enterprise version, SendKey format is `sctp{num}t{key}` | 企业版 SendKey 格式为 `sctp{num}t{key}`
- **Channel Selection**: You can specify multiple channels separated by `|` | 多渠道用 `|` 分隔
- **IP Hiding**: Use `.NoIP()` to hide the sender's IP address | 隐藏 IP 用 `.NoIP()`
- **Message History**: View message history on ServerChan website | 消息历史可在官网查看

---

## API Documentation | 官方文档

- [ServerChan API Documentation | Server 酱官方文档](https://sct.ftqq.com/)
- [ServerChan Demo | Server 酱官方示例](https://github.com/easychen/serverchan-demo)

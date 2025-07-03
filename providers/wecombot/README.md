[⬅️ Back to Main README](../../README.md)

# WeCom Bot Provider | 企业微信机器人推送组件

This provider supports sending messages to WeCom (企业微信) group robots via webhooks.

本组件支持通过 webhook 向企业微信（WeCom）群机器人发送消息。

---

## Features | 功能特性

- **Multiple Account Support 多账号支持**: Configure multiple accounts with different strategies (round-robin, random, weighted) | 支持多账号、灵活负载均衡（轮询、随机、加权）
- **Message Types 消息类型**: Text, Markdown, Image, News, Template Card, File, Voice | 文本、Markdown、图片、新闻、模板卡片、文件、语音
- **Media Upload 媒体上传**: Support for uploading media files (images, files, voice) | 支持自动上传图片、文件、语音等媒体
- **Mention Support @成员支持**: Support for mentioning users or @all | 支持@指定成员或@所有人

---

## Configuration | 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// English: Create WeCom Bot configuration
// 中文：创建企业微信机器人配置
config := wecombot.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机、加权等
    },
    Accounts: []core.Account{
        {
            Name:     "main",
            Key:      "YOUR_WEBHOOK_KEY", // webhook URL 末尾 /key= 后的部分
            Weight:   100,
            Disabled: false,
        },
        // ... more accounts
    },
}

provider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err) // 创建失败
}
```

---

## Message Types (Builder Style) | 消息类型（链式构建）

### 1. Text Message | 文本消息

```go
// English: Simple text message
// 中文：简单文本消息
msg := wecombot.Text().
    Content("Hello from go-sender! 你好，世界！").
    Build()

// English: Text message with mentions
// 中文：带@的文本消息
msg := wecombot.Text().
    Content("Hello @all, this is a test message").
    MentionUsers([]string{"@all"}).
    MentionMobiles([]string{"***REMOVED***"}).
    Build()
```

### 2. Markdown Message | Markdown 消息

```go
markdownContent := `# 系统通知 System Notice\n\n- **CPU**: 45%\n- **Memory**: 60%\n- **Disk**: 75%\n\n> 系统运行正常 System OK\n\n[查看详情 View Details](https://example.com)`

msg := wecombot.Markdown().
    Content(markdownContent).
    Build()
```

### 3. Image Message | 图片消息

```go
// English: Upload image and send
// 中文：上传图片并发送
file, err := os.Open("image.jpg")
if err != nil {
    log.Fatalf("Failed to open image: %v", err)
}
defer file.Close()
mediaId, account, err := provider.UploadMedia(ctx, "image.jpg", file)
if err != nil {
    log.Fatalf("Failed to upload image: %v", err)
}
msg := wecombot.Image().Base64(imgBase64).MD5(imgMD5).Build()
```

### 4. News Message | 新闻消息

```go
msg := wecombot.NewsMsg().
    AddArticle("重要通知 Important", "描述 Description", "https://example.com", "https://example.com/image.jpg").
    Build()
```

### 5. Template Card Message | 模板卡片消息

```go
msg := wecombot.Card(wecombot.CardTypeTextNotice).
    MainTitle("模板卡片标题 Main Title", "描述 Description").
    SubTitle("点击查看详情 Click for details").
    JumpURL("https://example.com").
    Build()
```

### 6. File Message (Auto Upload) | 文件消息（自动上传）

```go
// English: Send a file by local path (auto upload)
// 中文：通过本地路径发送文件（自动上传）
msg := wecombot.File().
    LocalPath("/path/to/report.pdf"). // 自动上传，无需手动获取 media_id
    Build()
```

### 7. Voice Message (Auto Upload) | 语音消息（自动上传）

```go
// English: Send a voice message by local path (auto upload, AMR format, ≤2MB, ≤60s)
// 中文：通过本地路径发送语音（自动上传，仅支持AMR格式，≤2MB，≤60秒）
msg := wecombot.Voice().
    LocalPath("/path/to/voice.amr"). // 自动上传，无需手动获取 media_id
    Build()
```

---

## Usage with Sender | 与 Sender 结合使用

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

s := gosender.NewSender(nil)
wecomProvider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

ctx := context.Background()
msg := wecombot.File().LocalPath("/path/to/report.pdf").Build() // 文件自动上传
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## Media Upload | 媒体上传说明

- **Auto Upload 自动上传**: For file/voice messages, you can use `.LocalPath("/path/to/file")` in the builder. The SDK will automatically upload the file/voice and fill in the media_id. No need to call UploadMedia manually.
- **Manual Upload 手动上传**: For image messages, you can use `provider.UploadMedia(ctx, "image.jpg", file)` to get media_id, then use it in the builder.
- **Constraints 限制**:
  - File: ≤20MB
  - Voice: ≤2MB, ≤60s, AMR format only
  - All media: >5 bytes
  - 媒体文件有效期 3 天，仅上传账号可用

---

## API Reference | API 参考

### Config | 配置

- `BaseConfig`: Common configuration fields | 通用配置字段
- `Accounts`: Array of account configurations | 账号配置数组

### Message Types | 消息类型

- `TextMessage`: Text | 文本
- `MarkdownMessage`: Markdown
- `ImageMessage`: Image | 图片
- `NewsMessage`: News | 新闻
- `TemplateCardMessage`: Template Card | 模板卡片
- `FileMessage`: File | 文件
- `VoiceMessage`: Voice | 语音

### Builder Options | 构建器选项

- `.MentionUsers(users []string)`: Mention users | @成员
- `.MentionMobiles(mobiles []string)`: Mention by mobile | @手机号
- `.LocalPath(path string)`: Auto upload file/voice | 自动上传文件/语音
- `.MediaID(id string)`: Use existing media_id | 使用已上传的 media_id

---

## Notes | 说明

- **Webhook Key**: Get from WeCom group robot settings | webhook key 请在企业微信机器人设置中获取
- **Media Upload**: Media files valid for 3 days, only usable by uploading account | 媒体文件 3 天有效，仅上传账号可用
- **Mention**: Use `@all` to mention all | 使用 `@all` 可@所有人
- **Voice**: Only AMR format, ≤2MB, ≤60s | 语音仅支持 AMR 格式，≤2MB，≤60 秒
- **File**: ≤20MB | 文件最大 20MB

---

## API Documentation | 官方文档

- [WeCom Bot API Documentation | 企业微信机器人官方文档](https://developer.work.weixin.qq.com/document/path/91770)

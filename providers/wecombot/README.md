# WeCom Bot Provider (企业微信群机器人)

企业微信群机器人Provider，通过Webhook方式向企业微信群聊发送消息。

[⬅️ 返回项目README](../../README.md) | [🔗 企业微信应用Provider](../wecomapp/README.md)

---

## 📝 与企业微信应用的区别

| 特性对比 | 群机器人 (WeCom Bot) | 企业应用 (WeCom App) |
|---------|---------------------|---------------------|
| **消息目标** | 群聊成员 | 企业内指定用户/部门/标签 |
| **认证方式** | Webhook Key | 企业ID + 应用Secret + 应用ID |
| **配置复杂度** | 简单（仅需Webhook URL） | 复杂（需管理后台配置） |
| **适用场景** | 群聊通知、告警推送 | 企业内部通信、工作流通知 |
| **用户范围** | 群聊成员 | 企业所有员工 |
| **权限控制** | 群管理员控制 | 企业管理员控制 |
| **消息送达** | 群聊推送 | 个人消息推送 |

**选择建议：**
- 🤖 **群机器人**: 适合团队协作、项目群告警、简单通知场景
- 🏢 **企业应用**: 适合正式的企业通信、HR通知、审批流程等场景

---

## 功能特性

- **多账号支持**: 配置多个机器人账号，支持轮询/随机/权重负载均衡策略
- **丰富消息类型**: 支持文本、Markdown、图片、图文、模板卡片、文件、语音消息
- **媒体文件上传**: 支持自动上传图片、文件、语音等媒体文件
- **@提醒功能**: 支持@指定用户或@所有人
- **简单易用**: 仅需群机器人Webhook Key即可快速接入

---

## 配置示例

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// 创建企业微信群机器人配置
config := wecombot.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // 轮询、随机或权重策略
    },
    Items: []*wecombot.Account{
        {
            Name:     "main",
            Key:      "YOUR_WEBHOOK_KEY", // Webhook URL中/key=后的部分
            Weight:   100,
            Disabled: false,
        },
        // ... 更多账号
    },
}

provider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("创建企业微信群机器人provider失败: %v", err)
}
```

### 获取Webhook Key

1. 在企业微信群聊中，点击群设置
2. 选择"群机器人" -> "添加机器人"
3. 完成设置后，复制Webhook地址
4. 提取URL中`key=`后面的部分作为配置中的Key值

例如：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa`
其中`693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa`就是需要的Key值。

---

## 消息类型（构建器风格）

### 1. Text Message

```go
// Simple text message
msg := wecombot.Text().
    Content("Hello from go-sender!").
    Build()

// Text message with mentions
msg := wecombot.Text().
    Content("Hello @all, this is a test message").
    MentionUsers([]string{"@all"}).
    MentionMobiles([]string{"***REMOVED***"}).
    Build()
```

### 2. Markdown Message

```go
markdownContent := `# System Notice\n\n- **CPU**: 45%\n- **Memory**: 60%\n- **Disk**: 75%\n\n> System OK\n\n[View Details](https://example.com)`

msg := wecombot.Markdown().
    Content(markdownContent).
    Build()
```

### 3. Image Message

```go
// Send image with base64 and MD5
msg := wecombot.Image().
    Base64(imgBase64).
    MD5(imgMD5).
    Build()
```

### 4. News Message

```go
msg := wecombot.News().
    AddArticle("Important", "Description", "https://example.com", "https://example.com/image.jpg").
    Build()
```

### 5. Template Card Message

```go
msg := wecombot.Card(wecombot.CardTypeTextNotice).
    MainTitle("Main Title", "Description").
    SubTitle("Click for details").
    JumpURL("https://example.com").
    Build()
```

### 6. File Message

```go
// Send a file by local path
msg := wecombot.File().
    LocalPath("/path/to/report.pdf").
    Build()
```

### 7. Voice Message

```go
// Send a voice message by local path (AMR format, ≤2MB, ≤60s)
msg := wecombot.Voice().
    LocalPath("/path/to/voice.amr").
    Build()
```

---

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

s := gosender.NewSender()
wecomProvider, err := wecombot.New(config)
if err != nil {
    log.Fatalf("Failed to create WeCom Bot provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeWecombot, wecomProvider, nil)

ctx := context.Background()
msg := wecombot.File().LocalPath("/path/to/report.pdf").Build()
err = s.Send(ctx, msg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

---

## Media Upload Notes

- **File/Voice Upload**: Use `.LocalPath("/path/to/file")` in the builder. The SDK will upload the file/voice and fill in the media_id.
- **Image Upload**: Use `provider.UploadMedia(ctx, "image.jpg", file)` to get media_id, then use it in the builder.
- **Constraints**:
  - File: ≤20MB
  - Voice: ≤2MB, ≤60s, AMR format only
  - All media: >5 bytes
  - Media files are valid for 3 days and usable only by the uploading account.

---

## API Reference

### Config

- `BaseConfig`: Common configuration fields.
- `Accounts`: Array of account configurations.

### Message Types

- `TextMessage`: Text
- `MarkdownMessage`: Markdown
- `ImageMessage`: Image
- `NewsMessage`: News
- `TemplateCardMessage`: Template Card
- `FileMessage`: File
- `VoiceMessage`: Voice

### Builder Options

- `.MentionUsers(users []string)`: Mention users.
- `.MentionMobiles(mobiles []string)`: Mention by mobile.
- `.LocalPath(path string)`: Auto upload file/voice.
- `.MediaID(id string)`: Use existing media_id.

---

## Notes

- **Webhook Key**: Obtain from WeCom group robot settings.
- **Media Upload**: Media files valid for 3 days, only usable by the uploading account.
- **Mention**: Use `@all` to mention all.
- **Voice**: Only AMR format, ≤2MB, ≤60s.
- **File**: ≤20MB.

---

## API Documentation

- [WeCom Bot API Documentation](https://developer.work.weixin.qq.com/document/path/91770)

# WeCom Bot Provider (企业微信群机器人)

> 通过 Webhook 向企业微信群聊发送消息

[⬅️ 返回项目README](../../README.md) | [📖 官方文档](https://developer.work.weixin.qq.com/document/path/91770) | [🔗 企业微信应用Provider](../wecomapp/README.md)

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
| **消息类型** | 基础消息类型 | 更丰富的消息类型 |

**选择建议：**
- 🤖 **群机器人**: 适合团队协作、项目群告警、简单通知场景
- 🏢 **企业应用**: 适合正式的企业通信、HR通知、审批流程等场景

---

## 🚀 快速开始

```go
import (
    "context"
    "github.com/shellvon/go-sender/providers/wecombot"
)

// 创建账号和 Provider
account := wecombot.NewAccount("your-webhook-key")
provider, _ := wecombot.NewProvider([]*wecombot.Account{account})

// 构建并发送消息
msg := wecombot.Text().Content("Hello from go-sender!").Build()
provider.Send(context.Background(), msg, nil)
```

---

## 💬 支持的消息类型

使用 Builder 模式轻松构建各种消息：

### 文本消息 (`wecombot.Text()`)
```go
// 简单文本
msg := wecombot.Text().
    Content("系统告警：CPU 使用率超过 90%").
    Build()

// 带 @提醒
msg := wecombot.Text().
    Content("紧急通知 @all").
    MentionUsers([]string{"@all"}).
    MentionMobiles([]string{"13800138000"}).
    Build()
```

### Markdown 消息 (`wecombot.Markdown()`)
```go
msg := wecombot.Markdown().
    Content("# 监控报告\n\n- **CPU**: 45%\n- **内存**: 60%").
    Build()
```

### 图片消息 (`wecombot.Image()`)
```go
msg := wecombot.Image().
    Base64(imgBase64).
    MD5(imgMD5).
    Build()
```

### 图文消息 (`wecombot.News()`)
```go
msg := wecombot.News().
    AddArticle("重要通知", "详细描述", "https://example.com", "image.jpg").
    Build()
```

### 模板卡片 (`wecombot.Card()`)
```go
msg := wecombot.Card(wecombot.CardTypeTextNotice).
    MainTitle("系统维护通知", "预计维护 2 小时").
    SubTitle("点击查看详情").
    JumpURL("https://example.com").
    Build()
```

### 文件消息 (`wecombot.File()`)
```go
// 自动上传本地文件
msg := wecombot.File().
    LocalPath("/path/to/report.pdf").
    Build()
```

### 语音消息 (`wecombot.Voice()`)
```go
// 自动上传本地语音文件
msg := wecombot.Voice().
    LocalPath("/path/to/voice.amr").
    Build()
```

---

## ⚙️ Provider 配置

### 基础配置

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecombot"
)

config := &wecombot.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin, // 轮询、随机或权重策略
    },
    Items: []*wecombot.Account{
        wecombot.NewAccount("webhook-key-1"),
        wecombot.NewAccount("webhook-key-2"),
    },
}

provider, err := wecombot.New(config)
```

### 获取 Webhook Key

1. 在企业微信群聊中，点击群设置
2. 选择"群机器人" → "添加机器人"
3. 完成设置后，复制 Webhook 地址
4. 提取 URL 中 `key=` 后面的部分

例如：`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa`
其中 `693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa` 就是需要的 Key 值。

---

## 🔧 高级功能

### 自动上传媒体文件

对于文件和语音消息，支持自动上传功能：

```go
// 文件消息 - 自动上传
msg := wecombot.File().LocalPath("/path/to/document.pdf").Build()

// 语音消息 - 自动上传  
msg := wecombot.Voice().LocalPath("/path/to/voice.amr").Build()

// 图片消息 - 手动上传后使用
mediaID, _ := provider.UploadMedia(ctx, "image.jpg", fileBytes)
msg := wecombot.Image().MediaID(mediaID).Build()
```

### 与 Sender 集成

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
)

sender := gosender.NewSender()
provider, _ := wecombot.New(config)
sender.RegisterProvider(core.ProviderTypeWecombot, provider, nil)

// 发送消息
msg := wecombot.Text().Content("Hello").Build()
sender.Send(context.Background(), msg)
```

---

## 📋 注意事项

### 基本使用要求
- 使用 `@all` 提醒所有人
- 支持通过用户ID或手机号提醒特定用户
- 文件和语音消息不支持 @提醒

### 重要提醒
关于以下具体限制和配置，请查阅企业微信官方文档：
- **媒体文件限制**：文件大小、格式要求、有效期等详细规定
- **频率限制**：群机器人消息发送频率限制的具体数值和规则
- **错误码说明**：各种错误情况的处理方式

详细信息请参考：[企业微信群机器人官方文档](https://developer.work.weixin.qq.com/document/path/91770)

---

## 相关链接

- [企业微信群机器人官方文档](https://developer.work.weixin.qq.com/document/path/91770)
- [企业微信应用 Provider](../wecomapp/README.md) - 更强大的企业内部通信

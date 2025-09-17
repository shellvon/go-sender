# WeCom App Provider (企业微信应用)

> 通过企业微信应用发送消息到用户、部门或标签

[⬅️ 返回项目README](../../README.md) | [📖 官方文档](https://developer.work.weixin.qq.com/document/path/90236) | [🔗 企业微信群机器人Provider](../wecombot/README.md)

---

## 📝 与企业微信群机器人的区别

| 特性对比 | 企业应用 (WeCom App) | 群机器人 (WeCom Bot) |
|---------|---------------------|---------------------|
| **消息目标** | 企业内指定用户/部门/标签 | 群聊成员 |
| **认证方式** | 企业ID + 应用Secret + 应用ID | Webhook Key |
| **配置复杂度** | 复杂（需管理后台配置） | 简单（仅需Webhook URL） |
| **适用场景** | 企业内部通信、工作流通知 | 群聊通知、告警推送 |
| **用户范围** | 企业所有员工 | 群聊成员 |
| **权限控制** | 企业管理员控制 | 群管理员控制 |
| **消息送达** | 个人消息推送 | 群聊推送 |
| **消息类型** | 更丰富的消息类型 | 基础消息类型 |
| **高级功能** | 模板卡片、安全模式、自动上传 | 基础功能 |

**选择建议：**
- 🏢 **企业应用**: 适合正式的企业通信、HR通知、审批流程等场景
- 🤖 **群机器人**: 适合团队协作、项目群告警、简单通知场景

---

## 🚀 快速开始

```go
import (
    "context"
    "github.com/shellvon/go-sender/providers/wecomapp"
)

// 创建账号和 Provider
account := wecomapp.NewAccount("corp-id", "app-secret", "agent-id")
provider, _ := wecomapp.NewProvider([]*wecomapp.Account{account})

// 构建并发送消息
msg := wecomapp.Text().
    Content("系统通知：部署已完成").
    ToUser("user1|user2").  // 发送给指定用户
    Build()

provider.Send(context.Background(), msg, nil)
```

---

## 💬 支持的消息类型

使用 Builder 模式轻松构建各种消息：

### 文本消息 (`wecomapp.Text()`)
```go
// 发送给指定用户
msg := wecomapp.Text().
    Content("系统告警：CPU使用率超过90%").
    ToUser("user1|user2").
    Build()

// 发送给所有人
msg := wecomapp.Text().
    Content("重要通知").
    ToUser("@all").
    Build()
```

### Markdown 消息 (`wecomapp.Markdown()`)
```go
msg := wecomapp.Markdown().
    Content("# 监控报告\n\n- **CPU**: 45%\n- **内存**: 60%").
    ToUser("admin1|admin2").
    Build()
```

### 媒体消息 (`wecomapp.Media()`)
```go
// 图片消息 - 自动上传
msg := wecomapp.Media().
    MediaType("image").
    LocalPath("/path/to/screenshot.png").
    ToUser("user1").
    Build()

// 文件消息 - 自动上传
msg := wecomapp.Media().
    MediaType("file").
    LocalPath("/path/to/report.pdf").
    ToUser("team@department").
    Build()

// 语音消息 - 自动上传
msg := wecomapp.Media().
    MediaType("voice").
    LocalPath("/path/to/voice.amr").
    ToUser("user1").
    Build()

// 视频消息 - 自动上传
msg := wecomapp.Media().
    MediaType("video").
    LocalPath("/path/to/video.mp4").
    ToUser("user1").
    Build()
```

### 图文消息 (`wecomapp.News()`)
```go
msg := wecomapp.News().
    AddArticle("重要公告", "请注意查看最新政策", "https://example.com/news", "pic.jpg").
    AddArticle("技术分享", "Go语言最佳实践", "https://example.com/tech", "tech.jpg").
    ToUser("@all").
    Build()
```

### 文本卡片 (`wecomapp.TextCard()`)
```go
msg := wecomapp.TextCard().
    Title("部署完成").
    Description("应用版本 v2.1.0 已成功部署到生产环境").
    URL("https://example.com/deployment").
    BTNText("查看详情").
    ToUser("devops@company").
    Build()
```

### 模板卡片 (`wecomapp.TemplateCard()`)
```go
msg := wecomapp.TemplateCard().
    CardType("text_notice").
    Source(wecomapp.CardSource{
        IconURL: "https://example.com/icon.png",
        Desc:    "企业微信",
    }).
    MainTitle(wecomapp.CardMainTitle{
        Title: "欢迎使用企业微信",
        Desc:  "您的好友正在邀请您加入企业微信",
    }).
    SubTitleText("下载企业微信还能抢红包！").
    ToUser("user1").
    Build()
```

### 图文消息 (`wecomapp.MPNews()`)
```go
msg := wecomapp.MPNews().
    AddArticle("标题", "作者", "内容摘要", "图片URL", "内容链接").
    ToUser("@all").
    Build()
```

### 小程序通知 (`wecomapp.MiniprogramNotice()`)
```go
msg := wecomapp.MiniprogramNotice().
    AppID("wx123456").
    Page("pages/index").
    Title("小程序通知").
    Description("点击查看详情").
    ToUser("user1").
    Build()
```

---

## ⚙️ Provider 配置

### 基础配置

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecomapp"
)

config := &wecomapp.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin, // 负载均衡策略
    },
    Items: []*wecomapp.Account{
        wecomapp.NewAccount("corp-id", "app-secret", "agent-id"),
    },
}

provider, _ := wecomapp.New(config, nil) // nil = 使用默认内存缓存
```

### 获取配置参数

1. **企业ID (CorpID)**: 登录企业微信管理后台，在"我的企业" → "企业信息"中获取
2. **应用Secret**: 在"应用管理" → 选择应用 → "Secret"中获取  
3. **应用ID (AgentID)**: 在"应用管理" → 选择应用中获取

---

## 🔧 高级功能

### 自动上传媒体文件

对于媒体消息，支持自动上传功能：

```go
// 自动上传本地文件
msg := wecomapp.Media().
    MediaType("file").
    LocalPath("/path/to/document.pdf").
    ToUser("user1").
    Build()

// 使用已有的 media_id
msg := wecomapp.Media().
    MediaType("image").
    MediaID("MEDIA_ID_FROM_UPLOAD").
    ToUser("user1").
    Build()
```

### 发送目标配置

```go
msg := wecomapp.Text().
    Content("多目标消息").
    ToUser("user1|user2").           // 指定用户
    ToParty("dept1|dept2").          // 指定部门  
    ToTag("tag1|tag2").              // 指定标签
    Build()
```

### 安全模式和消息检查

```go
msg := wecomapp.Text().
    Content("重要消息").
    ToUser("@all").
    Safe(1).                         // 开启安全模式
    EnableDuplicateCheck(1).         // 开启重复消息检查
    DuplicateCheckInterval(3600).    // 重复检查间隔(秒)
    Build()
```

### 自定义 Token 缓存

```go
import "github.com/shellvon/go-sender/core"

// 使用自定义缓存
customCache := core.NewMemoryCache[*wecomapp.AccessToken]()
provider, _ := wecomapp.New(config, customCache)

// 使用默认缓存（传入nil）
provider, _ := wecomapp.New(config, nil)
```

### 与 Sender 集成

```go
import (
    "context"
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecomapp"
)

// 创建 Sender 实例
sender := gosender.NewSender()

// 创建企业微信应用 Provider
account := wecomapp.NewAccount("corp-id", "app-secret", "agent-id")
provider, _ := wecomapp.NewProvider([]*wecomapp.Account{account})

// 注册 Provider
sender.RegisterProvider(core.ProviderTypeWecomApp, provider, nil)

// 发送消息
msg := wecomapp.Text().Content("Hello").ToUser("user1").Build()
result, err := sender.Send(context.Background(), msg)
```

---

## 📋 注意事项

### 基本使用要求
- ToUser 中的用户ID需要是企业微信中的用户ID
- ToParty 中需要使用企业微信中的部门ID
- ToTag 中需要使用企业微信中的标签ID

### 重要提醒
关于以下具体限制和配置，请查阅企业微信官方文档：
- **媒体文件限制**：文件大小、格式要求、有效期等详细规定
- **应用权限配置**：应用可见范围、发送权限设置方法
- **频率限制**：API调用频率限制的具体数值和规则
- **错误码说明**：各种错误情况的处理方式
- **IP白名单**：微信应用支持IP白名单配置，不在名单内的会报错

详细信息请参考：[企业微信应用官方文档](https://developer.work.weixin.qq.com/document/path/90236)

---

## 相关链接

- [企业微信应用官方文档](https://developer.work.weixin.qq.com/document/path/90236)
- [发送应用消息](https://developer.work.weixin.qq.com/document/path/90236)
- [上传多媒体文件](https://developer.work.weixin.qq.com/document/path/90253)
- [企业微信群机器人 Provider](../wecombot/README.md) - 更简单的群聊通知

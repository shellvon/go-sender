# WeCom App Provider (企业微信应用)

> 通过[企业微信应用](https://developer.work.weixin.qq.com/)发送消息到用户、部门或标签。

[⬅️ 返回项目README](../../README.md) | [🔗 企业微信群机器人Provider](../wecombot/README.md)

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
| **高级功能** | 模板卡片、文件上传、安全模式 | 基础消息类型 |

**选择建议：**
- 🏢 **企业应用**: 适合正式的企业通信、HR通知、审批流程等场景
- 🤖 **群机器人**: 适合团队协作、项目群告警、简单通知场景

---

## 功能特性

- 支持多应用账号配置，支持轮询/随机/权重负载均衡策略
- 自动获取和刷新访问令牌(access_token)，自动重试机制
- 多种消息类型支持:
  - 文本消息 (Text)
  - Markdown消息
  - 图片消息 (Image)
  - 语音消息 (Voice)
  - 视频消息 (Video)  
  - 文件消息 (File)
  - 图文消息 (News)
  - 文本卡片 (TextCard)
  - 模板卡片 (TemplateCard)
  - 图文消息 (MPNews)
  - 小程序通知 (MiniprogramNotice)
- 自动文件上传功能（图片、语音、视频、文件）
- 支持发送给指定用户、部门或标签
- 安全模式、重复消息检查等高级功能

---

## 配置说明

### 基本配置

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/wecomapp"
)

cfg := wecomapp.Config{
    ProviderMeta: core.ProviderMeta{
        Strategy: core.StrategyRoundRobin, // 负载均衡策略
    },
    Items: []*wecomapp.Account{
        {
            BaseAccount: core.BaseAccount{
                AccountMeta: core.AccountMeta{
                    Name: "main",
                },
                Credentials: core.Credentials{
                    APIKey:    "YOUR_CORP_ID",     // 企业ID
                    APISecret: "YOUR_CORP_SECRET", // 应用Secret
                    AppID:     "YOUR_AGENT_ID",    // 应用ID
                },
            },
        },
    },
}
```

### 获取配置参数

1. **企业ID (CorpID)**: 登录企业微信管理后台，在"我的企业" -> "企业信息"中获取
2. **应用Secret**: 在"应用管理" -> 选择应用 -> "Secret"中获取  
3. **应用ID (AgentID)**: 在"应用管理" -> 选择应用中获取

---

## 快速构建器

### 文本消息

```go
// 简单文本消息
msg := wecomapp.Text().
    Content("系统告警: CPU使用率超过90%").
    ToUser("user1|user2").  // 指定用户
    Build()

// 发送给所有人
msg := wecomapp.Text().
    Content("重要通知").
    ToUser("@all").  // 所有用户
    Build()
```

### Markdown消息

```go
markdownContent := `# 系统监控报告

## 服务器状态
- **CPU使用率**: 45%
- **内存使用率**: 60%  
- **磁盘使用率**: 75%

> 系统运行正常

[查看详情](https://example.com/dashboard)`

msg := wecomapp.Markdown().
    Content(markdownContent).
    ToUser("admin1|admin2").
    Build()
```

### 图片消息

```go
// 使用本地文件路径 (自动上传)
msg := wecomapp.Media().
    MediaType("image").
    LocalPath("/path/to/screenshot.png").
    ToUser("user1").
    Build()

// 使用已有的media_id
msg := wecomapp.Media().
    MediaType("image").
    MediaID("MEDIA_ID_FROM_UPLOAD").
    ToUser("user1").
    Build()
```

### 文件消息

```go
// 发送文件 (自动上传)
msg := wecomapp.Media().
    MediaType("file").
    LocalPath("/path/to/report.pdf").
    ToUser("team@department").
    Build()
```

### 语音消息

```go
// 发送语音 (仅支持AMR格式，≤2MB，≤60秒)
msg := wecomapp.Media().
    MediaType("voice").
    LocalPath("/path/to/voice.amr").
    ToUser("user1").
    Build()
```

### 图文消息

```go
msg := wecomapp.News().
    AddArticle(
        "重要公告",                    // 标题
        "请注意查看最新的公司政策",      // 描述  
        "https://example.com/news",   // 链接
        "https://example.com/pic.jpg", // 图片
    ).
    AddArticle(
        "技术分享",
        "Go语言最佳实践",
        "https://example.com/tech",
        "https://example.com/tech.jpg",
    ).
    ToUser("@all").
    Build()
```

### 文本卡片

```go
msg := wecomapp.TextCard().
    Title("部署完成").
    Description("应用版本 v2.1.0 已成功部署到生产环境").
    URL("https://example.com/deployment").
    BTNText("查看详情").
    ToUser("devops@company").
    Build()
```

### 模板卡片

```go
msg := wecomapp.TemplateCard().
    CardType("text_notice").
    Source(wecomapp.CardSource{
        IconURL:   "https://example.com/icon.png",
        Desc:      "企业微信",
        DescColor: 0,
    }).
    MainTitle(wecomapp.CardMainTitle{
        Title: "欢迎使用企业微信",
        Desc:  "您的好友正在邀请您加入企业微信",
    }).
    QuoteArea(wecomapp.CardQuoteArea{
        Type:      1,
        URL:       "https://example.com",
        Title:     "引用文本标题",
        QuoteText: "引用文本内容",
    }).
    SubTitleText("下载企业微信还能抢红包！").
    HorizontalContentList([]wecomapp.CardHorizontalContent{
        {
            KeyName: "邀请人",
            Value:   "张三",
        },
        {
            KeyName: "企业名称", 
            Value:   "腾讯",
        },
    }).
    JumpList([]wecomapp.CardJump{
        {
            Type:     1,
            URL:      "https://example.com",
            Title:    "企业微信官网",
            AppID:    "APPID",
            PagePath: "pages/index",
        },
    }).
    CardAction(wecomapp.CardAction{
        Type: 1,
        URL:  "https://example.com",
    }).
    ToUser("user1").
    Build()
```

---

## 使用方法

### 1. 直接使用Provider

```go
provider, err := wecomapp.NewWithDefaults(&cfg)
if err != nil {
    log.Fatalf("创建provider失败: %v", err)
}

ctx := context.Background()
result, err := provider.Send(ctx, msg, nil)
if err != nil {
    log.Printf("发送失败: %v", err)
}
```

### 2. 结合GoSender使用

```go
import (
    gosender "github.com/shellvon/go-sender"
)

sender := gosender.NewSender()
provider, _ := wecomapp.NewWithDefaults(&cfg)
sender.RegisterProvider(core.ProviderTypeWecomApp, provider, nil)

err := sender.Send(context.Background(), msg)
if err != nil {
    log.Printf("发送失败: %v", err)
}
```

---

## 高级功能

### 文件自动上传

对于媒体消息（图片、语音、视频、文件），支持自动上传功能:

```go
// 使用本地文件路径，SDK自动上传并获取media_id
msg := wecomapp.Media().
    MediaType("file").
    LocalPath("/path/to/document.pdf").
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

### 自定义Token缓存

```go
import "github.com/shellvon/go-sender/providers/wecomapp"

// 实现自定义缓存
type CustomTokenCache struct {
    // 你的缓存实现
}

func (c *CustomTokenCache) Get(key string) (*wecomapp.AccessToken, error) {
    // 获取token实现
}

func (c *CustomTokenCache) Set(key string, token *wecomapp.AccessToken) error {
    // 设置token实现  
}

func (c *CustomTokenCache) Delete(key string) error {
    // 删除token实现
}

// 使用自定义缓存
provider, err := wecomapp.New(&cfg, &CustomTokenCache{})

// 或者使用默认内存缓存
provider, err := wecomapp.NewWithDefaults(&cfg)
// 等价于: provider, err := wecomapp.New(&cfg, nil)
```

---

## 消息类型限制

| 消息类型 | 大小限制 | 格式要求 | 特殊说明 |
|---------|---------|---------|---------|
| 文本 | 最大2048字节 | UTF-8编码 | 支持换行符 |
| Markdown | 最大4096字节 | UTF-8编码 | 支持部分Markdown语法 |
| 图片 | ≤2MB | jpg/png格式 | 自动上传获取media_id |
| 语音 | ≤2MB，≤60秒 | AMR格式 | 企业微信录音格式 |
| 视频 | ≤10MB | MP4格式 | - |
| 文件 | ≤20MB | 任意格式 | 支持常见文档格式 |

---

## 错误处理

Provider内置了以下错误处理机制:

- **Token自动刷新**: 当access_token过期时自动获取新token并重试
- **网络重试**: 支持网络请求失败时的重试机制  
- **参数验证**: 发送前验证必需参数
- **错误码映射**: 将企业微信API错误码转换为友好的错误信息

常见错误码:
- `40001`: 不合法的secret参数
- `40014`: 不合法的access_token  
- `41001`: 缺少access_token参数
- `42001`: access_token超时
- `48001`: api接口未授权

---

## 便捷方法

### 快速创建Provider

```go
// 创建单账号provider
provider, err := wecomapp.NewProvider([]*wecomapp.Account{
    wecomapp.NewAccount("corpid", "secret", "agentid", 
        wecomapp.Name("main"),
        wecomapp.Weight(100),
    ),
})

// 创建多账号provider with策略
provider, err := wecomapp.NewProvider([]*wecomapp.Account{
    wecomapp.NewAccount("corpid1", "secret1", "agentid1", wecomapp.Name("app1")),
    wecomapp.NewAccount("corpid2", "secret2", "agentid2", wecomapp.Name("app2")),
}, wecomapp.Strategy(core.StrategyWeighted))
```

---

## API参考

### Provider配置

- `Config`: Provider配置结构
- `Account`: 企业微信应用账号配置
- `ProviderOption`: Provider配置选项
- `ConfigOption`: Provider实例配置选项

### 消息类型

- `TextMessage`: 文本消息
- `MarkdownMessage`: Markdown消息  
- `MediaMessage`: 媒体消息(图片/语音/视频/文件)
- `NewsMessage`: 图文消息
- `TextCardMessage`: 文本卡片
- `TemplateCardMessage`: 模板卡片
- `MPNewsMessage`: 图文消息(mpnews)
- `MiniprogramNoticeMessage`: 小程序通知

### 构建器

- `Text()`: 文本消息构建器
- `Markdown()`: Markdown消息构建器
- `Media()`: 媒体消息构建器
- `News()`: 图文消息构建器
- `TextCard()`: 文本卡片构建器
- `TemplateCard()`: 模板卡片构建器

---

## 注意事项

- **应用权限**: 确保应用具有发送消息的权限，在企业微信管理后台配置可见范围
- **用户ID**: ToUser中的用户ID需要是企业微信中的用户ID，不是微信昵称
- **部门ID**: ToParty中需要使用企业微信中的部门ID
- **文件上传**: 上传的媒体文件有效期为3天，仅上传账号可使用
- **频率限制**: 企业微信对API调用有频率限制，建议合理控制发送频率
- **消息加密**: 开启安全模式时，消息内容会被加密传输

---

## 官方文档

- [企业微信API文档](https://developer.work.weixin.qq.com/document/path/90236)
- [发送应用消息](https://developer.work.weixin.qq.com/document/path/90236)
- [上传多媒体文件](https://developer.work.weixin.qq.com/document/path/90253)

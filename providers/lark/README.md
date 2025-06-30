[⬅️ Back to Main README](../../README.md)

# Lark/Feishu Provider

This provider supports sending messages to Lark/Feishu group robots via webhooks.

## Features

- **Multiple Account Support**: Configure multiple accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for text, post (rich text), share chat, share user, image, and interactive card messages
- **Security**: Optional webhook signature verification
- **Internationalization**: Support for Chinese and English content in post messages and interactive cards

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/lark"
)

// Create Lark configuration
config := lark.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin, // or StrategyRandom, StrategyWeighted
    },
    Accounts: []core.Account{
        {
            Name:    "lark-account-1",
            Key:     "your-webhook-url", // The webhook URL key part
            Weight:  100,
            Disabled: false,
        },
        {
            Name:    "lark-account-2",
            Key:     "your-backup-webhook-url",
            Weight:  80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := lark.New(config)
if err != nil {
    log.Fatalf("Failed to create Lark provider: %v", err)
}
```

## Message Types

### 1. Text Message

```go
textMsg := lark.NewTextMessage("Hello from go-sender!")
```

### 2. Post Message (Rich Text)

```go
postMsg := lark.NewPostMessage().
    SetChineseContent("测试标题", [][]lark.PostElement{
        {
            {Tag: "text", Text: "这是一条测试消息"},
        },
        {
            {Tag: "a", Text: "点击这里", Href: "https://www.feishu.cn"},
        },
    }).
    SetEnglishContent("Test Title", [][]lark.PostElement{
        {
            {Tag: "text", Text: "This is a test message"},
        },
        {
            {Tag: "a", Text: "Click here", Href: "https://www.larksuite.com"},
        },
    })
```

### 3. Interactive Card Message

```go
cardMsg := lark.NewInteractiveMessage().
    SetHeader(&lark.CardHeader{
        Title: &lark.CardText{
            Tag:     "plain_text",
            Content: "Interactive Card",
        },
    }).
    AddElement(lark.CardElement{
        Tag: "div",
        Text: &lark.CardText{
            Tag:     "lark_md",
            Content: "This is an interactive card message!",
        },
    }).
    AddElement(lark.CardElement{
        Tag: "action",
        Action: &lark.CardAction{
            Tag:  "button",
            Text: &lark.CardText{Tag: "plain_text", Content: "Visit Website"},
            URL:  "https://www.larksuite.com",
        },
    })
```

### 4. Share Chat Message

```go
shareChatMsg := lark.NewShareChatMessage("oc_1234567890abcdef")
```

### 5. Share User Message

```go
shareUserMsg := lark.NewShareUserMessage("ou_1234567890abcdef")
```

### 6. Image Message

```go
imageMsg := lark.NewImageMessage("img_1234567890abcdef")
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/lark"
)

// Create sender
s := gosender.NewSender(nil)

// Register Lark provider
larkProvider, err := lark.New(config)
if err != nil {
    log.Fatalf("Failed to create Lark provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeLark, larkProvider, nil)

// Send message
ctx := context.Background()
textMsg := lark.NewTextMessage("Hello from go-sender!")
err = s.Send(ctx, textMsg)
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
- `Key`: Lark webhook URL key (the part after `/hook/`)
- `Weight`: Weight for weighted strategy (default: 1)
- `Disabled`: Whether this account is disabled
- `Webhook`: Optional webhook URL (if different from standard format)

### Message Types

All message types implement the `core.Message` interface and include:

- `Validate()`: Validates message content
- `ProviderType()`: Returns `core.ProviderTypeLark`
- `GetMsgType()`: Returns the specific message type

## Notes

- **Webhook URL**: Get your webhook URL from Lark/Feishu group robot settings
- **Key Format**: The provider automatically constructs the full webhook URL using the key
- **Image Key**: For image messages, you need to upload the image to Lark first and get the image_key
- **Chat ID**: For share chat messages, use the chat ID from Lark
- **User ID**: For share user messages, use the user ID from Lark

## API Documentation

For detailed API documentation, visit:

- [Lark Bot API](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)
- [Feishu Bot API](https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN)

## Security (Signature Verification)

Lark/Feishu group bots support optional security signature verification. This provider supports both standard and signature-protected webhooks:

- If the `webhook` field is set in the account, it will be used as the full webhook URL. If the `key` field is also set, the provider will automatically append `timestamp` and `sign` parameters using the key as the secret for signature calculation.
- If only the `webhook` field is set (and `key` is empty), the provider will use the webhook as-is (you can manually append signature parameters if needed).
- If the `webhook` field is empty, the provider will construct the webhook URL using the `key` field (no signature will be used).

**Signature Calculation:**

The signature is calculated as:

```
sign = base64(HMAC-SHA256(timestamp + "\n" + secret, secret))
```

**Example configuration with signature:**

```go
config := lark.Config{
    Accounts: []core.Account{
        {
            Name:    "lark-secure-bot",
            Webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-id",
            Key:     "your-signature-secret", // Secret for signature calculation
        },
    },
}
```

For more details on Lark/Feishu bot security settings, see:

- [Lark Bot Security Settings](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)

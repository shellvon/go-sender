package serverchan

import (
	"errors"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
)

const (
	maxTitleLength        = 32
	maxShortContentLength = 64
	maxContentLength      = 32 * 1024 // 32KB
)

// Channel mapping for easy-to-remember names.
func channelMap() map[string]string {
	return map[string]string{
		// Android
		"android": "98",
		// WeCom
		"wecom":     "66", // 企业微信应用消息
		"wecom_bot": "1",  // 企业微信群机器人
		"wecom_app": "66", // 企业微信应用消息 (alias)
		// DingTalk
		"dingtalk": "2",
		"ding":     "2",
		// Feishu/Lark
		"feishu": "3",
		"lark":   "3",
		// Bark
		"bark": "8",
		// Test
		"test": "0",
		// Custom
		"custom": "88",
		// PushDeer
		"pushdeer": "18",
		// Service
		"service": "9",
		"ft":      "9", // 方糖服务号
	}
}

// Message for ServerChan
// Reference: https://sct.ftqq.com/
type Message struct {
	*core.BaseMessage

	Title   string `json:"title"`   // Message title, required, max 32 chars
	Content string `json:"desp"`    // Message content (desp), optional, supports Markdown, max 32KB
	Short   string `json:"short"`   // Message card content, optional, max 64 chars
	NoIP    bool   `json:"noip"`    // Hide calling IP, optional
	Channel string `json:"channel"` // Dynamic message channel, optional, multiple channels separated by |
	OpenID  string `json:"openid"`  // Message copy openid, optional, multiple openids separated by comma or |
}

// NewMessage creates a new ServerChan message with required title and content.
func NewMessage(title, content string) *Message {
	return Text().Title(title).Content(content).Build()
}

// Compile-time assertion: Message implements Message interface.
var (
	_ core.Message     = (*Message)(nil)
	_ core.Validatable = (*Message)(nil)
)

// Validate validates the message.
func (m *Message) Validate() error {
	if m.Title == "" {
		return errors.New("title cannot be empty")
	}
	if len(m.Title) > maxTitleLength {
		return errors.New("title cannot exceed 32 characters")
	}
	if len(m.Content) > maxContentLength {
		return errors.New("content cannot exceed 32KB")
	}
	if len(m.Short) > maxShortContentLength {
		return errors.New("short content cannot exceed 64 characters")
	}
	return nil
}

// GetSupportedChannels returns a map of supported channel names and their codes.
func GetSupportedChannels() map[string]string {
	result := make(map[string]string)
	for name, code := range channelMap() {
		result[name] = code
	}
	return result
}

// --- Builder Pattern for ServerChan ---

// TextBuilder provides a chainable builder for ServerChan messages.
type TextBuilder struct {
	title   string
	content string
	short   string
	noIP    bool
	channel string
	openid  string
}

// Text returns a new TextBuilder as the entry point for builder style.
func Text() *TextBuilder {
	return &TextBuilder{}
}

// Title sets the message title (required).
func (b *TextBuilder) Title(title string) *TextBuilder {
	b.title = title
	return b
}

// Content sets the message content (supports Markdown).
func (b *TextBuilder) Content(content string) *TextBuilder {
	b.content = content
	return b
}

// Short sets the short content for message card.
func (b *TextBuilder) Short(short string) *TextBuilder {
	b.short = short
	return b
}

// Channel sets the message channel(s).
// all channels are supported, but only 2 channels are allowed to be used at the same time
// available channels:
//   - android(98): 安卓
//   - wecom(66): 企业微信
//   - wecom_bot(1): 企业微信群机器人
//   - wecom_app(66): 企业微信应用消息
//   - dingtalk(2): 钉钉
//   - feishu(3): 飞书
//   - lark(3): 飞书
//   - bark(8): Bark
//   - test(0): 测试
//   - custom(88): 自定义
//   - pushdeer(18): 推送宝
//   - service(9): 方糖服务号
//
// If not specified, the default channel(s) configured on the website will be used.
func (b *TextBuilder) Channel(channel string) *TextBuilder {
	b.channel = channel
	return b
}

// NoIP hides the calling IP.
func (b *TextBuilder) NoIP() *TextBuilder {
	b.noIP = true
	return b
}

// OpenID sets the openid for message copy.
//
// Only the test account and WeCom App Message channels are supported.
// The openid for the test account can be obtained from the test account page, and multiple openids should be separated by commas (,).
// For the WeCom App Message channel, the openid parameter should be the recipient's UID in WeCom
// (which can be viewed via a link after configuring the channel on the message channel page).
// To send to multiple people, separate UIDs with a vertical bar (|).
// If not specified, the message will be sent to the recipients configured on the channel settings page.
//
// Reference:
//   - https://sct.ftqq.com/
func (b *TextBuilder) OpenID(openid string) *TextBuilder {
	b.openid = openid
	return b
}

// Build constructs the Message from the builder.
func (b *TextBuilder) Build() *Message {
	msg := &Message{
		BaseMessage: core.NewBaseMessage(core.ProviderTypeServerChan),
		Title:       b.title,
		Content:     b.content,
		Short:       b.short,
		NoIP:        b.noIP,
		OpenID:      b.openid,
	}
	if b.channel != "" {
		msg.Channel = msg.translateChannel(b.channel)
	}
	return msg
}

// translateChannel translates channel names to numeric codes
// If the input is already a number, it returns as is
// If the input is a mapped name, it returns the corresponding number
// If the input is not found in the map, it returns the original input.
func (m *Message) translateChannel(channel string) string {
	if channel == "" {
		return ""
	}

	// Split by | for multiple channels
	channels := strings.Split(channel, "|")
	translated := make([]string, 0, len(channels))

	for _, ch := range channels {
		ch = strings.TrimSpace(ch)
		if ch == "" {
			continue
		}

		// Check if it's already a number
		if _, err := strconv.Atoi(ch); err == nil {
			translated = append(translated, ch)
			continue
		}

		// Try to translate from map
		if translatedCode, exists := channelMap()[strings.ToLower(ch)]; exists {
			translated = append(translated, translatedCode)
		} else {
			// If not found in map, use original
			translated = append(translated, ch)
		}
	}

	return strings.Join(translated, "|")
}

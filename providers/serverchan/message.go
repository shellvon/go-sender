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
//
//nolint:gochecknoglobals // Reason: channelMap is a global mapping for serverchan channels
var channelMap = map[string]string{
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

// Message for ServerChan
// Reference: https://sct.ftqq.com/
type Message struct {
	core.DefaultMessage

	Title   string `json:"title"`   // Message title, required, max 32 chars
	Content string `json:"desp"`    // Message content (desp), optional, supports Markdown, max 32KB
	Short   string `json:"short"`   // Message card content, optional, max 64 chars
	NoIP    bool   `json:"noip"`    // Hide calling IP, optional
	Channel string `json:"channel"` // Dynamic message channel, optional, multiple channels separated by |
	OpenID  string `json:"openid"`  // Message copy openid, optional, multiple openids separated by comma or |
}

// MessageOption defines a function type for configuring Message.
type MessageOption func(*Message)

// NewMessage creates a new ServerChan message with optional configurations.
func NewMessage(title, content string, opts ...MessageOption) *Message {
	msg := &Message{
		Title:   title,
		Content: content,
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(msg)
	}

	return msg
}

// WithShort sets the short content for message card.
func WithShort(short string) MessageOption {
	return func(m *Message) {
		m.Short = short
	}
}

// WithNoIP hides the calling IP.
func WithNoIP() MessageOption {
	return func(m *Message) {
		m.NoIP = true
	}
}

// WithChannel sets the message channel(s) (replaces existing channels)
// Supports both channel names and numeric codes
// Examples: "wecom", "66", "wecom|dingtalk", "66|2".
func WithChannel(channel string) MessageOption {
	return func(m *Message) {
		m.Channel = m.translateChannel(channel)
	}
}

// AddChannel adds a channel to existing channels
// If no channels exist, it sets the channel
// Examples: AddChannel("wecom"), AddChannel("66").
func (m *Message) AddChannel(channel string) *Message {
	translated := m.translateChannel(channel)
	if m.Channel == "" {
		m.Channel = translated
	} else {
		m.Channel = m.Channel + "|" + translated
	}
	return m
}

// ClearChannels removes all channels.
func (m *Message) ClearChannels() *Message {
	m.Channel = ""
	return m
}

// WithOpenID sets the openid for message copy.
func WithOpenID(openid string) MessageOption {
	return func(m *Message) {
		m.OpenID = openid
	}
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
		if translatedCode, exists := channelMap[strings.ToLower(ch)]; exists {
			translated = append(translated, translatedCode)
		} else {
			// If not found in map, use original
			translated = append(translated, ch)
		}
	}

	return strings.Join(translated, "|")
}

func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeServerChan
}

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
	for name, code := range channelMap {
		result[name] = code
	}
	return result
}

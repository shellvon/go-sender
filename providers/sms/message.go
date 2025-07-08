package sms

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// MessageType represents the type of SMS message.
type MessageType int

const (
	SMSText MessageType = iota // 普通文本短信
	MMS                        // 多媒体短信（彩信）
	Voice                      // 语音短信
)

// MessageCategory represents the category of SMS message.
type MessageCategory int

const (
	CategoryVerification MessageCategory = iota // 验证码
	CategoryNotification                        // 通知
	CategoryPromotion                           // 营销
)

// Minimum number of digits required for a valid mobile number (international standard).
const minMobileDigits = 7

// Message represents an SMS message.
type Message struct {
	core.DefaultMessage

	// Message type and category
	Type     MessageType     `json:"type"`     // 消息类型（文本/彩信/语音）
	Category MessageCategory `json:"category"` // 消息分类（验证码/营销/通知）

	// Sub-provider specification (required for SMS due to multiple platforms)
	SubProvider string `json:"sub_provider"` // 子提供商类型（aliyun, cl253, tencent等）

	// Basic fields
	Mobiles  []string `json:"mobiles"`   // 接收号码（单个或多个）
	Content  string   `json:"content"`   // 文本内容（模板短信时可空）
	SignName string   `json:"sign_name"` // 短信签名（国内平台必备）

	// Template related fields
	TemplateID     string            `json:"template_id"`           // 平台模板ID（国际平台如Twilio需此字段）
	TemplateParams map[string]string `json:"template_params"`       // 模板参数（key-value）
	ParamsOrder    []string          `json:"template_params_array"` // （有序数组，华为等平台）

	// International SMS support
	RegionCode int `json:"region_code"` // 地区代码（regionCode，E.164国际区号），如中国大陆为86，港澳台及海外为其他

	// Common callback and scheduling fields (unified interface)
	CallbackURL string     `json:"callback_url,omitempty"` // 统一回调地址 - 各平台内部适配
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"` // 统一发送时间 - 各平台内部适配
	Extend      string     `json:"extend,omitempty"`       // 统一扩展字段 - 各平台内部适配
	UID         string     `json:"uid,omitempty"`          // 统一用户ID - 各平台内部适配

	// Extensions for platform-specific parameters
	Extras map[string]interface{} `json:"extras"` // 扩展字段（平台特定参数）
}

// ProviderType returns the provider type for this message.
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeSMS
}

// Validate validates the SMS message.
func (m *Message) Validate() error {
	if len(m.Mobiles) == 0 {
		return core.NewParamError("mobiles cannot be empty")
	}

	for _, mobile := range m.Mobiles {
		if !isValidMobileNumber(mobile) {
			return core.NewParamError(fmt.Sprintf("invalid mobile number: %s", mobile))
		}
	}

	if m.Type == SMSText && m.Content == "" && m.TemplateID == "" {
		return core.NewParamError("content or template_id is required for text messages")
	}

	if m.Type == Voice && m.Content == "" {
		return core.NewParamError("content is required for voice messages")
	}

	return nil
}

// isValidMobileNumber validates if a mobile number is valid.
func isValidMobileNumber(mobile string) bool {
	if mobile == "" {
		return false
	}

	// Remove all non-digit and non-plus characters
	clean := regexp.MustCompile(`[^\d+]`).ReplaceAllString(mobile, "")

	// Must have at least minMobileDigits digits (international minimum)
	if len(regexp.MustCompile(`\d`).FindAllString(clean, -1)) < minMobileDigits {
		return false
	}

	// Must start with + or digit
	if len(clean) > 0 && clean[0] != '+' && (clean[0] < '0' || clean[0] > '9') {
		return false
	}

	return true
}

// String returns the string representation of MessageType.
func (mt MessageType) String() string {
	switch mt {
	case SMSText:
		return "SMS Text"
	case MMS:
		return "MMS"
	case Voice:
		return "Voice"
	default:
		return "Unknown"
	}
}

// String returns the string representation of MessageCategory.
func (mc MessageCategory) String() string {
	switch mc {
	case CategoryVerification:
		return "Verification"
	case CategoryNotification:
		return "Notification"
	case CategoryPromotion:
		return "Promotion"
	default:
		return "Unknown"
	}
}

// IsIntl 判断是否为国际/港澳台短信（regionCode != 0 且 != 86）.
func (m *Message) IsIntl() bool {
	return m.RegionCode != 0 && m.RegionCode != 86
}

// IsDomestic 判断是否为中国大陆短信（regionCode == 0 或 86）.
func (m *Message) IsDomestic() bool {
	return !m.IsIntl()
}

// HasMultipleRecipients returns true if the message has multiple recipients.
func (m *Message) HasMultipleRecipients() bool {
	return len(m.Mobiles) > 1
}

// SubProviderType returns the sub-provider type for this message.
func (m *Message) SubProviderType() SubProviderType {
	return SubProviderType(m.SubProvider)
}

// GetExtraString returns a string value from extras.
func (m *Message) GetExtraString(key string) (string, bool) {
	if m.Extras == nil {
		return "", false
	}
	if value, exists := m.Extras[key]; exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetExtraStringOrDefault returns a string value from extras, or the default value if not found.
func (m *Message) GetExtraStringOrDefault(key, defaultValue string) string {
	if value, exists := m.GetExtraString(key); exists {
		return value
	}
	return defaultValue
}

// GetExtraStringOrDefaultEmpty returns a string value from extras, or empty string if not found.
func (m *Message) GetExtraStringOrDefaultEmpty(key string) string {
	return m.GetExtraStringOrDefault(key, "")
}

// GetExtraInt returns an int value from extras.
func (m *Message) GetExtraInt(key string) (int, bool) {
	if m.Extras == nil {
		return 0, false
	}
	if value, exists := m.Extras[key]; exists {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i, true
			}
		}
	}
	return 0, false
}

// GetExtraIntOrDefault returns an int value from extras, or the default value if not found.
func (m *Message) GetExtraIntOrDefault(key string, defaultValue int) int {
	if value, exists := m.GetExtraInt(key); exists {
		return value
	}
	return defaultValue
}

// GetExtraBool returns a bool value from extras.
func (m *Message) GetExtraBool(key string) (bool, bool) {
	if m.Extras == nil {
		return false, false
	}
	if value, exists := m.Extras[key]; exists {
		switch v := value.(type) {
		case bool:
			return v, true
		case string:
			return strings.ToLower(v) == "true", true
		case int:
			return v != 0, true
		}
	}
	return false, false
}

// GetExtraBoolOrDefault returns a bool value from extras, or the default value if not found.
func (m *Message) GetExtraBoolOrDefault(key string, defaultValue bool) bool {
	if value, exists := m.GetExtraBool(key); exists {
		return value
	}
	return defaultValue
}

// GetSubProvider returns the sub-provider type.
func (m *Message) GetSubProvider() string {
	return m.SubProvider
}

// GetMsgType returns the string representation of Message.Type.
func (m *Message) GetMsgType() string {
	return m.Type.String()
}

// ApplyCommonDefaults applies common default values from account to message.
// This method handles the common defaults that are shared across all SMS providers:
// - SignName: use message's SignName if present, otherwise extract from content, otherwise use account's default
// - CallbackURL: use message's callback if present, otherwise use account's default
// - RegionCode: set to 86 (China) if not set
// - Extras: initialize if nil.
func (m *Message) ApplyCommonDefaults(account *Account) {
	// Setup SignName: use message's SignName if present, otherwise extract from content, otherwise use account's default
	if m.SignName == "" {
		// Try to extract signature from content
		extractedSignName := utils.GetSignatureFromContent(m.Content)
		if extractedSignName != "" {
			// Found signature in content, set it and remove from content
			m.SignName = extractedSignName
			// Remove the signature from content (【signName】)
			m.Content = strings.TrimPrefix(m.Content, "【"+extractedSignName+"】")
			m.Content = strings.TrimSpace(m.Content)
		} else {
			// No signature in content, use account's default
			m.SignName = account.SignName
		}
	}

	// Setup CallbackURL: use message's callback if present, otherwise use account's default
	if m.CallbackURL == "" && account.Callback != "" {
		m.CallbackURL = account.Callback
	}

	// Setup Extras for platform-specific fields
	if m.Extras == nil {
		m.Extras = make(map[string]interface{})
	}

	// Setup default region code
	if m.RegionCode == 0 {
		m.RegionCode = 86
	}
}

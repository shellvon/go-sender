package sms

import (
	"strconv"

	"github.com/shellvon/go-sender/core"
)

// MessageType represents the type of SMS message
type MessageType int

const (
	SMSText MessageType = iota // 文本短信
	MMS                        // 彩信（多媒体）
	Voice                      // 语音短信
)

// MessageCategory represents the category of SMS message
type MessageCategory int

const (
	CategoryVerification MessageCategory = iota // 验证码
	CategoryNotification                        // 通知
	CategoryPromotion                           // 营销
)

// Message represents an SMS message with enhanced type and category support
type Message struct {
	core.DefaultMessage

	// Message type and category
	Type     MessageType     `json:"type"`     // 消息类型（文本/彩信/语音）
	Category MessageCategory `json:"category"` // 消息分类（验证码/营销/通知）

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

	// Extensions
	Extras map[string]interface{} `json:"extras"` // 扩展字段（平台特定参数）
}

var (
	_ core.Message = (*Message)(nil)
)

// ProviderType returns the provider type for this message
func (m *Message) ProviderType() core.ProviderType {
	return core.ProviderTypeSMS
}

// Validate checks if the Message is valid
func (m *Message) Validate() error {
	if len(m.Mobiles) == 0 {
		return core.NewParamError("mobiles cannot be empty")
	}
	return nil
}

// MessageOption defines a function type for configuring Message
type MessageOption func(*Message)

// WithType sets the message type
func WithType(msgType MessageType) MessageOption {
	return func(m *Message) {
		m.Type = msgType
	}
}

// WithCategory sets the message category
func WithCategory(category MessageCategory) MessageOption {
	return func(m *Message) {
		m.Category = category
	}
}

// WithMobiles sets the mobile phone numbers
func WithMobiles(mobiles []string) MessageOption {
	return func(m *Message) {
		m.Mobiles = mobiles
	}
}

// WithMobile sets a single mobile phone number
func WithMobile(mobile string) MessageOption {
	return func(m *Message) {
		m.Mobiles = []string{mobile}
	}
}

// WithContent sets the SMS content
func WithContent(content string) MessageOption {
	return func(m *Message) {
		m.Content = content
	}
}

// WithTemplateID sets the template ID
func WithTemplateID(templateID string) MessageOption {
	return func(m *Message) {
		m.TemplateID = templateID
	}
}

// WithTemplateParams sets the template parameters
func WithTemplateParams(params map[string]string) MessageOption {
	return func(m *Message) {
		m.TemplateParams = params
	}
}

// WithParamsOrder sets the template parameters array (ordered parameters for some providers)
func WithParamsOrder(paramsArray []string) MessageOption {
	return func(m *Message) {
		m.ParamsOrder = paramsArray
	}
}

// WithSignName sets the SMS signature name
func WithSignName(signName string) MessageOption {
	return func(m *Message) {
		m.SignName = signName
	}
}

// WithRegionCode sets the region code (E.164 国际区号)
func WithRegionCode(regionCode int) MessageOption {
	return func(m *Message) {
		m.RegionCode = regionCode
	}
}

// WithExtras sets the extra fields
func WithExtras(extras map[string]interface{}) MessageOption {
	return func(m *Message) {
		m.Extras = extras
	}
}

// NewMessage creates a new Message with required fields and optional configurations
func NewMessage(mobile string, opts ...MessageOption) *Message {
	m := &Message{
		Mobiles: []string{mobile},
		Type:    SMSText, // Default to text SMS
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(m)
	}

	return m
}

// NewTextMessage creates a new text SMS message
func NewTextMessage(mobile string, content string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{WithType(SMSText), WithContent(content)}, opts...)
	return NewMessage(mobile, opts...)
}

// NewTemplateMessage creates a new template SMS message
func NewTemplateMessage(mobile string, templateID string, params map[string]string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{
		WithType(SMSText),
		WithTemplateID(templateID),
		WithTemplateParams(params),
	}, opts...)
	return NewMessage(mobile, opts...)
}

// NewVoiceMessage creates a new voice SMS message
func NewVoiceMessage(mobile string, content string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{WithType(Voice), WithContent(content)}, opts...)
	return NewMessage(mobile, opts...)
}

// NewMMSMessage creates a new MMS message
func NewMMSMessage(mobile string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{WithType(MMS)}, opts...)
	return NewMessage(mobile, opts...)
}

// NewVerificationMessage creates a new verification SMS message
func NewVerificationMessage(mobile string, content string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{
		WithType(SMSText),
		WithCategory(CategoryVerification),
		WithContent(content),
	}, opts...)
	return NewMessage(mobile, opts...)
}

// NewNotificationMessage creates a new notification SMS message
func NewNotificationMessage(mobile string, content string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{
		WithType(SMSText),
		WithCategory(CategoryNotification),
		WithContent(content),
	}, opts...)
	return NewMessage(mobile, opts...)
}

// NewPromotionMessage creates a new promotion SMS message
func NewPromotionMessage(mobile string, content string, opts ...MessageOption) *Message {
	opts = append([]MessageOption{
		WithType(SMSText),
		WithCategory(CategoryPromotion),
		WithContent(content),
	}, opts...)
	return NewMessage(mobile, opts...)
}

// GetExtraString safely gets a string value from Extras map
func (m *Message) GetExtraString(key string) (string, bool) {
	if m.Extras == nil {
		return "", false
	}
	if value, ok := m.Extras[key]; ok {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetExtraInt safely gets an int value from Extras map
func (m *Message) GetExtraInt(key string) (int, bool) {
	if m.Extras == nil {
		return 0, false
	}
	if value, ok := m.Extras[key]; ok {
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

// GetExtraBool safely gets a bool value from Extras map
func (m *Message) GetExtraBool(key string) (bool, bool) {
	if m.Extras == nil {
		return false, false
	}
	if value, ok := m.Extras[key]; ok {
		if b, ok := value.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// GetExtraFloat safely gets a float64 value from Extras map
func (m *Message) GetExtraFloat(key string) (float64, bool) {
	if m.Extras == nil {
		return 0, false
	}
	if value, ok := m.Extras[key]; ok {
		switch v := value.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

// GetExtraStringOrDefault gets a string value from Extras map with default value
func (m *Message) GetExtraStringOrDefault(key, defaultValue string) string {
	if value, ok := m.GetExtraString(key); ok && value != "" {
		return value
	}
	return defaultValue
}

// GetExtraIntOrDefault gets an int value from Extras map with default value
func (m *Message) GetExtraIntOrDefault(key string, defaultValue int) int {
	if value, ok := m.GetExtraInt(key); ok {
		return value
	}
	return defaultValue
}

// GetExtraBoolOrDefault gets a bool value from Extras map with default value
func (m *Message) GetExtraBoolOrDefault(key string, defaultValue bool) bool {
	if value, ok := m.GetExtraBool(key); ok {
		return value
	}
	return defaultValue
}

// GetExtraFloatOrDefault gets a float64 value from Extras map with default value
func (m *Message) GetExtraFloatOrDefault(key string, defaultValue float64) float64 {
	if value, ok := m.GetExtraFloat(key); ok {
		return value
	}
	return defaultValue
}

// String returns the string representation of MessageType
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

// String returns the string representation of MessageCategory
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

// IsIntl 判断是否为国际/港澳台短信（regionCode != 0 且 != 86）
func (m *Message) IsIntl() bool {
	return m.RegionCode != 0 && m.RegionCode != 86
}

// IsDomestic 判断是否为中国大陆短信（regionCode == 0 或 86）
func (m *Message) IsDomestic() bool {
	return m.RegionCode == 0 || m.RegionCode == 86
}

// HasMultipleRecipients returns true if the message has multiple recipients.
func (m *Message) HasMultipleRecipients() bool {
	return len(m.Mobiles) > 1
}

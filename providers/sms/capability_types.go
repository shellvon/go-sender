package sms

// Limits represents message limits
type Limits struct {
	MaxBatchSize  int    `json:"max_batch_size"`  // 群发时最多多少条
	MaxContentLen int    `json:"max_content_len"` // 最大内容长度
	RateLimit     string `json:"rate_limit"`      // 频率限制
	DailyLimit    string `json:"daily_limit"`     // 日限制
}

// RegionCapability represents capability for a specific region (domestic/international)
type RegionCapability struct {
	Single     bool              `json:"single"`     // 是否支持单发
	Batch      bool              `json:"batch"`      // 是否支持群发
	Types      []MessageType     `json:"types"`      // 支持的消息类型
	Categories []MessageCategory `json:"categories"` // 支持的消息分类
	Desc       string            `json:"desc"`       // 描述说明
}

// MessageCapability represents capability for a specific message type
type MessageCapability struct {
	International RegionCapability `json:"international"` // 国际短信能力
	Domestic      RegionCapability `json:"domestic"`      // 国内短信能力
	Limits        Limits           `json:"limits"`        // 限制信息
}

// Capabilities represents provider capabilities for all message types
type Capabilities struct {
	SMS   MessageCapability `json:"sms"`   // 文本短信能力
	Voice MessageCapability `json:"voice"` // 语音短信能力
	MMS   MessageCapability `json:"mms"`   // 彩信能力
}

// NewRegionCapability creates a new RegionCapability with default values
func NewRegionCapability(single, batch bool, types []MessageType, categories []MessageCategory, desc string) RegionCapability {
	return RegionCapability{
		Single:     single,
		Batch:      batch,
		Types:      types,
		Categories: categories,
		Desc:       desc,
	}
}

// NewMessageCapability creates a new MessageCapability with default values
func NewMessageCapability() MessageCapability {
	return MessageCapability{
		International: NewRegionCapability(
			false, false,
			[]MessageType{},
			[]MessageCategory{},
			"不支持国际短信",
		),
		Domestic: NewRegionCapability(
			true, true,
			[]MessageType{SMSText},
			[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
			"支持国内短信",
		),
		Limits: Limits{
			MaxBatchSize:  100,
			MaxContentLen: 500,
			RateLimit:     "1条/秒",
			DailyLimit:    "无限制",
		},
	}
}

// NewCapabilities creates a new Capabilities with default values
func NewCapabilities() *Capabilities {
	return &Capabilities{
		SMS:   NewMessageCapability(),
		Voice: NewMessageCapability(),
		MMS:   NewMessageCapability(),
	}
}

// GetLimits returns the limits for this message capability
func (mc MessageCapability) GetLimits() Limits {
	return mc.Limits
}

// SupportsMessageType checks if this region capability supports a specific message type
func (rc RegionCapability) SupportsMessageType(msgType MessageType) bool {
	for _, t := range rc.Types {
		if t == msgType {
			return true
		}
	}
	return false
}

// SupportsCategory checks if this region capability supports a specific message category
func (rc RegionCapability) SupportsCategory(category MessageCategory) bool {
	for _, c := range rc.Categories {
		if c == category {
			return true
		}
	}
	return false
}

// Supports checks if this region capability supports a specific message type and category
func (rc RegionCapability) Supports(msgType MessageType, category MessageCategory) bool {
	return rc.SupportsMessageType(msgType) && rc.SupportsCategory(category)
}

// SupportsSMS checks if this provider supports SMS
func (c *Capabilities) SupportsSMS() bool {
	return c.SMS.Domestic.Single || c.SMS.International.Single
}

// SupportsVoice checks if this provider supports voice messages
func (c *Capabilities) SupportsVoice() bool {
	return c.Voice.Domestic.Single || c.Voice.International.Single
}

// SupportsMMS checks if this provider supports MMS
func (c *Capabilities) SupportsMMS() bool {
	return c.MMS.Domestic.Single || c.MMS.International.Single
}

// SupportsMessageType checks if this provider supports a specific message type
func (c *Capabilities) SupportsMessageType(msgType MessageType) bool {
	switch msgType {
	case SMSText:
		return c.SupportsSMS()
	case Voice:
		return c.SupportsVoice()
	case MMS:
		return c.SupportsMMS()
	default:
		return false
	}
}

// SupportsCategory checks if this provider supports a specific message category
func (c *Capabilities) SupportsCategory(category MessageCategory) bool {
	return c.SMS.Domestic.SupportsCategory(category) ||
		c.SMS.International.SupportsCategory(category) ||
		c.Voice.Domestic.SupportsCategory(category) ||
		c.Voice.International.SupportsCategory(category) ||
		c.MMS.Domestic.SupportsCategory(category) ||
		c.MMS.International.SupportsCategory(category)
}

// Supports checks if this provider supports a specific message type and category
func (c *Capabilities) Supports(msgType MessageType, category MessageCategory) bool {
	switch msgType {
	case SMSText:
		return c.SMS.Domestic.Supports(msgType, category) ||
			c.SMS.International.Supports(msgType, category)
	case Voice:
		return c.Voice.Domestic.Supports(msgType, category) ||
			c.Voice.International.Supports(msgType, category)
	case MMS:
		return c.MMS.Domestic.Supports(msgType, category) ||
			c.MMS.International.Supports(msgType, category)
	default:
		return false
	}
}

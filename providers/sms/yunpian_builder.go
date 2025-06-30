package sms

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/official/document/sms/zh_CN/api_reference.html
//
// 官方文档:
//   - 短信API: https://www.yunpian.com/official/document/sms/zh_CN/api_reference.html
//
// builder 仅支持 text（普通短信）类型。

// YunpianSMSBuilder provides Yunpian-specific SMS message creation
type YunpianSMSBuilder struct {
	BaseSMSBuilder
}

// NewYunpianSMSBuilder creates a new Yunpian SMS builder
func NewYunpianSMSBuilder() *YunpianSMSBuilder {
	return &YunpianSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderYunpian},
	}
}

// NewTextMessage creates a new Yunpian text SMS message
func (b *YunpianSMSBuilder) NewTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(SMSText, mobiles, content, opts...)
}

// NewVoiceMessage creates a new Yunpian voice message
func (b *YunpianSMSBuilder) NewVoiceMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(Voice, mobiles, content, opts...)
}

// NewMMSMessage creates a new Yunpian MMS message
func (b *YunpianSMSBuilder) NewMMSMessage(mobiles []string, opts ...MessageOption) *Message {
	return b.newMessage(MMS, mobiles, "", opts...)
}

// newMessage is a generic method to create messages of any type
func (b *YunpianSMSBuilder) newMessage(msgType MessageType, mobiles []string, content string, opts ...MessageOption) *Message {
	// Build base options
	baseOpts := []MessageOption{
		WithSubProvider(string(b.subProvider)),
		WithType(msgType),
		WithMobiles(mobiles),
	}

	// Add content for text and voice messages
	if content != "" {
		baseOpts = append(baseOpts, WithContent(content))
	}

	// Add user options
	allOpts := append(baseOpts, opts...)

	// Create message with all options at once
	return NewMessageWithOptions(allOpts...)
}

// WithYunpianExtend 设置扩展字段
// 云片短信服务中，extend 字段用于扩展信息
// 文档地址: https://www.yunpian.com/official/document/sms/zh_cn/api-single_send
func WithYunpianExtend(extend string) MessageOption {
	return WithExtra(yunpianExtend, extend)
}

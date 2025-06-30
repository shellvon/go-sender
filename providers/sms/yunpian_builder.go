package sms

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/dev-doc
//
// 官方文档:
//   - 短信API: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
//   - 语音API: https://www.yunpian.com/official/document/sms/zh_CN/voice_send
//   - 超级短信API: https://www.yunpian.com/official/document/sms/zh_CN/super_sms_send
//
//
// builder 仅支持 text（普通短信）类型。

// YunpianSMSBuilder provides Yunpian-specific SMS message creation.
type YunpianSMSBuilder struct {
	BaseSMSBuilder
}

// NewYunpianSMSBuilder creates a new YunpianSMSBuilder instance.
func NewYunpianSMSBuilder() *YunpianSMSBuilder {
	return &YunpianSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderYunpian},
	}
}

// NewTextMessage creates a new Yunpian text SMS message.
func (b *YunpianSMSBuilder) NewTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(SMSText, mobiles, content, opts...)
}

// NewVoiceMessage creates a new Yunpian voice message.
func (b *YunpianSMSBuilder) NewVoiceMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(Voice, mobiles, content, opts...)
}

// NewMMSMessage creates a new Yunpian MMS message.
func (b *YunpianSMSBuilder) NewMMSMessage(mobiles []string, opts ...MessageOption) *Message {
	return b.newMessage(MMS, mobiles, "", opts...)
}

// newMessage is a generic method to create messages of any type.
func (b *YunpianSMSBuilder) newMessage(
	msgType MessageType,
	mobiles []string,
	content string,
	opts ...MessageOption,
) *Message {
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
	baseOpts = append(baseOpts, opts...)

	// Create message with all options at once
	return NewMessageWithOptions(baseOpts...)
}

// WithYunpianRegister sets the register field for Yunpian SMS.
// 是否为注册验证码短信，如果传入 true，则该条短信作为注册验证码短信统计注册成功率，需联系客服开通。
//   - https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
func WithYunpianRegister(register bool) MessageOption {
	return WithExtra(yunpianRegisterKey, register)
}

// WithYunpianMobileStat sets the mobile_stat field for Yunpian SMS.
// 若短信中包含云片短链接，此参数传入 true 将会把短链接替换为目标手机号的专属链接，用于统计哪些号码的机主点击了短信中的链接，可在云片后台查看。详情参考短信点击统计；
// 传false时，短信中包含的云片短链接将原样发送给终端手机号，不会替换专属链接。
// 该字段默认值为false。
//   - https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
func WithYunpianMobileStat(stat bool) MessageOption {
	return WithExtra(yunpianMobileStatKey, stat)
}

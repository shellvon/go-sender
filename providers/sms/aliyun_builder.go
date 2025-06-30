package sms

import (
	"strconv"
)

// AliyunSMSBuilder provides Aliyun-specific SMS message creation.
type AliyunSMSBuilder struct {
	BaseSMSBuilder
}

// NewAliyunSMSBuilder creates a new Aliyun SMS builder.
//   - SMS(普通文本): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//   - MMS(卡片消息): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Vocie(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
func NewAliyunSMSBuilder() *AliyunSMSBuilder {
	return &AliyunSMSBuilder{
		BaseSMSBuilder: BaseSMSBuilder{subProvider: SubProviderAliyun},
	}
}

// NewTextMessage creates a new Aliyun text SMS message.
func (b *AliyunSMSBuilder) NewTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(SMSText, mobiles, content, opts...)
}

// NewVoiceMessage creates a new Aliyun voice message.
func (b *AliyunSMSBuilder) NewVoiceMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	return b.newMessage(Voice, mobiles, content, opts...)
}

// NewMMSMessage creates a new Aliyun MMS (CardSMS) message.
func (b *AliyunSMSBuilder) NewMMSMessage(mobiles []string, opts ...MessageOption) *Message {
	return b.newMessage(MMS, mobiles, "", opts...)
}

// newMessage is a generic method to create messages of any type.
func (b *AliyunSMSBuilder) newMessage(
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

// WithAliyunCalledShowNumber 语音短信可选参数
// 发送语音通知的通话号码（被叫显号）。若此参数不填，则为公共模式通话；若传入真实号或服务实例 ID，则为专属模式通话。
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Vocie(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//
// 当发送语音短信时，优先使用消息本身指定的显示号码，若不存在，则会使用账号信息上配置的 From 字段。
func WithAliyunCalledShowNumber(calledShowNumber string) MessageOption {
	return WithExtra(aliyunCalledShowNumberKey, calledShowNumber)
}

// WithAliyunPlayTimes 语音短信可选参数
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Vocie(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//
// 语音通知文件的播放次数。取值范围：1~3。
func WithAliyunPlayTimes(playTimes int) MessageOption {
	return WithExtra(aliyunPlayTimesKey, strconv.Itoa(playTimes))
}

// WithAliyunVolume 语音短信可选参数
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Vocie(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//
// 语音通知文件播放的音量。取值范围：0~100，默认取值 100。
func WithAliyunVolume(volume int) MessageOption {
	return WithExtra(aliyunVolumeKey, strconv.Itoa(volume))
}

// WithAliyunSpeed 语音短信可选参数
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Vocie(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//
// 语音文件播放的语速。取值范围：-500~500。
func WithAliyunSpeed(speed int) MessageOption {
	return WithExtra(aliyunSpeedKey, strconv.Itoa(speed))
}

// WithAliyunOutID 可选参数:外部流水号
// 外部流水扩展字段，用于标识业务流水号，在状态报告中会原样返回。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func WithAliyunOutID(outID string) MessageOption {
	return WithExtra(aliyunOutIDKey, outID)
}

// WithAliyunFallbackType 卡片短信可选参数 回落类型。取值：
//   - SMS：不支持卡片短信的号码，回落文本短信。
//   - DIGITALSMS：不支持卡片短信的号码，回落数字短信。
//   - NONE：不需要回落。
//
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyunFallbackType(fallbackType string) MessageOption {
	return WithExtra(aliyunFallbackTypeKey, fallbackType)
}

// WithAliyunSmsTemplateCode 卡片短信可选参数，回落文本短信的模板 Code。FallbackType 选择 SMS 回落文本短信时，此参数必填
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyunSmsTemplateCode(smsTemplateCode string) MessageOption {
	return WithExtra(aliyunSmsTemplateCodeKey, smsTemplateCode)
}

// WithAliyunDigitalTemplateCode 卡片短信可选参数，回落数字短信的模板 Code。FallbackType 选择 DIGITALSMS 回落数字短信时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyunDigitalTemplateCode(digitalTemplateCode string) MessageOption {
	return WithExtra(aliyunDigitalTemplateCodeKey, digitalTemplateCode)
}

// WithAliyunSmsTemplateParam 卡片短信可选参数: 回落文本短信的模板变量对应的实际值。SmsTemplateCode 回落的文本短信模板内含有变量时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyunSmsTemplateParam(smsTemplateParam string) MessageOption {
	return WithExtra(aliyunSmsTemplateParamKey, smsTemplateParam)
}

// WithAliyuDigitalTemplateParam 卡片短信可选参数: 回落数字短信的模板变量对应的实际值。DigitalTemplateCode 回落的数字短信模板内含有变量时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyuDigitalTemplateParam(digitalTemplateParam string) MessageOption {
	return WithExtra(aliyunDigitalTemplateParamKey, digitalTemplateParam)
}

// WithAliyunSmsUpExtendCode 可选参数:上行短信扩展码。上行短信指发送给通信服务提供商的短信，用于定制某种服务、完成查询，或是办理某种业务等，需要收费，按运营商普通短信资费进行扣费。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func WithAliyunSmsUpExtendCode(smsUpExtendCode string) MessageOption {
	return WithExtra(aliyunSmsUpExtendCodeKey, smsUpExtendCode)
}

// WithAliyunCardObjects 卡片短信可选参数: 卡片消息的卡片对象。
// 用于定义卡片短信的具体内容和样式，包括标题、描述、图片、按钮等元素。
// 每个卡片对象包含卡片的配置信息，如卡片类型、内容、样式等。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func WithAliyunCardObjects(cardObjects string) MessageOption {
	return WithExtra(aliyunCardObjectsKey, cardObjects)
}

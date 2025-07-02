package sms

// AliyunSMSBuilder provides Aliyun-specific SMS message creation.
type AliyunSMSBuilder struct {
	*BaseBuilder

	region               string
	calledShowNumber     string
	playTimes            int
	volume               int
	speed                int
	outID                string
	fallbackType         string
	smsTemplateCode      string
	digitalTemplateCode  string
	smsTemplateParam     string
	digitalTemplateParam string
	smsUpExtendCode      string
	cardObjects          string
}

func newAliyunSMSBuilder() *AliyunSMSBuilder {
	return &AliyunSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderAliyun},
		volume:      aliyunDefaultVolume, // 默认音量
	}
}

// AliyunSMSBuilder inherits all methods from BaseSMSBuilder

// CalledShowNumber 设置阿里云语音短信的被叫显号。
// 发送语音通知的通话号码（被叫显号）。若此参数不填，则为公共模式通话；若传入真实号或服务实例 ID，则为专属模式通话。
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Voice(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
//
// 当发送语音短信时，优先使用消息本身指定的显示号码，若不存在，则会使用账号信息上配置的 From 字段。
func (b *AliyunSMSBuilder) CalledShowNumber(num string) *AliyunSMSBuilder {
	b.calledShowNumber = num
	return b
}

// PlayTimes 设置语音短信的播放次数。
// 语音通知文件的播放次数。取值范围：1~3。
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Voice(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
func (b *AliyunSMSBuilder) PlayTimes(times int) *AliyunSMSBuilder {
	b.playTimes = times
	return b
}

// Volume 设置语音短信的音量。
// 语音通知文件播放的音量。取值范围：0~100，默认取值 100。
func (b *AliyunSMSBuilder) Volume(volume int) *AliyunSMSBuilder {
	b.volume = volume
	return b
}

// Speed 设置语音短信的语速。
// 语音文件播放的语速。取值范围：-500~500。
func (b *AliyunSMSBuilder) Speed(speed int) *AliyunSMSBuilder {
	b.speed = speed
	return b
}

// OutID 设置外部流水号。
// 外部流水扩展字段，用于标识业务流水号，在状态报告中会原样返回。字符串类型，长度限制为 1~15 个字符。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func (b *AliyunSMSBuilder) OutID(outID string) *AliyunSMSBuilder {
	b.outID = outID
	return b
}

// FallbackType 设置卡片短信回落类型。
// 回落类型。取值：SMS、DIGITALSMS、NONE。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) FallbackType(t string) *AliyunSMSBuilder {
	b.fallbackType = t
	return b
}

// SmsTemplateCode 设置卡片短信回落文本短信的模板 Code。
// FallbackType 选择 SMS 回落文本短信时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) SmsTemplateCode(code string) *AliyunSMSBuilder {
	b.smsTemplateCode = code
	return b
}

// DigitalTemplateCode 设置卡片短信回落数字短信的模板 Code。
// FallbackType 选择 DIGITALSMS 回落数字短信时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) DigitalTemplateCode(code string) *AliyunSMSBuilder {
	b.digitalTemplateCode = code
	return b
}

// SmsTemplateParam 设置卡片短信回落文本短信的模板变量值。
// SmsTemplateCode 回落的文本短信模板内含有变量时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) SmsTemplateParam(param string) *AliyunSMSBuilder {
	b.smsTemplateParam = param
	return b
}

// DigitalTemplateParam 设置卡片短信回落数字短信的模板变量值。
// DigitalTemplateCode 回落的数字短信模板内含有变量时，此参数必填。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) DigitalTemplateParam(param string) *AliyunSMSBuilder {
	b.digitalTemplateParam = param
	return b
}

// SmsUpExtendCode 设置上行短信扩展码。
// 上行短信指发送给通信服务提供商的短信，用于定制某种服务、完成查询，或是办理某种业务等，需要收费，按运营商普通短信资费进行扣费。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
func (b *AliyunSMSBuilder) SmsUpExtendCode(code string) *AliyunSMSBuilder {
	b.smsUpExtendCode = code
	return b
}

// CardObjects 设置卡片短信的卡片对象。
// 用于定义卡片短信的具体内容和样式，包括标题、描述、图片、按钮等元素。
// 每个卡片对象包含卡片的配置信息，如卡片类型、内容、样式等。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
func (b *AliyunSMSBuilder) CardObjects(objs string) *AliyunSMSBuilder {
	b.cardObjects = objs
	return b
}

// Region 设置阿里云短信的区域, 默认值为 cn-hangzhou。
// https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-endpoint
// 该值影响当次发放的服务接入点
// 阿里云短信的区域。取值：
//
// 亚太:
//   - cn-hangzhou: 中国大陆(默认)
//   - cn-shanghai: 中国上海
//   - cn-beijing: 中国北京
//   - cn-guangzhou: 中国广州
//   - cn-chengdu: 中国成都
//   - ap-southeast-1 新加坡
//   - ap-southeast-5 印度尼西亚（雅加达）
//   - cn-hongkong 中国香港
//
// 欧洲:
//   - eu-central-1 德国
//
// 行业云:
//   - cn-beijing-finance-1	 华北2 金融云（邀测）
//   - cn-hangzhou-finance	华东1 金融云
//   - cn-shanghai-finance-1 华东2 金融云
//   - cn-shenzhen-finance-1 华南1 金融云
//
// 对于语音短信，无论设置哪一个区域，目前都是固定值。
func (b *AliyunSMSBuilder) Region(region string) *AliyunSMSBuilder {
	b.region = region
	return b
}

func (b *AliyunSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	// 阿里云专属参数写入Extras
	extra := map[string]interface{}{}
	if b.calledShowNumber != "" {
		extra[aliyunCalledShowNumberKey] = b.calledShowNumber
	}
	if b.playTimes != 0 {
		extra[aliyunPlayTimesKey] = b.playTimes
	}
	if b.volume != 0 {
		extra[aliyunVolumeKey] = b.volume
	}
	if b.speed != 0 {
		extra[aliyunSpeedKey] = b.speed
	}
	if b.outID != "" {
		extra[aliyunOutIDKey] = b.outID
	}
	if b.fallbackType != "" {
		extra[aliyunFallbackTypeKey] = b.fallbackType
	}
	if b.smsTemplateCode != "" {
		extra[aliyunSmsTemplateCodeKey] = b.smsTemplateCode
	}
	if b.digitalTemplateCode != "" {
		extra[aliyunDigitalTemplateCodeKey] = b.digitalTemplateCode
	}
	if b.smsTemplateParam != "" {
		extra[aliyunSmsTemplateParamKey] = b.smsTemplateParam
	}
	if b.digitalTemplateParam != "" {
		extra[aliyunDigitalTemplateParamKey] = b.digitalTemplateParam
	}
	if b.smsUpExtendCode != "" {
		extra[aliyunSmsUpExtendCodeKey] = b.smsUpExtendCode
	}
	if b.cardObjects != "" {
		extra[aliyunCardObjectsKey] = b.cardObjects
	}
	if b.region != "" {
		extra[aliyunRegionKey] = b.region
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}

package sms

// @ProviderName: Huawei Cloud / 华为云
// @Website: https://www.huaweicloud.com
// @APIDoc: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 官方文档:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 仅支持国内/国际短信，不支持彩信和语音。
//
// builder 仅支持 text（模板短信）类型。

// HuaweiSMSBuilder provides Huawei-specific SMS message creation.
type HuaweiSMSBuilder struct {
	*BaseBuilder

	from string
}

// newHuaweiSMSBuilder creates a new Huawei SMS builder.
func newHuaweiSMSBuilder() *HuaweiSMSBuilder {
	return &HuaweiSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderHuawei},
	}
}

// From sets the sender number for Huawei SMS.
// 华为云短信服务中，from 字段用于指定发送方号码。
//   - 文档地址: https://support.huaweicloud.com/api-msgsms/sms_05_0002.html
func (b *HuaweiSMSBuilder) From(from string) *HuaweiSMSBuilder {
	b.from = from
	return b
}

// Build constructs the Huawei SMS message with all the configured options.
func (b *HuaweiSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	if b.from != "" {
		extra[huaweiFromKey] = b.from
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}

package sms

import "fmt"

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6361/67380
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6361/67380
//
// builder 仅支持 text（普通短信）类型。

type VolcSMSBuilder struct {
	*BaseBuilder

	smsAccount string
	tag        string
}

func newVolcSMSBuilder() *VolcSMSBuilder {
	return &VolcSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderVolc},
	}
}

// Tag 设置火山引擎短信标签。
// 用于标识短信业务类型，便于后续统计分析和管理。
//   - 短信API: https://www.volcengine.com/docs/6361/67380
func (b *VolcSMSBuilder) Tag(tag string) *VolcSMSBuilder {
	b.tag = tag
	return b
}

// SmsAccount 消息组 ID, 对于火山云API必须设置。
//   - 短信API: https://www.volcengine.com/docs/6361/67380
func (b *VolcSMSBuilder) SmsAccount(smsAccount string) *VolcSMSBuilder {
	b.smsAccount = smsAccount

	fmt.Println("b.smsAccount", b.smsAccount)
	return b
}

func (b *VolcSMSBuilder) Build() *Message {

	fmt.Println("b.smsAccount--->", b.smsAccount)
	msg := b.BaseBuilder.Build()

	fmt.Println("b.smsAccount", b.smsAccount)
	fmt.Println("b.tag", b.tag)
	extra := map[string]interface{}{}
	if b.smsAccount != "" {
		extra[volcSmsAccountKey] = b.smsAccount
	}
	if b.tag != "" {
		extra[volcTagKey] = b.tag
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	fmt.Println("msg.Extras", msg.Extras)
	return msg
}

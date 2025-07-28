package sms

// @ProviderName: Submail / 赛邮
// @Website: https://www.mysubmail.com
// @APIDoc: https://www.mysubmail.com/documents
//
// 官方文档:
//   - 国内短信: https://www.mysubmail.com/documents/FppOR3
//   - 国际短信: https://www.mysubmail.com/documents/3UQA3
//   - 模板短信: https://www.mysubmail.com/documents/OOVyh
//   - 群发: https://www.mysubmail.com/documents/AzD4Z4
//   - 语音: https://www.mysubmail.com/documents/meE3C1
//   - 彩信: https://www.mysubmail.com/documents/N6ktR
//
// builder 支持 text（国内/国际，模板/非模板，单发/群发）、voice（语音）、mms（彩信）类型。

type SubmailSMSBuilder struct {
	*BaseBuilder[*SubmailSMSBuilder]
}

func newSubmailSMSBuilder() *SubmailSMSBuilder {
	b := &SubmailSMSBuilder{}
	b.BaseBuilder = &BaseBuilder[*SubmailSMSBuilder]{subProvider: SubProviderSubmail, self: b}
	return b
}

// Tag sets the tag for Submail SMS.
// 消息标签（用于消息追踪，最大32字符）。
//   - https://www.mysubmail.com/documents/FppOR3
func (b *SubmailSMSBuilder) Tag(tag string) *SubmailSMSBuilder {
	return b.meta(submailTagKey, tag)
}

// Sender sets the sender identifier for Submail SMS.
// 主要用于国际短信，可选字段。
//   - https://www.mysubmail.com/documents/3UQA3
func (b *SubmailSMSBuilder) Sender(sender string) *SubmailSMSBuilder {
	return b.meta(submailSenderKey, sender)
}

// SignType sets the signature type for Submail SMS.
//   - md5 (default)
//   - sha1
//   - normal
//
// Docs:
//   - https://www.mysubmail.com/documents/FppOR3
func (b *SubmailSMSBuilder) SignType(signType string) *SubmailSMSBuilder {
	return b.meta(submailSignTypeKey, signType)
}

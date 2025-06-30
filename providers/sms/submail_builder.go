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

// NewSubmailTextMessage 创建赛邮短信消息.
func NewSubmailTextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderSubmail)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// NewSubmailVoiceMessage 创建赛邮语音消息.
func NewSubmailVoiceMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderSubmail)),
		WithType(Voice),
		WithMobiles(mobiles),
		WithContent(content),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// NewSubmailMMSMessage 创建赛邮彩信消息.
func NewSubmailMMSMessage(mobiles []string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderSubmail)),
		WithType(MMS),
		WithMobiles(mobiles),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithSubmailTag 设置标签
// 用于消息追踪，最大32个字符.
func WithSubmailTag(tag string) MessageOption {
	return WithExtra("tag", tag)
}

// WithSubmailSender 设置发送方标识
// 主要用于国际短信，可选字段.
func WithSubmailSender(sender string) MessageOption {
	return WithExtra("sender", sender)
}

// WithSubmailSignType 设置签名类型
// 可选值：md5（默认）、sha1、normal.
func WithSubmailSignType(signType string) MessageOption {
	return WithExtra("signType", signType)
}

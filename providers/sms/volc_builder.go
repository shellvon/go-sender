package sms

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6348/70146
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6348/70146
//
// builder 仅支持 text（普通短信）类型。

// NewVolcTextMessage 创建火山引擎短信消息
//
// 示例：
//
//	msg := NewVolcTextMessage(
//	         []string{"13800138000"},
//	         "您的验证码是1234",
//	         "签名",
//	         WithExtend("12345"),
//	         WithCallbackURL("https://callback.example.com"),
//	      )
func NewVolcTextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderVolc)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithVolcTag 设置火山引擎短信标签.
//   - 短信API: https://www.volcengine.com/docs/6348/70146
func WithVolcTag(tag string) MessageOption {
	return WithExtra(volcTagKey, tag)
}

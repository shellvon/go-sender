package sms

// @ProviderName: Yuntongxun / 容联云
// @Website: https://www.yuntongxun.com
// @APIDoc: https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
//
// 官方文档:
//   - 短信API: https://www.yuntongxun.com/doc/rest/sms/3_2_2_1.html
//
// builder 仅支持 text（普通短信）类型。

// NewYuntongxunTextMessage creates a new Yuntongxun text SMS message.
func NewYuntongxunTextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderYuntongxun)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

package sms

// @ProviderName: Smsbao / 短信宝
// @Website: https://www.smsbao.com
// @APIDoc: https://www.smsbao.com/openapi
//
// 官方文档:
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//   - 国际短信: https://www.smsbao.com/openapi/299.html
//   - 语音验证码: https://www.smsbao.com/openapi/214.html
//
// 能力说明:
//   - 国内短信：支持单发和群发，最多99个号码/次。
//   - 国际短信：支持单发和群发，最多99个号码/次。
//   - 语音验证码：仅支持国内、仅验证码类型、仅单号码。
//
// builder 支持文本短信（国内/国际）和语音验证码。

// NewSMSBaoTextMessage 创建短信宝文本短信消息（国内/国际均可）
//
// 参数：
//   - mobiles: 接收手机号列表，支持单发/群发，最多99个号码
//   - content: 短信内容
//   - opts:    其他可选参数（如签名、msg_id、国际短信等）
//
// 示例：
//
//	msg := NewSMSBaoTextMessage(
//	         []string{"***REMOVED***"},
//	         "您的验证码是1234",
//	         WithSMSBaoSign("签名"),
//	         WithSMSBaoMsgID("custom_id"),
//	         WithSMSBaoIntl(),
//	      )
func NewSMSBaoTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderSmsbao)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
	}
	allOpts := append(baseOpts, opts...)
	return NewMessageWithOptions(allOpts...)
}

// NewSMSBaoVoiceMessage 创建短信宝语音验证码消息（仅支持国内单号码）
//
// 参数：
//   - mobile:  接收手机号（仅支持单个国内号码）
//   - code:    验证码内容
//   - opts:    其他可选参数（如msg_id）
//
// 示例：
//
//	msg := NewSMSBaoVoiceMessage(
//	         "***REMOVED***",
//	         "123456",
//	         WithSMSBaoMsgID("custom_id"),
//	      )
func NewSMSBaoVoiceMessage(mobile string, code string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderSmsbao)),
		WithType(Voice),
		WithMobiles([]string{mobile}),
		WithContent(code),
	}
	allOpts := append(baseOpts, opts...)
	return NewMessageWithOptions(allOpts...)
}

// 当客户使用专用通道产品时，需要指定产品ID
// 产品ID可在短信宝后台或联系客服获得,不填则默认使用通用短信产品
// 文档
//   - 国内短信: https://www.smsbao.com/openapi/213.html
//
// Alias: WithSMSBaoProductID
// Equivalent to WithTemplateID
func WithSMSBaoProductID(productID string) MessageOption {
	return WithTemplateID(productID)
}

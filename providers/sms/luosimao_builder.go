package sms

// @ProviderName: Luosimao / 螺丝帽
// @Website: https://luosimao.com
// @APIDoc: https://luosimao.com/docs/api
//
// 官方文档:
//   - 短信API: https://luosimao.com/docs/api
//
// builder 仅支持 text（普通短信）类型。

// NewLuosimaoTextMessage 创建螺丝帽短信消息
//
// 示例：
//
//	msg := NewLuosimaoTextMessage(
//	         []string{"13800138000"},
//	         "您的验证码是1234",
//	         WithLuosimaoCallbackUrl("https://callback.example.com"),
//	      )
func NewLuosimaoTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderLuosimao)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
	}
	allOpts := append(baseOpts, opts...)
	return NewMessageWithOptions(allOpts...)
}

// NewLuosimaoVoiceMessage 创建螺丝帽语音验证码消息
//
// 仅支持单发验证码，type=Voice，category=CategoryVerification。
// content 字段为验证码内容，签名无效。
//
// 语音API文档: https://luosimao.com/docs/api/51
//
// 示例：
//
//	msg := NewLuosimaoVoiceMessage("13800138000", "123456")
func NewLuosimaoVoiceMessage(mobile string, code string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderLuosimao)),
		WithType(Voice),
		WithCategory(CategoryVerification),
		WithMobiles([]string{mobile}),
		WithContent(code),
	}
	allOpts := append(baseOpts, opts...)
	return NewMessageWithOptions(allOpts...)
}

// WithLuosimaoSignName 设置短信签名（仅文本短信有效，语音短信不支持签名）
//
// 螺丝帽短信签名需放在内容前部，如：【签名】内容，亦可以直接通过指定Content字段，如：【签名】内容，然后signName字段为空
// 示例：WithLuosimaoSignName("您的品牌")
func WithLuosimaoSignName(sign string) MessageOption {
	return func(m *Message) {
		m.SignName = sign
	}
}

// WithLuosimaoSendTime 设置批量短信定时发送时间（仅收件人>1时有效）
//
// 格式：2025-06-27 10:00:00
// 批量短信API文档: https://luosimao.com/docs/api#send_batch
//
// 示例：WithLuosimaoSendTime("2025-06-27 10:00:00")
func WithLuosimaoSendTime(sendTime string) MessageOption {
	return WithExtra("time", sendTime)
}

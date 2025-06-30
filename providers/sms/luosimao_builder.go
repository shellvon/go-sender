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
//	         WithSignName("您的品牌"),
//	         WithCallbackURL("https://callback.example.com"),
//	      )
func NewLuosimaoTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderLuosimao)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
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
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

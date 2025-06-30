package sms

// @ProviderName: Juhe / 聚合数据
// @Website: https://www.juhe.cn
// @APIDoc: https://www.juhe.cn/docs/api/id/54
//
// 官方文档:
//   - 短信API文档: https://www.juhe.cn/docs/api/id/54
//   - 国内短信API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信API: https://www.juhe.cn/docs/api/id/357
//   - 视频短信API: https://www.juhe.cn/docs/api/id/363
//
// builder 仅支持 text（普通短信）类型。

// NewJuheTextMessage 创建聚合短信消息
//
// 参数：
//
//	mobiles: 接收短信的手机号列表，支持一个或多个手机号
//	content: 短信内容（如验证码等）
//	sign:    短信签名
//	opts:    其他可选参数（如模板ID、回调地址、扩展码等）
func NewJuheTextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderJuhe)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

package sms

// @ProviderName: UCP / 云之讯
// @Website: https://www.ucpaas.com
// @APIDoc: http://docs.ucpaas.com
//
// 官方文档:
//   - 短信API: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
//
// builder 仅支持 text（普通短信）类型。

// NewUcpTextMessage 创建UCP短信消息
//
// 参数：
//   - mobiles: 手机号列表
//   - templateID: 模板ID
//   - paramsOrder: 模板参数顺序
//   - opts: 其他可选参数
func NewUcpTextMessage(mobiles []string, templateID string, paramsOrder []string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderUcp)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithTemplateID(templateID),
		WithParamsOrder(paramsOrder),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithUcpUID 设置自定义UID（回调用）.
func WithUcpUID(uid string) MessageOption {
	return WithUID(uid)
}

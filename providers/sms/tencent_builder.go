package sms

// @ProviderName: Tencent / 腾讯云
// @Website: https://cloud.tencent.com
// @APIDoc: https://cloud.tencent.com/document/product/382/55981
//
// 官方文档:
//   - 短信API: https://cloud.tencent.com/document/product/382/55981
//   - 语音API: https://cloud.tencent.com/document/product/1128/51559
//
// builder 支持 text（普通短信）和 voice（语音短信）类型。

// NewTencentTextMessage 创建腾讯云短信消息
//
// 参数：
//   - mobiles: 手机号列表
//   - templateID: 模板ID
//   - sign: 签名
//   - paramsOrder: 模板参数顺序
//   - opts: 其他可选参数
func NewTencentTextMessage(
	mobiles []string,
	templateID, sign string,
	paramsOrder []string,
	opts ...MessageOption,
) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderTencent)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithTemplateID(templateID),
		WithSignName(sign),
		WithParamsOrder(paramsOrder),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// NewTencentVoiceMessage 创建腾讯云语音短信消息
//
// 参数：
//   - mobiles: 手机号列表（仅支持单个）
//   - templateID: 模板ID
//   - paramsOrder: 模板参数顺序
//   - opts: 其他可选参数
func NewTencentVoiceMessage(mobile string, templateID string, paramsOrder []string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderTencent)),
		WithType(Voice),
		WithMobiles([]string{mobile}),
		WithTemplateID(templateID),
		WithParamsOrder(paramsOrder),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithTencentSenderID 设置 SenderId（国际/港澳台短信专用）
//   - https://cloud.tencent.com/document/product/382/55981
//
// 国内短信无需填写该项；国际/港澳台短信已申请独立 SenderId 需要填写该字段，默认使用公共 SenderId，无需填写该字段。
//
//   - 注：月度使用量达到指定量级可申请独立 SenderId 使用，详情请联系 腾讯云短信小助手。示例值：Qsms
func WithTencentSenderID(senderID string) MessageOption {
	return WithExtra(tencentSenderID, senderID)
}

// WithTencentRegion 设置 Region
//   - https://cloud.tencent.com/document/product/382/55981
//
// 目前支持的区域可查看: https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.9F.E5.88.97.E8.A1.A8
//   - 华北地区（北京）	ap-beijing
//   - 华南地区（广州）	ap-guangzhou
//   - 华东地区（南京）	ap-nanjing
func WithTencentRegion(region string) MessageOption {
	return WithExtra("Region", region)
}

// WithTencentPlayTimes 设置语音播放次数
//
// 语音验证码/通知API参数：PlayTimes
// - 取值范围：1~3，默认值为2。
// - 仅语音短信有效，文本短信无效。
// - 超出范围将被自动修正为默认值。
//
// 官方文档：
//   - https://cloud.tencent.com/document/product/382/38778
//   - https://cloud.tencent.com/document/product/1128/51559
//
// 示例：
//
//	msg := NewTencentVoiceMessage("***REMOVED***", "1234", nil, WithTencentPlayTimes(3))
func WithTencentPlayTimes(times int) MessageOption {
	return WithExtra("PlayTimes", times)
}

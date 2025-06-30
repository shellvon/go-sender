package sms

// @ProviderName: Huawei Cloud / 华为云
// @Website: https://www.huaweicloud.com
// @APIDoc: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 官方文档:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 仅支持国内/国际短信，不支持彩信和语音。
//
// builder 仅支持 text（模板短信）类型。

// NewHuaweiTextMessage 创建华为云模板短信消息（国内/国际均可）
//
// 示例：
//
//	msg := NewHuaweiTextMessage(
//	         []string{"13800138000"},
//	         "模板ID",
//	         []string{"param1", "param2"},
//	         "签名",
//	         WithCallbackURL("https://callback.example.com"),
//	         WithExtend("12345"),
//	      )
func NewHuaweiTextMessage(
	mobiles []string,
	templateID string,
	paramsOrder []string,
	sign string,
	opts ...MessageOption,
) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderHuawei)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithTemplateID(templateID),
		WithParamsOrder(paramsOrder),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithHuaweiFrom 设置发送方号码
// 华为云短信服务中，from 字段用于指定发送方号码
// 文档地址: https://support.huaweicloud.com/api-msgsms/sms_05_0002.html
func WithHuaweiFrom(from string) MessageOption {
	return WithExtra(huaweiFromKey, from)
}

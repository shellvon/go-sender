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
//	         WithHuaweiCallbackUrl("https://callback.example.com"),
//	         WithHuaweiExtendCode("12345"),
//	      )
func NewHuaweiTextMessage(mobiles []string, templateId string, paramsOrder []string, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderHuawei)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithTemplateID(templateId),
		WithParamsOrder(paramsOrder),
		WithSignName(sign),
	}
	allOpts := append(baseOpts, opts...)
	return NewMessageWithOptions(allOpts...)
}

// WithHuaweiCallbackUrl 设置状态回执的回调地址
// 客户的回调地址，用于接收短信状态报告，如：http://my.com/receiveSMSReport。
//
//   - 如果设置了该字段，则该消息的状态报告将通过"接收状态报告"接口直接通知客户。
//   - 如果未设置该字段，则短信平台收到运营商短信中心返回的状态报告不会推送给客户，该状态报告将在短信平台中保存1个小时，超时后系统会自动删除。
//   - 回调地址推荐使用域名。
//
// 文档地址:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 当发送消息之中不存在此字段时，则使用账号信息上配置的 Webhook 字段。
func WithHuaweiStatusCallback(statusCallback string) MessageOption {
	return WithExtra("statusCallback", statusCallback)
}

// WithHuaweiExtendCode 设置扩展码
// 扩展参数，在状态报告中会原样返回。不允许赋空值，不允许携带以下字符："{"，"}"（即大括号）。
//
// 文档地址:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
func WithHuaweiExtendCode(extend string) MessageOption {
	return WithExtra("extend", extend)
}

// WithHuaweiFrom 设置发送方号码
// 华为云短信服务中，from 字段用于指定发送方号码
// 文档地址: https://support.huaweicloud.com/api-msgsms/sms_05_0002.html
func WithHuaweiFrom(from string) MessageOption {
	return WithExtra(huaweiFrom, from)
}

package sms

import (
	"strconv"
)

// @ProviderName: Cl253 (Chuanglan) / 创蓝253
// @Website: https://www.253.com
// @APIDoc: https://www.253.com/api
//
// 官方文档:
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// CL253 支持能力:
//   - 国内短信：支持验证码、通知、营销，单发/群发，签名自动拼接，需遵守工信部规范。
//   - 国际短信：支持验证码、通知、营销，仅单发，需带国际区号，内容需以签名开头。
//   - 彩信/语音短信：暂不支持。
//
// 本 builder 仅支持 text（普通短信）类型。

// NewCL253TextMessage 创建 CL253 普通短信消息
//
// 示例：
//
//	msg := NewCL253TextMessage(
//	         []string{"13800138000"},
//	         "您的验证码是1234",
//	         "签名",
//	         WithCL253Report(true),
//	         WithCL253CallbackUrl("https://callback.example.com"),
//	      )
func NewCL253TextMessage(mobiles []string, content, sign string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(SubProviderCl253)),
		WithType(SMSText),
		WithMobiles(mobiles),
		WithContent(content),
		WithSignName(sign),
	}
	baseOpts = append(baseOpts, opts...)
	return NewMessageWithOptions(baseOpts...)
}

// WithCL253Report 设置是否需要状态报告
// CL253短信服务中，report 字段用于指定是否需要状态报告
// 如您需要状态回执，则需要传"true",不传默认为"false"，则无法获取状态回执。
//
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// 文档地址: https://www.cl253.com/api/send
func WithCL253Report(report bool) MessageOption {
	return WithExtra(cl253Report, strconv.FormatBool(report))
}

// WithCL253SenderID 设置发件人ID
// 用户收到短信之后显示的发件人，国内不支持自定义，国外支持，但是需要提前和运营商沟通注册，具体请与 TIG 对接人员确定
//
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// 此参数仅国际短信支持.
func WithCL253SenderID(senderID string) MessageOption {
	return WithExtra(cl253SenderID, senderID)
}

// WithCL253TdFlag 设置国际短信的 TD 标志
// 退订开启标识，1 开启；0 或 null 关闭
//
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// 此参数仅国际短信支持.
func WithCL253TdFlag(tdFlag string) MessageOption {
	return WithExtra(cl253TDFlag, tdFlag)
}

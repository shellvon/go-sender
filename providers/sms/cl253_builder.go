package sms

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

// Cl253SMSBuilder provides CL253-specific SMS message creation.
type Cl253SMSBuilder struct {
	*BaseBuilder

	senderID string
	tdFlag   int    // 退订开启标识，1 开启；0 或 null 关闭
	report   string // 状态报告参数
}

// newCl253SMSBuilder creates a new CL253 SMS builder.
func newCl253SMSBuilder() *Cl253SMSBuilder {
	return &Cl253SMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderCl253},
	}
}

// SenderID sets the sender ID for CL253 international SMS.
// 用户收到短信之后显示的发件人，国内不支持自定义，国外支持，但是需要提前和运营商沟通注册，具体请与 TIG 对接人员确定。
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// 此参数仅国际短信支持.
func (b *Cl253SMSBuilder) SenderID(id string) *Cl253SMSBuilder {
	b.senderID = id
	return b
}

// TDFlag sets the international SMS TD flag for CL253.
// 退订开启标识，1 开启；0 或 null 关闭。
//   - 国际短信 API: https://doc.chuanglan.com/document/O58743GF76M7754H
//
// 此参数仅国际短信支持.
func (b *Cl253SMSBuilder) TDFlag(flag int) *Cl253SMSBuilder {
	b.tdFlag = flag
	return b
}

// Report sets the report parameter for CL253 SMS.
// 状态报告参数，用于接收短信发送状态回调。
//   - 国内短信 API: https://doc.chuanglan.com/document/HAQYSZKH9HT5Z50L
func (b *Cl253SMSBuilder) Report(report string) *Cl253SMSBuilder {
	b.report = report
	return b
}

func (b *Cl253SMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	if b.senderID != "" {
		extra[cl253SenderIDKey] = b.senderID
	}
	extra[cl253TDFlagKey] = b.tdFlag
	if b.report != "" {
		extra[cl253ReportKey] = b.report
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}

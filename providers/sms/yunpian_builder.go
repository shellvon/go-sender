package sms

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/dev-doc
//
// 官方文档:
//   - 短信API: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
//   - 语音API: https://www.yunpian.com/official/document/sms/zh_CN/voice_send
//   - 超级短信API: https://www.yunpian.com/official/document/sms/zh_CN/super_sms_send
//
// builder 仅支持 text（普通短信）类型。

// YunpianSMSBuilder provides Yunpian-specific SMS message creation.
type YunpianSMSBuilder struct {
	*BaseBuilder

	register   bool
	mobileStat bool
}

func newYunpianSMSBuilder() *YunpianSMSBuilder {
	return &YunpianSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderYunpian},
	}
}

// Register sets the register field for Yunpian SMS.
// 是否为注册验证码短信，如果传入 true，则该条短信作为注册验证码短信统计注册成功率，需联系客服开通。
//   - https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
func (b *YunpianSMSBuilder) Register(register bool) *YunpianSMSBuilder {
	b.register = register
	return b
}

// MobileStat sets the mobile_stat field for Yunpian SMS.
// 若短信中包含云片短链接，此参数传入 true 将会把短链接替换为目标手机号的专属链接，用于统计哪些号码的机主点击了短信中的链接，可在云片后台查看。详情参考短信点击统计；
// 传false时，短信中包含的云片短链接将原样发送给终端手机号，不会替换专属链接。
// 该字段默认值为false。
//   - https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
func (b *YunpianSMSBuilder) MobileStat(stat bool) *YunpianSMSBuilder {
	b.mobileStat = stat
	return b
}

func (b *YunpianSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	extra[yunpianRegisterKey] = b.register
	extra[yunpianMobileStatKey] = b.mobileStat
	msg.Extras = extra
	return msg
}

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

// TencentSMSBuilder provides Tencent-specific SMS message creation.
type TencentSMSBuilder struct {
	*BaseBuilder

	senderID      string
	region        string
	playTimes     int
	smsSdkAppID   string
	voiceSdkAppID string
}

// newTencentSMSBuilder creates a new Tencent SMS builder.
func newTencentSMSBuilder() *TencentSMSBuilder {
	return &TencentSMSBuilder{
		BaseBuilder: &BaseBuilder{subProvider: SubProviderTencent},
		region:      tencentDefaultRegion,
	}
}

// SenderID sets the SenderId for Tencent international/HK/Macau/Taiwan SMS.
// 国内短信无需填写该项；国际/港澳台短信已申请独立 SenderId 需要填写该字段，默认使用公共 SenderId，无需填写该字段。
//   - https://cloud.tencent.com/document/product/382/55981
//
// 示例值：Qsms.
func (b *TencentSMSBuilder) SenderID(id string) *TencentSMSBuilder {
	b.senderID = id
	return b
}

// Region sets the Region for Tencent SMS.
// 目前支持的区域可查看: https://cloud.tencent.com/document/api/382/52071#.E5.9C.B0.E5.9F.E5.88.E6.A1.A3
//   - 华北地区（北京） ap-beijing
//   - 华南地区（广州） ap-guangzhou (default)
//   - 华东地区（南京） ap-nanjing
func (b *TencentSMSBuilder) Region(region string) *TencentSMSBuilder {
	b.region = region
	return b
}

// PlayTimes sets the play times for Tencent voice SMS.
// 语音验证码/通知API参数：PlayTimes
// - 取值范围：1~3，默认值为2。
// - 仅语音短信有效，文本短信无效。
func (b *TencentSMSBuilder) PlayTimes(times int) *TencentSMSBuilder {
	b.playTimes = times
	return b
}

// SmsSdkAppID sets the SmsSdkAppId for Tencent SMS.
//   - This is required for most Tencent SMS API calls.
//   - See: https://cloud.tencent.com/document/product/382/55981.
func (b *TencentSMSBuilder) SmsSdkAppID(appID string) *TencentSMSBuilder {
	b.smsSdkAppID = appID
	return b
}

// VoiceSdkAppID sets the Voice SdkAppId for Tencent voice SMS.
// This should be used for voice SMS scenarios only.
//   - https://cloud.tencent.com/document/product/1128/51559.
func (b *TencentSMSBuilder) VoiceSdkAppID(appID string) *TencentSMSBuilder {
	b.voiceSdkAppID = appID
	return b
}

func (b *TencentSMSBuilder) Build() *Message {
	msg := b.BaseBuilder.Build()
	extra := map[string]interface{}{}
	if b.senderID != "" {
		extra[tencentSenderIDKey] = b.senderID
	}
	if b.region != "" {
		extra[tencentRegionKey] = b.region
	}
	if b.playTimes != 0 {
		extra[tencentPlayTimesKey] = b.playTimes
	}
	if b.smsSdkAppID != "" {
		extra[tencentSmsSdkAppIDKey] = b.smsSdkAppID
	}
	if b.voiceSdkAppID != "" {
		extra[tencentVoiceSdkAppIDKey] = b.voiceSdkAppID
	}
	if len(extra) > 0 {
		msg.Extras = extra
	}
	return msg
}

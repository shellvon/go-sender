package sms

// SubProviderType represents the specific SMS service provider.
type SubProviderType string

const (
	// SubProviderAliyun represents Aliyun SMS service.
	SubProviderAliyun SubProviderType = "aliyun"
	// SubProviderTencent represents Tencent SMS service.
	SubProviderTencent SubProviderType = "tencent"
	// SubProviderCl253 represents CL253 SMS service.
	SubProviderCl253 SubProviderType = "cl253"
	// SubProviderHuawei represents Huawei SMS service.
	SubProviderHuawei SubProviderType = "huawei"
	// SubProviderJuhe represents Juhe SMS service.
	SubProviderJuhe SubProviderType = "juhe"
	// SubProviderLuosimao represents Luosimao SMS service.
	SubProviderLuosimao SubProviderType = "luosimao"
	// SubProviderSmsbao represents Smsbao SMS service.
	SubProviderSmsbao SubProviderType = "smsbao"
	// SubProviderSubmail represents Submail SMS service.
	SubProviderSubmail SubProviderType = "submail"
	// SubProviderUcp represents UCP SMS service.
	SubProviderUcp SubProviderType = "ucp"
	// SubProviderVolc represents Volc SMS service.
	SubProviderVolc SubProviderType = "volc"
	// SubProviderYuntongxun represents Yuntongxun SMS service.
	SubProviderYuntongxun SubProviderType = "yuntongxun"
	// SubProviderYunpian represents Yunpian SMS service.
	SubProviderYunpian SubProviderType = "yunpian"
)

// MessageBuilder defines the interface for building SMS messages with various options.
type MessageBuilder interface {
	Build() *Message
}

// BaseSMSBuilder provides common functionality for all SMS builders.
type BaseSMSBuilder struct {
	subProvider SubProviderType
}

// NewTextMessage creates a new text SMS message.
func (b *BaseSMSBuilder) NewTextMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(b.subProvider)),
		WithType(SMSText),
		WithContent(content),
		WithMobiles(mobiles),
	}

	return NewMessage("", append(baseOpts, opts...)...)
}

// NewMMSMessage creates a new MMS (multimedia) message.
func (b *BaseSMSBuilder) NewMMSMessage(mobiles []string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(b.subProvider)),
		WithType(MMS),
		WithMobiles(mobiles),
	}

	return NewMessage("", append(baseOpts, opts...)...)
}

// NewVoiceMessage creates a new voice message.
func (b *BaseSMSBuilder) NewVoiceMessage(mobiles []string, content string, opts ...MessageOption) *Message {
	baseOpts := []MessageOption{
		WithSubProvider(string(b.subProvider)),
		WithType(Voice),
		WithContent(content),
		WithMobiles(mobiles),
	}

	return NewMessage("", append(baseOpts, opts...)...)
}

// Factory functions for each provider
// These functions will be implemented in their respective builder files

// Aliyun returns an Aliyun SMS builder.
func Aliyun() *AliyunSMSBuilder {
	return NewAliyunSMSBuilder()
}

// Tencent returns a Tencent SMS builder (to be implemented)
// func Tencent() *TencentSMSBuilder {
//     return NewTencentSMSBuilder()
// }

// Cl253 returns a CL253 SMS builder (to be implemented)
// func Cl253() *Cl253SMSBuilder {
//     return NewCl253SMSBuilder()
// }

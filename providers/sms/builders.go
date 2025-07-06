/*
Package sms provides a unified API for sending SMS messages across multiple providers.

Usage (链式Builder风格)：

	// 文本短信
	msg := sms.Aliyun().To("13800138000").Content("您的验证码是1234").SignName("签名").Build()
	msg := sms.Tencent().To("13800138000").Content("您的验证码是1234").SignName("签名").Build()
	msg := sms.Huawei().To("13800138000").Content("您的验证码是1234").SignName("签名").Build()

	// 模板短信（无序参数）


	msg := sms.Aliyun().To("13800138000").TemplateID("SMS_123456").Params(map[string]string{"code": "1234"}).SignName("签名").Build()
	msg := sms.Yunpian().To("13800138000").TemplateID("模板ID").Params(map[string]string{"code": "1234"}).Build()

	// 模板短信（有序参数）
	msg := sms.Tencent().To("13800138000").TemplateID("模板ID").ParamsOrder([]string{"1234"}).SignName("签名").Build()
	msg := sms.Huawei().To("13800138000").TemplateID("模板ID").ParamsOrder([]string{"1234"}).SignName("签名").Build()

	// 语音短信


	msg := sms.Aliyun().To("13800138000").Type(sms.Voice).TemplateID("TTS_123456").Params(map[string]string{"code": "1234"}).CalledShowNumber("400xxxxxxx").Build()
	msg := sms.Tencent().To("13800138000").Type(sms.Voice).TemplateID("模板ID").ParamsOrder([]string{"1234"}).Build()

	// 彩信
	msg := sms.Aliyun().To("13800138000").Type(sms.MMS).Build()

Provider Documentation:
  - Aliyun: https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
  - Tencent: https://cloud.tencent.com/document/product/382/55981
  - Huawei: https://support.huaweicloud.com/api-msgsms/sms_05_0001.html
  - CL253: https://doc.chuanglan.com/
  - Juhe: https://www.juhe.cn/docs
  - Luosimao: https://luosimao.com/docs/api
  - Smsbao: https://www.smsbao.com/openapi
  - Submail: https://www.mysubmail.com/documents
  - UCP: http://docs.ucpaas.com/doku.php
  - Volc: https://www.volcengine.com/product/cloud-sms
  - Yuntongxun: https://www.yuntongxun.com/developer-center
  - Yunpian: https://www.yunpian.com/dev-doc
*/
package sms

import (
	"time"
)

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

// BaseBuilder provides chainable methods for common SMS parameters.
// It is generic so that each concrete builder can embed a `*BaseBuilder[*ConcreteBuilder]`
// and all链式 methods将返回具体的 builder 类型，避免 "跳回" 基类的问题。
type BaseBuilder[T any] struct {
	subProvider    SubProviderType
	mobiles        []string
	content        string
	signName       string
	templateID     string
	templateParams map[string]string
	paramsOrder    []string
	category       MessageCategory
	msgType        MessageType
	regionCode     int
	callbackURL    string
	scheduledAt    *time.Time
	extend         string
	uid            string

	// self 持有实际的 builder 指针，用于链式调用时返回具体类型。
	self T
}

// To sets the mobile numbers to send the message to.
func (b *BaseBuilder[T]) To(mobiles ...string) T {
	b.mobiles = mobiles
	return b.self
}

// Content sets the content of the message.
func (b *BaseBuilder[T]) Content(content string) T {
	b.content = content
	return b.self
}

// SignName sets the SMS signature (sign name).
//   - For most domestic (China mainland) SMS providers this field is mandatory and the value must be a previously approved signature.
//   - A few international providers allow the signature to be embedded directly in the message body, making this field optional.
//
// When in doubt, always set a sign name to avoid provider-side validation errors.
func (b *BaseBuilder[T]) SignName(sign string) T {
	b.signName = sign
	return b.self
}

// TemplateID sets the template ID of the message.
func (b *BaseBuilder[T]) TemplateID(id string) T {
	b.templateID = id
	return b.self
}

// Params sets the template parameters of the message.
//   - For most SMS providers, this is a map of parameter names to their values.
//   - For some providers, the order of parameters is important and should be specified using ParamsOrder.
//   - For providers that support both ordered and unordered parameters, this method can be used for both.
//
// When in doubt, always set the parameters to avoid provider-side validation errors.
func (b *BaseBuilder[T]) Params(params map[string]string) T {
	b.templateParams = params
	return b.self
}

// ParamsOrder sets the order of template parameters.
//   - For some SMS providers, the order of parameters is important and should be specified using this method.
//   - For providers that support both ordered and unordered parameters, this method can be used for both.
//
// When in doubt, always set the parameters to avoid provider-side validation errors.
func (b *BaseBuilder[T]) ParamsOrder(order []string) T {
	b.paramsOrder = order
	return b.self
}

// Category sets the category of the message.
//   - For some SMS providers, the category of the message is important and should be specified using this method.
//   - For providers that support both ordered and unordered parameters, this method can be used for both.
//
// When in doubt, always set the parameters to avoid provider-side validation errors.
func (b *BaseBuilder[T]) Category(category MessageCategory) T {
	b.category = category
	return b.self
}

// Type sets the type of the message.
//   - For some SMS providers, the type of the message is important and should be specified using this method.
//   - For providers that support both ordered and unordered parameters, this method can be used for both.
//
// When in doubt, always set the parameters to avoid provider-side validation errors.
func (b *BaseBuilder[T]) Type(t MessageType) T {
	b.msgType = t
	return b.self
}

// RegionCode sets the destination country/region code (E.164).
//   - This only affects international SMS; for domestic (China mainland) SMS the value is ignored.
//   - A value of 0 or 86 is treated as mainland China.
//   - Any other value tells the provider to use its international-SMS API, if supported.
func (b *BaseBuilder[T]) RegionCode(code int) T {
	b.regionCode = code
	return b.self
}

// CallbackURL sets the callback URL of the message.
func (b *BaseBuilder[T]) CallbackURL(url string) T {
	b.callbackURL = url
	return b.self
}

// ScheduledAt sets the scheduled (delayed-send) time of the message.
// The feature is only effective if the underlying SMS provider supports scheduled delivery; the
// transformer will format the timestamp as required by that provider. Examples of supported
// providers:
//   - Luosimao: https://luosimao.com/docs/api#send_batch
//   - CL253   : https://doc.chuanglan.com/
//
// If the provider has no concept of scheduled SMS, this parameter is silently ignored.
func (b *BaseBuilder[T]) ScheduledAt(at time.Time) T {
	b.scheduledAt = &at
	return b.self
}

// Extend sets the vendor-specific extension field.
// Most SMS providers expose an optional "extend/extra" parameter that is echoed back in delivery
// receipts or callbacks. Validation rules vary by provider—refer to their respective docs for
// the exact format.
func (b *BaseBuilder[T]) Extend(ext string) T {
	b.extend = ext
	return b.self
}

// UID sets a user-defined unique identifier for the message.
// Providers known to honor this field include:
//   - Yunpian: https://www.yunpian.com/dev-doc
//   - UCP    : http://docs.ucpaas.com/doku.php
func (b *BaseBuilder[T]) UID(uid string) T {
	b.uid = uid
	return b.self
}

// Build assembles the final *Message with all已设置字段.
func (b *BaseBuilder[T]) Build() *Message {
	return &Message{
		Type:           b.msgType,
		Category:       b.category,
		SubProvider:    string(b.subProvider),
		Mobiles:        b.mobiles,
		Content:        b.content,
		SignName:       b.signName,
		TemplateID:     b.templateID,
		TemplateParams: b.templateParams,
		ParamsOrder:    b.paramsOrder,
		RegionCode:     b.regionCode,
		CallbackURL:    b.callbackURL,
		ScheduledAt:    b.scheduledAt,
		Extend:         b.extend,
		UID:            b.uid,
	}
}

// Aliyun creates a new Aliyun SMS message builder.
func Aliyun() *AliyunSMSBuilder {
	return newAliyunSMSBuilder()
}

// Tencent creates a new Tencent SMS message builder.
func Tencent() *TencentSMSBuilder {
	return newTencentSMSBuilder()
}

// Cl253 creates a new CL253 SMS message builder.
func Cl253() *Cl253SMSBuilder {
	return newCl253SMSBuilder()
}

// Huawei creates a new Huawei SMS message builder.
func Huawei() *HuaweiSMSBuilder {
	return newHuaweiSMSBuilder()
}

// Juhe creates a new Juhe SMS message builder.
func Juhe() *JuheSMSBuilder {
	return newJuheSMSBuilder()
}

// Luosimao creates a new Luosimao SMS message builder.
func Luosimao() *LuosimaoSMSBuilder {
	return newLuosimaoSMSBuilder()
}

// Smsbao creates a new Smsbao SMS message builder.
func Smsbao() *SmsbaoSMSBuilder {
	return newSmsbaoSMSBuilder()
}

// Submail creates a new Submail SMS message builder.
func Submail() *SubmailSMSBuilder {
	return newSubmailSMSBuilder()
}

// Ucp creates a new UCP SMS message builder.
func Ucp() *UcpSMSBuilder {
	return newUcpSMSBuilder()
}

// Volc creates a new Volc SMS message builder.
func Volc() *VolcSMSBuilder {
	return newVolcSMSBuilder()
}

// Yuntongxun creates a new Yuntongxun SMS message builder.
func Yuntongxun() *YuntongxunSMSBuilder {
	return newYuntongxunSMSBuilder()
}

// Yunpian creates a new Yunpian SMS message builder.
func Yunpian() *YunpianSMSBuilder {
	return newYunpianSMSBuilder()
}

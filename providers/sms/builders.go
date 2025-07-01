/*
Package sms provides a unified API for sending SMS messages across multiple providers.

Usage (链式Builder风格)：

	// 文本短信
	msg := sms.Aliyun().To("***REMOVED***").Content("您的验证码是1234").SignName("签名").Build()
	msg := sms.Tencent().To("***REMOVED***").Content("您的验证码是1234").SignName("签名").Build()
	msg := sms.Huawei().To("***REMOVED***").Content("您的验证码是1234").SignName("签名").Build()

	// 模板短信（无序参数）


	msg := sms.Aliyun().To("***REMOVED***").TemplateID("SMS_123456").Params(map[string]string{"code": "1234"}).SignName("签名").Build()
	msg := sms.Yunpian().To("***REMOVED***").TemplateID("模板ID").Params(map[string]string{"code": "1234"}).Build()

	// 模板短信（有序参数）
	msg := sms.Tencent().To("***REMOVED***").TemplateID("模板ID").ParamsOrder([]string{"1234"}).SignName("签名").Build()
	msg := sms.Huawei().To("***REMOVED***").TemplateID("模板ID").ParamsOrder([]string{"1234"}).SignName("签名").Build()

	// 语音短信


	msg := sms.Aliyun().To("***REMOVED***").Type(sms.Voice).TemplateID("TTS_123456").Params(map[string]string{"code": "1234"}).CalledShowNumber("400xxxxxxx").Build()
	msg := sms.Tencent().To("***REMOVED***").Type(sms.Voice).TemplateID("模板ID").ParamsOrder([]string{"1234"}).Build()

	// 彩信
	msg := sms.Aliyun().To("***REMOVED***").Type(sms.MMS).Build()

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
type BaseBuilder struct {
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
}

func (b *BaseBuilder) To(mobiles ...string) *BaseBuilder {
	b.mobiles = mobiles
	return b
}
func (b *BaseBuilder) Content(content string) *BaseBuilder {
	b.content = content
	return b
}
func (b *BaseBuilder) SignName(sign string) *BaseBuilder {
	b.signName = sign
	return b
}
func (b *BaseBuilder) TemplateID(id string) *BaseBuilder {
	b.templateID = id
	return b
}
func (b *BaseBuilder) Params(params map[string]string) *BaseBuilder {
	b.templateParams = params
	return b
}
func (b *BaseBuilder) ParamsOrder(order []string) *BaseBuilder {
	b.paramsOrder = order
	return b
}
func (b *BaseBuilder) Category(category MessageCategory) *BaseBuilder {
	b.category = category
	return b
}
func (b *BaseBuilder) Type(t MessageType) *BaseBuilder {
	b.msgType = t
	return b
}
func (b *BaseBuilder) RegionCode(code int) *BaseBuilder {
	b.regionCode = code
	return b
}
func (b *BaseBuilder) CallbackURL(url string) *BaseBuilder {
	b.callbackURL = url
	return b
}
func (b *BaseBuilder) ScheduledAt(at time.Time) *BaseBuilder {
	b.scheduledAt = &at
	return b
}
func (b *BaseBuilder) Extend(ext string) *BaseBuilder {
	b.extend = ext
	return b
}
func (b *BaseBuilder) UID(uid string) *BaseBuilder {
	b.uid = uid
	return b
}

func (b *BaseBuilder) Build() *Message {
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

// BaseSMSBuilder provides common functionality for all SMS builders.
type BaseSMSBuilder struct {
	subProvider SubProviderType
}

// Aliyun returns an Aliyun SMS builder.
//   - SMS(普通文本): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendsms
//   - MMS(卡片消息): https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-sendcardsms
//   - Voice(验证码或文本转语音): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbytts
//   - Voice(通知): https://help.aliyun.com/zh/vms/developer-reference/api-dyvmsapi-2017-05-25-singlecallbyvoice
func Aliyun() *AliyunSMSBuilder {
	return newAliyunSMSBuilder()
}

// Tencent returns a Tencent SMS builder.
//   - SMS(普通文本): https://cloud.tencent.com/document/product/382/55981
//   - Voice(语音验证码API): https://cloud.tencent.com/document/product/1128/51559
//   - Voice(语音通知API): https://cloud.tencent.com/document/product/1128/51558
func Tencent() *TencentSMSBuilder {
	return newTencentSMSBuilder()
}

// Cl253 returns a CL253 SMS builder.
//   - SMS(普通文本): https://doc.chuanglan.com/
//   - Voice(语音验证码API): https://doc.chuanglan.com/
func Cl253() *Cl253SMSBuilder {
	return newCl253SMSBuilder()
}

// Huawei returns a Huawei SMS builder.
//   - SMS(普通文本): https://support.huaweicloud.com/api-msgsms/sms_05_0001.html
func Huawei() *HuaweiSMSBuilder {
	return newHuaweiSMSBuilder()
}

// Juhe returns a Juhe SMS builder.
//   - SMS国内短信API: https://www.juhe.cn/docs/api/id/54
//   - SMS国际短信API: https://www.juhe.cn/docs/api/id/357
//   - MMS视频短信API: https://www.juhe.cn/docs/api/id/363
func Juhe() *JuheSMSBuilder {
	return newJuheSMSBuilder()
}

// Luosimao returns a Luosimao SMS builder.
//   - SMS(普通文本): https://luosimao.com/docs/api
//   - Voice(语音验证码API): https://luosimao.com/docs/api/51
func Luosimao() *LuosimaoSMSBuilder {
	return newLuosimaoSMSBuilder()
}

// Smsbao returns a Smsbao SMS builder.
//   - SMS(国内文本): https://www.smsbao.com/openapi
//   - SMS(国际短信): https://www.smsbao.com/openapi/299.html
//   - Voice(语音短信): https://www.smsbao.com/openapi/214.html
func Smsbao() *SmsbaoSMSBuilder {
	return newSmsbaoSMSBuilder()
}

// Submail returns a Submail SMS builder.
//   - SMS国内短信: https://www.mysubmail.com/documents/FppOR3
//   - SMS国际短信: https://www.mysubmail.com/documents/3UQA3
//   - SMS模板短信: https://www.mysubmail.com/documents/OOVyh
//   - SMS群发: https://www.mysubmail.com/documents/AzD4Z4
//   - Voice语音: https://www.mysubmail.com/documents/meE3C1
//   - MMS彩信: https://www.mysubmail.com/documents/N6ktR
func Submail() *SubmailSMSBuilder {
	return newSubmailSMSBuilder()
}

// Ucp returns a UCP SMS builder.
//   - SMS(普通文本): http://docs.ucpaas.com/doku.php
func Ucp() *UcpSMSBuilder {
	return newUcpSMSBuilder()
}

// Volc returns a Volc SMS builder.
//   - SMS(普通文本): https://www.volcengine.com/product/cloud-sms
func Volc() *VolcSMSBuilder {
	return newVolcSMSBuilder()
}

// Yuntongxun returns a Yuntongxun SMS builder.
//   - SMS(普通文本): https://www.yuntongxun.com/developer-center
func Yuntongxun() *YuntongxunSMSBuilder {
	return newYuntongxunSMSBuilder()
}

// Yunpian returns a Yunpian SMS builder.
//   - SMS(普通文本): https://www.yunpian.com/dev-doc
func Yunpian() *YunpianSMSBuilder {
	return newYunpianSMSBuilder()
}

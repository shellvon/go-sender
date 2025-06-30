package sms

// Internal constants for SMS message extra fields
// These constants are used internally by transformers and should not be exposed to users

// Aliyun specific constants
const (
	// aliyunOutID is the key for Aliyun external flow ID
	aliyunOutID = "OutId"

	// aliyunCalledShowNumber is the key for Aliyun called show number
	aliyunCalledShowNumber = "CalledShowNumber"

	// aliyunPlayTimes is the key for Aliyun play times
	aliyunPlayTimes = "PlayTimes"

	// aliyunVolume is the key for Aliyun volume
	aliyunVolume = "Volume"

	// aliyunSpeed is the key for Aliyun speed
	aliyunSpeed = "Speed"

	// aliyunSmsUpExtendCode is the key for Aliyun SMS up extend code
	aliyunSmsUpExtendCode = "SmsUpExtendCode"

	// aliyunFallbackType is the key for Aliyun fallback type
	aliyunFallbackType = "FallbackType"

	// aliyunSmsTemplateCode is the key for Aliyun SMS template code
	aliyunSmsTemplateCode = "SmsTemplateCode"

	// aliyunDigitalTemplateCode is the key for Aliyun digital template code
	aliyunDigitalTemplateCode = "DigitalTemplateCode"

	// aliyunSmsTemplateParam is the key for Aliyun SMS template param
	aliyunSmsTemplateParam = "SmsTemplateParam"

	// aliyunDigitalTemplateParam is the key for Aliyun digital template param
	aliyunDigitalTemplateParam = "DigitalTemplateParam"

	// aliyunCardObjects is the key for Aliyun card objects
	aliyunCardObjects = "CardObjects"
)

// Tencent specific constants
const (
	// tencentSmsSdkAppID is the key for Tencent SMS SDK app ID
	tencentSmsSdkAppID = "SmsSdkAppId"

	// tencentExtendCode is the key for Tencent extend code
	tencentExtendCode = "ExtendCode"

	// tencentSenderID is the key for Tencent sender ID
	tencentSenderID = "SenderId"

	// tencentRegion is the key for Tencent region
	tencentRegion = "Region"
)

// CL253 specific constants
const (
	// cl253Report is the key for CL253 report
	cl253Report = "report"

	// cl253SendTime is the key for CL253 send time
	cl253SendTime = "sendtime"

	// cl253Extend is the key for CL253 extend
	cl253Extend = "extend"

	// cl253SenderID is the key for CL253 sender ID
	cl253SenderID = "senderId"

	// cl253TemplateID is the key for CL253 template ID
	cl253TemplateID = "templateId"

	// cl253TDFlag is the key for CL253 TD flag
	cl253TDFlag = "tdFlag"
)

// Huawei specific constants
const (
	// huaweiFrom is the key for Huawei from field
	huaweiFrom = "from"

	// huaweiStatusCallback is the key for Huawei status callback
	huaweiStatusCallback = "statusCallback"

	// huaweiExtend is the key for Huawei extend
	huaweiExtend = "extend"
)

// Yunpian specific constants
const (
	// yunpianExtend is the key for Yunpian extend
	yunpianExtend = "extend"
)

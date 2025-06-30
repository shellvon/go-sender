package sms

// Internal constants for SMS message extra fields
// These constants are used internally by transformers and should not be exposed to users

// Aliyun specific constants.
const (
	// aliyunOutIDKey is the key for Aliyun external flow ID.
	aliyunOutIDKey = "OutId"

	// aliyunCalledShowNumberKey is the key for Aliyun called show number.
	aliyunCalledShowNumberKey = "CalledShowNumber"

	// aliyunPlayTimesKey is the key for Aliyun play times.
	aliyunPlayTimesKey = "PlayTimes"

	// aliyunVolumeKey is the key for Aliyun volume.
	aliyunVolumeKey = "Volume"

	// aliyunSpeedKey is the key for Aliyun speed.
	aliyunSpeedKey = "Speed"

	// aliyunSmsUpExtendCodeKey is the key for Aliyun SMS up extend code.
	aliyunSmsUpExtendCodeKey = "SmsUpExtendCode"

	// aliyunFallbackTypeKey is the key for Aliyun fallback type.
	aliyunFallbackTypeKey = "FallbackType"

	// aliyunSmsTemplateCodeKey is the key for Aliyun SMS template code.
	aliyunSmsTemplateCodeKey = "SmsTemplateCode"

	// aliyunDigitalTemplateCodeKey is the key for Aliyun digital template code.
	aliyunDigitalTemplateCodeKey = "DigitalTemplateCode"

	// aliyunSmsTemplateParamKey is the key for Aliyun SMS template param.
	aliyunSmsTemplateParamKey = "SmsTemplateParam"

	// aliyunDigitalTemplateParamKey is the key for Aliyun digital template param.
	aliyunDigitalTemplateParamKey = "DigitalTemplateParam"

	// aliyunCardObjectsKey is the key for Aliyun card objects.
	aliyunCardObjectsKey = "CardObjects"
)

// Tencent specific constants.
const (
	// tencentSmsSdkAppIDKey is the key for Tencent SMS SDK app ID.
	tencentSmsSdkAppIDKey = "SmsSdkAppId"

	// tencentExtendCodeKey is the key for Tencent extend code.
	tencentExtendCodeKey = "ExtendCode"

	// tencentSenderIDKey is the key for Tencent sender ID.
	tencentSenderIDKey = "SenderId"

	// tencentRegionKey is the key for Tencent region.
	tencentRegionKey = "Region"

	// tencentPlayTimesKey is the key for Tencent play times.
	tencentPlayTimesKey = "PlayTimes"

	// tencentVoiceSdkAppIDKey is the key for Tencent Voice SdkAppId (for voice SMS).
	tencentVoiceSdkAppIDKey = "VoiceSdkAppId"
)

// CL253 specific constants.
const (
	// cl253ReportKey is the key for CL253 report.
	cl253ReportKey = "report"

	// cl253SenderIDKey is the key for CL253 sender ID.
	cl253SenderIDKey = "senderId"

	// cl253TDFlagKey is the key for CL253 TD flag.
	cl253TDFlagKey = "tdFlag"
)

// Huawei specific constants.
const (
	// huaweiFromKey is the key for Huawei from field.
	huaweiFromKey = "from"
)

// SMSBao specific constants.
const (
	// smsbaoProductIDKey is the key for SMSBao product ID.
	smsbaoProductIDKey = "productId"
)

// Submail extra keys.
const (
	// submailTagKey is the key for Submail message tag (used for message tracking, max 32 chars).
	submailTagKey = "tag"
	// submailSenderKey is the key for Submail sender identifier (mainly for international SMS, optional).
	submailSenderKey = "sender"
	// submailSignTypeKey is the key for Submail signature type (optional: md5 (default), sha1, normal).
	submailSignTypeKey = "sign_type"
)

// Volc specific constants.
const (
	// volcTagKey is the key for Volc tag.
	volcTagKey = "Tag"
)

// Yunpian specific constants.
const (

	// yunpianRegisterKey is the key for Yunpian register.
	yunpianRegisterKey = "register"

	// yunpianMobileStatKey is the key for Yunpian mobile stat.
	yunpianMobileStatKey = "mobile_stat"
)

// YunTongXun specific constants.
const (
	// yuntongxunPlayTimesKey is the key for YunTongXun play times.
	yuntongxunPlayTimesKey = "playTimes"

	// yuntongxunMediaNameKey is the key for YunTongXun media name.
	yuntongxunMediaNameKey = "mediaName"

	// yuntongxunDisplayNumKey is the key for YunTongXun display num.
	yuntongxunDisplayNumKey = "displayNum"

	// yuntongxunUserDataKey is the key for YunTongXun user data.
	yuntongxunUserDataKey = "userData"

	// yuntongxunMaxCallTimeKey is the key for YunTongXun max call time.
	yuntongxunMaxCallTimeKey = "maxCallTime"
)

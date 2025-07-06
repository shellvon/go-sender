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

	// aliyunDefaultVolume is the default volume for Aliyun voice calls (0-100).
	aliyunDefaultVolume = 100

	// aliyunSpeedKey is the key for Aliyun speed.
	aliyunSpeedKey = "Speed"

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

	// aliyunRegionKey is the key for Aliyun region.
	aliyunRegionKey = "Region"
)

// Tencent specific constants.
const (
	// tencentSmsSdkAppIDKey is the key for Tencent SMS SDK app ID.
	tencentSmsSdkAppIDKey = "SmsSdkAppId"

	// tencentSenderIDKey is the key for Tencent sender ID.
	tencentSenderIDKey = "SenderId"

	// tencentRegionKey is the key for Tencent region.
	tencentRegionKey = "Region"

	// tencentPlayTimesKey is the key for Tencent play times.
	tencentPlayTimesKey = "PlayTimes"

	// tencentVoiceSdkAppIDKey is the key for Tencent Voice SdkAppId (for voice SMS).
	tencentVoiceSdkAppIDKey = "VoiceSdkAppid"

	// tencentDefaultRegion is the default region for Tencent SMS.
	tencentDefaultRegion = "ap-guangzhou"

	// tencentDefaultPlayTimes is the default play times for Tencent voice SMS.
	tencentDefaultPlayTimes = 2
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

	// huaweiRegionKey is the key for Huawei region.
	huaweiRegionKey = "region"
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

	// submailDefaultSignType is the default sign type for Submail SMS.
	submailDefaultSignType = "md5"
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

	// yuntongxunMediaNameTypeKey is the key for YunTongXun media name type.
	yuntongxunMediaNameTypeKey = "mediaNameType"

	// yuntongxunDisplayNumKey is the key for YunTongXun display num.
	yuntongxunDisplayNumKey = "displayNum"

	// yuntongxunUserDataKey is the key for YunTongXun user data.
	yuntongxunUserDataKey = "userData"

	// yuntongxunMaxCallTimeKey is the key for YunTongXun max call time.
	yuntongxunMaxCallTimeKey = "maxCallTime"

	// yuntongxunWelcomePromptKey is the key for YunTongXun welcome prompt.
	yuntongxunWelcomePromptKey = "welcomePrompt"

	// yuntongxunPlayVerifyCodeKey is the key for YunTongXun play verify code.
	yuntongxunPlayVerifyCodeKey = "playVerifyCode"

	// yuntongxunRegionKey is the key for YunTongXun region.
	yuntongxunRegionKey = "region"

	// yuntongxunTxtSpeedKey is the key for YunTongXun txt speed.
	yuntongxunTxtSpeedKey = "txtSpeed"

	// yuntongxunTxtPitchKey is the key for YunTongXun txt pitch.
	yuntongxunTxtPitchKey = "txtPitch"

	// yuntongxunTxtVolumeKey is the key for YunTongXun txt volume.
	yuntongxunTxtVolumeKey = "txtVolume"

	// yuntongxunTxtBgsoundKey is the key for YunTongXun txt bgsound.
	yuntongxunTxtBgsoundKey = "txtBgsound"

	// yuntongxunPlayModeKey is the key for YunTongXun play mode.
	yuntongxunPlayModeKey = "playMode"
)

// Luosimao specific constants.
const (
	// luosimaoScheduledAtKey is the key for Luosimao scheduled at.
	luosimaoScheduledAtKey = "time"
)

// ProviderType represents different SMS service providers.
type ProviderType string

const (
	ProviderTypeAliyun     ProviderType = "aliyun"     // 阿里云短信（支持国内和国际）
	ProviderTypeCl253      ProviderType = "cl253"      // 蓝创253
	ProviderTypeSmsbao     ProviderType = "smsbao"     // 短信宝
	ProviderTypeJuhe       ProviderType = "juhe"       // 聚合服务
	ProviderTypeLuosimao   ProviderType = "luosimao"   // 螺丝帽
	ProviderTypeHuawei     ProviderType = "huawei"     // 华为云短信
	ProviderTypeUcp        ProviderType = "ucp"        // 云之讯
	ProviderTypeYunpian    ProviderType = "yunpian"    // 云片短信（支持国内和国际）
	ProviderTypeSubmail    ProviderType = "submail"    // 赛邮短信（支持国内和国际）
	ProviderTypeVolc       ProviderType = "volc"       // 火山引擎短信
	ProviderTypeYuntongxun ProviderType = "yuntongxun" // 云讯通（容联云通讯）
	ProviderTypeTencent    ProviderType = "tencent"    // 腾讯云短信（支持国内和国际）
)

package sms

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Yuntongxun is a message provider for Yuntongxun.
// It supports sending text, voice message
//
// Reference:
//   - Official Website: https://www.yuntongxun.com/
//   - API Docs: https://www.yuntongxun.com/developer-center
//   - SMS China Mainland: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - SMS International: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d

const (
	cloopenEndpoint      = "app.cloopen.com:8883"
	cloopenHKEndpoint    = "hksms.cloopen.com:8883"
	cloopenDefaultRegion = "cn"
)

type yuntongxunTransformer struct {
	*BaseTransformer
}

func newYuntongxunTransformer() *yuntongxunTransformer {
	transformer := &yuntongxunTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(SubProviderYunpian),
		&core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: true,
			Path:      "statusCode",
			Expect:    "000000",
			MsgPath:   "statusMsg",
			Mode:      core.MatchEq,
		},
		nil,
		WithSMSHandler(transformer.transformSMS),
		WithVoiceHandler(transformer.transformVoice),
	)
	return transformer
}
func init() {
	RegisterTransformer(string(SubProviderYuntongxun), newYuntongxunTransformer())
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
func (t *yuntongxunTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, NewProviderError(string(SubProviderYuntongxun), "MISSING_PARAM", "mobiles is required")
	}
	if msg.IsIntl() {
		if msg.Content == "" {
			return nil, nil, NewProviderError(
				string(SubProviderYuntongxun),
				"MISSING_PARAM",
				"international sms requires content",
			)
		}
	} else {
		if msg.TemplateID == "" {
			return nil, nil, NewProviderError(string(SubProviderYuntongxun), "MISSING_PARAM", "domestic sms requires templateID")
		}
	}
	// 判断是否为国际短信
	if msg.IsIntl() {
		return t.transformIntlSMS(msg, account)
	}
	return t.transformDomesticSMS(msg, account)
}

// 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
func (t *yuntongxunTransformer) transformDomesticSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 构建请求体
	data := map[string]interface{}{
		"to":         strings.Join(msg.Mobiles, ","),
		"appId":      account.AppID,
		"templateId": msg.TemplateID,
		"datas":      msg.ParamsOrder,
	}

	// 构建完整URL
	endpoint := cloopenEndpoint
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/SMS/TemplateSMS?sig=%s",
		endpoint, "2013-12-26", account.APIKey, t.generateSignature(account))

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, NewProviderError(
			string(SubProviderYuntongxun),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal yuntongxun request body: %v", err),
		)
	}

	return &core.HTTPRequestSpec{
		Method:  http.MethodPost,
		URL:     url,
		Headers: t.buildHeaders(account),
		Body:    bodyData,
	}, nil, nil
}

// 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
func (t *yuntongxunTransformer) transformIntlSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 构建请求体
	data := map[string]interface{}{
		"mobile":  strings.Join(msg.Mobiles, ","),
		"content": utils.AddSignature(msg.Content, msg.SignName),
		"appId":   account.APIKey,
	}

	region := utils.FirstNonEmpty(
		msg.GetExtraStringOrDefault(yuntongxunRegionKey, ""),
		account.Region,
		cloopenDefaultRegion,
	)
	var endpoint string
	if region == cloopenDefaultRegion {
		endpoint = cloopenEndpoint
	} else {
		endpoint = cloopenHKEndpoint
	}
	url := fmt.Sprintf("https://%s/%s/account/%s/international/send?sig=%s",
		endpoint, "v2", account.APIKey, t.generateSignature(account))

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, NewProviderError(
			string(SubProviderYuntongxun),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal yuntongxun international request body: %v", err),
		)
	}

	return &core.HTTPRequestSpec{
		Method:  http.MethodPost,
		URL:     url,
		Headers: t.buildHeaders(account),
		Body:    bodyData,
	}, nil, nil
}

// transformVoice transforms voice message to HTTP request
// 语音验证码:
//   - 默认: https://doc.yuntongxun.com/pe/5a533de43b8496dd00dce07e
//   - 自定义: https://doc.yuntongxun.com/pe/5a533de53b8496dd00dce080
//
// 外呼通知
//   - 语音通知: https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
func (t *yuntongxunTransformer) transformVoice(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, NewProviderError(string(SubProviderYuntongxun), "MISSING_PARAM", "mobiles is required")
	}
	if msg.IsIntl() {
		return nil, nil, NewProviderError(
			string(SubProviderYuntongxun),
			"UNSUPPORTED_COUNTRY",
			"voice sms only supports domestic mobile",
		)
	}
	// 只支持国内
	if msg.IsIntl() {
		return nil, nil, NewUnsupportedInternationalError(string(SubProviderYuntongxun), "voice call")
	}

	// 构建请求体
	//  https://doc.yuntongxun.com/pe/5a5342c73b8496dd00dce139
	body := map[string]interface{}{
		"appId": account.AppID,
		"to":    strings.Join(msg.Mobiles, ","),
		// 可选 语音文件名称，格式 wav，播放多个文件用英文分号隔开。与mediaTxt不能同时为空。当不为空时mediaTxt属性失效。测试用默认语音：ccp_marketingcall.wav
		"mediaName": msg.GetExtraStringOrDefault(yuntongxunMediaNameKey, ""),
		// 可选 语音文件名的类型，默认值为0，表示用户语音文件；　值为1表示平台通用文件。
		"mediaNameType": msg.GetExtraStringOrDefault(yuntongxunMediaNameTypeKey, ""),
		// 可选 文本内容，文本中汉字要求utf8编码，默认值为空。当mediaName为空才有效。
		"mediaTxt": msg.Content,

		// 来电显示的号码，根据平台侧显号规则控制(有显号需求请联系云通讯商务，并且说明显号的方式)，不在平台规则内或空则显示云通讯平台默认号码。默认值空。
		// 注：来电显示的号码不能和呼叫的号码相同，否则显示云通讯平台默认号码。
		"displayNum": msg.GetExtraStringOrDefault(yuntongxunDisplayNumKey, ""),

		// 循环播放次数，1－3次，默认播放1次。
		"playTimes": msg.GetExtraStringOrDefault(yuntongxunPlayTimesKey, ""),
		// 云通讯平台将向该Url地址发送呼叫结果通知。
		"respUrl": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
		// 可选 用户数据，透传字段，可填入任意字符串，如：用户id，用户名等。
		"userData": msg.GetExtraStringOrDefault(yuntongxunUserDataKey, ""),
		// 可选 文本转语音的语速，默认值为空。文本转语音后的发音速度，取值范围：-50至50，当mediaTxt有效才生效,默认值为0。
		"txtSpeed": msg.GetExtraStringOrDefault(yuntongxunTxtSpeedKey, ""),
		// 可选 文本转语音后的音调，取值范围：-500至500，当mediaTxt有效才生效，默认值为0。
		"txtPitch": msg.GetExtraStringOrDefault(yuntongxunTxtPitchKey, ""),
		// 可选 文本转语音后的音量大小，取值范围：-20至20，当mediaTxt有效才生效，默认值为0。
		"txtVolume": msg.GetExtraStringOrDefault(yuntongxunTxtVolumeKey, ""),
		// 文本转语音后的背景音编号，目前云通讯平台支持6种背景音，1到6的六种背景音编码，0为不需要背景音。暂时不支持第三方自定义背景音。当mediaTxt有效才生效。
		"txtBgsound": msg.GetExtraStringOrDefault(yuntongxunTxtBgsoundKey, ""),
		// 可选 是否同时播放文本和语音文件 , 0、否 1、是，默认0。优先播放文本。
		"playMode": msg.GetExtraStringOrDefault(yuntongxunPlayModeKey, ""),
	}

	// 构建完整URL
	endpoint := cloopenEndpoint
	action := "LandingCalls"

	if msg.Category == CategoryVerification {
		action = "VoiceVerify"

		body["maxCallTime"] = msg.GetExtraStringOrDefault(yuntongxunMaxCallTimeKey, "")
		body["welcomePrompt"] = msg.GetExtraStringOrDefault(yuntongxunWelcomePromptKey, "")
		body["respUrl"] = utils.FirstNonEmpty(msg.CallbackURL, account.Callback)
		body["verifyCode"] = msg.Content
		// 对于自定义的语音验证码，需要设置playVerifyCode
		body["playVerifyCode"] = msg.GetExtraStringOrDefault(yuntongxunPlayVerifyCodeKey, "")
	}
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/Calls/%s?sig=%s",
		endpoint, "2013-12-26", account.APIKey, action, t.generateSignature(account))

	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, nil, NewProviderError(
			string(SubProviderYuntongxun),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal yuntongxun voice request body: %v", err),
		)
	}

	return &core.HTTPRequestSpec{
		Method:  http.MethodPost,
		URL:     url,
		Headers: t.buildHeaders(account),
		Body:    bodyData,
	}, nil, nil
}

// generateSignature 生成云讯通API签名.
func (t *yuntongxunTransformer) generateSignature(account *Account) string {
	datetime := time.Now().Format("20060102150405")
	return strings.ToUpper(utils.HashHex(md5.New, []byte(account.APIKey+account.APISecret+datetime)))
}

// buildHeaders 构建云讯通API请求头.
func (t *yuntongxunTransformer) buildHeaders(account *Account) map[string]string {
	datetime := time.Now().Format("20060102150405")
	auth := fmt.Sprintf("%s:%s", account.APIKey, datetime)
	return map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64EncodeBytes([]byte(auth)),
	}
}

package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// huaweiTransformer implements HTTPRequestTransformer for Huawei Cloud SMS.
// It supports sending text message.
//
// Reference:
//   - Official Website: https://www.huaweicloud.com
//   - API Docs: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//   - SMS API(Domestic): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//   - SMS API(International): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html

const (
	huaweiAPIPath       = "https://api.rtc.huaweicloud.com:10443/sms/batchSendSms/v1"
	huaweiSuccessCode   = "000000"
	huaweiDefaultRegion = "cn-north-1"
)

// huaweiTransformer implements HTTPRequestTransformer for Huawei Cloud SMS.
type huaweiTransformer struct {
	*BaseTransformer
}

// init automatically registers the Huawei transformer.
func init() {
	RegisterTransformer(string(SubProviderHuawei), newHuaweiTransformer())
}

func newHuaweiTransformer() *huaweiTransformer {
	transformer := &huaweiTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(SubProviderHuawei),
		&core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: true,
			Path:      "code",
			Expect:    "000000",
			Mode:      core.MatchEq,
		},
		nil,
		WithSMSHandler(transformer.transformSMS),
	)
	return transformer
}

// transformSMS transforms SMS message to HTTP request
//
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 如果"to"参数携带的号码中包含除数字和+之外的其他字符，则无法向该参数携带的所有号码发送短信。如果"to"参数携带的所有号码只包含数字和+，
// 但部分号码不符合号码规则要求，则在响应消息中会通过状态码标识发送失败的号码，不影响其他正常号码的短信发送。号码之间以英文逗号分隔，
// 每个号码最大长度为21位，最多允许携带500个号码。如果携带超过500个号码，则全部号码都会发送失败。
func (t *huaweiTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 格式化手机号
	formattedMobiles := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		formattedMobiles[i] = t.formatHuaweiPhoneNumber(mobile, msg.RegionCode)
	}
	paras := ""
	if len(msg.ParamsOrder) > 0 {
		b, _ := json.Marshal(msg.ParamsOrder)
		paras = string(b)
	}
	// 构建参数，注释与每个参数紧密对应
	params := url.Values{}
	params.Set("from", msg.GetExtraStringOrDefault(huaweiFromKey, ""))
	params.Set("to", strings.Join(formattedMobiles, ","))
	params.Set("templateId", msg.TemplateID)
	params.Set("templateParas", paras)

	// 签名名称，必须是已审核通过的，与模板类型一致的签名名称。
	// 仅中国大陆短信可携带此参数。
	// 仅在templateId指定的模板类型为通用模板时生效且必填，用于指定在通用模板短信内容前面补充的签名
	if msg.IsDomestic() {
		params.Set("signature", msg.SignName)
	}
	// 客户的回调地址，用于接收短信状态报告，如：http://my.com/receiveSMSReport。
	// 如果设置了该字段，则该消息的状态报告将通过"接收状态报告"接口直接通知客户。
	// 如果未设置该字段，则短信平台收到运营商短信中心返回的状态报告不会推送给客户，该状态报告将在短信平台中保存1个小时，超时后系统会自动删除。
	// 回调地址推荐使用域名。
	params.Set("statusCallback", msg.CallbackURL)
	// 扩展参数，在状态报告中会原样返回。
	// 不允许赋空值，不允许携带以下字符："{","}"（即大括号）。
	params.Set("extend", msg.Extend)

	body := []byte(params.Encode())

	appKey := account.APIKey
	appSecret := account.APISecret

	headers := t.buildHeaders(appKey, appSecret)

	region := utils.FirstNonEmpty(msg.GetExtraStringOrDefault(huaweiRegionKey, ""), account.Region, huaweiDefaultRegion)
	reqSpec := &core.HTTPRequestSpec{
		Method:  http.MethodPost,
		URL:     fmt.Sprintf("https://%s.myhuaweicloud.com/sms/batchSendSms/v1", region),
		Headers: headers,
		Body:    body,
	}
	return reqSpec, nil, nil
}

// formatHuaweiPhoneNumber formats phone number for Huawei Cloud API
// 华为云要求：标准号码格式为：+{国家码}{地区码}{终端号码}.
func (t *huaweiTransformer) formatHuaweiPhoneNumber(mobile string, regionCode int) string {
	if regionCode == 0 {
		regionCode = 86
	}
	return fmt.Sprintf("+%d%s", regionCode, mobile)
}

// buildHeaders 构建华为云短信请求头.
func (t *huaweiTransformer) buildHeaders(appKey, appSecret string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"",
		"X-WSSE":        t.buildHuaweiWsseHeader(appKey, appSecret),
	}
}

// buildHuaweiWsseHeader 构建 X-WSSE 认证头.
func (t *huaweiTransformer) buildHuaweiWsseHeader(appKey, appSecret string) string {
	now := time.Now().UTC().Format(time.RFC3339)
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	passwordDigest := utils.Base64EncodeBytes(utils.SHA256Sum([]byte(nonce + now + appSecret)))
	return fmt.Sprintf(
		"UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"",
		appKey, passwordDigest, nonce, now,
	)
}

package sms

// @ProviderName: Huawei Cloud / 华为云
// @Website: https://www.huaweicloud.com
// @APIDoc: https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// # HuaweiProvider implements SMSProviderInterface for Huawei Cloud SMS
//
// 官方文档:
//   - 短信API(国内/国际): https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
//
// 前置条件
// 已创建短信应用，获取Application Key、Application Secret、APP接入地址和通道号（中国大陆短信除外）。
// 仅中国大陆短信）已申请短信签名，获取签名通道号。
// 已申请短信模板，获取模板ID。
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID。
//   - 国际短信：支持单发和群发，需模板ID。
//   - 彩信/语音：暂不支持。
//
// endpoint 只需配置主机名（如 api.rtc.huaweicloud.com:10443），具体 path 由服务方法追加。
//
// 群发短信时，如果"to"参数携带的号码中包含除数字和+之外的其他字符，则无法向该参数携带的所有号码发送短信。如果"to"参数携带的所有号码只包含数字和+，但部分号码不符合号码规则要求，则在响应消息中会通过状态码标识发送失败的号码，不影响其他正常号码的短信发送。号码之间以英文逗号分隔，每个号码最大长度为21位，最多允许携带500个号码。如果携带超过500个号码，则全部号码都会发送失败。
// 根据短信内容的长度，一条长短信可能会被拆分为多条短信发送，拆分规则详见短信发送规则。
// 通过X-WSSE方式鉴权时，生成随机数的时间与发送请求时的本地时间的差值不能超过24小时，否则会导致鉴权失败。
// 请求body中的参数需要进行urlencode。
// 通过特殊AK/SK鉴权时，生成随机数的时间与发送请求时的本地时间的差值不能超过15分钟，否则会导致鉴权失败。
// Example endpoint: https://api.rtc.huaweicloud.com:10443
import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/utils"
)

// 华为云短信服务实现
// 仅支持模板短信发送，API文档：https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
const (
	huaweiEndpoint     = "https://api.rtc.huaweicloud.com:10443"
	huaweiEndpointPath = "/sms/batchSendSms/v1"
	huaweiSuccessCode  = "000000"
)

type HuaweiProvider struct {
	config SMSProvider
}

// NewHuaweiProvider creates a new Huawei Cloud SMS provider
func NewHuaweiProvider(config SMSProvider) *HuaweiProvider {
	return &HuaweiProvider{config: config}
}

// formatHuaweiPhoneNumber formats phone number for Huawei Cloud API
// 华为云要求：标准号码格式为：+{国家码}{地区码}{终端号码}
// 发送中国大陆短信，如果"+"不存在，则默认为+86
func formatHuaweiPhoneNumber(mobile string, regionCode int) string {
	if regionCode == 0 {
		regionCode = 86
	}
	return fmt.Sprintf("+%d%s", regionCode, mobile)
}

// Send sends an SMS message via Huawei Cloud
func (provider *HuaweiProvider) Send(ctx context.Context, msg *Message) error {
	if err := ValidateForSend(provider, msg); err != nil {
		return err
	}
	if msg.Type != SMSText {
		return NewUnsupportedMessageTypeError(string(ProviderTypeHuawei), msg.Type.String(), msg.Category.String())
	}
	return provider.sendSMS(ctx, msg)
}

// https://support.huaweicloud.com/intl/zh-cn/api-msgsms/sms_05_0001.html
// sendSMS sends SMS message via Huawei Cloud API
func (provider *HuaweiProvider) sendSMS(ctx context.Context, msg *Message) error {
	appKey := provider.config.AppID
	appSecret := provider.config.AppSecret
	//
	//	短信发送方的号码。
	//	全球短信填写为创建短信应用时分配的通道号，如：isms100000001。
	//	中国大陆短信填写为短信平台为短信签名分配的通道号码，请在申请短信签名时获取签名通道号，如：csms100000001，且通道号码对应的签名类型和模板ID对应的模板类型必须相同
	channel := provider.config.Channel
	templateId := msg.TemplateID
	paras := "[]"
	if len(msg.ParamsOrder) > 0 {
		b, _ := json.Marshal(msg.ParamsOrder)
		paras = string(b)
	}
	// 短信接收方的号码，标准号码格式为：+{国家码}{地区码}{终端号码}。
	// 发送全球短信：不区分接收号码类型，所填号码都必须符合标准号码格式。示例：+2412000000（加蓬号码）
	// 发送中国大陆短信，如果"+"不存在，则默认为+86，如果接收方号码为手机号码，则{地区码}可选。如：+8613112345678。
	// 如果携带多个接收方号码，则以英文逗号分隔。每个号码最大长度为21位，最多允许携带500个号码。

	// 格式化手机号
	formattedMobiles := make([]string, len(msg.Mobiles))
	for i, mobile := range msg.Mobiles {
		formattedMobiles[i] = formatHuaweiPhoneNumber(mobile, msg.RegionCode)
	}

	params := map[string]string{
		"from":          channel,
		"to":            strings.Join(formattedMobiles, ","),
		"templateId":    templateId,
		"templateParas": paras,
		// 签名名称，必须是已审核通过的，与模板类型一致的签名名称。
		// 仅中国大陆短信可携带此参数。
		// 签名名称，必须是已审核通过的，与模板类型一致的签名名称。
		// 仅在templateId指定的模板类型为通用模板时生效且必填，用于指定在通用模板短信内容前面补充的签名
		"signature": msg.SignName,
	}
	// 客户的回调地址，用于接收短信状态报告，如：http://my.com/receiveSMSReport。
	// 如果设置了该字段，则该消息的状态报告将通过"接收状态报告"接口直接通知客户。
	// 如果未设置该字段，则短信平台收到运营商短信中心返回的状态报告不会推送给客户，该状态报告将在短信平台中保存1个小时，超时后系统会自动删除。
	// 回调地址推荐使用域名。
	if callbackUrl := msg.GetExtraStringOrDefault("callbackUrl", provider.config.Callback); callbackUrl != "" {
		params["statusCallback"] = callbackUrl
	}
	//扩展参数，在状态报告中会原样返回。
	// 不允许赋空值，不允许携带以下字符："{""，""}"（即大括号）。
	if extend := msg.GetExtraStringOrDefault("extend", ""); extend != "" {
		params["extend"] = extend
	}

	headers := getHuaweiHeaders(appKey, appSecret)
	endpoint := getHuaweiEndpoint(&provider.config)
	resp, _, err := utils.DoRequest(ctx, endpoint, utils.RequestOptions{
		Method:  "POST",
		Headers: headers,
		Data:    params,
	})
	if err != nil {
		return fmt.Errorf("huawei SMS request failed: %w", err)
	}
	return parseHuaweiResponse(resp)
}

func getHuaweiEndpoint(provider *SMSProvider) string {
	endpoint := huaweiEndpoint
	if provider.Endpoint != "" {
		endpoint = provider.Endpoint
	}
	return strings.TrimRight(endpoint, "/") + huaweiEndpointPath
}

func getHuaweiHeaders(appKey, appSecret string) map[string]string {
	return map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "WSSE realm=\"SDP\",profile=\"UsernameToken\",type=\"Appkey\"",
		"X-WSSE":        buildHuaweiWsseHeader(appKey, appSecret),
	}
}

// X-WSSE 认证方式:
// 取值为UsernameToken Username="app_key的值", PasswordDigest="PasswordDigest的值", Nonce="随机数", Created="随机数生成时间"。
// PasswordDigest：根据PasswordDigest = Base64 (SHA256 (Nonce + Created + Password))生成，直接使用Nonce、Created、Password拼接后的字符串进行SHA256加密即可，字符串中无需包含+号和空格。其中，Password为app_secret的值。
// Nonce：客户发送请求时生成的一个随机数，长度为1~128位，可包含数字和大小写字母。例如：66C92B11FF8A425FB8D4CCFE0ED9ED1F。
// Created：随机数生成时间。采用标准UTC格式，例如：2018-02-12T15:30:20Z。不同编程语言中的时间格式转换方式不同，
func buildHuaweiWsseHeader(appKey, appSecret string) string {
	now := time.Now().UTC().Format(time.RFC3339)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	passwordDigest := utils.Base64EncodeBytes(utils.SHA256Sum([]byte(nonce + now + appSecret)))
	return fmt.Sprintf(
		"UsernameToken Username=\"%s\",PasswordDigest=\"%s\",Nonce=\"%s\",Created=\"%s\"",
		appKey, passwordDigest, nonce, now,
	)
}
func parseHuaweiResponse(resp []byte) error {
	var result struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse huawei response: %w", err)
	}
	if result.Code != huaweiSuccessCode {
		return &SMSError{
			Code:     result.Code,
			Message:  result.Description,
			Provider: string(ProviderTypeHuawei),
		}
	}
	return nil
}

func (p *HuaweiProvider) GetCapabilities() *Capabilities {
	capabilities := NewCapabilities()
	// 国内/国际短信均支持单发/群发
	capabilities.SMS.Domestic = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国内短信，需模板ID，支持单发和群发",
	)
	capabilities.SMS.International = NewRegionCapability(
		true, true,
		[]MessageType{SMSText},
		[]MessageCategory{CategoryVerification, CategoryNotification, CategoryPromotion},
		"支持国际短信，需模板ID，支持单发和群发",
	)
	capabilities.SMS.Limits.MaxBatchSize = 1000
	capabilities.SMS.Limits.MaxContentLen = 500
	capabilities.SMS.Limits.RateLimit = "未知"
	capabilities.SMS.Limits.DailyLimit = "未知"
	// 彩信、语音均不支持
	capabilities.MMS.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内彩信")
	capabilities.MMS.International = NewRegionCapability(false, false, nil, nil, "不支持国际彩信")
	capabilities.Voice.Domestic = NewRegionCapability(false, false, nil, nil, "不支持国内语音")
	capabilities.Voice.International = NewRegionCapability(false, false, nil, nil, "不支持国际语音")
	return capabilities
}
func (p *HuaweiProvider) CheckCapability(msg *Message) error {
	return DefaultCheckCapability(p, msg)
}
func (p *HuaweiProvider) GetLimits(msgType MessageType) Limits {
	capabilities := p.GetCapabilities()
	switch msgType {
	case SMSText:
		return capabilities.SMS.GetLimits()
	default:
		return Limits{}
	}
}
func (p *HuaweiProvider) GetName() string {
	return p.config.Name
}
func (p *HuaweiProvider) GetType() string {
	return string(p.config.Type)
}
func (p *HuaweiProvider) IsEnabled() bool {
	return !p.config.Disabled
}
func (p *HuaweiProvider) GetWeight() int {
	return p.config.GetWeight()
}
func (p *HuaweiProvider) CheckConfigured() error {
	if p.config.AppID == "" || p.config.AppSecret == "" || p.config.Channel == "" {
		return fmt.Errorf("huawei SMS provider requires AppID, AppSecret, and Channel")
	}
	return nil
}

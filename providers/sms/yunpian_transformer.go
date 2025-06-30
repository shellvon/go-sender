package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/official/document/sms/zh_CN/api_reference.html
//
// 官方文档:
//   - 短信API: https://www.yunpian.com/official/document/sms/zh_CN/api_reference.html
//   - 语音API: https://www.yunpian.com/official/document/sms/zh_CN/voice_send
//   - 超级短信API: https://www.yunpian.com/official/document/sms/zh_CN/super_sms_send
//
// transformer 支持 text（国内/国际/模板/非模板/群发）、voice（语音验证码）、mms（超级短信）类型。

type yunpianTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderYunpian), &yunpianTransformer{})
}

func (t *yunpianTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderYunpian)
}

func (t *yunpianTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for yunpianTransformer")
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}
	switch smsMsg.Type {
	case SMSText:
		return t.transformTextSMS(ctx, smsMsg, account)
	case Voice:
		return t.transformVoiceSMS(ctx, smsMsg, account)
	case MMS:
		return t.transformMMSSMS(ctx, smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported yunpian message type: %v", smsMsg.Type)
	}
}

func (t *yunpianTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.Type == SMSText && msg.Content == "" && msg.TemplateID == "" {
		return errors.New("content or templateID is required")
	}
	if msg.Type == Voice && msg.Content == "" {
		return errors.New("voice content is required")
	}
	if msg.Type == MMS && msg.TemplateID == "" {
		return errors.New("mms requires templateID")
	}
	return nil
}

func (t *yunpianTransformer) transformTextSMS(
	ctx context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 国际短信
	if msg.IsIntl() {
		if len(msg.Mobiles) > 1 {
			return nil, nil, errors.New("yunpian international SMS only supports single send")
		}
		if msg.TemplateID != "" {
			return nil, nil, errors.New("yunpian international SMS does not support template")
		}
		return t.transformIntlSMS(ctx, msg, account)
	}

	// 国内短信
	if msg.TemplateID != "" {
		if len(msg.Mobiles) > 1 {
			return t.transformTplBatchSMS(ctx, msg, account)
		}
		return t.transformTplSMS(ctx, msg, account)
	}
	if len(msg.Mobiles) > 1 {
		return t.transformBatchSMS(ctx, msg, account)
	}
	return t.transformSingleSMS(ctx, msg, account)
}

func (t *yunpianTransformer) transformSingleSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/single_send.json", account.Endpoint)
	params := map[string]string{
		"apikey": account.Secret,
		"mobile": msg.Mobiles[0],
		"text":   utils.AddSignature(msg.Content, msg.SignName),
	}

	// 处理统一的接口字段 - 适配到云片特定的key
	if msg.Extend != "" {
		params["extend"] = msg.Extend // 云片使用 extend
	} else if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if msg.UID != "" {
		params["uid"] = msg.UID // 云片使用 uid
	} else if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if msg.CallbackURL != "" {
		params["callback_url"] = msg.CallbackURL // 云片使用 callback_url
	} else if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}
	if reg, ok := msg.GetExtraBool("register"); ok {
		params["register"] = strconv.FormatBool(reg)
	}
	if stat, ok := msg.GetExtraBool("mobile_stat"); ok {
		params["mobile_stat"] = strconv.FormatBool(stat)
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformBatchSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/batch_send.json", account.Endpoint)
	params := map[string]string{
		"apikey": account.Secret,
		"mobile": strings.Join(msg.Mobiles, ","),
		"text":   utils.AddSignature(msg.Content, msg.SignName),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}
	if stat, ok := msg.GetExtraBool("mobile_stat"); ok {
		params["mobile_stat"] = strconv.FormatBool(stat)
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformTplSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/tpl_single_send.json", account.Endpoint)
	params := map[string]string{
		"apikey":    account.Secret,
		"mobile":    msg.Mobiles[0],
		"tpl_id":    msg.TemplateID,
		"tpl_value": t.buildTemplateValue(msg.TemplateParams),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformTplBatchSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/tpl_batch_send.json", account.Endpoint)
	params := map[string]string{
		"apikey":    account.Secret,
		"mobile":    strings.Join(msg.Mobiles, ","),
		"tpl_id":    msg.TemplateID,
		"tpl_value": t.buildTemplateValue(msg.TemplateParams),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformIntlSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/single_send.json", account.IntlEndpoint)
	params := map[string]string{
		"apikey": account.Secret,
		"mobile": fmt.Sprintf("+%d%s", msg.RegionCode, msg.Mobiles[0]),
		"text":   utils.AddSignature(msg.Content, msg.SignName),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformVoiceSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("voice", "/v2/voice/send.json", account.Endpoint)
	params := map[string]string{
		"apikey": account.Secret,
		"mobile": msg.Mobiles[0],
		"code":   msg.Content,
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformMMSSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("vsms", "/v2/vsms/tpl_batch_send.json", account.Endpoint)
	params := map[string]string{
		"apikey":    account.Secret,
		"mobile":    strings.Join(msg.Mobiles, ","),
		"tpl_id":    msg.TemplateID,
		"tpl_value": t.buildTemplateValue(msg.TemplateParams),
	}
	if ext := msg.GetExtraStringOrDefault("extend", ""); ext != "" {
		params["extend"] = ext
	}
	if uid := msg.GetExtraStringOrDefault("uid", ""); uid != "" {
		params["uid"] = uid
	}
	if cb := msg.GetExtraStringOrDefault("callback_url", account.Webhook); cb != "" {
		params["callback_url"] = cb
	}

	return t.buildRequest(endpoint, params)
}

// yunpianEndpoint 统一生成云片 API endpoint，支持可选自定义域名.
func (t *yunpianTransformer) yunpianEndpoint(service, path string, domainOverride string) string {
	domain := fmt.Sprintf("%s.yunpian.com", service)
	if domainOverride != "" {
		domain = domainOverride
	}
	return fmt.Sprintf("https://%s%s", domain, path)
}

// buildRequest 构建HTTP请求.
func (t *yunpianTransformer) buildRequest(
	endpoint string,
	params map[string]string,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	body := []byte(values.Encode())

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      endpoint,
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: "form",
	}, t.handleYunpianResponse, nil
}

// buildTemplateValue 构建模板参数值.
func (t *yunpianTransformer) buildTemplateValue(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	var pairs []string
	for key, value := range params {
		// Format: urlencode("#key#") + "=" + urlencode("value")
		encodedKey := url.QueryEscape("#" + key + "#")
		encodedValue := url.QueryEscape(value)
		pairs = append(pairs, encodedKey+"="+encodedValue)
	}
	sort.Strings(pairs) // Sort for consistent ordering
	return strings.Join(pairs, "&")
}

// handleYunpianResponse 处理云片API响应.
func (t *yunpianTransformer) handleYunpianResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse yunpian response: %w", err)
	}
	if result.Code != 0 {
		return &Error{
			Code:     strconv.Itoa(result.Code),
			Message:  result.Msg,
			Provider: string(SubProviderYunpian),
		}
	}
	return nil
}

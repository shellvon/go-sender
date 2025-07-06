package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Yunpian / 云片
// @Website: https://www.yunpian.com
// @APIDoc: https://www.yunpian.com/dev-doc
//
// 官方文档:
//   - 短信API: https://www.yunpian.com/official/document/sms/zh_CN/domestic_list
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
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Yunpian: %T", msg)
	}

	// Apply Yunpian-specific defaults
	t.applyYunpianDefaults(smsMsg, account)

	switch smsMsg.Type {
	case SMSText:
		return t.transformSMS(smsMsg, account)
	case Voice:
		return t.transformVoice(smsMsg, account)
	case MMS:
		return t.transformMMS(smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", smsMsg.Type)
	}
}

// applyYunpianDefaults applies Yunpian-specific defaults to the message.
func (t *yunpianTransformer) applyYunpianDefaults(msg *Message, account *Account) {
	// Apply common defaults first
	msg.ApplyCommonDefaults(account)
}

// transformSMS transforms SMS message to HTTP request.
func (t *yunpianTransformer) transformSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.Content == "" && msg.TemplateID == "" {
		return nil, nil, errors.New("content or templateID is required")
	}
	// 国际短信
	if msg.IsIntl() {
		if len(msg.Mobiles) > 1 {
			return nil, nil, errors.New("yunpian international SMS only supports single send")
		}
		if msg.TemplateID != "" {
			return nil, nil, errors.New("yunpian international SMS does not support template")
		}
		return t.transformIntlSMS(msg, account)
	}

	// 国内短信
	if msg.TemplateID != "" {
		if len(msg.Mobiles) > 1 {
			return t.transformTplBatchSMS(msg, account)
		}
		return t.transformTplSMS(msg, account)
	}
	if len(msg.Mobiles) > 1 {
		return t.transformBatchSMS(msg, account)
	}
	return t.transformSingleSMS(msg, account)
}

func (t *yunpianTransformer) transformSingleSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/single_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       msg.Mobiles[0],
		"text":         utils.AddSignature(msg.Content, msg.SignName),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
		"register":     strconv.FormatBool(msg.GetExtraBoolOrDefault(yunpianRegisterKey, false)),
		"mobile_stat":  strconv.FormatBool(msg.GetExtraBoolOrDefault(yunpianMobileStatKey, false)),
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformBatchSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/batch_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       strings.Join(msg.Mobiles, ","),
		"text":         utils.AddSignature(msg.Content, msg.SignName),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
		"register":     strconv.FormatBool(msg.GetExtraBoolOrDefault(yunpianRegisterKey, false)),
		"mobile_stat":  strconv.FormatBool(msg.GetExtraBoolOrDefault(yunpianMobileStatKey, false)),
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformTplSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/tpl_single_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       msg.Mobiles[0],
		"tpl_id":       msg.TemplateID,
		"tpl_value":    t.buildTemplateValue(msg.TemplateParams),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformTplBatchSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/tpl_batch_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       strings.Join(msg.Mobiles, ","),
		"tpl_id":       msg.TemplateID,
		"tpl_value":    t.buildTemplateValue(msg.TemplateParams),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
	}

	return t.buildRequest(endpoint, params)
}

func (t *yunpianTransformer) transformIntlSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	endpoint := t.yunpianEndpoint("sms", "/v2/sms/single_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       fmt.Sprintf("+%d%s", msg.RegionCode, msg.Mobiles[0]),
		"text":         utils.AddSignature(msg.Content, msg.SignName),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
	}

	return t.buildRequest(endpoint, params)
}

// transformVoice transforms voice message to HTTP request.
func (t *yunpianTransformer) transformVoice(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.Content == "" {
		return nil, nil, errors.New("voice content is required")
	}
	endpoint := t.yunpianEndpoint("voice", "/v2/voice/send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       msg.Mobiles[0],
		"code":         msg.Content,
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
	}

	return t.buildRequest(endpoint, params)
}

// transformMMS transforms MMS message to HTTP request.
func (t *yunpianTransformer) transformMMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.TemplateID == "" {
		return nil, nil, errors.New("mms requires templateID")
	}
	endpoint := t.yunpianEndpoint("vsms", "/v2/vsms/tpl_batch_send.json")
	params := map[string]string{
		"apikey":       account.APISecret,
		"mobile":       strings.Join(msg.Mobiles, ","),
		"tpl_id":       msg.TemplateID,
		"tpl_value":    t.buildTemplateValue(msg.TemplateParams),
		"extend":       msg.Extend,
		"uid":          msg.UID,
		"callback_url": utils.FirstNonEmpty(msg.CallbackURL, account.Callback),
	}

	return t.buildRequest(endpoint, params)
}

// yunpianEndpoint 统一生成云片 API.
func (t *yunpianTransformer) yunpianEndpoint(service, path string) string {
	domain := fmt.Sprintf("%s.yunpian.com", service)
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
		Method:   http.MethodPost,
		URL:      endpoint,
		Body:     body,
		BodyType: core.BodyTypeForm,
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
	if !utils.IsAcceptableStatus(statusCode) {
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

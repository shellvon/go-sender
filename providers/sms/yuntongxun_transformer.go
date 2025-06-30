package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Yuntongxun / 云讯通
// @Website: https://www.yuntongxun.com
// @APIDoc: https://www.yuntongxun.com/developer-center
//
// 官方文档:
//   - 国内短信: https://doc.yuntongxun.com/pe/5a533de33b8496dd00dce07c
//   - 国际短信: https://doc.yuntongxun.com/pe/604f29eda80948a1006e928d
//
// transformer 支持 text（国内/国际，国内模板，国际内容）和 voice（仅国内）类型。

type yuntongxunTransformer struct{}

func newYuntongxunTransformer() core.HTTPTransformer[*core.Account] {
	return &yuntongxunTransformer{}
}

func init() {
	RegisterTransformer(string(SubProviderYuntongxun), &yuntongxunTransformer{})
}

func (t *yuntongxunTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderYuntongxun)
}

func (t *yuntongxunTransformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for yuntongxun: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}
	switch smsMsg.Type {
	case SMSText:
		return t.transformTextSMS(ctx, smsMsg, account)
	case Voice:
		return t.transformVoiceSMS(ctx, smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported yuntongxun message type: %s", smsMsg.Type)
	}
}

func (t *yuntongxunTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.IsIntl() {
		if msg.Content == "" {
			return errors.New("international sms requires content")
		}
	} else {
		if msg.TemplateID == "" {
			return errors.New("domestic sms requires templateID")
		}
	}
	if msg.Type == Voice && msg.IsIntl() {
		return errors.New("voice sms only supports domestic mobile")
	}
	return nil
}

func (t *yuntongxunTransformer) transformTextSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 判断是否为国际短信
	if msg.IsIntl() {
		return t.transformIntlSMS(ctx, msg, account)
	}
	return t.transformDomesticSMS(ctx, msg, account)
}

func (t *yuntongxunTransformer) transformDomesticSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	datetime := time.Now().Format("20060102150405")
	accountSid := account.From
	accountToken := account.Secret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求头
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	// 构建请求体
	data := map[string]interface{}{
		"to":         strings.Join(msg.Mobiles, ","),
		"appId":      account.Key,
		"templateId": msg.TemplateID,
		"datas":      msg.ParamsOrder,
	}

	// 构建完整URL
	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = "app.cloopen.com:8883"
	}
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/SMS/TemplateSMS?sig=%s",
		endpoint, "2013-12-26", accountSid, sig)

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal yuntongxun request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  headers,
		Body:     bodyData,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}, t.handleYuntongxunResponse, nil
}

func (t *yuntongxunTransformer) transformIntlSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	datetime := time.Now().Format("20060102150405")
	accountSid := account.From
	accountToken := account.Secret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求头
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	// 构建请求体
	data := map[string]interface{}{
		"mobile":  strings.Join(msg.Mobiles, ","),
		"content": utils.AddSignature(msg.Content, msg.SignName),
		"appId":   account.Key,
	}

	// 构建完整URL
	endpoint := account.IntlEndpoint
	if endpoint == "" {
		endpoint = "app.cloopen.com:8883"
	}
	url := fmt.Sprintf("https://%s/%s/account/%s/international/send?sig=%s",
		endpoint, "v2", accountSid, sig)

	bodyData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal yuntongxun international request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  headers,
		Body:     bodyData,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}, t.handleYuntongxunResponse, nil
}

func (t *yuntongxunTransformer) transformVoiceSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 只支持国内
	if msg.IsIntl() {
		return nil, nil, NewUnsupportedInternationalError(string(SubProviderYuntongxun), "voice call")
	}

	datetime := time.Now().Format("20060102150405")
	accountSid := account.From
	accountToken := account.Secret

	// 生成签名
	sig := strings.ToUpper(utils.MD5Hex(accountSid + accountToken + datetime))

	// 构建请求头
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json;charset=utf-8",
		"Authorization": utils.Base64Encode(fmt.Sprintf("%s:%s", accountSid, datetime)),
	}

	// 构建请求体
	body := map[string]interface{}{
		"to":        strings.Join(msg.Mobiles, ","),
		"appId":     account.Key,
		"mediaTxt":  msg.Content,
		"playTimes": msg.GetExtraStringOrDefault("playTimes", "3"),
	}

	// 可选参数
	if v := msg.GetExtraStringOrDefault("mediaName", ""); v != "" {
		body["mediaName"] = v
	}
	if v := msg.GetExtraStringOrDefault("displayNum", ""); v != "" {
		body["displayNum"] = v
	}
	if v := msg.GetExtraStringOrDefault("respUrl", account.Webhook); v != "" {
		body["respUrl"] = v
	}
	if v := msg.GetExtraStringOrDefault("userData", ""); v != "" {
		body["userData"] = v
	}
	if v := msg.GetExtraStringOrDefault("maxCallTime", ""); v != "" {
		body["maxCallTime"] = v
	}

	// 构建完整URL
	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = "app.cloopen.com:8883"
	}
	url := fmt.Sprintf("https://%s/%s/Accounts/%s/Calls/VoiceNotify?sig=%s",
		endpoint, "2013-12-26", accountSid, sig)

	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal yuntongxun voice request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  headers,
		Body:     bodyData,
		BodyType: "json",
		Timeout:  30 * time.Second,
	}, t.handleYuntongxunResponse, nil
}

// handleYuntongxunResponse 处理云讯通API响应
func (t *yuntongxunTransformer) handleYuntongxunResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var result struct {
		StatusCode string `json:"statusCode"`
		StatusMsg  string `json:"statusMsg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse yuntongxun response: %w", err)
	}
	if result.StatusCode != "000000" {
		return &SMSError{
			Code:     result.StatusCode,
			Message:  result.StatusMsg,
			Provider: string(SubProviderYuntongxun),
		}
	}
	return nil
}

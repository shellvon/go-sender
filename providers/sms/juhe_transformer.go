package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
)

// @ProviderName: Juhe / 聚合数据
// @Website: https://www.juhe.cn
// @APIDoc: https://www.juhe.cn/docs/api/id/54
//
// 官方文档:
//   - 短信API文档: https://www.juhe.cn/docs/api/id/54
//   - 国内短信API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信API: https://www.juhe.cn/docs/api/id/357
//   - 视频短信API: https://www.juhe.cn/docs/api/id/363
//
// transformer 支持国内短信、国际短信、彩信/视频短信。

const (
	juheDefaultEndpoint     = "v.juhe.cn"
	juheDefaultIntlEndpoint = "v.juhe.cn"
)

type juheTransformer struct{}

func newJuheTransformer() core.HTTPTransformer[*core.Account] {
	return &juheTransformer{}
}

func init() {
	RegisterTransformer(string(SubProviderJuhe), &juheTransformer{})
}

func (t *juheTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == string(SubProviderJuhe)
}

func (t *juheTransformer) Transform(ctx context.Context, msg core.Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Juhe: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}

	if smsMsg.Type == MMS {
		return t.transformMMSSMS(ctx, smsMsg, account)
	}
	if smsMsg.IsIntl() {
		return t.transformIntlSMS(ctx, smsMsg, account)
	}
	return t.transformDomesticSMS(ctx, smsMsg, account)
}

func (t *juheTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("at least one mobile number is required")
	}
	return nil
}

// transformDomesticSMS 处理国内短信
//   - API: https://www.juhe.cn/docs/api/id/54
func (t *juheTransformer) transformDomesticSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("tpl_id", msg.TemplateID)
	params.Set("tpl_value", t.buildTemplateValue(msg.TemplateParams))
	params.Set("key", account.Key)
	if extend := msg.GetExtraStringOrDefault("ext", ""); extend != "" {
		params.Set("ext", extend)
	}
	body := []byte(params.Encode())
	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = juheDefaultEndpoint
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      "http://" + endpoint + "/sms/send",
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: "form",
		Timeout:  30 * time.Second,
	}
	return reqSpec, t.handleJuheResponse, nil
}

// transformIntlSMS 处理国际短信
//   - API: https://www.juhe.cn/docs/api/id/357
func (t *juheTransformer) transformIntlSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("areaNum", fmt.Sprintf("%d", msg.RegionCode))
	params.Set("tpl_id", msg.TemplateID)
	params.Set("tpl_value", t.buildTemplateValue(msg.TemplateParams))
	params.Set("key", account.Key)
	body := []byte(params.Encode())
	endpoint := account.IntlEndpoint
	if endpoint == "" {
		endpoint = juheDefaultIntlEndpoint
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      "http://" + endpoint + "/smsInternational/send",
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: "form",
		Timeout:  30 * time.Second,
	}
	return reqSpec, t.handleJuheResponse, nil
}

// transformMMSSMS 处理彩信/视频短信
//   - API: https://www.juhe.cn/docs/api/id/363
func (t *juheTransformer) transformMMSSMS(ctx context.Context, msg *Message, account *core.Account) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", strings.Join(msg.Mobiles, ","))
	params.Set("tpl_id", msg.TemplateID)
	params.Set("key", account.Key)

	body := []byte(params.Encode())
	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = juheDefaultEndpoint
	}
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      "http://" + endpoint + "/caixinv2/send",
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: "form",
		Timeout:  30 * time.Second,
	}
	return reqSpec, t.handleJuheResponse, nil
}

func (t *juheTransformer) handleJuheResponse(statusCode int, body []byte) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}
	var result struct {
		ErrorCode int    `json:"error_code"`
		Reason    string `json:"reason"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse juhe response: %w", err)
	}
	if result.ErrorCode != 0 {
		return &SMSError{
			Code:     fmt.Sprintf("%d", result.ErrorCode),
			Message:  result.Reason,
			Provider: string(SubProviderJuhe),
		}
	}
	return nil
}

// 如果您的模板里面有变量则需要提交此参数,如:#code#=123456,参数需要urlencode
// -> https://www.juhe.cn/docs/api/id/357
func (t *juheTransformer) buildTemplateValue(params map[string]string) string {
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

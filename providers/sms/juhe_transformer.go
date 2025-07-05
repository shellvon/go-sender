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
	juheDefaultBaseURI = "https//v.juhe.cn"

	// API路径常量.
	juheDomesticAPIPath = "/sms/send"
	juheIntlAPIPath     = "/smsInternational/send"
	juheMMSAPIPath      = "/caixinv2/send"
)

type juheTransformer struct{}

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

func (t *juheTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for juheTransformer")
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}

	switch smsMsg.Type {
	case SMSText:
		return t.transformSMS(smsMsg, account)
	case Voice:
		return nil, nil, errors.New("Juhe does not support voice messages")
	case MMS:
		return t.transformMMS(smsMsg, account)
	default:
		return nil, nil, fmt.Errorf("unsupported message type: %v", smsMsg.Type)
	}
}

func (t *juheTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("at least one mobile number is required")
	}
	return nil
}

// transformSMS 处理短信（国内/国际）
//   - 国内短信 API: https://www.juhe.cn/docs/api/id/54
//   - 国际短信 API: https://www.juhe.cn/docs/api/id/357
func (t *juheTransformer) transformSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("tpl_id", msg.TemplateID)
	params.Set("tpl_value", t.buildTemplateValue(msg.TemplateParams))
	params.Set("key", account.APIKey)

	// 根据是否为国际短信设置不同的参数和URL
	var apiPath string
	if msg.IsIntl() {
		// 国际短信特有参数
		params.Set("areaNum", strconv.Itoa(msg.RegionCode))
		apiPath = juheIntlAPIPath
	} else {
		// 国内短信特有参数
		params.Set("ext", msg.Extend)
		apiPath = juheDomesticAPIPath
	}

	body := []byte(params.Encode())
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      juheDefaultBaseURI + apiPath,
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: core.BodyTypeRaw,
	}
	return reqSpec, t.handleJuheResponse, nil
}

// transformMMS 处理彩信/视频短信
//   - API: https://www.juhe.cn/docs/api/id/363
func (t *juheTransformer) transformMMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", strings.Join(msg.Mobiles, ","))
	params.Set("tpl_id", msg.TemplateID)
	params.Set("key", account.APIKey)

	body := []byte(params.Encode())
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      juheDefaultBaseURI + juheMMSAPIPath,
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: core.BodyTypeForm,
	}
	return reqSpec, t.handleJuheResponse, nil
}

func (t *juheTransformer) handleJuheResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
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
		return &Error{
			Code:     strconv.Itoa(result.ErrorCode),
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

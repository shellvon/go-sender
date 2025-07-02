package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: UCP / 云之讯
// @Website: https://www.ucpaas.com
// @APIDoc: http://docs.ucpaas.com
//
// 官方文档:
//   - 短信API文档: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
//   - 国内外指定模版单发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms
//   - 国内外指定模版群发: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:sendsms_batch
//
// 能力说明:
//   - 国内短信：支持单发和群发，需模板ID。
//   - 国际短信：支持单发和群发，需模板ID。
//   - 彩信/语音：暂不支持。
//
// 注意：支持国内外手机号码，需模板ID。
type ucpTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderUcp), &ucpTransformer{})
}

func (t *ucpTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderUcp)
}

func (t *ucpTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for ucpTransformer")
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}
	return t.transformTextSMS(ctx, smsMsg, account)
}

func (t *ucpTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.TemplateID == "" {
		return errors.New("templateID is required")
	}
	return nil
}

// transformTextSMS 构造 UCP 短信 HTTP 请求
//   - 短信API: http://docs.ucpaas.com/doku.php?id=%E7%9F%AD%E4%BF%A1:about_sms
func (t *ucpTransformer) transformTextSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	// 根据手机号数量选择API
	var apiPath string
	if len(msg.Mobiles) > 1 {
		apiPath = "templatesms"
	} else {
		apiPath = "variablesms"
	}

	params := map[string]interface{}{
		"clientid":   account.Key,
		"password":   account.Secret,
		"templateid": msg.TemplateID,
		"mobile":     strings.Join(msg.Mobiles, ","),
		"uid":        msg.UID,
	}

	// 模板参数处理
	if len(msg.ParamsOrder) > 0 {
		// 模板中的替换参数，如该模板不存在参数则无需传该参数或者参数为空，如果有多个参数则需要写在同一个字符串中，以分号分隔 （如："a;b;c"），参数中不能含有特殊符号"【】"和","
		params["param"] = strings.Join(msg.ParamsOrder, ";")
	}

	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = "open.ucpaas.com"
	}
	if msg.IsIntl() && account.IntlEndpoint != "" {
		endpoint = account.IntlEndpoint
	}
	url := "https://" + endpoint + "/sms-server/" + apiPath

	bodyData, err := json.Marshal(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal ucp request body: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     bodyData,
		BodyType: "json",
	}, t.handleUcpResponse, nil
}

// handleUcpResponse 处理云之讯API响应.
func (t *ucpTransformer) handleUcpResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse ucp response: %w", err)
	}
	if result.Code != 0 {
		return &Error{
			Code:     strconv.Itoa(result.Code),
			Message:  result.Msg,
			Provider: string(SubProviderUcp),
		}
	}
	return nil
}

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

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6348/70146
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6348/70146
//
// transformer 仅支持 text（普通短信）类型。

type volcTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderVolc), &volcTransformer{})
}

func (t *volcTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderVolc)
}

func (t *volcTransformer) Transform(
	ctx context.Context,
	msg core.Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for volcTransformer")
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}
	return t.transformTextSMS(ctx, smsMsg, account)
}

func (t *volcTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}
	if msg.Content == "" {
		return errors.New("content is required")
	}
	if msg.SignName == "" {
		return errors.New("sign name is required")
	}
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(SubProviderVolc), "sendSMS")
	}
	return nil
}

func (t *volcTransformer) transformTextSMS(
	_ context.Context,
	msg *Message,
	account *core.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	body := map[string]interface{}{
		"SmsAccount":   account.Key,
		"Sign":         msg.SignName,
		"TemplateID":   msg.TemplateID,
		"PhoneNumbers": strings.Join(msg.Mobiles, ","),
	}
	if len(msg.TemplateParams) > 0 {
		body["TemplateParam"] = utils.ToJSONString(msg.TemplateParams)
	}
	if tag, ok := msg.GetExtraString(volcTagKey); ok && tag != "" {
		body[volcTagKey] = tag
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal volc request body: %w", err)
	}

	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = "sms.volcengineapi.com"
	}
	url := "https://" + endpoint + "/?Action=SendSms&Version=2020-01-01"

	headers := t.buildVolcHeaders(account, bodyJSON)

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      url,
		Headers:  headers,
		Body:     bodyJSON,
		BodyType: "json",
	}, t.handleVolcResponse, nil
}

// buildVolcHeaders 构建火山引擎TOP网关签名头
// 签名文档: https://www.volcengine.com/docs/6361/1205061
func (t *volcTransformer) buildVolcHeaders(account *core.Account, body []byte) map[string]string {
	ak := account.Key
	sk := account.Secret
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	xDate := timestamp[:10]

	endpoint := account.Endpoint
	if endpoint == "" {
		endpoint = "sms.volcengineapi.com"
	}

	headers := map[string]string{
		"Content-Type": "application/json;charset=utf-8",
		"Host":         endpoint,
		"X-Date":       timestamp,
	}
	canonicalHeaders := "content-type:application/json;charset=utf-8\nhost:" + endpoint + "\nx-date:" + timestamp + "\n"
	canonicalRequest := "POST\n/\n\n" + canonicalHeaders + "\n" + string(body)
	stringToSign := "HMAC-SHA256\n" + xDate + "\n" + utils.SHA256Hex([]byte(canonicalRequest))
	signingKey := utils.HMACSHA256([]byte(sk), []byte(xDate))
	signature := utils.HMACSHA256(signingKey, []byte(stringToSign))
	signatureBase64 := utils.Base64EncodeBytes(signature)
	authHeader := "HMAC-SHA256 Credential=" + ak + ", SignedHeaders=content-type;host;x-date, Signature=" + signatureBase64
	headers["Authorization"] = authHeader
	return headers
}

// handleVolcResponse 处理火山引擎API响应.
func (t *volcTransformer) handleVolcResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var result struct {
		ResponseMetadata struct {
			RequestID string `json:"RequestId"`
			Action    string `json:"Action"`
			Version   string `json:"Version"`
			Service   string `json:"Service"`
			Region    string `json:"Region"`
			Error     *struct {
				CodeN   int    `json:"CodeN"`
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error,omitempty"`
		} `json:"ResponseMetadata"`
		Result struct {
			MessageID string `json:"MessageID"`
		} `json:"Result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse volc response: %w", err)
	}
	if result.ResponseMetadata.Error != nil {
		return &Error{
			Code:     result.ResponseMetadata.Error.Code,
			Message:  result.ResponseMetadata.Error.Message,
			Provider: string(SubProviderVolc),
		}
	}
	return nil
}

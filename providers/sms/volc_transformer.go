package sms

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Volc / 火山引擎
// @Website: https://www.volcengine.com
// @APIDoc: https://www.volcengine.com/docs/6361/67380
//
// 官方文档:
//   - 短信API: https://www.volcengine.com/docs/6361/67380
//
// transformer 仅支持 text（普通短信）类型。
const (
	volcDefaultSmsEndpoint = "sms.volcengineapi.com"
	volcDefaultSmsAction   = "SendSms"
	volcDefaultSmsVersion  = "2020-01-01"
	volcDefaultSmsRegion   = "cn-beijing"
	volcDefaultSmsService  = "volcSMS"
)

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
	account *Account,
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
	if msg.SignName == "" {
		return errors.New("sign name is required")
	}
	if msg.GetExtraStringOrDefault(volcSmsAccountKey, "") == "" {
		return errors.New("sms account is required")
	}
	if msg.IsIntl() {
		return NewUnsupportedInternationalError(string(SubProviderVolc), "sendSMS")
	}
	return nil
}

func (t *volcTransformer) transformTextSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	body := map[string]interface{}{
		"SmsAccount":   msg.GetExtraStringOrDefault(volcSmsAccountKey, ""),
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

	// 只维护一份qs
	qs := url.Values{}
	qs.Set("Action", volcDefaultSmsAction)
	qs.Set("Version", volcDefaultSmsVersion)

	url := fmt.Sprintf("https://%s/?%s", volcDefaultSmsEndpoint, qs.Encode())
	headers := t.buildVolcHeaders(account, bodyJSON, qs)

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      url,
		Headers:  headers,
		Body:     bodyJSON,
		BodyType: core.BodyTypeRaw,
	}, t.handleVolcResponse, nil
}

// buildVolcHeaders constructs VolcEngine signature headers.
func (t *volcTransformer) buildVolcHeaders(account *Account, body []byte, qs url.Values) map[string]string {
	// 1. Basic parameters
	xDate := time.Now().UTC().Format("20060102T150405Z")
	authDate := xDate[:8]
	service := volcDefaultSmsService
	region := volcDefaultSmsRegion
	path := "/"
	host := volcDefaultSmsEndpoint
	contentType := "application/json"

	// 2. Payload: hex-encoded SHA256
	payload := utils.SHA256Hex(body)

	// 3. Canonicalize query string
	canonicalQueryString := t.normalizeQueryString(qs.Encode())

	// 4. Signed headers in fixed order
	signedHeaders := []string{"host", "x-date", "x-content-sha256", "content-type"}

	// 5. Canonical headers in the same order as signedHeaders
	headerMap := map[string]string{
		"host":             host,
		"x-date":           xDate,
		"x-content-sha256": payload,
		"content-type":     contentType,
	}
	var canonicalHeaders []string
	for _, k := range signedHeaders {
		canonicalHeaders = append(canonicalHeaders, k+":"+headerMap[k])
	}
	headerString := strings.Join(canonicalHeaders, "\n") + "\n"

	// 6. Build canonical request
	canonicalRequest := strings.Join([]string{
		http.MethodPost,
		path,
		canonicalQueryString,
		headerString,
		strings.Join(signedHeaders, ";"),
		payload,
	}, "\n")

	// 7. Build string to sign
	credentialScope := authDate + "/" + region + "/" + service + "/request"
	stringToSign := strings.Join([]string{
		"HMAC-SHA256",
		xDate,
		credentialScope,
		utils.SHA256Hex([]byte(canonicalRequest)),
	}, "\n")

	// 8. Derive signing key
	kDate := utils.HMACSHA256([]byte(account.APISecret), []byte(authDate))
	kRegion := utils.HMACSHA256(kDate, []byte(region))
	kService := utils.HMACSHA256(kRegion, []byte(service))
	kSigning := utils.HMACSHA256(kService, []byte("request"))
	signature := utils.HMACSHA256(kSigning, []byte(stringToSign))
	signatureHex := hex.EncodeToString(signature)

	// 9. Build Authorization header
	authHeader := "HMAC-SHA256 Credential=" + account.APIKey + "/" + credentialScope +
		", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signatureHex

	// 10. Return headers with capitalized keys
	return map[string]string{
		"Content-Type":     contentType,
		"Host":             host,
		"X-Date":           xDate,
		"X-Content-Sha256": payload,
		"Authorization":    authHeader,
	}
}

// normalizeQueryString canonicalizes the query string: sorts, URL-encodes, and replaces spaces with %20.
func (t *volcTransformer) normalizeQueryString(queryString string) string {
	if queryString == "" {
		return ""
	}
	values, _ := url.ParseQuery(queryString)
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var canonicalQS []string
	for _, k := range keys {
		for _, v := range values[k] {
			// url.QueryEscape replaces spaces with '+', so we replace '+' with '%20'
			keyEsc := strings.ReplaceAll(url.QueryEscape(k), "+", "%20")
			valEsc := strings.ReplaceAll(url.QueryEscape(v), "+", "%20")
			canonicalQS = append(canonicalQS, keyEsc+"="+valEsc)
		}
	}
	return strings.Join(canonicalQS, "&")
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
			// Array of String
			MessageID []string `json:"MessageID"`
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

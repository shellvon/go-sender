package sms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// volcTransformer implements HTTPRequestTransformer for VolcEngine SMS.
// It supports sending text message.
//
// Reference:
//   - Official Website: https://www.volcengine.com/
//   - API Docs: https://www.volcengine.com/docs/6361/67380
//   - SMS API: https://www.volcengine.com/docs/6361/67380

const (
	volcDefaultSmsEndpoint = "sms.volcengineapi.com"
	volcDefaultSmsAction   = "SendSms"
	volcDefaultSmsVersion  = "2020-01-01"
	volcDefaultSmsRegion   = "cn-beijing"
	volcDefaultSmsService  = "volcSMS"
)

type volcTransformer struct {
	*BaseTransformer
}

func init() {
	RegisterTransformer(string(SubProviderVolc), newVolcTransformer())
}

func newVolcTransformer() *volcTransformer {
	transformer := &volcTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(SubProviderVolc),
		&core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: true,
			Path:      "ResponseMetadata.Error.Code",
			CodePath:  "ResponseMetadata.Error.Code",
			MsgPath:   "ResponseMetadata.Error.Message",
			Expect:    "",
			Mode:      core.MatchEq,
		},
		nil,
		WithSMSHandler(transformer.transformSMS),
	)
	return transformer
}

// transformSMS transforms SMS message to HTTP request
//   - 短信API: https://www.volcengine.com/docs/6361/67380
func (t *volcTransformer) transformSMS(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	// 验证参数
	if len(msg.Mobiles) == 0 {
		return nil, nil, NewProviderError(string(SubProviderVolc), "MISSING_PARAM", "mobiles is required")
	}
	if msg.SignName == "" {
		return nil, nil, NewProviderError(string(SubProviderVolc), "MISSING_PARAM", "sign name is required")
	}
	if msg.IsIntl() {
		return nil, nil, NewUnsupportedInternationalError(string(SubProviderVolc), "sendSMS")
	}
	body := map[string]interface{}{
		"SmsAccount":   account.AppID,
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
		return nil, nil, NewProviderError(
			string(SubProviderVolc),
			"JSON_MARSHAL_ERROR",
			fmt.Sprintf("failed to marshal volc request body: %v", err),
		)
	}

	// 只维护一份qs
	qs := url.Values{}
	qs.Set("Action", volcDefaultSmsAction)
	qs.Set("Version", volcDefaultSmsVersion)

	url := fmt.Sprintf("https://%s/?%s", volcDefaultSmsEndpoint, qs.Encode())
	headers := t.buildVolcHeaders(account, bodyJSON, qs)

	return &core.HTTPRequestSpec{
		Method:  http.MethodPost,
		URL:     url,
		Headers: headers,
		Body:    bodyJSON,
	}, nil, nil
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
	payload := utils.HashHex(sha256.New, body)

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
		utils.HashHex(sha256.New, []byte(canonicalRequest)),
	}, "\n")

	// 8. Derive signing key
	kDate := utils.HMACSum(sha256.New, []byte(account.APISecret), []byte(authDate))
	kRegion := utils.HMACSum(sha256.New, kDate, []byte(region))
	kService := utils.HMACSum(sha256.New, kRegion, []byte(service))
	kSigning := utils.HMACSum(sha256.New, kService, []byte("request"))
	signature := utils.HMACSum(sha256.New, kSigning, []byte(stringToSign))
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

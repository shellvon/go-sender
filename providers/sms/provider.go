package sms

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// Provider supports multiple SMS service providers and strategy selection
type Provider struct {
	providers []*SMSProvider
	selector  *utils.Selector[*SMSProvider]
}

var _ core.Provider = (*Provider)(nil)

// New creates a new SMS provider instance
func New(config Config) (*Provider, error) {
	if !config.IsConfigured() {
		return nil, errors.New("SMS provider is not configured or is disabled")
	}

	// Convert to pointer slice
	providers := make([]*SMSProvider, len(config.Providers))
	for i := range config.Providers {
		providers[i] = &config.Providers[i]
	}

	// Use common initialization logic
	enabledProviders, selector, err := utils.InitProvider(&config, providers)
	if err != nil {
		return nil, errors.New("no enabled SMS providers found")
	}

	return &Provider{
		providers: enabledProviders,
		selector:  selector,
	}, nil
}

// Send sends an SMS message
func (p *Provider) Send(ctx context.Context, message core.Message) error {
	smsMsg, ok := message.(*Message)
	if !ok {
		return core.NewParamError(fmt.Sprintf("invalid message type: expected *sms.Message, got %T", message))
	}

	if err := smsMsg.Validate(); err != nil {
		return err
	}

	provider := p.selectProvider(ctx)
	if provider == nil {
		return errors.New("no available provider")
	}
	return p.doSendSMS(ctx, provider, smsMsg)
}

// selectProvider selects a provider based on context
func (p *Provider) selectProvider(ctx context.Context) *SMSProvider {
	return p.selector.Select(ctx)
}

func (p *Provider) doSendSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	switch provider.Type {
	case ProviderTypeTencent:
		return p.sendTencentSMS(ctx, provider, msg)
	case ProviderTypeAliyun:
		return p.sendAliyunSMS(ctx, provider, msg)
	case ProviderTypeYunpian:
		return p.sendYunpianSMS(ctx, provider, msg)
	case ProviderTypeUcp:
		return p.sendUcpSMS(ctx, provider, msg)
	case ProviderTypeCl253:
		return p.sendCl253SMS(ctx, provider, msg)
	case ProviderTypeSmsbao:
		return p.sendSmsbaoSMS(ctx, provider, msg)
	case ProviderTypeJuhe:
		return p.sendJuheSMS(ctx, provider, msg)
	case ProviderTypeLuosimao:
		return p.sendLuosimaoSMS(ctx, provider, msg)
	case ProviderTypeNetease:
		return p.sendNeteaseSMS(ctx, provider, msg)
	case ProviderTypeNormal:
		return p.sendNormalSMS(ctx, provider, msg)
	default:
		return fmt.Errorf("unsupported SMS provider type: %s", provider.Type)
	}
}

// Tencent SMS implementation
func (p *Provider) sendTencentSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	nonce := strconv.Itoa(int(time.Now().UnixNano()))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Build signature
	data := map[string]string{
		"appkey": provider.AppSecret,
		"random": nonce,
		"time":   timestamp,
		"mobile": msg.Mobile,
	}

	plainText := fmt.Sprintf("appkey=%s&random=%s&time=%s&mobile=%s",
		data["appkey"], data["random"], data["time"], data["mobile"])

	sig := p.tencentChecksum(plainText)

	// Build request body
	requestBody := map[string]interface{}{
		"ext":    "",
		"extend": "",
		"sig":    sig,
		"tel": map[string]string{
			"mobile":     msg.Mobile,
			"nationcode": "86",
		},
		"time": timestamp,
		"type": 0,
	}

	// Set message content
	if msg.TemplateCode != "" {
		requestBody["tpl_id"] = msg.TemplateCode
		requestBody["params"] = msg.TemplateParams
	} else {
		requestBody["msg"] = msg.Content
	}

	// Build URL
	apiURL := fmt.Sprintf("https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s",
		provider.AppID, nonce)

	// Send request
	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/json"},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// Aliyun SMS implementation
func (p *Provider) sendAliyunSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Build parameters
	params := map[string]string{
		"AccessKeyId":      provider.AppID,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     msg.Mobile,
		"RegionId":         "cn-hangzhou",
		"SignName":         msg.SignName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   nonce,
		"SignatureVersion": "1.0",
		"TemplateCode":     msg.TemplateCode,
		"TemplateParam":    p.mapToJSON(msg.TemplateParams),
		"Timestamp":        timestamp,
		"Version":          "2017-05-25",
	}

	// Build signature
	canonicalizedQueryString := p.aliyunCanonicalize(params)
	stringToSign := "POST&%2F&" + p.urlEncode(canonicalizedQueryString)
	signature := p.aliyunChecksum(stringToSign, provider.AppSecret)

	// Build request body
	queryString := p.aliyunStringify(params)
	requestBody := "Signature=" + signature + "&" + queryString

	// Send request
	_, _, err := utils.DoRequest(ctx, "https://dysmsapi.aliyuncs.com", utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:        []byte(requestBody),
		ContentType: "application/x-www-form-urlencoded",
	})

	return err
}

// Yunpian SMS implementation
func (p *Provider) sendYunpianSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		// For template SMS, content should be the template content with parameters
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	requestBody := map[string]string{
		"apikey": provider.AppSecret,
		"mobile": msg.Mobile,
		"text":   content,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, "https://sms.yunpian.com/v2/sms/single_send.json", utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/json"},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// UCP SMS implementation
func (p *Provider) sendUcpSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	requestBody := map[string]interface{}{
		"sid":        provider.AppID,
		"token":      provider.AppSecret,
		"appid":      provider.ExtraConfig["app_id"],
		"mobile":     msg.Mobile,
		"templateid": msg.TemplateCode,
		"param":      strings.Join(p.mapValues(msg.TemplateParams), ","),
	}

	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, "https://open.ucpaas.com/ol/sms/sendsms", utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/json"},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// CL253 SMS implementation
func (p *Provider) sendCl253SMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	requestBody := map[string]string{
		"account":  provider.AppID,
		"password": provider.AppSecret,
		"phone":    msg.Mobile,
		"msg":      content,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, "http://smssh1.253.com/msg/send/json", utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/json"},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// Smsbao SMS implementation
func (p *Provider) sendSmsbaoSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	// MD5 password
	passwordMD5 := p.md5Hash(provider.AppSecret)

	// Build query parameters
	params := url.Values{}
	params.Set("u", provider.AppID)
	params.Set("p", passwordMD5)
	params.Set("m", msg.Mobile)
	params.Set("c", content)

	// Send GET request
	apiURL := "http://api.smsbao.com/sms?" + params.Encode()
	_, _, err := utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method: "GET",
	})

	return err
}

// Helper functions for signature calculation
func (p *Provider) tencentChecksum(plainText string) string {
	hash := sha256.Sum256([]byte(plainText))
	return hex.EncodeToString(hash[:])
}

func (p *Provider) aliyunChecksum(stringToSign, secretKey string) string {
	key := (secretKey + "&")
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(stringToSign))
	digest := h.Sum(nil)
	return p.urlEncode(base64.StdEncoding.EncodeToString(digest))
}

func (p *Provider) aliyunCanonicalize(params map[string]string) string {
	var pairs []string
	for k, v := range params {
		pair := p.urlEncode(k) + "=" + p.urlEncode(v)
		pairs = append(pairs, pair)
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

func (p *Provider) aliyunStringify(params map[string]string) string {
	var pairs []string
	for k, v := range params {
		pairs = append(pairs, k+"="+v)
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "&")
}

func (p *Provider) urlEncode(s string) string {
	return url.QueryEscape(s)
}

func (p *Provider) mapToJSON(m map[string]string) string {
	bytes, _ := json.Marshal(m)
	return string(bytes)
}

func (p *Provider) mapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// MD5 hash function
func (p *Provider) md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Juhe SMS implementation
func (p *Provider) sendJuheSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	requestBody := map[string]string{
		"mobile":    msg.Mobile,
		"tpl_id":    msg.TemplateCode,
		"tpl_value": p.buildTemplateValue(msg.TemplateParams),
		"key":       provider.AppSecret,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, "http://v.juhe.cn/sms/send", utils.RequestOptions{
		Method:      "POST",
		Headers:     map[string]string{"Content-Type": "application/json"},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// Luosimao SMS implementation
func (p *Provider) sendLuosimaoSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	requestBody := map[string]string{
		"mobile":  msg.Mobile,
		"message": content,
	}

	// Basic auth with API key
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("api:key-"+provider.AppSecret))

	bodyBytes, _ := json.Marshal(requestBody)
	_, _, err := utils.DoRequest(ctx, "http://sms-api.luosimao.com/v1/send.json", utils.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": authHeader,
		},
		Body:        bodyBytes,
		ContentType: "application/json",
	})

	return err
}

// Netease SMS implementation
func (p *Provider) sendNeteaseSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	curtime := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Build checksum: SHA1(app_secret + nonce + curtime)
	plainText := provider.AppSecret + nonce + curtime
	checksum := p.sha1Hash(plainText)

	requestBody := map[string]string{
		"mobile": msg.Mobile,
	}

	if msg.TemplateCode != "" {
		requestBody["templateid"] = msg.TemplateCode
		requestBody["params"] = p.buildTemplateValue(msg.TemplateParams)
	} else {
		requestBody["code"] = msg.Content
	}

	headers := map[string]string{
		"AppKey":       provider.AppID,
		"CurTime":      curtime,
		"Nonce":        nonce,
		"CheckSum":     checksum,
		"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
	}

	// Convert to form data
	formData := url.Values{}
	for k, v := range requestBody {
		formData.Set(k, v)
	}

	_, _, err := utils.DoRequest(ctx, "https://api.netease.im/sms/sendcode.action", utils.RequestOptions{
		Method:      "POST",
		Headers:     headers,
		Body:        []byte(formData.Encode()),
		ContentType: "application/x-www-form-urlencoded",
	})

	return err
}

// Normal SMS implementation (generic)
func (p *Provider) sendNormalSMS(ctx context.Context, provider *SMSProvider, msg *Message) error {
	content := msg.Content
	if msg.TemplateCode != "" {
		content = msg.TemplateCode
		for key, value := range msg.TemplateParams {
			content = strings.ReplaceAll(content, "#"+key+"#", value)
		}
	}

	// Build request body with auth and message data
	requestBody := map[string]string{
		"mobile":  msg.Mobile,
		"content": content,
	}

	// Add auth credentials
	if provider.AppID != "" {
		requestBody["app_id"] = provider.AppID
	}
	if provider.AppSecret != "" {
		requestBody["app_secret"] = provider.AppSecret
	}

	// Add extra config as additional parameters
	for k, v := range provider.ExtraConfig {
		requestBody[k] = v
	}

	// Replace placeholders
	for k, v := range requestBody {
		v = strings.ReplaceAll(v, "{{mobile}}", msg.Mobile)
		v = strings.ReplaceAll(v, "{{content}}", content)
		requestBody[k] = v
	}

	// Determine HTTP method and URL from extra config
	method := "POST"
	if m, ok := provider.ExtraConfig["method"]; ok {
		method = m
	}

	apiURL := "http://api.example.com/sms" // Default URL
	if u, ok := provider.ExtraConfig["url"]; ok {
		apiURL = u
	}

	// Determine content type
	contentType := "application/json"
	if ct, ok := provider.ExtraConfig["content_type"]; ok {
		contentType = ct
	}

	var bodyBytes []byte
	var err error

	if contentType == "application/json" {
		bodyBytes, err = json.Marshal(requestBody)
	} else {
		// Form data
		formData := url.Values{}
		for k, v := range requestBody {
			formData.Set(k, v)
		}
		bodyBytes = []byte(formData.Encode())
	}

	if err != nil {
		return err
	}

	_, _, err = utils.DoRequest(ctx, apiURL, utils.RequestOptions{
		Method:      method,
		Headers:     map[string]string{"Content-Type": contentType},
		Body:        bodyBytes,
		ContentType: contentType,
	})

	return err
}

// Helper function to build template value string
func (p *Provider) buildTemplateValue(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	var pairs []string
	for k, v := range params {
		pairs = append(pairs, k+"="+v)
	}
	return strings.Join(pairs, ",")
}

// SHA1 hash function
func (p *Provider) sha1Hash(text string) string {
	hash := sha1.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (p *Provider) Name() string {
	return string(core.ProviderTypeSMS)
}

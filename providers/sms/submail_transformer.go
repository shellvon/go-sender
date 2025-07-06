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
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Submail / 赛邮
// @Website: https://www.mysubmail.com
// @APIDoc: https://www.mysubmail.com/documents
//
// 官方文档:
//   - 国内短信: https://www.mysubmail.com/documents/FppOR3
//   - 国际短信: https://www.mysubmail.com/documents/3UQA3
//   - 模板短信: https://www.mysubmail.com/documents/OOVyh
//   - 群发: https://www.mysubmail.com/documents/AzD4Z4
//   - 语音: https://www.mysubmail.com/documents/meE3C1
//   - 彩信: https://www.mysubmail.com/documents/N6ktR
//
// transformer 支持 text（国内/国际，模板/非模板，单发/群发）、voice（语音）、mms（彩信）类型。

// API endpoint paths.
const (
	// 国际短信-模版单发.
	intlTemplateSingle = "/internationalsms/xsend" // https://www.mysubmail.com/documents/87QTB2
	// 国际短信-批量群发.
	intlBatch = "internationalsms/batchsend" // https://www.mysubmail.com/documents/yD46O
	// 国际短信-单发.
	intlSingle = "/internationalsms/send" // https://www.mysubmail.com/documents/3UQA3
	// 国际短信-模版一对多(没有找到批量的API).
	intlTemplateBatch = "internationalsms/multixsend" // https://www.mysubmail.com/documents/B70hy
	// 国内短信-模版单发.
	domesticTemplateSingle = "/sms/xsend" // https://www.mysubmail.com/documents/OOVyh
	// 国内短信-模版-群发.
	domesticTemplateBatch = "/sms/multixsend" // https://www.mysubmail.com/documents/G5KBR
	// 国内短信-单发.
	domesticSingle = "/sms/send" // https://www.mysubmail.com/documents/FppOR3
	// 国内短信-群发.
	domesticBatch = "/sms/multisend" // https://www.mysubmail.com/documents/AzD4Z4
	// 彩信-单发.
	mmsSingle = "/mms/send" // https://www.mysubmail.com/documents/N6ktR
	// 语音-单发.
	voiceSingle = "/voice/send" // https://www.mysubmail.com/documents/meE3C1
	// 语音-模版单发.
	voiceTemplateSingle = "/voice/xsend" // https://www.mysubmail.com/documents/KbG03
	// 语音-模版群发.
	voiceTemplateBatch   = "/voice/multixsend" // https://www.mysubmail.com/documents/FkgkM2
	submailDefaultDomain = "https://api-v4.mysubmail.com"
)

type submailTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderSubmail), &submailTransformer{})
}

func (t *submailTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	return ok && smsMsg.SubProvider == string(SubProviderSubmail)
}

func (t *submailTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, fmt.Errorf("unsupported message type for Submail: %T", msg)
	}

	// Apply Submail-specific defaults
	t.applySubmailDefaults(smsMsg, account)

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

// applySubmailDefaults applies Submail-specific defaults to the message.
func (t *submailTransformer) applySubmailDefaults(msg *Message, account *Account) {
	// Apply common defaults first
	msg.ApplyCommonDefaults(account)
}

// 通用的submail请求构造方法.
func (t *submailTransformer) buildSubmailRequest(
	msg *Message,
	account *Account,
	apiPath string,
	buildParams func(*Message, *Account) map[string]string,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := buildParams(msg, account)
	body := t.encodeParams(params)
	endpoint := fmt.Sprintf("%s%s", submailDefaultDomain, apiPath)
	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      endpoint,
		Body:     body,
		BodyType: core.BodyTypeForm,
	}, t.handleSubmailResponse, nil
}

// transformSMS transforms SMS message to HTTP request
//   - 国内短信: https://www.mysubmail.com/documents/FppOR3
//   - 国际短信: https://www.mysubmail.com/documents/3UQA3
//   - 模板短信: https://www.mysubmail.com/documents/OOVyh
//   - 群发: https://www.mysubmail.com/documents/AzD4Z4
func (t *submailTransformer) transformSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.Content == "" && msg.TemplateID == "" {
		return nil, nil, errors.New("content or templateID is required")
	}
	if msg.IsIntl() && len(msg.Mobiles) > 1000 {
		return nil, nil, errors.New("international sms: at most 1000 mobiles per request")
	}
	if !msg.IsIntl() && len(msg.Mobiles) > 10000 {
		return nil, nil, errors.New("domestic sms: at most 10000 mobiles per request")
	}
	apiPath := t.getSMSPath(msg)
	return t.buildSubmailRequest(msg, account, apiPath, t.buildSMSParams)
}

func (t *submailTransformer) getSMSPath(msg *Message) string {
	isIntl := msg.IsIntl()
	isTemplate := msg.TemplateID != ""
	isBatch := len(msg.Mobiles) > 1

	if isIntl {
		if isTemplate {
			if isBatch {
				return intlTemplateBatch
			}
			return intlTemplateSingle
		}
		if isBatch {
			return intlBatch
		}
		return intlSingle
	}

	// Domestic SMS
	if isTemplate {
		if isBatch {
			return domesticTemplateBatch
		}
		return domesticTemplateSingle
	}
	if isBatch {
		return domesticBatch
	}
	return domesticSingle
}

// transformVoice transforms voice message to HTTP request
//   - 语音: https://www.mysubmail.com/documents/meE3C1
//   - 语音模板: https://www.mysubmail.com/documents/KbG03
func (t *submailTransformer) transformVoice(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.Content == "" && msg.TemplateID == "" {
		return nil, nil, errors.New("voice content or templateID is required")
	}
	if len(msg.Mobiles) > 1 {
		return nil, nil, errors.New("voice only supports single send")
	}
	apiPath := t.getVoicePath(msg)
	return t.buildSubmailRequest(msg, account, apiPath, t.buildVoiceParams)
}

func (t *submailTransformer) getVoicePath(msg *Message) string {
	if msg.TemplateID == "" {
		return voiceSingle
	}
	if len(msg.Mobiles) > 1 {
		return voiceTemplateBatch
	}
	return voiceTemplateSingle
}

// transformMMS transforms MMS message to HTTP request
//   - 彩信: https://www.mysubmail.com/documents/N6ktR
func (t *submailTransformer) transformMMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if len(msg.Mobiles) == 0 {
		return nil, nil, errors.New("mobiles is required")
	}
	if msg.TemplateID == "" {
		return nil, nil, errors.New("mms requires templateID")
	}
	if len(msg.Mobiles) > 1 {
		return nil, nil, errors.New("mms only supports single send")
	}
	apiPath := mmsSingle
	return t.buildSubmailRequest(msg, account, apiPath, t.buildMMSParams)
}

func (t *submailTransformer) buildSMSParams(msg *Message, account *Account) map[string]string {
	params := map[string]string{
		"appid": account.APIKey,
	}

	// 添加接收者
	t.addRecipients(params, msg)

	// 添加内容或模板
	t.addContentOrTemplate(params, msg)

	// 添加通用参数
	t.addCommonParams(params, msg, account)

	return params
}

func (t *submailTransformer) buildVoiceParams(msg *Message, account *Account) map[string]string {
	params := map[string]string{
		"appid": account.APIKey,
	}

	// 添加接收者
	t.addRecipients(params, msg)

	// 添加内容或模板
	t.addContentOrTemplate(params, msg)

	// 添加通用参数
	t.addCommonParams(params, msg, account)

	return params
}

func (t *submailTransformer) buildMMSParams(msg *Message, account *Account) map[string]string {
	params := map[string]string{
		"appid": account.APIKey,
	}

	// 添加接收者
	t.addRecipients(params, msg)

	// 添加内容或模板
	t.addContentOrTemplate(params, msg)

	// 添加通用参数
	t.addCommonParams(params, msg, account)

	return params
}

func (t *submailTransformer) addRecipients(params map[string]string, msg *Message) {
	// 语音和彩信仅支持单发
	if msg.Type == Voice || msg.Type == MMS {
		params["to"] = msg.Mobiles[0]
		return
	}

	// 短信支持群发
	if msg.IsIntl() && msg.TemplateID != "" && len(msg.Mobiles) > 1 {
		// 国际批量模板发送使用multi参数
		t.addMultiRecipients(params, msg)
	} else {
		params["to"] = strings.Join(msg.Mobiles, ",")
	}
}

func (t *submailTransformer) addMultiRecipients(params map[string]string, msg *Message) {
	multi := make([]map[string]interface{}, 0, len(msg.Mobiles))
	for _, mobile := range msg.Mobiles {
		item := map[string]interface{}{
			"to": mobile,
		}
		if len(msg.TemplateParams) > 0 {
			item["vars"] = msg.TemplateParams
		}
		multi = append(multi, item)
	}
	multiJSON, _ := json.Marshal(multi)
	params["multi"] = string(multiJSON)
}

func (t *submailTransformer) addContentOrTemplate(params map[string]string, msg *Message) {
	params["project"] = msg.TemplateID
	params["vars"] = utils.ToJSONString(msg.TemplateParams)
	params["content"] = utils.AddSignature(msg.Content, msg.SignName)
	// 或者使用extras中的sender
	if sender := msg.GetExtraStringOrDefault("sender", ""); sender != "" {
		params["sender"] = sender
	}
}

func (t *submailTransformer) addCommonParams(params map[string]string, msg *Message, account *Account) {
	if tag := msg.GetExtraStringOrDefault(submailTagKey, ""); tag != "" {
		params[submailTagKey] = tag
	}
	if signType := msg.GetExtraStringOrDefault(submailSignTypeKey, ""); signType != "" {
		params[submailSignTypeKey] = signType
	}

	// 添加时间戳
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params["timestamp"] = timestamp

	// 添加签名
	params["signature"] = t.calculateSignature(account, params)
}

// calculateSignature 计算签名.
//   - https://www.mysubmail.com/documents/pdxzv1
func (t *submailTransformer) calculateSignature(account *Account, params map[string]string) string {
	// 获取签名类型，默认为md5
	signType := submailDefaultSignType
	// 或者从消息的extras中获取
	if msgSignType := params[submailTagKey]; msgSignType != "" {
		signType = msgSignType
	}

	// 构建签名字符串
	var keys []string
	for k := range params {
		if k != "signature" && k != submailSignTypeKey && k != "sign_version" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, k+"="+params[k])
	}
	stringToSign := strings.Join(pairs, "&")

	// 根据签名类型计算签名
	switch signType {
	case "sha1":
		return utils.SHA1Hex(stringToSign + account.APISecret)
	case "normal":
		return account.APISecret
	default: // md5
		return utils.MD5Hex(stringToSign + account.APISecret)
	}
}

func (t *submailTransformer) encodeParams(params map[string]string) []byte {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return []byte(values.Encode())
}

func (t *submailTransformer) handleSubmailResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}

	var result struct {
		Status string `json:"status"`  // Request status: success/error
		SendID string `json:"send_id"` // Unique send ID
		Fee    int    `json:"fee"`     // Billing count
		Code   string `json:"code"`    // Error code (when status is error)
		Msg    string `json:"msg"`     // Error message (when status is error)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse submail response: %w", err)
	}

	if result.Status != "success" {
		return &Error{
			Code:     result.Code,
			Message:  result.Msg,
			Provider: string(SubProviderSubmail),
		}
	}

	return nil
}

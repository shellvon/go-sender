package sms

import (
	"context"
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
		return nil, nil, fmt.Errorf("unsupported message type for submail: %T", msg)
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, err
	}

	endpoint := t.buildEndpoint(smsMsg, account)
	params := t.buildParams(smsMsg, account)
	body := t.encodeParams(params)

	return &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      endpoint,
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Body:     body,
		BodyType: "form",
	}, t.handleSubmailResponse, nil
}

func (t *submailTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("mobiles is required")
	}

	switch msg.Type {
	case SMSText:
		if msg.Content == "" && msg.TemplateID == "" {
			return errors.New("content or templateID is required")
		}
		if msg.IsIntl() && len(msg.Mobiles) > 1000 {
			return errors.New("international sms: at most 1000 mobiles per request")
		}
		if !msg.IsIntl() && len(msg.Mobiles) > 10000 {
			return errors.New("domestic sms: at most 10000 mobiles per request")
		}
	case Voice:
		if msg.Content == "" && msg.TemplateID == "" {
			return errors.New("voice content or templateID is required")
		}
		if len(msg.Mobiles) > 1 {
			return errors.New("voice only supports single send")
		}
	case MMS:
		if msg.TemplateID == "" {
			return errors.New("mms requires templateID")
		}
		if len(msg.Mobiles) > 1 {
			return errors.New("mms only supports single send")
		}
	default:
		return fmt.Errorf("unsupported submail message type: %s", msg.Type)
	}
	return nil
}

func (t *submailTransformer) buildEndpoint(msg *Message, _ *Account) string {
	var apiPath string

	switch msg.Type {
	case SMSText:
		apiPath = t.getSMSPath(msg)
	case Voice:
		apiPath = t.getVoicePath(msg)
	case MMS:
		apiPath = mmsSingle
	default:
		apiPath = "" // 或 panic/return error
	}
	return fmt.Sprintf("%s%s", smsbaoDefaultBaseURI, apiPath)
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

func (t *submailTransformer) getVoicePath(msg *Message) string {
	if msg.TemplateID == "" {
		return voiceSingle
	}
	if len(msg.Mobiles) > 1 {
		return voiceTemplateBatch
	}
	return voiceTemplateSingle
}

func (t *submailTransformer) buildParams(msg *Message, account *Account) map[string]string {
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

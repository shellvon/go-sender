package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// @ProviderName: Luosimao / 螺丝帽
// @Website: https://luosimao.com
// @APIDoc: https://luosimao.com/docs/api
//
// 官方文档:
//   - 短信API: https://luosimao.com/docs/api
//
// transformer 仅支持 text（普通短信）类型。

const (
	luosimaoSmsDefaultBaseURI   = "https://sms-api.luosimao.com"
	luosimaoVoiceDefaultBaseURI = "https://voice-api.luosimao.com"
)

type luosimaoTransformer struct{}

func init() {
	RegisterTransformer(string(SubProviderLuosimao), &luosimaoTransformer{})
}

func (t *luosimaoTransformer) CanTransform(msg core.Message) bool {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return false
	}
	return smsMsg.SubProvider == string(SubProviderLuosimao)
}

// Transform 构造螺丝帽短信/语音 HTTP 请求
//   - 短信单发: https://luosimao.com/docs/api#send
//   - 短信批量: https://luosimao.com/docs/api#send_batch
//   - 语音验证码: https://luosimao.com/docs/api/51
//
// 规则:
//   - type==Voice 且 category==CategoryVerification：走语音接口（仅支持单发）
//   - type==SMSText：手机号数量大于1走批量，否则单发
//   - endpoint 只用域名，transformer 拼接协议和路径
//
// 参数:
//   - ctx: 上下文
//   - msg: 消息体
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *luosimaoTransformer) Transform(
	_ context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, NewProviderError(string(SubProviderLuosimao), "INVALID_MESSAGE_TYPE", fmt.Sprintf("unsupported message type for Luosimao: %T", msg))
	}

	// Apply Luosimao-specific defaults
	t.applyLuosimaoDefaults(smsMsg, account)

	switch smsMsg.Type {
	case SMSText:
		return t.transformSMS(smsMsg, account)
	case Voice:
		return t.transformVoice(smsMsg, account)
	case MMS:
		return nil, nil, NewProviderError(string(SubProviderLuosimao), "UNSUPPORTED_MESSAGE_TYPE", fmt.Sprintf("unsupported message type: %v", smsMsg.Type))
	default:
		return nil, nil, NewProviderError(string(SubProviderLuosimao), "UNSUPPORTED_MESSAGE_TYPE", fmt.Sprintf("unsupported message type: %v", smsMsg.Type))
	}
}

// applyLuosimaoDefaults applies Luosimao-specific defaults to the message.
func (t *luosimaoTransformer) applyLuosimaoDefaults(msg *Message, account *Account) {
	// Apply common defaults first
	msg.ApplyCommonDefaults(account)
}

// buildLuosimaoRequestSpec 构造 Luosimao HTTPRequestSpec
//   - 负责拼接 body、设置 Authorization header、构造 HTTPRequestSpec 和响应处理器
//   - 由 transformSingleSMS/transformBatchSMS/transformVoiceSMS 统一调用
//
// 参数:
//   - params: url.Values 请求参数
//   - requestURL: 完整请求 URL
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *luosimaoTransformer) buildLuosimaoRequestSpec(
	params url.Values,
	requestURL string,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	body := []byte(params.Encode())
	authHeader := "Basic " + utils.Base64EncodeBytes([]byte("api:key-"+account.APISecret))
	reqSpec := &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      requestURL,
		Headers:  map[string]string{"Authorization": authHeader},
		Body:     body,
		BodyType: core.BodyTypeForm,
	}
	return reqSpec, t.handleLuosimaoResponse, nil
}

// transformSMS transforms SMS message to HTTP request
//   - 根据手机号数量决定走单发还是批量
//   - 单发API: https://luosimao.com/docs/api#send
//   - 批量API: https://luosimao.com/docs/api#send_batch
//
// 对于批量多发或者只有一个手机号但指定了定时发送的任务的，都采用批量发送，否则使用单个发送API.
func (t *luosimaoTransformer) transformSMS(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	// 对于手机号多余1个的或者只有一个手机号却有定时任务时，则采用批量发送API
	// 因为单个手机号的不支持定时任务
	isBatch := len(msg.Mobiles) > 1 || msg.ScheduledAt != nil
	if isBatch {
		// 批量发送
		params.Set("mobile_list", strings.Join(msg.Mobiles, ","))
		params.Set("message", utils.AddSignature(msg.Content, msg.SignName))
		if msg.ScheduledAt != nil {
			// 定时发送的时间，定时的发送任务可以在发送前10分钟在发送历史界面进行取消（仅限提交当天）, 格式为 YYYY-MM-DD HH:MM:SS
			params.Set(luosimaoScheduledAtKey, msg.ScheduledAt.Format(time.DateTime))
		}
		return t.buildLuosimaoRequestSpec(
			params,
			fmt.Sprintf("%s/v1/send_batch.json", luosimaoSmsDefaultBaseURI),
			account,
		)
	}
	// 单发
	params.Set("mobile", msg.Mobiles[0])
	params.Set("message", utils.AddSignature(msg.Content, msg.SignName))
	return t.buildLuosimaoRequestSpec(params, fmt.Sprintf("%s/v1/send.json", luosimaoSmsDefaultBaseURI), account)
}

// transformVoice transforms voice message to HTTP request
// 目前语音短信仅支持验证码，即检查category是否为CategoryVerification
//   - API: https://luosimao.com/docs/api/51
func (t *luosimaoTransformer) transformVoice(
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if msg.Category != CategoryVerification {
		return nil, nil, NewProviderError(string(SubProviderLuosimao), "UNSUPPORTED_CATEGORY", fmt.Sprintf("unsupported voice category: %v", msg.Category))
	}
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("code", msg.Content)
	return t.buildLuosimaoRequestSpec(
		params,
		fmt.Sprintf("%s/v1/verify.json", luosimaoVoiceDefaultBaseURI),
		account,
	)
}

// handleLuosimaoResponse 处理螺丝帽 API 响应
//   - 统一处理单发、批量、语音接口返回
func (t *luosimaoTransformer) handleLuosimaoResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return NewProviderError(string(SubProviderLuosimao), strconv.Itoa(statusCode), string(body))
	}
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return NewProviderError(string(SubProviderLuosimao), "PARSE_ERROR", err.Error())
	}
	if response["errorno"] != float64(0) {
		return NewProviderError(string(SubProviderLuosimao), strconv.Itoa(int(response["errorno"].(float64))), response["msg"].(string))
	}
	return nil
}

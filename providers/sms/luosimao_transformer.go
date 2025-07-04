package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	ctx context.Context,
	msg core.Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	smsMsg, ok := msg.(*Message)
	if !ok {
		return nil, nil, errors.New("invalid message type for luosimaoTransformer")
	}
	if err := t.validateMessage(smsMsg); err != nil {
		return nil, nil, fmt.Errorf("message validation failed: %w", err)
	}
	if smsMsg.Type == Voice && smsMsg.Category == CategoryVerification {
		return t.transformVoiceSMS(ctx, smsMsg, account)
	}
	if len(smsMsg.Mobiles) > 1 {
		return t.transformBatchSMS(ctx, smsMsg, account)
	}
	return t.transformSingleSMS(ctx, smsMsg, account)
}

func (t *luosimaoTransformer) validateMessage(msg *Message) error {
	if len(msg.Mobiles) == 0 {
		return errors.New("at least one mobile number is required")
	}
	if msg.Content == "" {
		return errors.New("content is required for Luosimao SMS")
	}
	return nil
}

// buildLuosimaoRequestSpec 构造 Luosimao HTTPRequestSpec
//   - 负责拼接 body、设置 Authorization header、构造 HTTPRequestSpec 和响应处理器
//   - 由 transformSingleSMS/transformBatchSMS/transformVoiceSMS 统一调用
//
// 参数:
//   - ctx: 上下文
//   - params: url.Values 请求参数
//   - requestURL: 完整请求 URL
//   - account: 账号配置
//
// 返回:
//   - HTTPRequestSpec: HTTP 请求规范
//   - ResponseHandler: 响应处理器
//   - error: 错误信息
func (t *luosimaoTransformer) buildLuosimaoRequestSpec(
	_ context.Context,
	params url.Values,
	requestURL string,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	body := []byte(params.Encode())
	authHeader := "Basic " + utils.Base64EncodeBytes([]byte("api:key-"+account.APISecret))
	reqSpec := &core.HTTPRequestSpec{
		Method:   "POST",
		URL:      requestURL,
		Headers:  map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Authorization": authHeader},
		Body:     body,
		BodyType: "form",
	}
	return reqSpec, t.handleLuosimaoResponse, nil
}

// transformSingleSMS 构造单发短信 HTTP 请求
//   - API: https://luosimao.com/docs/api#send
func (t *luosimaoTransformer) transformSingleSMS(
	ctx context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("message", utils.AddSignature(msg.Content, msg.SignName))
	return t.buildLuosimaoRequestSpec(ctx, params, fmt.Sprintf("%s/v1/send.json", luosimaoSmsDefaultBaseURI), account)
}

// transformBatchSMS 构造批量短信 HTTP 请求
//   - API: https://luosimao.com/docs/api#send_batch
//
// 对于批量短信，额外多一个定时发送的能力，可通过 ScheduledAt 字段设置
// 定时的发送任务可以在发送前10分钟在发送历史界面进行取消（仅限提交当天）.
func (t *luosimaoTransformer) transformBatchSMS(
	ctx context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile_list", strings.Join(msg.Mobiles, ","))
	params.Set("message", utils.AddSignature(msg.Content, msg.SignName))
	if msg.ScheduledAt != nil {
		// 定时发送的时间，定时的发送任务可以在发送前10分钟在发送历史界面进行取消（仅限提交当天）, 格式为 YYYY-MM-DD HH:MM:SS
		params.Set(luosimaoScheduledAtKey, msg.ScheduledAt.Format(time.DateTime))
	}
	return t.buildLuosimaoRequestSpec(
		ctx,
		params,
		fmt.Sprintf("%s/v1/send_batch.json", luosimaoSmsDefaultBaseURI),
		account,
	)
}

// transformVoiceSMS 构造语音验证码 HTTP 请求
//   - API: https://luosimao.com/docs/api/51
func (t *luosimaoTransformer) transformVoiceSMS(
	ctx context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	params := url.Values{}
	params.Set("mobile", msg.Mobiles[0])
	params.Set("code", msg.Content)
	return t.buildLuosimaoRequestSpec(
		ctx,
		params,
		fmt.Sprintf("%s/v1/verify.json", luosimaoVoiceDefaultBaseURI),
		account,
	)
}

// handleLuosimaoResponse 处理螺丝帽 API 响应
//   - 统一处理单发、批量、语音接口返回
func (t *luosimaoTransformer) handleLuosimaoResponse(statusCode int, body []byte) error {
	if !utils.IsAcceptableStatus(statusCode) {
		return fmt.Errorf("HTTP request failed with status %d: %s", statusCode, string(body))
	}
	var result struct {
		Error int    `json:"error"`
		Msg   string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse luosimao response: %w", err)
	}
	if result.Error != 0 {
		return &Error{
			Code:     strconv.Itoa(result.Error),
			Message:  result.Msg,
			Provider: string(SubProviderLuosimao),
		}
	}
	return nil
}

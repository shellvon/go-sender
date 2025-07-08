package sms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// luosimaoTransformer implements HTTPRequestTransformer for Luosimao SMS.
// It only supports sending text message.
//
// Reference:
//   - Official Website: https://luosimao.com
//   - API Docs: https://luosimao.com/docs/api
//   - SMS API: https://luosimao.com/docs/api#send
const (
	luosimaoSmsDefaultBaseURI   = "https://sms-api.luosimao.com"
	luosimaoVoiceDefaultBaseURI = "https://voice-api.luosimao.com"
)

type luosimaoTransformer struct {
	*BaseTransformer
}

func newLuosimaoTransformer() *luosimaoTransformer {
	transformer := &luosimaoTransformer{}
	transformer.BaseTransformer = NewBaseTransformer(
		string(SubProviderLuosimao),
		&core.ResponseHandlerConfig{
			BodyType:  core.BodyTypeJSON,
			CheckBody: true,
			Path:      "errorno",
			Expect:    "0",
			Mode:      core.MatchEq,
		},
		nil,
		WithSMSHandler(transformer.transformSMS),
		WithVoiceHandler(transformer.transformVoice),
	)
	return transformer
}

func init() {
	RegisterTransformer(string(SubProviderLuosimao), newLuosimaoTransformer())
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
	return reqSpec, nil, nil
}

// transformSMS transforms SMS message to HTTP request
//   - 根据手机号数量决定走单发还是批量
//   - 单发API: https://luosimao.com/docs/api#send
//   - 批量API: https://luosimao.com/docs/api#send_batch
//
// 对于批量多发或者只有一个手机号但指定了定时发送的任务的，都采用批量发送，否则使用单个发送API.
func (t *luosimaoTransformer) transformSMS(
	_ context.Context,
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
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	if msg.Category != CategoryVerification {
		return nil, nil, NewProviderError(
			string(SubProviderLuosimao),
			"UNSUPPORTED_CATEGORY",
			fmt.Sprintf("unsupported voice category: %v", msg.Category),
		)
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

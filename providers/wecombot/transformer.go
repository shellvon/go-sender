package wecombot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// WeComBot 是企业微信机器人的消息提供者。
// 它支持发送文本、图片、文件、Markdown、新闻等多种类型的消息。
// 它还支持向不同部门和用户发送消息。
//
// 参考：
//   - 官方网站：https://work.weixin.qq.com/
//   - API 文档：https://developer.work.weixin.qq.com/document/path/91770

const wecombotAPIURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

// wecombotTransformer 利用共享的 BaseHTTPTransformer 为企业微信机器人提供转换功能。
type wecombotTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]
}

// transform 为企业微信机器人消息构建 HTTPRequestSpec。
// 参数：
//   - _ context.Context - 请求的上下文（未使用）。
//   - msg Message - 要发送的消息。
//   - account *Account - 用于发送消息的账户信息。
//
// 返回值：
//   - *core.HTTPRequestSpec - HTTP 请求规格。
//   - core.SendResultHandler - 发送结果处理器（当前为 nil）。
//   - error - 如果转换失败，返回错误。
func (wt *wecombotTransformer) transform(
	_ context.Context,
	msg Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	if account == nil {
		return nil, nil, errors.New("no account provided")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal wecombot message: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method: http.MethodPost,
		URL:    wecombotAPIURL,
		QueryParams: url.Values{
			"key": {account.APIKey},
		},
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

// newWecombotTransformer 创建一个新的转换器实例。
// 返回值：core.HTTPTransformer[*Account] - 新创建的转换器实例，用于处理企业微信机器人消息。
func newWecombotTransformer() core.HTTPTransformer[*Account] {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "errcode",
		Expect:    "0",
		Mode:      core.MatchEq,
		CodePath:  "errcode",
		MsgPath:   "errmsg",
	}

	wt := &wecombotTransformer{}
	wt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeWecombot,
		"",
		respCfg,
		wt.transform,
	)

	return wt
}

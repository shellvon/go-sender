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

// WeComBot is a message provider for WeCom Bot.
// It supports sending text, image, file, markdown, news, and other types of messages.
// It also supports sending messages to different departments and users.
//
// Reference:
//   - Official Website: https://work.weixin.qq.com/
//   - API Docs: https://developer.work.weixin.qq.com/document/path/91770

const wecombotAPIURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

// wecombotTransformer leverages the shared BaseHTTPTransformer for WeCom Bot.
type wecombotTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]
}

// transform builds the HTTPRequestSpec for a WeCom Bot message.
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

// newWecombotTransformer creates a new transformer instance.
func newWecombotTransformer() core.HTTPTransformer[*Account] {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "errcode",
		Expect:    "0",
		Mode:      core.MatchEq,
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

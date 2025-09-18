package lark

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
	"github.com/shellvon/go-sender/utils"
)

// larkTransformer implements HTTPRequestTransformer for Lark.
// It supports sending text, image, file, markdown, news, and other types of messages.
//
// Reference:
//   - Official Website: https://open.feishu.cn
//   - API Docs: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot

const (
	larkWebhookURLTemplate = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"
)

type larkTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]
}

func (lt *larkTransformer) transform(
	_ context.Context,
	msg Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	qs := url.Values{}
	if account.APISecret != "" {
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		sign := utils.HMACBase64(sha256.New, []byte(account.APISecret), []byte(ts+"\n"+account.APISecret))
		qs.Add("timestamp", ts)
		qs.Add("sign", sign)
	}

	payload := map[string]interface{}{
		"msg_type": msg.GetMsgType(),
		"content":  msg,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal lark payload: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:      http.MethodPost,
		URL:         fmt.Sprintf(larkWebhookURLTemplate, account.APIKey),
		QueryParams: qs,
		Body:        body,
		BodyType:    core.BodyTypeJSON,
	}, nil, nil
}

func newLarkTransformer() *larkTransformer {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "code",
		Expect:    "0",
		Mode:      core.MatchEq,
		CodePath:  "code",
		MsgPath:   "msg",
	}

	lt := &larkTransformer{}
	lt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeLark,
		"",
		respCfg,
		lt.transform,
	)
	return lt
}

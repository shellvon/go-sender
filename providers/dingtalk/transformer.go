package dingtalk

import (
	"context"
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

// dingTalkTransformer implements HTTPRequestTransformer for DingTalk.
// It supports sending text, image, file, markdown, news, and other types of messages.
//
// Reference:
//   - Official Website: https://dingtalk.com/
//   - API Docs: https://open.dingtalk.com/document/robots/custom-robot-access

const dingtalkRobotURL = "https://oapi.dingtalk.com/robot/send"

type dingTalkTransformer struct {
	*transformer.BaseHTTPTransformer[Message, *Account]
}

func (dt *dingTalkTransformer) transform(
	_ context.Context,
	msg Message,
	account *Account,
) (*core.HTTPRequestSpec, core.SendResultHandler, error) {
	qs := url.Values{}
	qs.Add("access_token", account.APIKey)

	if account.APISecret != "" {
		ts := time.Now().UnixMilli()
		stringToSign := fmt.Sprintf("%d\n%s", ts, account.APISecret)
		sign := utils.HMACSHA256Base64(account.APISecret, stringToSign)
		qs.Add("timestamp", strconv.FormatInt(ts, 10))
		qs.Add("sign", sign)
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:      http.MethodPost,
		URL:         dingtalkRobotURL,
		QueryParams: qs,
		Body:        body,
		BodyType:    core.BodyTypeJSON,
	}, nil, nil
}

func newDingTalkTransformer() *dingTalkTransformer {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "errcode",
		Expect:    "0",
		Mode:      core.MatchEq,
		CodePath:  "errcode",
		MsgPath:   "errmsg",
	}

	dt := &dingTalkTransformer{}
	dt.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeDingtalk,
		"",
		respCfg,
		dt.transform,
	)
	return dt
}

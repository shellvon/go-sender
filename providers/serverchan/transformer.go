package serverchan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/transformer"
)

// serverChanTransformer implements HTTPRequestTransformer for ServerChan.
// It supports sending text message.
//
// Reference:
//   - Official Website: https://sct.ftqq.com
//   - API Docs: https://sct.ftqq.com/sendkey

type serverChanTransformer struct {
	*transformer.BaseHTTPTransformer[*Message, *Account]
}

func (st *serverChanTransformer) transform(
	_ context.Context,
	msg *Message,
	account *Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	apiURL := buildAPIURL(account.APIKey)

	body, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal serverchan payload: %w", err)
	}

	return &core.HTTPRequestSpec{
		Method:   http.MethodPost,
		URL:      apiURL,
		Body:     body,
		BodyType: core.BodyTypeJSON,
	}, nil, nil
}

func newServerChanTransformer() *serverChanTransformer {
	respCfg := &core.ResponseHandlerConfig{
		BodyType:  core.BodyTypeJSON,
		CheckBody: true,
		Path:      "code",
		Expect:    "0",
		Mode:      core.MatchEq,
	}

	st := &serverChanTransformer{}
	st.BaseHTTPTransformer = transformer.NewSimpleHTTPTransformer(
		core.ProviderTypeServerChan,
		"",
		respCfg,
		st.transform,
	)

	return st
}

// buildAPIURL 构造 ServerChan API 地址
// 支持企业版和专业版的 key 格式
// 参考: https://github.com/easychen/serverchan-demo/blob/481eccec5c3a6da71d6e744207a145d31a92b5a5/golang/main.go#L57C2-L69C3
func buildAPIURL(key string) string {
	// Check if key starts with 'sctp' (enterprise/professional version)
	if strings.HasPrefix(key, "sctp") {
		// Extract number from sctp{num}t format
		re := regexp.MustCompile(`^sctp(\d+)t`)
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 {
			num := matches[1]
			return fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, key)
		}
	}
	// Standard Turbo version
	return fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
}

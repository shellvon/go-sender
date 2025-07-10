package core_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/utils"
)

// helper to build *http.Response quickly.
//
//nolint:unparam // ignore unused parameters.
func buildResp(status int, ct, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{ct}},
	}
}

func TestResponseHandler(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *core.ResponseHandlerConfig
		resp        *http.Response
		wantErr     bool
		errContains string
	}{
		{
			name: "JSON eq success",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "data.success",
				Expect:    "true",
			},
			resp:    buildResp(200, "application/json", `{"data":{"success":"true"}}`),
			wantErr: false,
		},
		{
			name: "XML eq success",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeXML,
				Path:      "root.code",
				Expect:    "0",
			},
			resp:    buildResp(200, "application/xml", `<root><code>0</code></root>`),
			wantErr: false,
		},
		{
			name: "Auto detect JSON by header when BodyTypeNone",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeNone,
				Path:      "ok",
				Expect:    "yes",
			},
			resp:    buildResp(200, "application/json", `{"ok":"yes"}`),
			wantErr: false,
		},
		{
			name: "Contains match",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeText,
				Path:      "", // whole body
				Expect:    "world",
				Mode:      core.MatchContains,
			},
			resp:    buildResp(200, "text/plain", "hello world"),
			wantErr: false,
		},
		{
			name: "Regex match",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeText,
				Path:      "", // whole body
				Expect:    `h.*o`,
				Mode:      core.MatchRegex,
			},
			resp:    buildResp(200, "text/plain", "hello"),
			wantErr: false,
		},
		{
			name: "Numeric gt success",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "score",
				Expect:    "80",
				Mode:      core.MatchGt,
			},
			resp:    buildResp(200, "application/json", `{"score":90}`),
			wantErr: false,
		},
		{
			name: "Numeric gt failure",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "score",
				Expect:    "80",
				Mode:      core.MatchGt,
			},
			resp:        buildResp(200, "application/json", `{"score":70}`),
			wantErr:     true,
			errContains: "api error",
		},
		{
			name: "CodeMap translation",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "success",
				Expect:    "true",
				CodePath:  "error.code",
				MsgPath:   "error.msg",
				CodeMap:   map[string]string{"E01": "network busy"},
			},
			resp: buildResp(
				200,
				"application/json",
				`{"success":"false","error":{"code":"E01","msg":"server"}}`,
			),
			wantErr:     true,
			errContains: "network busy",
		},
		{
			name: "Escaped dot key",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "data\\.value",
				Expect:    "ok",
			},
			resp:    buildResp(200, "application/json", `{"data.value":"ok"}`),
			wantErr: false,
		},
		{
			name: "Array path access",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "list[1].id",
				Expect:    "b",
			},
			resp:    buildResp(200, "application/json", `{"list":[{"id":"a"},{"id":"b"}]}`),
			wantErr: false,
		},
		{
			name: "XML array path access",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeXML,
				Path:      "items.item[2].name",
				Expect:    "third",
			},
			resp: buildResp(
				200,
				"application/xml",
				`<items><item><name>first</name></item><item><name>second</name></item><item><name>third</name></item></items>`,
			),
			wantErr: false,
		},
		{
			name: "XML regex whole body",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeXML,
				Path:      "", // whole body regex
				Expect:    `third`,
				Mode:      core.MatchRegex,
			},
			resp:    buildResp(200, "application/xml", `<root>first second third</root>`),
			wantErr: false,
		},
		{
			name: "Regex case-insensitive field match",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "code.message",
				Expect:    `(?i)ok`,
				Mode:      core.MatchRegex,
			},
			resp:    buildResp(200, "application/json", `{"code":{"message":"OK"}}`),
			wantErr: false,
		},
		{
			name: "Whole body regex case-insensitive",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeNone, // auto detect
				Path:      "",
				Expect:    `(?i)hello`,
				Mode:      core.MatchRegex,
			},
			resp:    buildResp(200, "text/plain", "HeLLo world"),
			wantErr: false,
		},
		{
			name: "Escaped dot path with regex",
			cfg: &core.ResponseHandlerConfig{
				CheckBody: true,
				BodyType:  core.BodyTypeJSON,
				Path:      "data\\.msg",
				Expect:    `(?i)ok`,
				Mode:      core.MatchRegex,
			},
			resp:    buildResp(200, "application/json", `{"data.msg":"OK"}`),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := core.NewSendResultHandler(tc.cfg)
			body, _, err := utils.ReadAndClose(tc.resp)
			if err != nil {
				t.Fatalf("ReadAndClose failed: %v", err)
			}
			err = h(&core.SendResult{StatusCode: tc.resp.StatusCode, Body: body, Headers: tc.resp.Header})
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.errContains != "" && (err == nil || !strings.Contains(err.Error(), tc.errContains)) {
				t.Errorf("error should contain %q, got %v", tc.errContains, err)
			}
		})
	}
}

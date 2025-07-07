package sms_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/providers/sms"
	"github.com/shellvon/go-sender/utils"
)

func TestAliyunProvider_Send_Success(t *testing.T) {
	// 启动 httptest.Server 作为 mock HTTP 服务
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Code":"OK"}`))
	}))
	defer ts.Close()

	fakeTransformer := &fakeAliyunTransformer{url: ts.URL}
	config := &sms.Config{
		Items: []*sms.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{Provider: "sms", Name: "test", Disabled: false, SubType: "aliyun"},
					Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"},
				},
			},
		},
	}
	httpProvider, err := providers.NewHTTPProvider(
		"aliyun",
		fakeTransformer,
		config,
	)
	if err != nil {
		t.Fatalf("failed to create HTTP provider: %v", err)
	}
	p := &sms.Provider{HTTPProvider: httpProvider}
	msg := sms.Aliyun().To("***REMOVED***").Content("hi").SignName("sign").Build()
	err = p.Send(context.Background(), msg, &core.ProviderSendOptions{})
	if err != nil {
		t.Errorf("Send should succeed: %v", err)
	}
}

type fakeAliyunTransformer struct{ url string }

func (f *fakeAliyunTransformer) CanTransform(_ core.Message) bool { return true }

func (f *fakeAliyunTransformer) Transform(
	_ context.Context,
	_ core.Message,
	_ *sms.Account,
) (*core.HTTPRequestSpec, core.ResponseHandler, error) {
	return &core.HTTPRequestSpec{
			Method:   "POST",
			URL:      f.url,
			Headers:  map[string]string{},
			Body:     []byte("{}"),
			BodyType: core.BodyTypeJSON,
		}, func(resp *http.Response) error {
			body, _, err := utils.ReadAndClose(resp)
			if err != nil {
				return err
			}
			status := resp.StatusCode
			if status < 200 || status >= 300 {
				return fmt.Errorf("HTTP request failed with status %d: %s", status, string(body))
			}
			return nil
		}, nil
}

func TestAliyunProvider_Send_Error(t *testing.T) {
	// 启动 httptest.Server 返回 500
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Code":"ERROR"}`))
	}))
	defer ts.Close()

	fakeTransformer := &fakeAliyunTransformer{url: ts.URL}
	config := &sms.Config{
		Items: []*sms.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{Provider: "sms", SubType: "aliyun", Name: "test", Disabled: false},
					Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"},
				},
			},
		},
	}
	httpProvider, err := providers.NewHTTPProvider(
		"aliyun",
		fakeTransformer,
		config,
	)
	if err != nil {
		t.Fatalf("failed to create HTTP provider: %v", err)
	}
	p := &sms.Provider{HTTPProvider: httpProvider}
	msg := sms.Aliyun().To("***REMOVED***").Content("hi").SignName("sign").Build()
	err = p.Send(context.Background(), msg, &core.ProviderSendOptions{})
	if err == nil {
		t.Error("Send should fail when HTTPProvider returns error")
	}
}

package sms_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	gosender "github.com/shellvon/go-sender"
	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers"
	"github.com/shellvon/go-sender/providers/sms"
)

func TestAliyun_Integration_SendTextSMS(t *testing.T) {
	// mock http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Code":"OK"}`))
	}))
	defer ts.Close()

	acc := &sms.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{Provider: "sms", Name: "test"},
			Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"},
		},
	}
	msg := sms.Aliyun()
	msg.To("13800138000")
	msg.Content("hi")
	msg.SignName("sign")
	msg.Type(sms.Voice)
	msg.TemplateID("TTS_123456") // 语音短信需要 TTS_ 开头的模板
	// 可选：设置语音专属参数
	msg.CalledShowNumber("4000000000")
	msg.PlayTimes(2)
	msg.Volume(80)

	tr, ok := sms.GetTransformer("aliyun")
	if !ok {
		t.Fatal("Aliyun transformer not registered")
	}
	spec, handler, err := tr.Transform(context.Background(), msg.Build(), acc)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if spec == nil || handler == nil {
		t.Fatal("Transform should return spec and handler")
	}
	// 模拟 HTTP 请求
	client := &http.Client{}
	req, err := http.NewRequest(spec.Method, ts.URL, nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
	if handleErr := handler(resp); handleErr != nil {
		t.Errorf("handler failed: %v", handleErr)
	}
}

// rewriteRoundTripper 重写请求的 Host + Scheme，让所有请求都转到指定 mock 服务器。
type rewriteRoundTripper struct {
	base       http.RoundTripper
	targetHost string
}

func (rt rewriteRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	req.URL.Host = rt.targetHost
	return rt.base.RoundTrip(req)
}

func TestSender_DispatchToAliyun_EndToEnd(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm failed: %v", err)
		}
		if got := r.FormValue("PhoneNumbers"); got != "13800138000" {
			t.Errorf("PhoneNumbers = %s, want 13800138000", got)
		}
		if got := r.FormValue("SignName"); got != "sign" {
			t.Errorf("SignName = %s, want sign", got)
		}
		if got := r.FormValue("TemplateCode"); got != "SMS_123456" {
			t.Errorf("TemplateCode = %s, want SMS_123456", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Code":"OK"}`))
	}))
	defer ts.Close()

	acc := &sms.Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{Provider: "sms", SubType: "aliyun", Name: "test"},
			Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"},
		},
	}
	tr, ok := sms.GetTransformer("aliyun")
	if !ok {
		t.Fatal("Aliyun transformer not registered")
	}
	httpProvider, err := providers.NewHTTPProvider(
		"aliyun",
		tr,
		&sms.Config{
			Items: []*sms.Account{acc},
			ProviderMeta: core.ProviderMeta{
				Strategy: core.StrategyRoundRobin,
			},
		},
	)
	if err != nil {
		t.Fatalf("failed to create HTTP provider: %v", err)
	}
	aliyunProvider := &sms.Provider{HTTPProvider: httpProvider}

	sender := gosender.NewSender()
	sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, nil)

	// 使用自定义 RoundTripper 将请求重写到 httptest server。
	targetHost := ts.Listener.Addr().String()
	baseTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: rewriteRoundTripper{base: baseTransport, targetHost: targetHost}}
	sender.SetDefaultHTTPClient(client)

	msg := sms.Aliyun().
		To("13800138000").
		Content("hi").
		SignName("sign").
		Type(sms.SMSText).
		TemplateID("SMS_123456").
		Build()

	err = sender.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

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
	"github.com/shellvon/go-sender/utils"
)

func TestAliyun_Integration_SendTextSMS(t *testing.T) {
	// mock http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Code":"OK"}`))
	}))
	defer ts.Close()

	acc := &core.Account{
		APIKey:    "ak",
		APISecret: "sk",
	}
	msg := sms.Aliyun()
	msg.To("***REMOVED***")
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
		t.Fatalf("client.Do failed: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	if handleErr := handler(resp.StatusCode, body); handleErr != nil {
		t.Errorf("handler failed: %v", handleErr)
	}
}

func TestSender_DispatchToAliyun_EndToEnd(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm failed: %v", err)
		}
		if got := r.FormValue("PhoneNumbers"); got != "***REMOVED***" {
			t.Errorf("PhoneNumbers = %s, want ***REMOVED***", got)
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

	acc := &core.Account{
		Name:      "test",
		APIKey:    "ak",
		APISecret: "sk",
	}
	tr, ok := sms.GetTransformer("aliyun")
	if !ok {
		t.Fatal("Aliyun transformer not registered")
	}
	httpProvider := providers.NewHTTPProvider(
		"aliyun",
		[]*core.Account{acc},
		tr,
		utils.GetStrategy(core.StrategyRoundRobin),
	)
	aliyunProvider := &sms.Provider{HTTPProvider: httpProvider}

	sender := gosender.NewSender()
	sender.RegisterProvider(core.ProviderTypeSMS, aliyunProvider, nil)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	sender.SetDefaultHTTPClient(client)

	msg := sms.Aliyun().
		To("***REMOVED***").
		Content("hi").
		SignName("sign").
		Type(sms.SMSText).
		TemplateID("SMS_123456").
		Build()

	err := sender.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

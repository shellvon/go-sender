package sms_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

func mustGetAliyunTransformer(t *testing.T) core.HTTPTransformer[*sms.Account] {
	tr, ok := sms.GetTransformer("aliyun")
	if !ok {
		t.Fatal("Aliyun transformer not registered")
	}
	return tr
}

func TestAliyunTransformer_CanTransform(t *testing.T) {
	tr := mustGetAliyunTransformer(t)
	msg := sms.Aliyun().To("13800138000").Content("hi").SignName("sign").Build()
	if !tr.CanTransform(msg) {
		t.Error("CanTransform should return true for Aliyun message")
	}
	msg.SubProvider = "tencent"
	if tr.CanTransform(msg) {
		t.Error("CanTransform should return false for non-Aliyun message")
	}
}

func TestAliyunTransformer_Transform_Text(t *testing.T) {
	tr := mustGetAliyunTransformer(t)
	msg := sms.Aliyun().To("13800138000").Content("hi").SignName("sign").Build()
	msg.Type = sms.SMSText
	acc := &sms.Account{BaseAccount: core.BaseAccount{Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"}}}
	spec, handler, err := tr.Transform(context.Background(), msg, acc)
	if err != nil {
		t.Fatalf("Transform failed: %v", err)
	}
	if spec == nil || handler == nil {
		t.Fatal("Transform should return spec and handler")
	}
}

func TestAliyunTransformer_Transform_UnsupportedType(t *testing.T) {
	tr := mustGetAliyunTransformer(t)
	msg := sms.Aliyun().To("13800138000").Content("hi").SignName("sign").Build()
	msg.Type = 99 // 非法类型
	acc := &sms.Account{BaseAccount: core.BaseAccount{Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"}}}
	_, _, err := tr.Transform(context.Background(), msg, acc)
	if err == nil {
		t.Error("Transform should fail for unsupported type")
	}
}

func TestAliyunTransformer_Transform_TextWithVoiceParam(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.Type(sms.SMSText).To("13800138000").Content("hi").SignName("sign")
	msg := aliyunBuilder.Volume(80).Build()
	tr := mustGetAliyunTransformer(t)
	acc := &sms.Account{BaseAccount: core.BaseAccount{Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"}}}
	spec, handler, err := tr.Transform(context.Background(), msg, acc)
	if err != nil {
		t.Fatalf("Transform should succeed even with voice params in text message: %v", err)
	}
	if spec == nil || handler == nil {
		t.Fatal("Transform should return spec and handler")
	}
	// 验证语音参数被忽略（不包含在请求中）
	if spec.QueryParams.Get("Volume") != "" {
		t.Error("Volume parameter should be ignored for text messages")
	}
}

func TestAliyunTransformer_Transform_VoiceWithVoiceParam(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.Type(sms.Voice).To("13800138000").Content("hi").SignName("sign")
	msg := aliyunBuilder.Volume(80).Build()
	tr := mustGetAliyunTransformer(t)
	acc := &sms.Account{BaseAccount: core.BaseAccount{Credentials: core.Credentials{APIKey: "ak", APISecret: "sk"}}}
	spec, handler, err := tr.Transform(context.Background(), msg, acc)
	if err != nil {
		t.Fatalf("Transform failed for voice: %v", err)
	}
	if spec == nil || handler == nil {
		t.Fatal("Transform should return spec and handler for voice")
	}
}

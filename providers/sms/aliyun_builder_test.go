package sms_test

import (
	"reflect"
	"testing"

	"github.com/shellvon/go-sender/providers/sms"
)

func TestAliyunSMSBuilder_AllFields_Text(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.To("***REMOVED***", "13900139000")
	aliyunBuilder.Content("test content")
	aliyunBuilder.SignName("test-sign")
	aliyunBuilder.TemplateID("TEMPLATE_ID")
	aliyunBuilder.Params(map[string]string{"code": "1234"})
	aliyunBuilder.
		OutID("out-id").
		FallbackType("SMS").
		SmsTemplateCode("SMS_CODE").
		DigitalTemplateCode("DIGITAL_CODE").
		SmsTemplateParam("sms-param").
		DigitalTemplateParam("digital-param").
		SmsUpExtendCode("up-ext").
		CardObjects("card-obj")

	msg := aliyunBuilder.Build()

	if msg == nil {
		t.Fatal("Build() returned nil")
	}
	if got, want := msg.SignName, "test-sign"; got != want {
		t.Errorf("SignName = %q, want %q", got, want)
	}
	if got, want := msg.TemplateID, "TEMPLATE_ID"; got != want {
		t.Errorf("TemplateID = %q, want %q", got, want)
	}
	if got, want := msg.Content, "test content"; got != want {
		t.Errorf("Content = %q, want %q", got, want)
	}
	if got, want := msg.Mobiles, []string{"***REMOVED***", "13900139000"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Mobiles = %v, want %v", got, want)
	}
	if got, want := msg.TemplateParams["code"], "1234"; got != want {
		t.Errorf("TemplateParams[code] = %q, want %q", got, want)
	}
	// 检查 Extras 字段（只断言文本 key 存在，不关心语音 key）
	extra := msg.Extras
	if extra == nil {
		t.Fatal("Extras should not be nil")
	}
	textKeys := []string{
		"OutId", "FallbackType", "SmsTemplateCode", "DigitalTemplateCode", "SmsTemplateParam", "DigitalTemplateParam", "SmsUpExtendCode", "CardObjects",
	}
	for _, k := range textKeys {
		if _, ok := extra[k]; !ok {
			t.Errorf("Extras missing key: %s", k)
		}
	}
}

func TestAliyunSMSBuilder_AllFields_Voice(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.To("***REMOVED***")
	aliyunBuilder.Content("test voice content")
	aliyunBuilder.SignName("voice-sign")
	aliyunBuilder.TemplateID("VOICE_TEMPLATE")
	aliyunBuilder.Type(sms.Voice)
	aliyunBuilder.
		CalledShowNumber("4000000000").
		PlayTimes(2).
		Volume(80).
		Speed(100)

	msg := aliyunBuilder.Build()

	if msg == nil {
		t.Fatal("Build() returned nil")
	}
	if got, want := msg.SignName, "voice-sign"; got != want {
		t.Errorf("SignName = %q, want %q", got, want)
	}
	if got, want := msg.TemplateID, "VOICE_TEMPLATE"; got != want {
		t.Errorf("TemplateID = %q, want %q", got, want)
	}
	if got, want := msg.Content, "test voice content"; got != want {
		t.Errorf("Content = %q, want %q", got, want)
	}
	if got, want := msg.Mobiles, []string{"***REMOVED***"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Mobiles = %v, want %v", got, want)
	}
	// 检查 Extras 字段（语音专属 key）
	extra := msg.Extras
	if extra == nil {
		t.Fatal("Extras should not be nil")
	}
	voiceKeys := []string{"CalledShowNumber", "PlayTimes", "Volume", "Speed"}
	for _, k := range voiceKeys {
		if _, ok := extra[k]; !ok {
			t.Errorf("Extras missing key: %s", k)
		}
	}
}

func TestAliyunSMSBuilder_AllFields_Extras(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.To("***REMOVED***", "13900139000")
	aliyunBuilder.Content("test content")
	aliyunBuilder.SignName("test-sign")
	aliyunBuilder.TemplateID("TEMPLATE_ID")
	aliyunBuilder.Params(map[string]string{"code": "1234"})
	aliyunBuilder.
		CalledShowNumber("4000000000").
		PlayTimes(2).
		OutID("out-id").
		Volume(80).
		Speed(100).
		FallbackType("SMS").
		SmsTemplateCode("SMS_CODE").
		DigitalTemplateCode("DIGITAL_CODE").
		SmsTemplateParam("sms-param").
		DigitalTemplateParam("digital-param").
		SmsUpExtendCode("up-ext").
		CardObjects("card-obj")

	msg := aliyunBuilder.Build()

	if msg == nil {
		t.Fatal("Build() returned nil")
	}
	if got, want := msg.SignName, "test-sign"; got != want {
		t.Errorf("SignName = %q, want %q", got, want)
	}
	if got, want := msg.TemplateID, "TEMPLATE_ID"; got != want {
		t.Errorf("TemplateID = %q, want %q", got, want)
	}
	if got, want := msg.Content, "test content"; got != want {
		t.Errorf("Content = %q, want %q", got, want)
	}
	if got, want := msg.Mobiles, []string{"***REMOVED***", "13900139000"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Mobiles = %v, want %v", got, want)
	}
	if got, want := msg.TemplateParams["code"], "1234"; got != want {
		t.Errorf("TemplateParams[code] = %q, want %q", got, want)
	}
	// 检查 Extras 字段（只断言所有设置的 key 都在）
	extra := msg.Extras
	if extra == nil {
		t.Fatal("Extras should not be nil")
	}
	allKeys := []string{
		"CalledShowNumber", "PlayTimes", "Volume", "Speed", "OutId", "FallbackType", "SmsTemplateCode", "DigitalTemplateCode", "SmsTemplateParam", "DigitalTemplateParam", "SmsUpExtendCode", "CardObjects",
	}
	for _, k := range allKeys {
		if _, ok := extra[k]; !ok {
			t.Errorf("Extras missing key: %s", k)
		}
	}
}

func TestAliyunSMSBuilder_RequiredFields(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.To("***REMOVED***")
	aliyunBuilder.Content("hi")
	aliyunBuilder.SignName("sign")
	msg := aliyunBuilder.Build()
	if msg == nil {
		t.Fatal("Build() returned nil")
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}
}

func TestAliyunSMSBuilder_EmptyMobiles(t *testing.T) {
	aliyunBuilder := sms.Aliyun()
	aliyunBuilder.Content("hi")
	aliyunBuilder.SignName("sign")
	msg := aliyunBuilder.Build()
	if err := msg.Validate(); err == nil {
		t.Error("Validate() should fail if mobiles is empty")
	}
}

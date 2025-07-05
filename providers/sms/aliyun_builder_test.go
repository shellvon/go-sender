package sms_test

import (
	"reflect"
	"testing"

	"github.com/shellvon/go-sender/providers/sms"
)

func TestAliyunSMSBuilder_AllFields_Text(t *testing.T) {
	msg := sms.Aliyun().
		To("***REMOVED***", "13900139000").
		Content("test content").
		SignName("test-sign").
		TemplateID("TEMPLATE_ID").
		Params(map[string]string{"code": "1234"}).
		OutID("out-id").
		FallbackType("SMS").
		SmsTemplateCode("SMS_CODE").
		DigitalTemplateCode("DIGITAL_CODE").
		SmsTemplateParam("sms-param").
		DigitalTemplateParam("digital-param").
		SmsUpExtendCode("up-ext").
		CardObjects("card-obj").
		Build()

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
		"OutId", "FallbackType", "SmsTemplateCode", "DigitalTemplateCode", "SmsTemplateParam", "DigitalTemplateParam", "CardObjects",
	}
	for _, k := range textKeys {
		if _, ok := extra[k]; !ok {
			t.Errorf("Extras missing key: %s", k)
		}
	}
	// 检查 Extend 字段
	if got, want := msg.Extend, "up-ext"; got != want {
		t.Errorf("Extend = %q, want %q", got, want)
	}
}

func TestAliyunSMSBuilder_AllFields_Voice(t *testing.T) {
	msg := sms.Aliyun().
		To("***REMOVED***").
		Content("test voice content").
		SignName("voice-sign").
		TemplateID("VOICE_TEMPLATE").
		Type(sms.Voice).
		CalledShowNumber("4000000000").
		PlayTimes(2).
		Volume(80).
		Speed(100).
		Build()

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
	msg := sms.Aliyun().
		To("***REMOVED***", "13900139000").
		Content("test content").
		SignName("test-sign").
		TemplateID("TEMPLATE_ID").
		Params(map[string]string{"code": "1234"}).
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
		CardObjects("card-obj").
		Build()

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
		"CalledShowNumber", "PlayTimes", "Volume", "Speed", "OutId", "FallbackType", "SmsTemplateCode", "DigitalTemplateCode", "SmsTemplateParam", "DigitalTemplateParam", "CardObjects",
	}
	for _, k := range allKeys {
		if _, ok := extra[k]; !ok {
			t.Errorf("Extras missing key: %s", k)
		}
	}
	// 检查 Extend 字段
	if got, want := msg.Extend, "up-ext"; got != want {
		t.Errorf("Extend = %q, want %q", got, want)
	}
}

func TestAliyunSMSBuilder_RequiredFields(t *testing.T) {
	msg := sms.Aliyun().
		To("***REMOVED***").
		Content("hi").
		SignName("sign").
		Build()
	if msg == nil {
		t.Fatal("Build() returned nil")
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}
}

func TestAliyunSMSBuilder_EmptyMobiles(t *testing.T) {
	msg := sms.Aliyun().
		Content("hi").
		SignName("sign").
		Build()
	if err := msg.Validate(); err == nil {
		t.Error("Validate() should fail if mobiles is empty")
	}
}

package sms_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/sms"
)

func TestSMSBuilderAndValidate(t *testing.T) {
	msg := sms.Aliyun().
		To("***REMOVED***").
		Content("Test message").
		SignName("TestSign").
		Build()

	if msg == nil {
		t.Fatal("Builder returned nil message")
	}
	if err := msg.Validate(); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

// FakeProvider implements core.Provider for testing.
type FakeProvider struct{}

func (f *FakeProvider) Send(
	_ context.Context,
	msg core.Message,
	_ ...*core.ProviderSendOptions,
) (*core.SendResult, error) {
	if msg == nil {
		return nil, core.NewParamError("msg is nil")
	}
	return &core.SendResult{}, nil
}
func (f *FakeProvider) Name() string { return "fake" }

func TestSMSSendWithFakeProvider(t *testing.T) {
	provider := &FakeProvider{}
	msg := sms.Aliyun().To("***REMOVED***").Content("Test").SignName("Sign").Build()
	_, err := provider.Send(context.Background(), msg)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

package sms

import (
	"context"
	"reflect"
	"testing"

	"github.com/shellvon/go-sender/core"
)

// dummyTransformer is a minimal HTTPTransformer for test purposes.
// It echoes SubProvider back in the URL so we can assert the output.
type dummyTransformer struct{}

func (d *dummyTransformer) CanTransform(msg core.Message) bool { return true }

func (d *dummyTransformer) Transform(_ context.Context, _ core.Message, _ *Account) (
	*core.HTTPRequestSpec, core.ResponseHandler, error,
) {
	return &core.HTTPRequestSpec{Method: "GET", URL: "https://example.com"}, nil, nil
}

func TestSMSRegistryTransform_Success(t *testing.T) {
	const sub = "dummy"
	// register transformer
	RegisterTransformer(sub, &dummyTransformer{})

	msg := &Message{
		Type:        SMSText,
		Mobiles:     []string{"13800138000"},
		Content:     "hello",
		SubProvider: sub,
	}

	acc := &Account{
		BaseAccount: core.BaseAccount{
			AccountMeta: core.AccountMeta{
				Provider: string(core.ProviderTypeSMS),
				SubType:  sub,
				Name:     "acc1",
				Weight:   1,
			},
			Credentials: core.Credentials{APIKey: "x"},
		},
	}

	tr := &smsTransformer{}
	spec, _, err := tr.Transform(context.Background(), msg, acc)
	if err != nil {
		t.Fatalf("Transform() unexpected error: %v", err)
	}
	if spec == nil || spec.URL != "https://example.com" {
		t.Fatalf("unexpected spec: %+v", spec)
	}
}

func TestSMSRegistryTransform_Missing(t *testing.T) {
	msg := &Message{
		Type:        SMSText,
		Mobiles:     []string{"13800138000"},
		Content:     "hello",
		SubProvider: "unknown",
	}
	acc := &Account{}
	tr := &smsTransformer{}
	_, _, err := tr.Transform(context.Background(), msg, acc)
	if err == nil {
		t.Fatal("expected error for unsupported sub-provider, got nil")
	}
}

func TestBaseConfigSelect_FilterBySubProvider(t *testing.T) {
	cfg := &core.BaseConfig[*Account]{
		Items: []*Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{Name: "aliyun", SubType: "aliyun"},
					Credentials: core.Credentials{APIKey: "x"},
				},
			},
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{Name: "tencent", SubType: "tencent"},
					Credentials: core.Credentials{APIKey: "y"},
				},
			},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() failed: %v", err)
	}

	// filter by SubProvider == tencent
	sel, err := cfg.Select(context.Background(), func(a *Account) bool {
		return a.SubType == "tencent"
	})
	if err != nil {
		t.Fatalf("Select() error: %v", err)
	}
	if sel == nil || sel.GetName() != "tencent" {
		t.Fatalf("unexpected selected account: %+v", sel)
	}
}

// Ensure RegisterTransformer is thread safe (basic check).
func TestSMSRegistry_IdempotentRegister(t *testing.T) {
	const sub = "dummy2"
	tr1 := &dummyTransformer{}
	tr2 := &dummyTransformer{}
	RegisterTransformer(sub, tr1)
	RegisterTransformer(sub, tr2)

	got, ok := GetTransformer(sub)
	if !ok {
		t.Fatalf("GetTransformer returned false for %s", sub)
	}
	// It should hold the last registered one (overwrites).
	if !reflect.DeepEqual(got, tr2) {
		t.Fatalf("registry did not overwrite transformer, got %#v want %#v", got, tr2)
	}
}

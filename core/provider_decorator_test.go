package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

type fakeProvider struct {
	sendErr error
	name    string
}

func (f *fakeProvider) Send(_ context.Context, _ core.Message, _ *core.ProviderSendOptions) error {
	return f.sendErr
}
func (f *fakeProvider) Name() string { return f.name }

type fakeMessage struct{}

func (f *fakeMessage) Validate() error                 { return nil }
func (f *fakeMessage) ProviderType() core.ProviderType { return core.ProviderTypeSMS }
func (f *fakeMessage) MsgID() string                   { return "id" }
func (f *fakeMessage) GetMsgType() string              { return "" }
func (f *fakeMessage) GetSubProvider() string          { return "" }

func TestNewProviderDecorator(t *testing.T) {
	p := &fakeProvider{name: "p1"}
	pd := core.NewProviderDecorator(p, nil, &core.NoOpLogger{})
	if pd.Provider.Name() != "p1" {
		t.Error("Provider not set correctly")
	}
}

func TestProviderDecorator_Send_Sync(t *testing.T) {
	p := &fakeProvider{name: "p2"}
	pd := core.NewProviderDecorator(p, nil, &core.NoOpLogger{})
	msg := &fakeMessage{}
	err := pd.Send(context.Background(), msg)
	if err != nil {
		t.Errorf("Send should succeed, got %v", err)
	}
}

func TestProviderDecorator_Send_Async(t *testing.T) {
	p := &fakeProvider{name: "p3"}
	pd := core.NewProviderDecorator(p, nil, &core.NoOpLogger{})
	msg := &fakeMessage{}
	ch := make(chan bool, 1)
	err := pd.Send(
		context.Background(),
		msg,
		core.WithSendAsync(true),
		core.WithSendCallback(func(*core.SendResult, error) { ch <- true }),
	)
	if err != nil {
		t.Errorf("Async Send should succeed, got %v", err)
	}
	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Error("Async callback not called")
	}
}

func TestProviderDecorator_Send_Error(t *testing.T) {
	p := &fakeProvider{name: "p4", sendErr: errors.New("fail")}
	pd := core.NewProviderDecorator(p, nil, &core.NoOpLogger{})
	msg := &fakeMessage{}
	err := pd.Send(context.Background(), msg)
	if err == nil {
		t.Error("Send should return error")
	}
}

func TestProviderDecorator_Close_Idempotent(_ *testing.T) {
	p := &fakeProvider{name: "p5"}
	pd := core.NewProviderDecorator(p, nil, &core.NoOpLogger{})
	_ = pd.Close()
	_ = pd.Close()
}

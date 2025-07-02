package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

func TestDefaultMessage_MsgIDAndExtras(t *testing.T) {
	m := core.NewDefaultMessage()
	id1 := m.MsgID()
	id2 := m.MsgID()
	if id1 != id2 || id1 == "" {
		t.Errorf("MsgID should be stable and non-empty")
	}
	m.Extras = map[string]interface{}{
		"str":   "v",
		"int":   1,
		"float": 1.5,
		"bool":  true,
	}
	if v, ok := m.GetExtraString("str"); !ok || v != "v" {
		t.Error("GetExtraString failed")
	}
	if v := m.GetExtraStringOrDefault("none", "d"); v != "d" {
		t.Error("GetExtraStringOrDefault failed")
	}
	if v, ok := m.GetExtraInt("int"); !ok || v != 1 {
		t.Error("GetExtraInt failed")
	}
	if v := m.GetExtraIntOrDefault("none", 2); v != 2 {
		t.Error("GetExtraIntOrDefault failed")
	}
	if v, ok := m.GetExtraBool("bool"); !ok || !v {
		t.Error("GetExtraBool failed")
	}
	if v := m.GetExtraBoolOrDefault("none", true); !v {
		t.Error("GetExtraBoolOrDefault failed")
	}
	if v, ok := m.GetExtraFloat("float"); !ok || v != 1.5 {
		t.Error("GetExtraFloat failed")
	}
	if v := m.GetExtraFloatOrDefault("none", 2.5); v != 2.5 {
		t.Error("GetExtraFloatOrDefault failed")
	}
}

type fakeHealth struct{}

func (f *fakeHealth) HealthCheck(_ context.Context) *core.HealthCheck {
	return &core.HealthCheck{Status: core.HealthStatusHealthy, Message: "ok", Timestamp: time.Now()}
}

func TestHealthCheckerInterface(t *testing.T) {
	var h core.HealthChecker = &fakeHealth{}
	res := h.HealthCheck(context.Background())
	if res == nil || res.Status != core.HealthStatusHealthy {
		t.Error("HealthChecker interface failed")
	}
}

func TestDefaultMessage_SerializeDeserialize(t *testing.T) {
	ser := &core.DefaultSendOptionsSerializer{}
	opts := &core.SendOptions{
		Priority:              1,
		Timeout:               2,
		DisableCircuitBreaker: true,
		DisableRateLimiter:    true,
		Metadata:              map[string]interface{}{"k": "v"},
	}
	b, err := ser.Serialize(opts)
	if err != nil || len(b) == 0 {
		t.Fatalf("Serialize failed: %v", err)
	}
	opts2, err := ser.Deserialize(b)
	if err != nil || opts2.Priority != 1 || opts2.Timeout != 2 || !opts2.DisableCircuitBreaker ||
		!opts2.DisableRateLimiter ||
		opts2.Metadata["k"] != "v" {
		t.Errorf("Deserialize failed or wrong: %+v, %v", opts2, err)
	}
}

type cmpItem struct{ v int }

func (c *cmpItem) Compare(other *cmpItem) bool { return c.v < other.v }

func TestCompare(t *testing.T) {
	a := &cmpItem{v: 1}
	b := &cmpItem{v: 2}
	if !a.Compare(b) || b.Compare(a) {
		t.Error("Compare logic error")
	}
}

type schedItem struct{ tm time.Time }

func (s *schedItem) SetScheduledAt(t time.Time) { s.tm = t }
func (s *schedItem) GetScheduledAt() *time.Time { return &s.tm }

func TestSchedulable(t *testing.T) {
	item := &schedItem{}
	tm := time.Now().Add(1 * time.Hour)
	item.SetScheduledAt(tm)
	if item.GetScheduledAt() == nil || !item.GetScheduledAt().Equal(tm) {
		t.Error("Schedulable Get/Set failed")
	}
}

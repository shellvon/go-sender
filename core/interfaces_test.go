package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/shellvon/go-sender/core"
)

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

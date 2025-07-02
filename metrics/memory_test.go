package metrics_test

import (
	"sync"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/metrics"
)

func TestMemoryMetricsCollector_RecordAndGetStats(t *testing.T) {
	m := metrics.NewMemoryMetricsCollector()
	provider := "test"
	m.RecordSendResult(core.MetricsData{Provider: provider, Success: true})
	m.RecordSendResult(core.MetricsData{Provider: provider, Success: false})
	total, succ, fail := m.GetStats(provider)
	if total != 2 || succ != 1 || fail != 1 {
		t.Errorf("unexpected stats: total=%d succ=%d fail=%d", total, succ, fail)
	}
	// 未知 provider
	t2, s2, f2 := m.GetStats("unknown")
	if t2 != 0 || s2 != 0 || f2 != 0 {
		t.Error("unknown provider should return zero stats")
	}
}

func TestMemoryMetricsCollector_Concurrency(t *testing.T) {
	m := metrics.NewMemoryMetricsCollector()
	provider := "concurrent"
	wg := sync.WaitGroup{}
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.RecordSendResult(core.MetricsData{Provider: provider, Success: i%2 == 0})
		}(i)
	}
	wg.Wait()
	total, _, _ := m.GetStats(provider)
	if total != 100 {
		t.Errorf("expected 100 total, got %d", total)
	}
}

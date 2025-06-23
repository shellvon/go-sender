package metrics

import (
	"sync"

	"github.com/shellvon/go-sender/core"
)

// MemoryMetricsCollector is an in-memory implementation of MetricsCollector
// Only keep sendResult related statistics
type MemoryMetricsCollector struct {
	mu      sync.RWMutex
	metrics map[string]*providerMetrics
}

type providerMetrics struct {
	totalRequests   int64
	successRequests int64
	failedRequests  int64
}

// NewMemoryMetricsCollector creates a new in-memory metrics collector
func NewMemoryMetricsCollector() *MemoryMetricsCollector {
	return &MemoryMetricsCollector{
		metrics: make(map[string]*providerMetrics),
	}
}

// RecordSendResult records the result of a send operation
func (m *MemoryMetricsCollector) RecordSendResult(data core.MetricsData) {
	m.mu.Lock()
	provider, exists := m.metrics[data.Provider]
	if !exists {
		provider = &providerMetrics{}
		m.metrics[data.Provider] = provider
	}
	m.mu.Unlock()

	if data.Success {
		provider.successRequests++
	} else {
		provider.failedRequests++
	}
	provider.totalRequests++
}

// GetStats returns the metrics for a given provider (for testing)
func (m *MemoryMetricsCollector) GetStats(provider string) (total, success, failed int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if pm, ok := m.metrics[provider]; ok {
		return pm.totalRequests, pm.successRequests, pm.failedRequests
	}
	return 0, 0, 0
}

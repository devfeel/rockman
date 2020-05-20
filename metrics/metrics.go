package metrics

import (
	"sync"
)

const (
	CounterKeyNode    = "node"
	CounterKeyCluster = "cluster"
	CounterKeyRuntime = "runtime"
	CounterKeyDefault = "default"
)

type (
	Metrics interface {
		GetCounter(key string) Counter
		GetNodeCounter() Counter
		GetClusterCounter() Counter
		GetRuntimeCounter() Counter
		GetDefaultCounter() Counter
	}

	StandardMetrics struct {
		counters *sync.Map
	}
)

func NewMetrics() Metrics {
	metrics := new(StandardMetrics)
	metrics.counters = new(sync.Map)
	return metrics
}

// GetCounter get Counter by key
func (m *StandardMetrics) GetCounter(key string) Counter {
	var counter Counter
	loadCounter, exists := m.counters.Load(key)
	if !exists {
		counter = NewCounter()
		m.counters.Store(key, counter)
	} else {
		counter = loadCounter.(Counter)
	}
	return counter
}

func (m *StandardMetrics) GetNodeCounter() Counter {
	return m.GetCounter(CounterKeyNode)
}

func (m *StandardMetrics) GetClusterCounter() Counter {
	return m.GetCounter(CounterKeyCluster)
}

func (m *StandardMetrics) GetRuntimeCounter() Counter {
	return m.GetCounter(CounterKeyRuntime)
}

func (m *StandardMetrics) GetDefaultCounter() Counter {
	return m.GetCounter(CounterKeyDefault)
}

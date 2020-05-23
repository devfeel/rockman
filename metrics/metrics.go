package metrics

import (
	"sync"
)

const (
	LabelNodeStart          = "NodeStart"
	LabelStartTheWorld      = "StartTheWorld"
	LabelStopTheWorld       = "StopTheWorld"
	LabelTaskExec           = "TaskExec"
	LabelLeaderChange       = "LeaderChange"
	LabelLeaderChangeFailed = "LeaderChangeFailed"
	LabelWorkerNodeOffline  = "WorkerNodeOffline"
	LabelRegistryOnline     = "RegistryOnline"
	LabelRegistryOffline    = "RegistryOffline"
)

type (
	Metrics interface {
		GetCounter(key string) Counter
	}

	StandardMetrics struct {
		counters *sync.Map
	}
)

var (
	defaultMetrics Metrics
	labelMap       map[string]string
)

func init() {
	defaultMetrics = NewMetrics()
	labelMap = make(map[string]string)
	labelMap[LabelNodeStart] = LabelNodeStart
	labelMap[LabelStartTheWorld] = LabelStartTheWorld
	labelMap[LabelStopTheWorld] = LabelStopTheWorld
	labelMap[LabelTaskExec] = LabelTaskExec
	labelMap[LabelLeaderChange] = LabelLeaderChange
	labelMap[LabelLeaderChangeFailed] = LabelLeaderChangeFailed
	labelMap[LabelWorkerNodeOffline] = LabelWorkerNodeOffline
	labelMap[LabelRegistryOnline] = LabelRegistryOnline
	labelMap[LabelRegistryOffline] = LabelRegistryOffline
}

func Default() Metrics {
	return defaultMetrics
}

// GetCounter is a shortcut for Default().GetCounter
func GetCounter(key string) Counter {
	return Default().GetCounter(key)
}

// GetAllCountInfo get all counter's count for map[string]int64
func GetAllCountInfo() map[string]int64 {
	m := make(map[string]int64)
	for _, label := range labelMap {
		m[label] = GetCounter(label).Count()
	}
	return m
}

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

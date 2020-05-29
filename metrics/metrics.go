package metrics

import (
	"strings"
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

	Opts struct {
		// Namespace, Subsystem, and Name are components of the fully-qualified
		// name of the Metric (created by joining these components with
		// "_").
		Namespace string
		Subsystem string
		Name      string
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

// GetCounterByOpts return Counter by Opts
func GetCounterByOpts(opts *Opts) Counter {
	return Default().GetCounter(buildFQName(opts.Namespace, opts.Subsystem, opts.Name))
}

// GetCounter is a shortcut for Default().GetCounter
func GetCounter(label string) Counter {
	return Default().GetCounter(label)
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
func (m *StandardMetrics) GetCounter(label string) Counter {
	var counter Counter
	loadCounter, exists := m.counters.Load(label)
	if !exists {
		counter = NewCounter()
		m.counters.Store(label, counter)
	} else {
		counter = loadCounter.(Counter)
	}
	return counter
}

func buildFQName(namespace, subsystem, name string) string {
	if name == "" {
		return ""
	}
	switch {
	case namespace != "" && subsystem != "":
		return strings.Join([]string{namespace, subsystem, name}, "_")
	case namespace != "":
		return strings.Join([]string{namespace, name}, "_")
	case subsystem != "":
		return strings.Join([]string{subsystem, name}, "_")
	}
	return name
}

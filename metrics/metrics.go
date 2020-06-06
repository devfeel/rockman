package metrics

import (
	"github.com/devfeel/rockman/metrics/prometheus"
	"github.com/devfeel/rockman/metrics/standard"
	promclient "github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
	"time"
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
		GetStandardCounter(label string) Counter
		Inc(label string)
		Dec(label string)
		Add(label string, value int64)
	}

	// Counter incremented and decremented base on int64 value.
	Counter interface {
		StartTime() time.Time
		Clear()
		Count() int64
		Dec()
		Inc()
		Add(int64)
	}

	StandardMetrics struct {
		counters    *sync.Map
		promCounter *promclient.CounterVec
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

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	return standard.NewStandardCounter()
}

// GetAllCountInfo get all counter's count for map[string]int64
func GetAllCountInfo() map[string]int64 {
	m := make(map[string]int64)
	for _, label := range labelMap {
		m[label] = Default().GetStandardCounter(label).Count()
	}
	return m
}

func NewMetrics() Metrics {
	metrics := new(StandardMetrics)
	metrics.counters = new(sync.Map)
	metrics.promCounter = prometheus.InitCounter()
	return metrics
}

// GetCounter get Counter by key
func (m *StandardMetrics) GetStandardCounter(label string) Counter {
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

// Inc increments the counter by 1.
func (m *StandardMetrics) Inc(label string) {
	m.GetStandardCounter(label).Inc()
	m.promCounter.WithLabelValues(label).Inc()
}

//Dec decrements the counter by 1.
func (m *StandardMetrics) Dec(label string) {
	m.GetStandardCounter(label).Dec()
	m.promCounter.WithLabelValues(label).Add(-1)
}

// Add increments the counter by the given value.
func (m *StandardMetrics) Add(label string, value int64) {
	m.GetStandardCounter(label).Add(value)
	m.promCounter.WithLabelValues(label).Add(float64(value))
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

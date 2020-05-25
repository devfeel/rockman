package prometheus

import "github.com/prometheus/client_golang/prometheus"

var NodeStartCounter *prometheus.CounterVec

func init() {
	NodeStartCounter = createCounterVec("NodeStart")
	prometheus.MustRegister(NodeStartCounter)
}

func createCounterVec(name string) *prometheus.CounterVec {
	opt := prometheus.CounterOpts{Name: name}
	labelName := []string{}
	return prometheus.NewCounterVec(opt, labelName)
}

package prometheus

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(createCounterVec("NodeStart"))
}

func createCounterVec(name string) *prometheus.CounterVec {
	opt := prometheus.CounterOpts{Name: name}
	labelName := []string{}
	return prometheus.NewCounterVec(opt, labelName)
}

package prometheus

import (
	"github.com/devfeel/rockman/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func StartMetricsWeb(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	logger.Default().Debug("MetricsWeb begin listen " + addr)
	return http.ListenAndServe(addr, nil)
}

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

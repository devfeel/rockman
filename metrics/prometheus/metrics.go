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

var DefaultCounter *prometheus.CounterVec

func init() {
	DefaultCounter = createCounterVec("Default")
	prometheus.MustRegister(DefaultCounter)
}

func createCounterVec(name string) *prometheus.CounterVec {
	opt := prometheus.CounterOpts{
		Namespace: "Rockman",
		Subsystem: "",
		Name:      name}
	labelName := []string{"Label"}
	return prometheus.NewCounterVec(opt, labelName)
}

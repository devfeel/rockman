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

func InitCounter() *prometheus.CounterVec {
	defaultCounter := createCounterVec("Default")
	prometheus.MustRegister(defaultCounter)
	return defaultCounter
}

func createCounterVec(name string) *prometheus.CounterVec {
	opt := prometheus.CounterOpts{
		Namespace: "Rockman",
		Subsystem: "",
		Name:      name}
	labelName := []string{"Label"}
	return prometheus.NewCounterVec(opt, labelName)
}

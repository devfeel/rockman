package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestNodeStartCounter(t *testing.T) {
	DefaultCounter.WithLabelValues("NodeStart").Inc()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", nil))
	time.Sleep(time.Hour)
}

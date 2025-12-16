package metrics

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	Hits *prometheus.CounterVec
}

func NewMetric(addr string, name string) (*PrometheusMetrics, error) {
	var metr PrometheusMetrics

	metr.Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_hits_total",
	},
		[]string{"code", "method", "path"},
	)

	if err := prometheus.Register(metr.Hits); err != nil {
		return nil, err
	}

	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	go func() {
		m := http.NewServeMux()
		m.Handle("GET /metrics", promhttp.Handler())
		log.Printf("Metrics server is started on %s addres\n", addr)
		if err := http.ListenAndServe(addr, m); err != nil {
			log.Fatalf("Metrics server is stopped\n")
		}
	}()
	return &metr, nil
}

// IncHits with information about request
func (m *PrometheusMetrics) IncHits(status int, method, path string) {
	m.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

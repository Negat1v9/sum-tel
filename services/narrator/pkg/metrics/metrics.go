package metrics

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	Tokens *prometheus.CounterVec
}

func NewMetric(addr string, name string) (*PrometheusMetrics, error) {
	var metr PrometheusMetrics

	metr.Tokens = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_tokens_total",
	},
		[]string{"success"},
	)

	if err := prometheus.Register(metr.Tokens); err != nil {
		return nil, err
	}

	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	go func() {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics server is started on %s address\n", addr)
		if err := http.ListenAndServe(addr, m); err != nil {
			log.Fatalf("Metrics server is stopped\n")
		}
	}()
	return &metr, nil
}

// AddTokens with information about success and number of tokens used
func (m *PrometheusMetrics) AddTokens(success bool, tokens int) {
	m.Tokens.WithLabelValues(strconv.FormatBool(success)).Add(float64(tokens))
}

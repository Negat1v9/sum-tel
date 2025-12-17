package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	ParsedChannels prometheus.Counter
	ParsedMessages *prometheus.CounterVec
}

func NewMetric(addr string, name string) (*PrometheusMetrics, error) {
	var metr PrometheusMetrics

	metr.ParsedMessages = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name + "_parsed_messages_total",
	},
		[]string{"username"},
	)

	if err := prometheus.Register(metr.ParsedMessages); err != nil {
		return nil, err
	}

	metr.ParsedChannels = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_parsed_channels_total",
	})

	if err := prometheus.Register(metr.ParsedChannels); err != nil {
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

// IncParsedChannels increments the counter for parsed channels
func (m *PrometheusMetrics) IncParsedChannels() {
	m.ParsedChannels.Inc()
}

// AddParsedMessages adds value to the counter for parsed messages with username label
func (m *PrometheusMetrics) AddParsedMessages(username string, value int) {
	m.ParsedMessages.WithLabelValues(username).Add(float64(value))
}

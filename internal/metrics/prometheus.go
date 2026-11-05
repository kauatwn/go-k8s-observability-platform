package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// httpRequestsTotal conta as requisições, separadas por rota e método HTTP
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests received.",
		},
		[]string{"path", "method"},
	)
)

// Register registra as métricas no registry global do Prometheus
func Register() {
	prometheus.MustRegister(httpRequestsTotal)
}

// RecordHTTPRequest incrementa o contador de requisições HTTP
func RecordHTTPRequest(path, method string) {
	httpRequestsTotal.WithLabelValues(path, method).Inc()
}

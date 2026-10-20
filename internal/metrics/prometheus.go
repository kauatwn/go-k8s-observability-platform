package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HTTPRequestsTotal conta as requisições, separadas por rota e método HTTP
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests received.",
		},
		[]string{"path", "method"},
	)
)

// Register inicializa os contadores no registro global do Prometheus
func Register() {
	prometheus.MustRegister(HTTPRequestsTotal)
}

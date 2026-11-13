// Package metrics centraliza a definição, o registro e a atualização
// das métricas customizadas da aplicação exportadas para o Prometheus.
package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// httpRequestsTotal contabiliza o volume total de requisições HTTP recebidas (Rate)
	// e permite calcular a taxa de falhas (Errors) por meio do rótulo "status".
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests received.",
		},
		[]string{"path", "method", "status"},
	)

	// httpRequestDuration registra a distribuição do tempo de resposta das requisições
	// em segundos (Duration), permitindo calcular percentis de latência (ex: p95, p99).
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets, // Buckets padrão do Prometheus (de 5ms até 10s)
		},
		[]string{"path", "method"},
	)
)

// Register inscreve todos os coletores customizados no registry global do Prometheus,
// tornando-os visíveis no endpoint /metrics.
func Register() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

// RecordHTTPRequest atualiza o contador de requisições e observa a duração
// de uma requisição HTTP concluída.
func RecordHTTPRequest(path, method string, statusCode int, durationSeconds float64) {
	status := strconv.Itoa(statusCode)
	httpRequestsTotal.WithLabelValues(path, method, status).Inc()
	httpRequestDuration.WithLabelValues(path, method).Observe(durationSeconds)
}

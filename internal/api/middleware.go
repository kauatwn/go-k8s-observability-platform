package api

import (
	"go-k8s-observability-platform/internal/metrics"
	"log/slog"
	"net/http"
)

// RecordMetricsMiddleware conta as requisições usando a lib oficial do Prometheus
func RecordMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Incrementa o contador passando o Path (ex: "/healthz") e o Method (ex: "GET")
		metrics.HTTPRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()

		slog.Info("HTTP request received",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})
}

package api

import (
	"go-k8s-observability-platform/internal/metrics"
	"log/slog"
	"net/http"
)

// RecordMetricsMiddleware registra métricas e logs das requisições HTTP usando a lib oficial do Prometheus.
func RecordMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			path := r.Pattern

			if path == "" {
				path = "unmatched"
			}

			// Incrementa o contador passando o Path (ex: "/healthz") e o Method (ex: "GET")
			metrics.RecordHTTPRequest(path, r.Method)
		}()

		slog.Info("HTTP request received",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})
}

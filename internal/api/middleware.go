package api

import (
	"go-k8s-observability-platform/internal/metrics"
	"log/slog"
	"net/http"
	"time"
)

// statusResponseWriter encapsula o http.ResponseWriter padrão para capturar
// o primeiro código de status HTTP enviado pelo handler.
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

// WriteHeader intercepta o envio explícito do status code.
// A trava wroteHeader garante que apenas o primeiro status enviado seja registrado,
// refletindo o comportamento real do protocolo HTTP no net/http.
func (w *statusResponseWriter) WriteHeader(code int) {
	if !w.wroteHeader {
		w.statusCode = code
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

// Write intercepta a escrita do corpo da resposta.
// Caso o handler escreva dados sem chamar WriteHeader antes, o net/http assume
// automaticamente o status 200 OK e trava o cabeçalho.
func (w *statusResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(b)
}

// Unwrap expõe o http.ResponseWriter original para o http.ResponseController,
// preservando recursos opcionais da conexão (como http.Flusher ou http.Hijacker).
func (w *statusResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// RecordMetricsMiddleware envolve um http.Handler para coletar métricas RED
// (Rate, Errors, Duration) e emitir logs estruturados em JSON.
func RecordMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Inicializa o wrapper assumindo 200 OK como status padrão
		sw := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// O defer garante que a coleta ocorra após o ServeMux resolver a rota
		// e o handler finalizar a escrita da resposta.
		defer func() {
			// r.Pattern contém o padrão registrado no ServeMux (ex: "/{$}", "/healthz").
			// Para rotas inexistentes (404), r.Pattern fica vazio; usar "unmatched"
			// evita explosão de cardinalidade nas séries temporais do Prometheus.
			path := r.Pattern
			if path == "" {
				path = "unmatched"
			}

			duration := time.Since(start).Seconds()
			metrics.RecordHTTPRequest(path, r.Method, sw.statusCode, duration)

			slog.Info("HTTP request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", sw.statusCode),
				slog.Float64("duration_seconds", duration),
				slog.String("remote_addr", r.RemoteAddr),
			)
		}()

		next.ServeHTTP(sw, r)
	})
}

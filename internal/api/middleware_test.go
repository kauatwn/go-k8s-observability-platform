package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordMetricsMiddleware(t *testing.T) {
	called := false

	// Handler simulado que retorna um status diferente do padrão (201 Created)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	handler := RecordMetricsMiddleware(next)
	req := httptest.NewRequest(http.MethodPost, "/test-path", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
}

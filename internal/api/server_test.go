package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	handler := SetupRouter()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Home endpoint",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Healthz endpoint",
			path:           "/healthz",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Ready endpoint initially unavailable",
			path:           "/ready",
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "Metrics endpoint",
			path:           "/metrics",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unknown route returns 404",
			path:           "/rota-inexistente",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("path %s: expected status %d, got %d", tc.path, tc.expectedStatus, rr.Code)
			}
		})
	}
}

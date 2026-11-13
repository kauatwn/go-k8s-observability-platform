package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzHandler(t *testing.T) {
	handler := &HealthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler.Healthz(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("expected status 'UP', got '%s'", resp["status"])
	}
}

func TestReadyHandler_NotReady(t *testing.T) {
	handler := &HealthHandler{}
	// isReady defaults to false

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()

	handler.Ready(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, status)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["status"] != "NOT READY" {
		t.Errorf("expected status 'NOT READY', got '%s'", resp["status"])
	}
}

func TestReadyHandler_Ready(t *testing.T) {
	handler := &HealthHandler{}
	handler.isReady.Store(true)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rr := httptest.NewRecorder()

	handler.Ready(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, status)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["status"] != "READY" {
		t.Errorf("expected status 'READY', got '%s'", resp["status"])
	}
}

package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_GetMetricsByName(t *testing.T) {
	defaultMemStorage := storage.NewMemStorage()
	presetMemStorage := func() Storage {
		mem := storage.NewMemStorage()
		mem.AddCounter(agent.Alloc, 22)
		mem.SetGauge(agent.Frees, 22)
		return mem
	}
	tests := []struct {
		name    string
		tpMet   string
		nameMet string
		mem     Storage
		status  int
	}{
		{name: "gauge not found", tpMet: agent.Gauge, nameMet: agent.Alloc, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "counter not found", tpMet: agent.Counter, nameMet: agent.Frees, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "response 200 for gauge", tpMet: agent.Gauge, nameMet: agent.Frees, mem: presetMemStorage(), status: http.StatusOK},
		{name: "response 200 for counter", tpMet: agent.Counter, nameMet: agent.Alloc, mem: presetMemStorage(), status: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/value/%s/%s", tt.tpMet, tt.nameMet)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				mem: tt.mem,
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", u.GetMetricsByName)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/internal/storage"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler_GetJSONMetricsByName(t *testing.T) {
	defaultMemStorage := storage.NewMemStorage()
	presetMemStorage := func() Storage {
		mem := storage.NewMemStorage()
		mem.AddCounter(agent.Alloc, 22)
		mem.SetGauge(agent.Frees, 22)
		return mem
	}
	tests := []struct {
		name   string
		mem    Storage
		body   metrics.MetricName
		status int
	}{
		{name: "gauge not found", body: metrics.MetricName{ID: agent.Alloc, MType: agent.Gauge}, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "counter not found", body: metrics.MetricName{ID: agent.Frees, MType: agent.Counter}, mem: defaultMemStorage, status: http.StatusNotFound},
		{name: "response 200 for gauge", body: metrics.MetricName{ID: agent.Frees, MType: agent.Gauge}, mem: presetMemStorage(), status: http.StatusOK},
		{name: "response 200 for counter", body: metrics.MetricName{ID: agent.Alloc, MType: agent.Counter}, mem: presetMemStorage(), status: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.body)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/value/", bytes.NewBuffer(out))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				mem: tt.mem,
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/value/", u.PostJSONMetricsByName)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

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

func TestUpdateMetrics_updateMetrics(t *testing.T) {
	tests := []struct {
		name    string
		tpMet   string
		nameMet string
		valMet  string
		status  int
	}{
		{name: "response 200 for gauge", tpMet: agent.Gauge, nameMet: agent.Alloc, valMet: "11.1", status: http.StatusOK},
		{name: "response 200 for counter", tpMet: agent.Counter, nameMet: agent.Frees, valMet: "11", status: http.StatusOK},
		{name: "wrong type", tpMet: "wrong", nameMet: agent.Alloc, valMet: "11", status: http.StatusBadRequest},
		{name: "without name", tpMet: agent.Gauge, nameMet: "", valMet: "11", status: http.StatusNotFound},
		{name: "value for gauge not float", tpMet: agent.Gauge, nameMet: agent.Frees, valMet: "str", status: http.StatusBadRequest},
		{name: "value for counter not int", tpMet: agent.Counter, nameMet: agent.Alloc, valMet: "11.1", status: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", tt.tpMet, tt.nameMet, tt.valMet)
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			u := &MetricsHandler{
				storage: storage.NewMemStorage(),
			}

			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", u.UpdateMetric)
			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

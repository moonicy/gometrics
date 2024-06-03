package handlers

import (
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
		{name: "response 200 for gauge", tpMet: agent.Gauge, nameMet: agent.Alloc, valMet: "11.1", status: 200},
		{name: "response 200 for counter", tpMet: agent.Counter, nameMet: agent.Frees, valMet: "11", status: 200},
		{name: "wrong type", tpMet: "wrong", nameMet: agent.Alloc, valMet: "11", status: 400},
		{name: "without name", tpMet: agent.Gauge, nameMet: "", valMet: "11", status: 404},
		{name: "value for gauge not float", tpMet: agent.Gauge, nameMet: agent.Frees, valMet: "str", status: 400},
		{name: "value for counter not int", tpMet: agent.Counter, nameMet: agent.Alloc, valMet: "11.1", status: 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/update/{type}/{name}/{value}/", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.SetPathValue("type", tt.tpMet)
			req.SetPathValue("name", tt.nameMet)
			req.SetPathValue("value", tt.valMet)

			u := &UpdateMetrics{
				mem: storage.NewMemStorage(),
			}
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(u.UpdateMetrics)
			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}
		})
	}
}

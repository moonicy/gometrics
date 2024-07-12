package client

import (
	"github.com/moonicy/gometrics/internal/agent"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if len(data) == 0 {
			t.Errorf("expected not empty body")
		}
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	report := agent.NewReport()
	report.Gauge[agent.Alloc] = 11
	report.Gauge[agent.Frees] = 22
	report.Gauge[agent.GCSys] = 33
	report.Counter[agent.PollCount] = 33

	cl := &Client{
		httpClient: http.DefaultClient,
		host:       server.URL,
	}
	cl.SendReport(report)
}

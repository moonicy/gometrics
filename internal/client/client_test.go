package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonicy/gometrics/internal/agent"
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
	report.SetGauge(agent.Alloc, 11)
	report.SetGauge(agent.Frees, 22)
	report.SetGauge(agent.GCSys, 33)
	report.AddCounter(agent.PollCount, 33)

	cl := &Client{
		httpClient: http.DefaultClient,
		host:       server.URL,
	}
	cl.SendReport(report)
}

func BenchmarkClient_makeResponseData(b *testing.B) {
	client := &Client{}
	report := agent.NewReport()
	reader := agent.NewMetricsReader()
	reader.Read(report)
	for i := 0; i < b.N; i++ {
		_, _ = client.makeRequestData(report)
	}
}

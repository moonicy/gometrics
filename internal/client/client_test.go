package client

import (
	"github.com/moonicy/gometrics/internal/agent"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_SendReport(t *testing.T) {
	counter := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
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
	l := len(report.Gauge) + len(report.Counter)
	time.Sleep(1 * time.Second)
	if counter != l {
		t.Errorf("expected count %d, got %d", l, counter)
	}
}

func TestClient_sendMetrics(t *testing.T) {
	type fields struct {
		httpClient     *http.Client
		reportInterval time.Duration
	}
	type args struct {
		tp    string
		name  string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "response 200 for gauge", fields: fields{
			httpClient:     &http.Client{},
			reportInterval: 10 * time.Millisecond,
		}, args: args{
			tp:    agent.Gauge,
			name:  agent.Alloc,
			value: "11",
		}},
		{name: "response 200 for counter", fields: fields{
			httpClient:     &http.Client{},
			reportInterval: 2 * time.Second,
		}, args: args{
			tp:    agent.Counter,
			name:  agent.PollCount,
			value: "22",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			defer server.Close()

			cl := &Client{
				httpClient: tt.fields.httpClient,
				host:       server.URL,
			}
			status := cl.sendMetrics(tt.args.tp, tt.args.name, tt.args.value)
			if status != "200 OK" {
				t.Errorf("expected status 200 OK, got %s", status)
			}
		})
	}
}

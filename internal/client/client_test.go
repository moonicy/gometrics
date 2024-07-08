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

func TestClient_sendGaugeMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	cl := &Client{
		httpClient: &http.Client{},
		host:       server.URL,
	}
	status := cl.sendGaugeMetrics(agent.Gauge, agent.Alloc, 11)
	if status != "200 OK" {
		t.Errorf("expected status 200 OK, got %s", status)
	}
}

func TestClient_sendCounterMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	cl := &Client{
		httpClient: &http.Client{},
		host:       server.URL,
	}
	status := cl.sendGaugeMetrics(agent.Counter, agent.Alloc, 11)
	if status != "200 OK" {
		t.Errorf("expected status 200 OK, got %s", status)
	}
}

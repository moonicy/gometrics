package client

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonicy/gometrics/internal/agent"
)

func TestClient_SendReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(r.Body)
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
	cl.SendReport(context.TODO(), report)
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

func TestNewClient(t *testing.T) {
	host := "http://example.com"
	key := "testHashKey"
	cryptoKey := "testCryptoKey"

	client := NewClient(host, key, cryptoKey)

	if client.host != host {
		t.Errorf("Expected host to be %v, got %v", host, client.host)
	}
	if client.hashKey != key {
		t.Errorf("Expected hashKey to be %v, got %v", key, client.hashKey)
	}
	if client.cryptoKey != cryptoKey {
		t.Errorf("Expected cryptoKey to be %v, got %v", cryptoKey, client.cryptoKey)
	}
	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized, but it was nil")
	}
}

func TestMakeRequestData(t *testing.T) {
	client := &Client{}

	report := agent.NewReport()

	report.SetGauge("gauge1", 10.5)
	report.AddCounter("counter1", 100)

	data, err := client.makeRequestData(report)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var metrics []map[string]interface{}
	err = jsoniter.Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Failed to unmarshal request data: %v", err)
	}

	if len(metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(metrics))
	}
}

func TestExternalIP(t *testing.T) {
	client := &Client{}

	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("Failed to get network interfaces: %v", err)
	}

	var hasIP bool
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				ip, _, _ := net.ParseCIDR(addr.String())
				if ip != nil && ip.To4() != nil {
					hasIP = true
					break
				}
			}
			if hasIP {
				break
			}
		}
	}

	_, err = client.externalIP()
	if hasIP && err != nil {
		t.Fatalf("Expected IP, got error: %v", err)
	} else if !hasIP && err == nil {
		t.Error("Expected error due to no external IP, but got none")
	}
}

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/metrics"
	"log"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	host       string
}

func NewClient(host string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
	}
}

func (cl *Client) sendGaugeMetrics(tp string, name string, value float64) string {
	body := metrics.Metrics{ID: name, MType: tp, Value: &value}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/update", cl.host)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return resp.Status
}

func (cl *Client) sendCounterMetrics(tp string, name string, delta int64) string {
	body := metrics.Metrics{ID: name, MType: tp, Delta: &delta}
	out, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/update", cl.host)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return resp.Status
}

func (cl *Client) SendReport(mem *agent.Report) {
	for k, v := range mem.Gauge {
		go cl.sendGaugeMetrics(agent.Gauge, k, v)
	}
	for k, v := range mem.Counter {
		go cl.sendCounterMetrics(agent.Counter, k, v)
	}
}

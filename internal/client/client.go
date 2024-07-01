package client

import (
	"bytes"
	"compress/gzip"
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
	body := metrics.Metrics{MetricName: metrics.MetricName{ID: name, MType: tp}, Value: &value}
	out, err := json.Marshal(body)
	if err != nil {
		log.Print(err)
		return ""
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, _ = zb.Write(out)

	zb.Close()

	url := fmt.Sprintf("%s/update/", cl.host)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Print(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()
	return resp.Status
}

func (cl *Client) sendCounterMetrics(tp string, name string, delta int64) string {
	body := metrics.Metrics{MetricName: metrics.MetricName{ID: name, MType: tp}, Delta: &delta}
	out, err := json.Marshal(body)
	if err != nil {
		log.Print(err)
		return ""
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, _ = zb.Write(out)

	zb.Close()

	url := fmt.Sprintf("%s/update/", cl.host)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Print(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	resp, err := cl.httpClient.Do(req)
	if err != nil {
		log.Print(err)
		return ""
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

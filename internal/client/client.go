package client

import (
	"fmt"
	"github.com/moonicy/gometrics/internal/agent"
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

func (cl *Client) sendMetrics(tp string, name string, value string) string {
	url := fmt.Sprintf("%s/update/%s/%s/%s", cl.host, tp, name, value)
	req, err := http.NewRequest("POST", url, nil)
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
		go cl.sendMetrics(agent.Gauge, k, fmt.Sprintf("%f", v))
	}
	for k, v := range mem.Counter {
		go cl.sendMetrics(agent.Counter, k, fmt.Sprintf("%d", v))
	}
}

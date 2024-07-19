package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moonicy/gometrics/internal/agent"
	m "github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/pkg/gzip"
	sign "github.com/moonicy/gometrics/pkg/hash"
	"github.com/moonicy/gometrics/pkg/retry"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	host       string
	hashKey    string
}

func NewClient(host string, key string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
		hashKey:    key,
	}
}

func (cl *Client) SendReport(report *agent.Report) {
	metrics := make([]m.Metric, 0, len(report.Gauge)+len(report.Counter))
	for k, v := range report.Counter {
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    k,
				MType: m.Counter,
			},
			Delta: &v,
			Value: nil,
		})
	}
	for k, v := range report.Gauge {
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    k,
				MType: m.Gauge,
			},
			Delta: nil,
			Value: &v,
		})
	}
	out, err := json.Marshal(metrics)
	if err != nil {
		log.Print(err)
		return
	}

	buf, err := gzip.Compress(out)
	if err != nil {
		log.Print(err)
		return
	}

	uri := fmt.Sprintf("%s/updates/", cl.host)
	req, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	if cl.hashKey != "" {
		hash := sign.CalcHash(out, cl.hashKey)
		req.Header.Add("HashSHA256", hash)
	}

	var resp *http.Response
	err = retry.RetryHandle(func() error {
		resp, err = cl.httpClient.Do(req)
		if err != nil {
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				return retry.NewRetryableError(urlErr.Error())
			}
			return err
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.Print(err)
			}
		}()
		if resp.StatusCode >= http.StatusInternalServerError {
			return retry.NewRetryableError("Server is not available")
		}
		return nil
	})
	if err != nil {
		log.Print(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Print("Wrong status code", resp.Status)
	}
}

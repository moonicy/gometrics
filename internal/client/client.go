package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"github.com/moonicy/gometrics/internal/agent"
	m "github.com/moonicy/gometrics/internal/metrics"
	"github.com/moonicy/gometrics/pkg/gzip"
	sign "github.com/moonicy/gometrics/pkg/hash"
	"github.com/moonicy/gometrics/pkg/retry"
)

// Client представляет клиента для отправки метрик на сервер.
type Client struct {
	httpClient *http.Client
	host       string
	hashKey    string
}

// NewClient создаёт и возвращает новый экземпляр Client с заданным хостом и ключом хеширования.
func NewClient(host string, key string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
		hashKey:    key,
	}
}

// SendReport отправляет отчет с метриками на сервер.
// Он собирает данные метрик, сжимает их, добавляет необходимые заголовки и отправляет HTTP-запрос.
// В случае ошибок выполняет повторные попытки с помощью механизма retry.
func (cl *Client) SendReport(report *agent.Report) {
	out, err := cl.makeRequestData(report)
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

func (cl *Client) makeRequestData(report *agent.Report) ([]byte, error) {
	metrics := make([]m.Metric, 0, report.GetCommonCount())
	report.Counter.Range(func(key, value any) bool {
		v := value.(int64)
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    key.(string),
				MType: m.Counter,
			},
			Delta: &v,
			Value: nil,
		})
		return true
	})
	report.Gauge.Range(func(key, value any) bool {
		v := value.(float64)
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    key.(string),
				MType: m.Gauge,
			},
			Delta: nil,
			Value: &v,
		})
		return true
	})
	out, err := jsoniter.Marshal(metrics)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return out, nil
}

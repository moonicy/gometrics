package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/moonicy/gometrics/pkg/crypt"
	"io"
	"log"
	"net"
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
	cryptoKey  string
}

// NewClient создаёт и возвращает новый экземпляр Client с заданным хостом и ключом хеширования.
func NewClient(host string, key string, cryptoKey string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
		hashKey:    key,
		cryptoKey:  cryptoKey,
	}
}

// SendReport отправляет отчет с метриками на сервер.
// Он собирает данные метрик, сжимает их, добавляет необходимые заголовки и отправляет HTTP-запрос.
// В случае ошибок выполняет повторные попытки с помощью механизма retry.
func (cl *Client) SendReport(_ context.Context, report *agent.Report) {
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

	compressedData, err := io.ReadAll(buf)
	if err != nil {
		log.Print(err)
		return
	}

	if cl.cryptoKey != "" {
		compressedData, err = crypt.Encrypt(cl.cryptoKey, compressedData)
		if err != nil {
			log.Print(err)
			return
		}
	}

	ip, err := cl.externalIP()
	if err != nil {
		log.Print(err)
		return
	}

	uri := fmt.Sprintf("%s/updates/", cl.host)
	req, err := http.NewRequest("POST", uri, bytes.NewReader(compressedData))
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("X-Real-IP", ip)

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
	counter := report.GetCounter()
	for k, v := range counter {
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    k,
				MType: m.Counter,
			},
			Delta: &v,
			Value: nil,
		})
	}
	gauges := report.GetGauge()
	for k, v := range gauges {
		metrics = append(metrics, m.Metric{
			MetricName: m.MetricName{
				ID:    k,
				MType: m.Gauge,
			},
			Delta: nil,
			Value: &v,
		})
	}

	report.Clean()

	out, err := jsoniter.Marshal(metrics)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return out, nil
}

func (cl *Client) externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

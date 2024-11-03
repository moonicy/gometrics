package client

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/url"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/pkg/retry"
	pb "github.com/moonicy/gometrics/proto"
)

// GRPCClient представляет клиента для отправки метрик на сервер.
type GRPCClient struct {
	metricsClient pb.MetricsClient
}

// NewGRPCClient создаёт и возвращает новый экземпляр Client с заданным хостом и ключом хеширования.
func NewGRPCClient() (*GRPCClient, error) {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewMetricsClient(conn)
	return &GRPCClient{
		metricsClient: c,
	}, nil
}

// SendReport отправляет отчет с метриками на сервер.
// Он собирает данные метрик, сжимает их, добавляет необходимые заголовки и отправляет HTTP-запрос.
// В случае ошибок выполняет повторные попытки с помощью механизма retry.
func (cl *GRPCClient) SendReport(ctx context.Context, report *agent.Report) {
	out := cl.makeRequestData(report)

	err := retry.RetryHandle(func() error {
		resp, err := cl.metricsClient.UpdateMetrics(ctx, out)
		if err != nil {
			var urlErr *url.Error
			if errors.As(err, &urlErr) {
				return retry.NewRetryableError(urlErr.Error())
			}
			return err
		}
		if resp.Error != "" {
			return fmt.Errorf("grpc error: %s", resp.Error)
		}
		return nil
	})
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Sent report")
}

func (cl *GRPCClient) makeRequestData(report *agent.Report) *pb.UpdateMetricsRequest {
	req := &pb.UpdateMetricsRequest{}
	counter := report.GetCounter()
	for k, v := range counter {
		req.Counters = append(req.Counters, &pb.Counter{
			Id:    k,
			Delta: v,
		})
	}
	gauges := report.GetGauge()
	for k, v := range gauges {
		req.Gauges = append(req.Gauges, &pb.Gauge{
			Id:    k,
			Value: v,
		})
	}

	report.Clean()

	return req
}

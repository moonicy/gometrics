package server

import (
	"context"
	"fmt"
	pb "github.com/moonicy/gometrics/proto"
)

// Storage определяет интерфейс для операций с хранилищем метрик.
type Storage interface {
	SetMetrics(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
}

type GRPCServer struct {
	pb.UnimplementedMetricsServer
	storage Storage
}

func NewGRPCServer(storage Storage) *GRPCServer {
	return &GRPCServer{storage: storage}
}

// UpdateMetrics реализует интерфейс добавления метрик.
func (s *GRPCServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse

	mtGauge := make(map[string]float64)
	for _, m := range in.Gauges {
		mtGauge[m.GetId()] = m.GetValue()
		fmt.Printf("mtGauge[%s] = %f\n", m.GetId(), m.GetValue())
	}
	mtCounter := make(map[string]int64)
	for _, m := range in.Counters {
		mtCounter[m.GetId()] = m.GetDelta()
		fmt.Printf("mtCounter[%s] = %d\n", m.GetId(), m.GetDelta())
	}
	err := s.storage.SetMetrics(ctx, mtCounter, mtGauge)
	if err != nil {
		response.Error = fmt.Sprintf("error adding metrics: %v", err)
	}

	fmt.Println("Got new metrics")

	return &response, nil
}

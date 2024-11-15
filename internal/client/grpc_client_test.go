package client

import (
	"context"
	"errors"
	"github.com/moonicy/gometrics/internal/agent"
	"testing"

	"github.com/moonicy/gometrics/pkg/retry"
	pb "github.com/moonicy/gometrics/proto"
	"google.golang.org/grpc"
)

type MockMetricsClient struct {
	updateMetricsCount int
	resp               *pb.UpdateMetricsResponse
	err                error
}

func (m *MockMetricsClient) UpdateMetrics(_ context.Context, _ *pb.UpdateMetricsRequest, _ ...grpc.CallOption) (*pb.UpdateMetricsResponse, error) {
	m.updateMetricsCount++
	return m.resp, m.err
}

func TestNewGRPCClient(t *testing.T) {
	client, err := NewGRPCClient()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.metricsClient == nil {
		t.Error("Expected metricsClient to be initialized, but it was nil")
	}
}

func TestSendReport(t *testing.T) {
	ctx := context.Background()
	mockMetricsClient := &MockMetricsClient{
		resp: &pb.UpdateMetricsResponse{},
	}
	client := &GRPCClient{metricsClient: mockMetricsClient}

	report := agent.NewReport()
	report.SetGauge("gauge1", 10.5)
	report.AddCounter("counter1", 100)

	client.SendReport(ctx, report)
	if mockMetricsClient.updateMetricsCount != 1 {
		t.Errorf("Expected UpdateMetrics to be called once, got %d", mockMetricsClient.updateMetricsCount)
	}

	mockMetricsClient.err = errors.New("network error")
	err := retry.RetryHandle(func() error {
		client.SendReport(ctx, report)
		return mockMetricsClient.err
	})

	if err == nil {
		t.Error("Expected error due to network issue, but got none")
	}
}

func TestMakeRequestDataGrpc(t *testing.T) {
	client := &GRPCClient{}
	report := agent.NewReport()
	report.SetGauge("gauge1", 10.5)
	report.AddCounter("counter1", 100)

	data := client.makeRequestData(report)

	if len(data.Counters) != 1 || data.Counters[0].Id != "counter1" || data.Counters[0].Delta != 100 {
		t.Errorf("Unexpected counter data: %+v", data.Counters)
	}

	if len(data.Gauges) != 1 || data.Gauges[0].Id != "gauge1" || data.Gauges[0].Value != 10.5 {
		t.Errorf("Unexpected gauge data: %+v", data.Gauges)
	}
}

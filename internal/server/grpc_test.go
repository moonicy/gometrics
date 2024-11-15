package server

import (
	"context"
	"errors"
	"testing"

	pb "github.com/moonicy/gometrics/proto"
	"github.com/stretchr/testify/assert"
)

// MockStorage - мок реализации интерфейса Storage
type MockStorage struct {
	setMetricsCalled bool
	setMetricsError  error
	lastCounter      map[string]int64
	lastGauge        map[string]float64
}

func (m *MockStorage) SetMetrics(_ context.Context, counter map[string]int64, gauge map[string]float64) error {
	m.setMetricsCalled = true
	m.lastCounter = counter
	m.lastGauge = gauge
	return m.setMetricsError
}

func TestUpdateMetrics_Success(t *testing.T) {
	mockStorage := &MockStorage{}
	server := NewGRPCServer(mockStorage)

	request := &pb.UpdateMetricsRequest{
		Gauges: []*pb.Gauge{
			{Id: "gauge1", Value: 10.5},
		},
		Counters: []*pb.Counter{
			{Id: "counter1", Delta: 100},
		},
	}

	resp, err := server.UpdateMetrics(context.Background(), request)

	assert.NoError(t, err)
	assert.Empty(t, resp.Error)
	assert.True(t, mockStorage.setMetricsCalled, "Expected SetMetrics to be called")

	assert.Equal(t, map[string]float64{"gauge1": 10.5}, mockStorage.lastGauge)
	assert.Equal(t, map[string]int64{"counter1": 100}, mockStorage.lastCounter)
}

func TestUpdateMetrics_StorageError(t *testing.T) {
	mockStorage := &MockStorage{
		setMetricsError: errors.New("storage error"),
	}
	server := NewGRPCServer(mockStorage)

	request := &pb.UpdateMetricsRequest{
		Gauges: []*pb.Gauge{
			{Id: "gauge1", Value: 10.5},
		},
		Counters: []*pb.Counter{
			{Id: "counter1", Delta: 100},
		},
	}

	resp, err := server.UpdateMetrics(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, "error adding metrics: storage error", resp.Error)
	assert.True(t, mockStorage.setMetricsCalled, "Expected SetMetrics to be called")
}

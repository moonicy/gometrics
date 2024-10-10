package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
)

type MockConsumer struct {
	Events []*file.Event
	Opened bool
	Closed bool
}

func (m *MockConsumer) Open() error {
	m.Opened = true
	return nil
}

func (m *MockConsumer) ReadEvent() (*file.Event, error) {
	if len(m.Events) == 0 {
		return nil, errors.New("EOF")
	}
	event := m.Events[0]
	m.Events = m.Events[1:]
	return event, nil
}

func (m *MockConsumer) Close() error {
	m.Closed = true
	return nil
}

type MockProducer struct {
	Events  []*file.Event
	Opened  bool
	Closed  bool
	FailOn  string
	WriteFn func(event *file.Event) error
}

func (m *MockProducer) Open() error {
	if m.FailOn == "Open" {
		return errors.New("open error")
	}
	m.Opened = true
	return nil
}

func (m *MockProducer) WriteEvent(event *file.Event) error {
	if m.WriteFn != nil {
		return m.WriteFn(event)
	}
	if m.FailOn == "WriteEvent" {
		return errors.New("write error")
	}
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockProducer) Close() error {
	if m.FailOn == "Close" {
		return errors.New("close error")
	}
	m.Closed = true
	return nil
}

func TestNewFileStorage(t *testing.T) {
	cfg := config.ServerConfig{
		Host:            "localhost:8080",
		StoreInternal:   300,
		FileStoragePath: "/tmp/storage.log",
		Restore:         true,
		DatabaseDsn:     "user:pass@/dbname",
		HashKey:         "hashkey",
	}

	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := NewFileStorage(cfg, mockConsumer, mockProducer)

	if fs.cfg != cfg {
		t.Errorf("Expected cfg to be %+v, got %+v", cfg, fs.cfg)
	}
	if fs.consumer != mockConsumer {
		t.Errorf("Expected consumer to be %+v, got %+v", mockConsumer, fs.consumer)
	}
	if fs.producer != mockProducer {
		t.Errorf("Expected producer to be %+v, got %+v", mockProducer, fs.producer)
	}
	if fs.mem == nil {
		t.Errorf("Expected mem to be initialized, got nil")
	}
}

func TestFileStorage_SetGauge(t *testing.T) {
	cfg := config.ServerConfig{
		StoreInternal: 0,
	}
	mockMem := NewMemStorage()
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	eventWritten := false
	mockProducer.WriteFn = func(event *file.Event) error {
		eventWritten = true
		return nil
	}

	err := fs.SetGauge(ctx, "cpu", 0.75)
	if err != nil {
		t.Fatalf("SetGauge returned error: %v", err)
	}

	if val, _ := fs.mem.GetGauge(ctx, "cpu"); val != 0.75 {
		t.Errorf("Expected gauge 'cpu' to be 0.75, got %v", val)
	}

	if !eventWritten {
		t.Errorf("Expected WriteEvent to be called")
	}

	if !mockProducer.Opened {
		t.Errorf("Expected Producer.Open to be called")
	}
	if !mockProducer.Closed {
		t.Errorf("Expected Producer.Close to be called")
	}
}

func TestFileStorage_AddCounter(t *testing.T) {
	cfg := config.ServerConfig{
		StoreInternal: 0,
	}
	mockMem := NewMemStorage()
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	eventWritten := false
	mockProducer.WriteFn = func(event *file.Event) error {
		eventWritten = true
		return nil
	}

	err := fs.AddCounter(ctx, "requests", 10)
	if err != nil {
		t.Fatalf("AddCounter returned error: %v", err)
	}

	if val, _ := fs.mem.GetCounter(ctx, "requests"); val != 10 {
		t.Errorf("Expected counter 'requests' to be 10, got %v", val)
	}

	if !eventWritten {
		t.Errorf("Expected WriteEvent to be called")
	}

	if !mockProducer.Opened {
		t.Errorf("Expected Producer.Open to be called")
	}
	if !mockProducer.Closed {
		t.Errorf("Expected Producer.Close to be called")
	}
}

func TestFileStorage_GetCounter(t *testing.T) {
	mockMem := NewMemStorage()
	mockMem.counter["errors"] = 5
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
	}

	val, err := fs.GetCounter(ctx, "errors")
	if err != nil {
		t.Fatalf("GetCounter returned error: %v", err)
	}
	if val != 5 {
		t.Errorf("Expected counter 'errors' to be 5, got %v", val)
	}
}

func TestFileStorage_GetGauge(t *testing.T) {
	mockMem := NewMemStorage()
	mockMem.gauge["memory"] = 512.5
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
	}

	val, err := fs.GetGauge(ctx, "memory")
	if err != nil {
		t.Fatalf("GetGauge returned error: %v", err)
	}
	if val != 512.5 {
		t.Errorf("Expected gauge 'memory' to be 512.5, got %v", val)
	}
}

func TestFileStorage_GetMetrics(t *testing.T) {
	mockMem := NewMemStorage()
	mockMem.counter["requests"] = 100
	mockMem.gauge["cpu"] = 0.75
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
	}

	counter, gauge, err := fs.GetMetrics(ctx)
	if err != nil {
		t.Fatalf("GetMetrics returned error: %v", err)
	}

	if counter["requests"] != 100 {
		t.Errorf("Expected counter 'requests' to be 100, got %v", counter["requests"])
	}
	if gauge["cpu"] != 0.75 {
		t.Errorf("Expected gauge 'cpu' to be 0.75, got %v", gauge["cpu"])
	}
}

func TestFileStorage_SetMetrics(t *testing.T) {
	mockMem := NewMemStorage()
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      config.ServerConfig{},
	}

	counter := map[string]int64{"errors": 3}
	gauge := map[string]float64{"disk": 256.0}

	err := fs.SetMetrics(ctx, counter, gauge)
	if err != nil {
		t.Fatalf("SetMetrics returned error: %v", err)
	}

	if val, _ := fs.mem.GetCounter(ctx, "errors"); val != 3 {
		t.Errorf("Expected counter 'errors' to be 3, got %v", val)
	}
	if val, _ := fs.mem.GetGauge(ctx, "disk"); val != 256.0 {
		t.Errorf("Expected gauge 'disk' to be 256.0, got %v", val)
	}
}

func TestFileStorage_Init_Restore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.ServerConfig{
		Restore:         true,
		StoreInternal:   0,
		FileStoragePath: "dummy_path",
	}
	mockMem := NewMemStorage()
	mockConsumer := &MockConsumer{
		Events: []*file.Event{
			{
				Gauge:     map[string]float64{"cpu": 0.80},
				Counter:   map[string]int64{"requests": 50},
				Timestamp: time.Now().Unix(),
			},
		},
	}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	err := fs.Init(ctx)
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if val, _ := fs.mem.GetGauge(ctx, "cpu"); val != 0.80 {
		t.Errorf("Expected gauge 'cpu' to be 0.80 after Restore, got %v", val)
	}
	if val, _ := fs.mem.GetCounter(ctx, "requests"); val != 50 {
		t.Errorf("Expected counter 'requests' to be 50 after Restore, got %v", val)
	}
}

func TestFileStorage_RunSync(t *testing.T) {
	cfg := config.ServerConfig{
		StoreInternal: 1,
	}
	mockMem := NewMemStorage()
	mockMem.gauge["cpu"] = 0.90
	mockMem.counter["requests"] = 200
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	fs.RunSync()

	time.Sleep(2 * time.Second)

	if len(mockProducer.Events) != 1 {
		t.Errorf("Expected 1 event to be written, got %d", len(mockProducer.Events))
	}

	event := mockProducer.Events[0]
	if event.Gauge["cpu"] != 0.90 {
		t.Errorf("Expected gauge 'cpu' to be 0.90, got %v", event.Gauge["cpu"])
	}
	if event.Counter["requests"] != 200 {
		t.Errorf("Expected counter 'requests' to be 200, got %v", event.Counter["requests"])
	}
}

func TestFileStorage_uploadToFile_ErrorOnProducerOpen(t *testing.T) {
	cfg := config.ServerConfig{
		StoreInternal: 0,
	}
	mockMem := NewMemStorage()
	mockMem.gauge["cpu"] = 0.85
	mockMem.counter["requests"] = 150
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{
		FailOn: "Open",
	}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	err := fs.uploadToFile(ctx)
	if err == nil {
		t.Errorf("Expected error when Producer.Open fails, got nil")
	}
}

func TestFileStorage_WaitShutDown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.ServerConfig{
		StoreInternal: 0,
	}
	mockMem := NewMemStorage()
	mockMem.gauge["memory"] = 1024.0
	mockMem.counter["errors"] = 25
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}
	fs := &FileStorage{
		mem:      mockMem,
		consumer: mockConsumer,
		producer: mockProducer,
		cfg:      cfg,
	}

	go fs.Init(ctx)

	cancel()

	time.Sleep(500 * time.Millisecond)

	if len(mockProducer.Events) != 1 {
		t.Errorf("Expected 1 event to be written on shutdown, got %d", len(mockProducer.Events))
	}

	event := mockProducer.Events[0]
	if event.Gauge["memory"] != 1024.0 {
		t.Errorf("Expected gauge 'memory' to be 1024.0, got %v", event.Gauge["memory"])
	}
	if event.Counter["errors"] != 25 {
		t.Errorf("Expected counter 'errors' to be 25, got %v", event.Counter["errors"])
	}
}

package file

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func TestNewProducer(t *testing.T) {
	filename := "test_producer.log"
	producer := NewProducer(filename)

	if producer.filename != filename {
		t.Errorf("Expected filename to be %s, got %s", filename, producer.filename)
	}
	if producer.file != nil {
		t.Errorf("Expected file to be nil, got %v", producer.file)
	}
	if producer.writer != nil {
		t.Errorf("Expected writer to be nil, got %v", producer.writer)
	}
}

func TestProducer_Open_Success(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_producer_open_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	producer := NewProducer(tmpfile.Name())
	err = producer.Open()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if producer.file == nil {
		t.Errorf("Expected file to be opened, got nil")
	}
	if producer.writer == nil {
		t.Errorf("Expected writer to be initialized, got nil")
	}
}

func TestProducer_Open_RetryableError(t *testing.T) {
	producer := NewProducer("/nonexistent_dir/test.log")
	err := producer.Open()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestProducer_WriteEvent_Success(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_producer_write_event_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	producer := NewProducer(tmpfile.Name())
	err = producer.Open()
	if err != nil {
		t.Fatalf("Failed to open producer: %v", err)
	}
	defer producer.Close()

	event := &Event{
		Gauge:     map[string]float64{"cpu": 0.85, "memory": 1024.0},
		Counter:   map[string]int64{"requests": 200, "errors": 10},
		Timestamp: 1625080200,
	}

	err = producer.WriteEvent(event)
	if err != nil {
		t.Fatalf("Expected no error on WriteEvent, got %v", err)
	}

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file for reading: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		t.Fatalf("Expected to read a line from the file")
	}

	var readEvent Event
	err = json.Unmarshal(scanner.Bytes(), &readEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if !compareEvents(&readEvent, event) {
		t.Errorf("Expected event %+v, got %+v", event, readEvent)
	}

	if scanner.Scan() {
		t.Errorf("Expected only one event in the file")
	}
}

func TestProducer_WriteEvent_InvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_producer_invalid_json_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	producer := NewProducer(tmpfile.Name())
	err = producer.Open()
	if err != nil {
		t.Fatalf("Failed to open producer: %v", err)
	}
	defer producer.Close()

	event := &Event{
		Gauge:     map[string]float64{"cpu": 0.95},
		Counter:   map[string]int64{"requests": 300},
		Timestamp: 1625080800,
	}

	err = producer.WriteEvent(event)
	if err != nil {
		t.Fatalf("Expected no error on WriteEvent, got %v", err)
	}

	file, err := os.OpenFile(tmpfile.Name(), os.O_WRONLY, 0666)
	if err != nil {
		t.Fatalf("Failed to open file for corruption: %v", err)
	}
	defer file.Close()

	_, err = file.WriteAt([]byte("invalid_json"), 0)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON: %v", err)
	}

	consumer := NewConsumer(tmpfile.Name())
	err = consumer.Open()
	if err != nil {
		t.Fatalf("Failed to open consumer: %v", err)
	}
	defer consumer.Close()

	_, err = consumer.ReadEvent()
	if err == nil {
		t.Errorf("Expected JSON unmarshal error, got nil")
	}
}

func compareEvents(a, b *Event) bool {
	if a.Timestamp != b.Timestamp {
		return false
	}

	if len(a.Gauge) != len(b.Gauge) {
		return false
	}
	for k, v := range a.Gauge {
		if bv, ok := b.Gauge[k]; !ok || bv != v {
			return false
		}
	}

	if len(a.Counter) != len(b.Counter) {
		return false
	}
	for k, v := range a.Counter {
		if bv, ok := b.Counter[k]; !ok || bv != v {
			return false
		}
	}

	return true
}

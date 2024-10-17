package file

import (
	"os"
	"testing"
)

func TestNewConsumer(t *testing.T) {
	filename := "test_consumer.log"
	consumer := NewConsumer(filename)

	if consumer.filename != filename {
		t.Errorf("Expected filename to be %s, got %s", filename, consumer.filename)
	}
	if consumer.file != nil {
		t.Errorf("Expected file to be nil, got %v", consumer.file)
	}
	if consumer.scanner != nil {
		t.Errorf("Expected scanner to be nil, got %v", consumer.scanner)
	}
}

func TestConsumer_Open_Success(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_consumer_open_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(tmpfile.Name())

	consumer := NewConsumer(tmpfile.Name())
	err = consumer.Open()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if consumer.file == nil {
		t.Errorf("Expected file to be opened, got nil")
	}
	if consumer.scanner == nil {
		t.Errorf("Expected scanner to be initialized, got nil")
	}
}

func TestConsumer_Open_RetryableError(t *testing.T) {
	consumer := NewConsumer("/nonexistent_dir/test.log")
	err := consumer.Open()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestConsumer_ReadEvent_InvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test_consumer_invalid_json_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			t.Errorf("Failed to remove temp file: %v", err)
		}
	}(tmpfile.Name())

	_, err = tmpfile.WriteString("invalid_json\n")
	if err != nil {
		t.Errorf("Failed to write to temp file: %v", err)
	}

	consumer := NewConsumer(tmpfile.Name())
	err = consumer.Open()
	if err != nil {
		t.Fatalf("Failed to open consumer: %v", err)
	}
	defer func(consumer *Consumer) {
		err = consumer.Close()
		if err != nil {
			t.Fatalf("Failed to close consumer: %v", err)
		}
	}(consumer)

	_, err = consumer.ReadEvent()
	if err == nil {
		t.Errorf("Expected JSON unmarshal error, got nil")
	}
}

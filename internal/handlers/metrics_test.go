package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
)

type MockDB struct{}

func (m *MockDB) ExecContext(_ context.Context, _ string, _ ...any) (sql.Result, error) {
	return nil, nil
}
func (m *MockDB) QueryContext(_ context.Context, _ string, _ ...any) (*sql.Rows, error) {
	return nil, nil
}
func (m *MockDB) QueryRowContext(_ context.Context, _ string, _ ...any) *sql.Row {
	return &sql.Row{}
}
func (m *MockDB) Begin() (*sql.Tx, error) {
	return nil, nil
}

type MockConsumer struct{}

func (m *MockConsumer) Open() error {
	return nil
}
func (m *MockConsumer) ReadEvent() (*file.Event, error) {
	return &file.Event{}, nil
}
func (m *MockConsumer) Close() error {
	return nil
}

type MockProducer struct{}

func (m *MockProducer) Open() error {
	return nil
}
func (m *MockProducer) WriteEvent(_ *file.Event) error {
	return nil
}
func (m *MockProducer) Close() error {
	return nil
}

func TestNewStorage(t *testing.T) {
	mockDB := &MockDB{}
	mockConsumer := &MockConsumer{}
	mockProducer := &MockProducer{}

	tests := []struct {
		name     string
		cfg      config.ServerConfig
		wantType string
	}{
		{
			name:     "Database storage",
			cfg:      config.ServerConfig{DatabaseDsn: "test_dsn"},
			wantType: "*storage.DBStorage",
		},
		{
			name:     "File storage",
			cfg:      config.ServerConfig{FileStoragePath: "/path/to/file"},
			wantType: "*storage.FileStorage",
		},
		{
			name:     "Memory storage",
			cfg:      config.ServerConfig{},
			wantType: "*storage.MemStorage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewStorage(tt.cfg, mockDB, mockConsumer, mockProducer)

			if gotType := fmt.Sprintf("%T", storage); gotType != tt.wantType {
				t.Errorf("Expected type %s, got %s", tt.wantType, gotType)
			}
		})
	}
}

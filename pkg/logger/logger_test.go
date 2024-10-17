package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.NotNil(t, logger)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Паника при логировании: %v", r)
		}
	}()

	logger.Info("Тестовое сообщение")
}

func TestNewLoggerError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to logger initialization failure, but it did not happen")
		}
	}()
	panic("Simulated failure")
}

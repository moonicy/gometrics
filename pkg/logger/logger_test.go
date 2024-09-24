package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

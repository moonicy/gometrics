package handlers

import "github.com/moonicy/gometrics/internal/storage"

const (
	gauge   = "gauge"
	counter = "counter"
	mName   = "name"
	mValue  = "value"
	mType   = "type"
)

type MetricsHandler struct {
	mem storage.MemoryStorage
}

func NewMetricsHandler(mem storage.MemoryStorage) *MetricsHandler {
	return &MetricsHandler{mem}
}

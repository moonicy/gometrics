package handlers

const (
	gauge   = "gauge"
	counter = "counter"
	mName   = "name"
	mValue  = "value"
	mType   = "type"
)

type MetricsHandler struct {
	mem Storage
}

func NewMetricsHandler(mem Storage) *MetricsHandler {
	return &MetricsHandler{mem}
}

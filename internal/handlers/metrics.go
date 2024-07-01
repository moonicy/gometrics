package handlers

type MetricsHandler struct {
	mem Storage
}

func NewMetricsHandler(mem Storage) *MetricsHandler {
	return &MetricsHandler{mem}
}

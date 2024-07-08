package handlers

type Storage interface {
	SetGauge(key string, value float64)
	AddCounter(key string, value int64)
	GetCounter(key string) (value int64, ok bool)
	GetGauge(key string) (value float64, ok bool)
	GetMetrics() (counter map[string]int64, gauge map[string]float64)
}

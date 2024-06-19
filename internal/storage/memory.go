package storage

import "sync"

type MemoryStorage interface {
	SetGauge(key string, value float64)
	AddCounter(key string, value int64)
	GetCounter(key string) (value int64, ok bool)
	GetGauge(key string) (value float64, ok bool)
	GetMetrics() (counter map[string]int64, gauge map[string]float64)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mx      sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) SetGauge(key string, value float64) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.gauge[key] = value
}

func (ms *MemStorage) AddCounter(key string, value int64) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.counter[key] += value
}

func (ms *MemStorage) GetCounter(key string) (value int64, ok bool) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok = ms.counter[key]
	return value, ok
}

func (ms *MemStorage) GetGauge(key string) (value float64, ok bool) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok = ms.gauge[key]
	return value, ok
}

func (ms *MemStorage) GetMetrics() (counter map[string]int64, gauge map[string]float64) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	counter = make(map[string]int64, len(ms.gauge))
	gauge = make(map[string]float64, len(ms.counter))
	for k, v := range ms.gauge {
		gauge[k] = v
	}
	for k, v := range ms.counter {
		counter[k] = v
	}
	return counter, gauge
}

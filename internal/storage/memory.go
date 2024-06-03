package storage

import "sync"

type MemoryStorage interface {
	SetGauge(key string, value float64)
	AddCounter(key string, value int64)
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

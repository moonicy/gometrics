package storage

import (
	"context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

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

func (ms *MemStorage) Init(_ context.Context) error {
	return nil
}

func (ms *MemStorage) SetGauge(_ context.Context, key string, value float64) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.gauge[key] = value
	return nil
}

func (ms *MemStorage) AddCounter(_ context.Context, key string, value int64) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.counter[key] += value
	return nil
}

func (ms *MemStorage) GetCounter(_ context.Context, key string) (value int64, err error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok := ms.counter[key]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

func (ms *MemStorage) GetGauge(_ context.Context, key string) (value float64, err error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok := ms.gauge[key]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

func (ms *MemStorage) GetMetrics(_ context.Context) (counter map[string]int64, gauge map[string]float64, err error) {
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
	return counter, gauge, nil
}

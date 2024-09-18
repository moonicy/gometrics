package storage

import (
	"context"
	"sync"
)

// MemStorage представляет хранилище метрик в памяти.
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mx      sync.Mutex
}

// NewMemStorage создаёт и возвращает новое хранилище метрик в памяти.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// Init инициализирует хранилище метрик в памяти.
func (ms *MemStorage) Init(_ context.Context) error {
	return nil
}

// SetGauge устанавливает значение метрики типа gauge с заданным именем и значением.
func (ms *MemStorage) SetGauge(_ context.Context, key string, value float64) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.gauge[key] = value
	return nil
}

// AddCounter увеличивает значение метрики типа counter с заданным именем на указанное значение.
func (ms *MemStorage) AddCounter(_ context.Context, key string, value int64) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	ms.counter[key] += value
	return nil
}

// GetCounter возвращает текущее значение метрики типа counter с заданным именем.
func (ms *MemStorage) GetCounter(_ context.Context, key string) (value int64, err error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok := ms.counter[key]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

// GetGauge возвращает текущее значение метрики типа gauge с заданным именем.
func (ms *MemStorage) GetGauge(_ context.Context, key string) (value float64, err error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	value, ok := ms.gauge[key]
	if !ok {
		return 0, ErrNotFound
	}
	return value, nil
}

// GetMetrics возвращает все сохранённые метрики типа counter и gauge.
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

// SetMetrics устанавливает переданные метрики в хранилище.
func (ms *MemStorage) SetMetrics(_ context.Context, counter map[string]int64, gauge map[string]float64) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	for k, v := range gauge {
		ms.gauge[k] = v
	}
	for k, v := range counter {
		ms.counter[k] += v
	}
	return nil
}

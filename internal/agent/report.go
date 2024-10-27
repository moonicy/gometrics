package agent

import (
	"sync"
)

// Константы, представляющие названия метрик, используемые агентом.
const (
	Alloc          = "Alloc"
	BuckHashSys    = "BuckHashSys"
	Frees          = "Frees"
	GCCPUFraction  = "GCCPUFraction"
	GCSys          = "GCSys"
	HeapAlloc      = "HeapAlloc"
	HeapIdle       = "HeapIdle"
	HeapInuse      = "HeapInuse"
	HeapObjects    = "HeapObjects"
	HeapReleased   = "HeapReleased"
	HeapSys        = "HeapSys"
	LastGC         = "LastGC"
	Lookups        = "Lookups"
	MCacheInuse    = "MCacheInuse"
	MCacheSys      = "MCacheSys"
	MSpanInuse     = "MSpanInuse"
	MSpanSys       = "MSpanSys"
	Mallocs        = "Mallocs"
	NextGC         = "NextGC"
	NumForcedGC    = "NumForcedGC"
	NumGC          = "NumGC"
	OtherSys       = "OtherSys"
	PauseTotalNs   = "PauseTotalNs"
	StackInuse     = "StackInuse"
	StackSys       = "StackSys"
	Sys            = "Sys"
	TotalAlloc     = "TotalAlloc"
	PollCount      = "PollCount"
	RandomValue    = "RandomValue"
	Gauge          = "gauge"
	Counter        = "counter"
	TotalMemory    = "TotalMemory"
	FreeMemory     = "FreeMemory"
	CPUutilization = "CPUutilization"
)

// Report хранит собранные метрики типа gauge и counter.
type Report struct {
	gauge   map[string]float64 // Map для хранения gauge-метрик.
	counter map[string]int64   // Map для хранения counter-метрик.
	mx      sync.Mutex
}

// NewReport создаёт и возвращает новый экземпляр Report с инициализированными Map.
func NewReport() *Report {
	return &Report{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// GetCommonCount возвращает общее количество собранных метрик.
func (r *Report) GetCommonCount() int {
	return len(r.gauge) + len(r.counter)
}

// SetGauge сохраняет gauge-метрику с указанным именем и значением.
func (r *Report) SetGauge(name string, value float64) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.gauge[name] = value
}

// AddCounter увеличивает counter-метрику с указанным именем на заданное значение.
// Если метрика не существует, она инициализируется.
func (r *Report) AddCounter(name string, value int64) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.counter[name] += value
}

func (r *Report) Clean() {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.gauge = make(map[string]float64)
	r.counter = make(map[string]int64)
}

func (r *Report) GetGauge() map[string]float64 {
	r.mx.Lock()
	defer r.mx.Unlock()
	gauges := make(map[string]float64)
	for k, v := range r.gauge {
		gauges[k] = v
	}
	return gauges
}
func (r *Report) GetCounter() map[string]int64 {
	r.mx.Lock()
	defer r.mx.Unlock()
	counters := make(map[string]int64)
	for k, v := range r.counter {
		counters[k] = v
	}
	return counters
}

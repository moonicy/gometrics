package agent

import (
	"sync"
	"sync/atomic"
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
	Gauge        sync.Map // Map для хранения gauge-метрик.
	Counter      sync.Map // Map для хранения counter-метрик.
	gaugeCount   int      // Количество gauge-метрик.
	counterCount int      // Количество counter-метрик.
}

// NewReport создаёт и возвращает новый экземпляр Report с инициализированными Map.
func NewReport() *Report {
	return &Report{
		Gauge:   sync.Map{},
		Counter: sync.Map{},
	}
}

// GetCommonCount возвращает общее количество собранных метрик.
func (r *Report) GetCommonCount() int {
	return r.gaugeCount + r.counterCount
}

// SetGauge сохраняет gauge-метрику с указанным именем и значением.
func (r *Report) SetGauge(name string, value float64) {
	r.Gauge.Store(name, value)
	r.gaugeCount++
}

// AddCounter увеличивает counter-метрику с указанным именем на заданное значение.
// Если метрика не существует, она инициализируется.
func (r *Report) AddCounter(name string, value int64) {
	if v, ok := r.Counter.Load(name); ok {
		ptr := v.(*int64)
		atomic.AddInt64(ptr, value)
	} else {
		r.Counter.Store(name, &value)
		r.counterCount++
	}
}

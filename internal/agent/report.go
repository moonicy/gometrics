package agent

import (
	"sync"
	"sync/atomic"
)

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

type Report struct {
	Gauge        sync.Map
	Counter      sync.Map
	gaugeCount   int
	counterCount int
}

func NewReport() *Report {
	return &Report{
		Gauge:   sync.Map{},
		Counter: sync.Map{},
	}
}

func (r *Report) GetCommonCount() int {
	return r.gaugeCount + r.counterCount
}

func (r *Report) SetGauge(name string, value float64) {
	r.Gauge.Store(name, value)
	r.gaugeCount++
}

func (r *Report) AddCounter(name string, value int64) {
	if v, ok := r.Counter.Load(name); ok {
		ptr := v.(*int64)
		atomic.AddInt64(ptr, value)
	} else {
		r.Counter.Store(name, value)
		r.counterCount++
	}
}

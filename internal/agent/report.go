package agent

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
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewReport() *Report {
	return &Report{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

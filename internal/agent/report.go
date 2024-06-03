package agent

const (
	Alloc         = "alloc"
	BuckHashSys   = "buckhashsys"
	Frees         = "frees"
	GCCPUFraction = "gccpufraction"
	GCSys         = "gcsys"
	HeapAlloc     = "heapalloc"
	HeapIdle      = "heapidle"
	HeapInuse     = "heapinuse"
	HeapObjects   = "heapobjects"
	HeapReleased  = "heapreleased"
	HeapSys       = "heapsys"
	LastGC        = "lastgc"
	Lookups       = "lookups"
	MCacheInuse   = "mcacheinuse"
	MCacheSys     = "mcachesys"
	MSpanInuse    = "mspaninuse"
	MSpanSys      = "mspansys"
	Mallocs       = "mallocs"
	NextGC        = "nextgc"
	NumForcedGC   = "numforcedgc"
	NumGC         = "numgc"
	OtherSys      = "othersys"
	PauseTotalNs  = "pausetotalns"
	StackInuse    = "stackinuse"
	StackSys      = "stacksys"
	Sys           = "sys"
	TotalAlloc    = "totalalloc"
	PollCount     = "pollcount"
	RandomValue   = "randomvalue"
	Gauge         = "gauge"
	Counter       = "counter"
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

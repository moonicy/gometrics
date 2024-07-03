package file

type Event struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	Timestamp int64
}

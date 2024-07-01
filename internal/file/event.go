package file

import "time"

type Event struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	Timestamp int64
}

func NewEvent(gauge map[string]float64, counter map[string]int64) *Event {
	return &Event{
		Gauge:     gauge,
		Counter:   counter,
		Timestamp: time.Now().Unix(),
	}
}

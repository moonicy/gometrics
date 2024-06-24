package metrics

import (
	"errors"
)

var ErrUnknownMetric = errors.New("unknown metric")
var ErrNotFound = errors.New("not found")
var ErrWrongValue = errors.New("wrong value")

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) Validate() error {

	if m.ID == "" {
		return ErrNotFound
	}
	if m.MType != Gauge && m.MType != Counter {
		return ErrUnknownMetric
	}
	if m.MType == Gauge && m.Value == nil {
		return ErrWrongValue
	}
	if m.MType == Counter && m.Delta == nil {
		return ErrWrongValue
	}
	return nil
}

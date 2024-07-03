package metrics

import (
	"errors"
)

var ErrUnknownMetric = errors.New("unknown metric")
var ErrNotFound = errors.New("not found")
var ErrWrongValue = errors.New("wrong value")

type MetricName struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

type Metric struct {
	MetricName
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (mn MetricName) Validate() error {
	if mn.ID == "" {
		return ErrNotFound
	}
	if mn.MType != Gauge && mn.MType != Counter {
		return ErrUnknownMetric
	}
	return nil
}

func (m Metric) Validate() error {
	if err := m.MetricName.Validate(); err != nil {
		return err
	}
	if m.MType == Gauge && m.Value == nil {
		return ErrWrongValue
	}
	if m.MType == Counter && m.Delta == nil {
		return ErrWrongValue
	}
	return nil
}

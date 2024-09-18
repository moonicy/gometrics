package metrics

import (
	"errors"
)

// ErrUnknownMetric возвращается, когда тип метрики неизвестен.
var ErrUnknownMetric = errors.New("unknown metric")

// ErrNotFound возвращается, когда метрика не найдена.
var ErrNotFound = errors.New("not found")

// ErrWrongValue возвращается, когда значение метрики некорректно.
var ErrWrongValue = errors.New("wrong value")

// MetricName представляет имя и тип метрики.
type MetricName struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

// Metric содержит данные метрики, включая её значение.
type Metric struct {
	MetricName
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// Validate проверяет корректность полей структуры MetricName.
func (mn MetricName) Validate() error {
	if mn.ID == "" {
		return ErrNotFound
	}
	if mn.MType != Gauge && mn.MType != Counter {
		return ErrUnknownMetric
	}
	return nil
}

// Validate проверяет корректность полей структуры Metric.
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

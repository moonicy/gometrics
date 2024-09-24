package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricName_Validate(t *testing.T) {
	tests := []struct {
		name       string
		metricName MetricName
		wantErr    error
	}{
		{
			name:       "Valid gauge metric name",
			metricName: MetricName{ID: "testMetric", MType: "gauge"},
			wantErr:    nil,
		},
		{
			name:       "Valid counter metric name",
			metricName: MetricName{ID: "testMetric", MType: "counter"},
			wantErr:    nil,
		},
		{
			name:       "Empty ID",
			metricName: MetricName{ID: "", MType: "gauge"},
			wantErr:    ErrNotFound,
		},
		{
			name:       "Unknown metric type",
			metricName: MetricName{ID: "testMetric", MType: "unknown"},
			wantErr:    ErrUnknownMetric,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.metricName.Validate()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestMetric_Validate(t *testing.T) {
	gaugeValue := 3.14
	counterValue := int64(42)

	tests := []struct {
		name    string
		metric  Metric
		wantErr error
	}{
		{
			name:    "Valid gauge metric",
			metric:  Metric{MetricName: MetricName{ID: "testMetric", MType: "gauge"}, Value: &gaugeValue},
			wantErr: nil,
		},
		{
			name:    "Valid counter metric",
			metric:  Metric{MetricName: MetricName{ID: "testMetric", MType: "counter"}, Delta: &counterValue},
			wantErr: nil,
		},
		{
			name:    "Gauge metric with nil value",
			metric:  Metric{MetricName: MetricName{ID: "testMetric", MType: "gauge"}, Value: nil},
			wantErr: ErrWrongValue,
		},
		{
			name:    "Counter metric with nil delta",
			metric:  Metric{MetricName: MetricName{ID: "testMetric", MType: "counter"}, Delta: nil},
			wantErr: ErrWrongValue,
		},
		{
			name:    "Unknown metric type",
			metric:  Metric{MetricName: MetricName{ID: "testMetric", MType: "unknown"}},
			wantErr: ErrUnknownMetric,
		},
		{
			name:    "Empty metric ID",
			metric:  Metric{MetricName: MetricName{ID: "", MType: "gauge"}, Value: &gaugeValue},
			wantErr: ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.metric.Validate()
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

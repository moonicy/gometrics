package storage

import (
	"context"
	"fmt"
	"github.com/moonicy/gometrics/internal/agent"
	"testing"
)

func TestMemStorage_AddCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key   string
		value int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wait   string
	}{
		{
			name: "create",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				key:   agent.Alloc,
				value: 11,
			},
			wait: "11",
		},
		{
			name: "update",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: map[string]int64{agent.Alloc: 11},
			},
			args: args{
				key:   agent.Alloc,
				value: 11,
			},
			wait: "22",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			ms.AddCounter(context.Background(), tt.args.key, tt.args.value)

			got := fmt.Sprintf("%d", ms.counter[tt.args.key])
			if tt.wait != got {
				t.Errorf("expected count %s, got %s", tt.wait, got)
			}
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		key   string
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wait   string
	}{
		{
			name: "create",
			fields: fields{
				gauge:   make(map[string]float64),
				counter: make(map[string]int64),
			},
			args: args{
				key:   agent.Alloc,
				value: 11,
			},
			wait: "11.000000",
		},
		{
			name: "update",
			fields: fields{
				gauge:   map[string]float64{agent.Alloc: 11},
				counter: make(map[string]int64),
			},
			args: args{
				key:   agent.Alloc,
				value: 22,
			},
			wait: "22.000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				gauge:   tt.fields.gauge,
				counter: tt.fields.counter,
			}
			ms.SetGauge(context.Background(), tt.args.key, tt.args.value)

			got := fmt.Sprintf("%f", ms.gauge[tt.args.key])
			if tt.wait != got {
				t.Errorf("expected gauge %s, got %s", tt.wait, got)
			}
		})
	}
}

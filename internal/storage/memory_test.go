package storage

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moonicy/gometrics/internal/agent"
)

var ctx = context.Background()

func TestNewMemStorage(t *testing.T) {
	ms := NewMemStorage()
	assert.NotNil(t, ms)
	assert.NotNil(t, ms.gauge)
	assert.NotNil(t, ms.counter)
}

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
			err := ms.AddCounter(ctx, tt.args.key, tt.args.value)
			if err != nil {
				log.Fatal(err)
			}

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
			err := ms.SetGauge(ctx, tt.args.key, tt.args.value)
			if err != nil {
				log.Fatal(err)
			}

			got := fmt.Sprintf("%f", ms.gauge[tt.args.key])
			if tt.wait != got {
				t.Errorf("expected gauge %s, got %s", tt.wait, got)
			}
		})
	}
}

func TestMemStorage_GetGauge_NotFound(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetGauge(ctx, "nonexistentGauge")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestMemStorage_GetCounter(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetCounter(ctx, "nonexistentCounter")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestMemStorage_GetMetrics(t *testing.T) {
	ms := NewMemStorage()

	err := ms.SetGauge(ctx, "testGauge1", 1.23)
	assert.NoError(t, err)
	err = ms.AddCounter(ctx, "testCounter1", 42)
	assert.NoError(t, err)

	counters, gauges, err := ms.GetMetrics(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1.23, gauges["testGauge1"])
	assert.Equal(t, int64(42), counters["testCounter1"])
}

func TestMemStorage_SetMetrics(t *testing.T) {
	ms := NewMemStorage()

	counters := map[string]int64{"counter1": 10}
	gauges := map[string]float64{"gauge1": 2.71}

	err := ms.SetMetrics(ctx, counters, gauges)
	assert.NoError(t, err)

	val, err := ms.GetCounter(ctx, "counter1")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), val)

	valGauge, err := ms.GetGauge(ctx, "gauge1")
	assert.NoError(t, err)
	assert.Equal(t, 2.71, valGauge)
}

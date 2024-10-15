package agent

import (
	"sync/atomic"
	"testing"
)

func TestNewReport(t *testing.T) {
	report := NewReport()

	if report == nil {
		t.Fatalf("NewReport() returned nil")
	}

	if report.GetCommonCount() != 0 {
		t.Errorf("Expected common count to be 0, got %d", report.GetCommonCount())
	}
}

func TestSetGauge(t *testing.T) {
	report := NewReport()
	report.SetGauge(Alloc, 100.5)

	val, ok := report.Gauge.Load(Alloc)
	if !ok {
		t.Fatalf("Expected to find gauge %s", Alloc)
	}

	if val != 100.5 {
		t.Errorf("Expected value 100.5, got %v", val)
	}

	if report.gaugeCount != 1 {
		t.Errorf("Expected gaugeCount to be 1, got %d", report.gaugeCount)
	}
}

func TestAddCounter_NewCounter(t *testing.T) {
	var initialValue int64 = 10
	report := NewReport()
	report.AddCounter(Mallocs, initialValue)

	val, ok := report.Counter.Load(Mallocs)
	if !ok {
		t.Fatalf("Expected to find counter %s", Mallocs)
	}

	if v := val.(*int64); *v != initialValue {
		t.Errorf("Expected value 10, got %d", v)
	}

	if report.counterCount != 1 {
		t.Errorf("Expected counterCount to be 1, got %d", report.counterCount)
	}
}

func TestAddCounter_ExistingCounter(t *testing.T) {
	report := NewReport()
	var initialValue int64 = 10
	report.Counter.Store(Mallocs, &initialValue)

	report.AddCounter(Mallocs, 5)

	val, ok := report.Counter.Load(Mallocs)
	if !ok {
		t.Fatalf("Expected to find counter %s", Mallocs)
	}

	if v := atomic.LoadInt64(val.(*int64)); v != 15 {
		t.Errorf("Expected value 15, got %d", v)
	}

	// counterCount shouldn't increase if we're adding to an existing counter
	if report.counterCount != 0 {
		t.Errorf("Expected counterCount to be 0 for existing counter, got %d", report.counterCount)
	}
}

func TestGetCommonCount(t *testing.T) {
	report := NewReport()

	report.SetGauge(HeapAlloc, 123.4)
	report.AddCounter(Mallocs, 10)

	expected := 2
	if report.GetCommonCount() != expected {
		t.Errorf("Expected common count to be %d, got %d", expected, report.GetCommonCount())
	}
}

package agent

import (
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

	val, ok := report.gauge[Alloc]
	if !ok {
		t.Fatalf("Expected to find gauge %s", Alloc)
	}

	if val != 100.5 {
		t.Errorf("Expected value 100.5, got %v", val)
	}

	if len(report.gauge) != 1 {
		t.Errorf("Expected len(gauge) to be 1, got %d", len(report.gauge))
	}
}

func TestAddCounter_NewCounter(t *testing.T) {
	var initialValue int64 = 10
	report := NewReport()
	report.AddCounter(Mallocs, initialValue)

	val, ok := report.counter[Mallocs]
	if !ok {
		t.Fatalf("Expected to find counter %s", Mallocs)
	}

	if val != initialValue {
		t.Errorf("Expected value 10, got %d", val)
	}

	if len(report.counter) != 1 {
		t.Errorf("Expected len(counter) to be 1, got %d", len(report.counter))
	}
}

func TestAddCounter_ExistingCounter(t *testing.T) {
	report := NewReport()
	var initialValue int64 = 10
	report.counter[Mallocs] = initialValue

	report.AddCounter(Mallocs, 5)

	val, ok := report.counter[Mallocs]
	if !ok {
		t.Fatalf("Expected to find counter %s", Mallocs)
	}

	if val != 15 {
		t.Errorf("Expected value 15, got %d", val)
	}

	// len(counter) shouldn't increase if we're adding to an existing counter
	if len(report.counter) != 1 {
		t.Errorf("Expected len(counter) to be 1 for existing counter, got %d", len(report.counter))
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

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

func TestReportClean(t *testing.T) {
	report := &Report{
		gauge:   map[string]float64{"metric1": 1.23, "metric2": 4.56},
		counter: map[string]int64{"count1": 10, "count2": 20},
	}

	report.Clean()

	if len(report.GetGauge()) != 0 {
		t.Error("expected gauge map to be empty after Clean()")
	}
	if len(report.GetCounter()) != 0 {
		t.Error("expected counter map to be empty after Clean()")
	}
}

func TestReportGetGauge(t *testing.T) {
	report := &Report{
		gauge: map[string]float64{"metric1": 1.23, "metric2": 4.56},
	}

	gauge := report.GetGauge()

	if len(gauge) != 2 {
		t.Errorf("expected gauge map length 2, got %d", len(gauge))
	}
	if gauge["metric1"] != 1.23 {
		t.Errorf("expected gauge[\"metric1\"] to be 1.23, got %f", gauge["metric1"])
	}
	if gauge["metric2"] != 4.56 {
		t.Errorf("expected gauge[\"metric2\"] to be 4.56, got %f", gauge["metric2"])
	}
}

func TestReportGetCounter(t *testing.T) {
	report := &Report{
		counter: map[string]int64{"count1": 10, "count2": 20},
	}

	counter := report.GetCounter()

	if len(counter) != 2 {
		t.Errorf("expected counter map length 2, got %d", len(counter))
	}
	if counter["count1"] != 10 {
		t.Errorf("expected counter[\"count1\"] to be 10, got %d", counter["count1"])
	}
	if counter["count2"] != 20 {
		t.Errorf("expected counter[\"count2\"] to be 20, got %d", counter["count2"])
	}
}

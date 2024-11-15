package workerpool

import (
	"testing"
	"time"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/config"
)

type MockMetricsReader struct {
	ReadCount int
}

func (m *MockMetricsReader) Read(_ *agent.Report) {
	m.ReadCount++
}

func TestRunReadMetrics(t *testing.T) {
	cfg := config.AgentConfig{
		PollInterval: 50 * time.Millisecond,
	}

	reader := &MockMetricsReader{}
	mem := agent.NewReport()
	callbackCalled := false

	callback := func() {
		callbackCalled = true
	}

	stopFunc := RunReadMetrics(cfg, reader, mem, callback)

	time.Sleep(200 * time.Millisecond)

	if reader.ReadCount == 0 {
		t.Error("Expected Read to be called at least once, but it wasn't")
	}

	stopFunc()

	time.Sleep(200 * time.Millisecond)

	if !callbackCalled {
		t.Error("Expected callback to be called on stop, but it wasn't")
	}

	finalReadCount := reader.ReadCount
	time.Sleep(100 * time.Millisecond)

	if reader.ReadCount != finalReadCount {
		t.Errorf("Expected Read to stop being called after stopFunc, but it continued")
	}
}

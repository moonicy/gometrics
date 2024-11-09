package workerpool

import (
	"context"
	"testing"
	"time"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/config"
)

type MockClient struct {
	SendReportCount int
}

func (m *MockClient) SendReport(_ context.Context, _ *agent.Report) {
	m.SendReportCount++
}

func TestRunSendReport(t *testing.T) {
	cfg := config.AgentConfig{
		ReportInterval: 50 * time.Millisecond,
		RateLimit:      1,
	}

	client := &MockClient{}
	mem := agent.NewReport()
	callbackCalled := false

	callback := func() {
		callbackCalled = true
	}

	stopFunc := RunSendReport(cfg, client, mem, callback)

	time.Sleep(200 * time.Millisecond)

	if client.SendReportCount == 0 {
		t.Error("Expected SendReport to be called at least once, but it wasn't")
	}

	stopFunc()

	time.Sleep(200 * time.Millisecond)

	if !callbackCalled {
		t.Error("Expected callback to be called on stop, but it wasn't")
	}

	finalSendReportCount := client.SendReportCount
	time.Sleep(100 * time.Millisecond)

	if client.SendReportCount != finalSendReportCount {
		t.Errorf("Expected SendReport to stop being called after stopFunc, but it continued")
	}
}

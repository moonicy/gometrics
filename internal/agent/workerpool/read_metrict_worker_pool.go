package workerpool

import (
	"time"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/workerpool"
)

func RunReadMetrics(cfg config.AgentConfig, reader *agent.MetricsReader, mem *agent.Report, callback func()) func() {
	var pollInterval = time.Duration(cfg.PollInterval) * time.Second

	rwp := workerpool.NewWorkerPool(1, 1)
	rwp.Run()

	go func() {
		defer callback()
		for {
			rwp.AddJob(func() error {
				reader.Read(mem)
				return nil
			})
			time.Sleep(pollInterval)
		}
	}()
	return rwp.Close
}

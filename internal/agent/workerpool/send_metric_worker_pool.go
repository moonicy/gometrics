package workerpool

import (
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/workerpool"
	"time"
)

func RunSendReport(cfg config.AgentConfig, client *client.Client, mem *agent.Report, callback func()) func() {
	var reportInterval = time.Duration(cfg.ReportInterval) * time.Second

	cwp := workerpool.NewWorkerPool(5, cfg.RateLimit)
	cwp.Run()

	go func() {
		defer callback()
		for {
			cwp.AddJob(func() error {
				client.SendReport(mem)
				return nil
			})
			time.Sleep(reportInterval)
		}
	}()

	return cwp.Close
}

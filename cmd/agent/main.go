package main

import (
	"github.com/moonicy/gometrics/internal/agent"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/workerpool"
	"sync"
	"time"
)

func main() {
	cfg := config.NewAgentConfig()

	cfg.Host = config.ParseURI(cfg.Host)

	var pollInterval = time.Duration(cfg.PollInterval) * time.Second
	var reportInterval = time.Duration(cfg.ReportInterval) * time.Second

	mem := agent.NewReport()
	client := metricsClient.NewClient(cfg.Host, cfg.HashKey)
	reader := agent.NewMetricsReader()
	var wg sync.WaitGroup
	wg.Add(2)

	cwp := workerpool.NewWorkerPool(5, cfg.RateLimit)
	cwp.Run()

	rwp := workerpool.NewWorkerPool(1, 1)
	rwp.Run()

	go func() {
		for {
			rwp.AddJob(func() error {
				reader.Read(mem)
				return nil
			})
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			cwp.AddJob(func() error {
				client.SendReport(mem)
				return nil
			})
			time.Sleep(reportInterval)
		}
	}()
	wg.Wait()
	cwp.Close()
	rwp.Close()
}

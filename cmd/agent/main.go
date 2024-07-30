package main

import (
	"github.com/moonicy/gometrics/internal/agent"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
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

	go func() {
		for {
			reader.Read(mem)
			time.Sleep(pollInterval)
		}
	}()

	go func() {
		for {
			client.SendReport(mem)
			time.Sleep(reportInterval)
		}
	}()
	wg.Wait()
}

package main

import (
	"github.com/moonicy/gometrics/internal/agent"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"sync"
	"time"
)

func main() {
	var pollInterval = 2 * time.Second
	var reportInterval = 10 * time.Second
	mem := agent.NewReport()
	client := metricsClient.NewClient("http://localhost:8080")
	reader := agent.NewMetricsReader()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for {
			time.Sleep(pollInterval)
			reader.Read(mem)
		}
	}()

	go func() {
		for {
			time.Sleep(reportInterval)
			client.SendReport(mem)
		}
	}()
	wg.Wait()
}

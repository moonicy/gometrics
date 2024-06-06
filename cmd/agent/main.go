package main

import (
	"github.com/moonicy/gometrics/internal/agent"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"strconv"
	"sync"
	"time"
)

func main() {
	parseFlags()
	flagRunAddr = parseURI(flagRunAddr)

	intPollI, _ := strconv.Atoi(flagPollInterval)
	intRepI, _ := strconv.Atoi(flagReportInterval)
	var pollInterval = time.Duration(intPollI) * time.Second
	var reportInterval = time.Duration(intRepI) * time.Second

	mem := agent.NewReport()
	client := metricsClient.NewClient(flagRunAddr)
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

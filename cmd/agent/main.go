package main

import (
	"github.com/moonicy/gometrics/internal/agent"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	parseFlags()
	flagRunAddr = http.ParseURI(flagRunAddr)

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

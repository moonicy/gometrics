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
	cfg := parseFlags()
	cfg.Host = http.ParseURI(cfg.Host)

	intPollI, _ := strconv.Atoi(cfg.PollInterval)
	intRepI, _ := strconv.Atoi(cfg.ReportInterval)
	var pollInterval = time.Duration(intPollI) * time.Second
	var reportInterval = time.Duration(intRepI) * time.Second

	mem := agent.NewReport()
	client := metricsClient.NewClient(cfg.Host)
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

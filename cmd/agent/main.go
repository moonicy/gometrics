package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/agent/workerpool"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
)

func main() {
	cfg := config.NewAgentConfig()

	cfg.Host = config.ParseURI(cfg.Host)

	mem := agent.NewReport()
	client := metricsClient.NewClient(cfg.Host, cfg.HashKey)
	reader := agent.NewMetricsReader()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	closeReadFn := workerpool.RunReadMetrics(cfg, reader, mem, wg.Done)
	closeSendFn := workerpool.RunSendReport(cfg, client, mem, wg.Done)

	wg.Wait()
	closeSendFn()
	closeReadFn()
}

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/agent/workerpool"
	metricsClient "github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	cfg := config.NewAgentConfig()

	cfg.Host = config.ParseURI(cfg.Host)

	mem := agent.NewReport()
	var client workerpool.Client
	client = metricsClient.NewClient(cfg.Host, cfg.HashKey, cfg.CryptoKey)
	if cfg.Grpc {
		grpcClient, err := metricsClient.NewGRPCClient()
		if err != nil {
			log.Fatal(err)
		}
		client = grpcClient
	}
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

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-exit

	wg.Wait()
	closeReadFn()
	closeSendFn()
}

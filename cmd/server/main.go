package main

import (
	"context"
	"errors"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/handlers"
	"github.com/moonicy/gometrics/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := config.NewServerConfig()

	sugar := logger.NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	route := handlers.NewRoute(ctx, sugar, cfg)

	var wg sync.WaitGroup
	wg.Add(1)

	sugar.Infow(
		"Starting server",
		"addr", cfg.Host,
	)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}

	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalw(err.Error(), "event", "start server")
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		sugar.Fatalw("Server shutdown error", "error", err)
	}

	cancel()
	wg.Wait()
	time.Sleep(1 * time.Second)
}

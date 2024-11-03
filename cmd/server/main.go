package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
	"github.com/moonicy/gometrics/internal/handlers"
	grpcserver "github.com/moonicy/gometrics/internal/server"
	database2 "github.com/moonicy/gometrics/pkg/database"
	"github.com/moonicy/gometrics/pkg/logger"
	pb "github.com/moonicy/gometrics/proto"
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

	cfg := config.NewServerConfig()

	sugar := logger.NewLogger()
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(sugar)

	ctx, cancel := context.WithCancel(context.Background())

	database, closeFn, err := database2.NewDatabase(sugar, cfg)
	if err != nil {
		sugar.Error(err)
	}
	defer func() {
		if err = closeFn(); err != nil {
			sugar.Error(err)
		}
	}()

	cr := file.NewConsumer(cfg.FileStoragePath)
	pr := file.NewProducer(cfg.FileStoragePath)
	storage := handlers.NewStorage(cfg, database, cr, pr)
	err = storage.Init(ctx)
	if err != nil {
		sugar.Error(err)
	}

	metricsHandler := handlers.NewMetricsHandler(storage, database, sugar)

	route := handlers.NewRoute(metricsHandler, sugar, cfg)

	gserver := grpcserver.NewGRPCServer(storage)

	sugar.Infow(
		"Starting server",
		"addr", cfg.Host,
	)

	AttachProfiler(route)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: route,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalw(err.Error(), "event", "start server")
		}
	}()

	go func() {
		defer wg.Done()
		// определяем порт для сервера
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}
		// создаём gRPC-сервер без зарегистрированной службы
		s := grpc.NewServer()
		// регистрируем сервис
		pb.RegisterMetricsServer(s, gserver)

		fmt.Println("Сервер gRPC начал работу")
		// получаем запрос gRPC
		if err = s.Serve(listen); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Сервер gRPC завершил работу")
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-exit

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err = server.Shutdown(ctxShutdown); err != nil {
		sugar.Fatalw("Server shutdown error", "error", err)
	}

	cancel()
	wg.Wait()
	time.Sleep(1 * time.Second)
}

func AttachProfiler(router *chi.Mux) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}

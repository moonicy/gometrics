package handlers

import (
	"context"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/storage"
	"go.uber.org/zap"
)

type Pingable interface {
	Ping() error
}

type Initable interface {
	Init(ctx context.Context) error
}

type Storage interface {
	SetGauge(ctx context.Context, key string, value float64) error
	AddCounter(ctx context.Context, key string, value int64) error
	GetCounter(ctx context.Context, key string) (value int64, err error)
	GetGauge(ctx context.Context, key string) (value float64, err error)
	GetMetrics(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
	SetMetrics(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
}

type MetricsHandler struct {
	storage Storage
	pinger  Pingable
	logger  *zap.SugaredLogger
}

func NewMetricsHandler(storage Storage, pinger Pingable, logger *zap.SugaredLogger) *MetricsHandler {
	return &MetricsHandler{storage, pinger, logger}
}

func NewStorage(cfg config.ServerConfig, db storage.DB, cr storage.Consumer, pr storage.Producer) interface {
	Storage
	Initable
} {
	if cfg.DatabaseDsn != "" {
		return storage.NewDBStorage(db)
	} else if cfg.FileStoragePath != "" {
		return storage.NewFileStorage(cfg, cr, pr)
	}
	return storage.NewMemStorage()
}

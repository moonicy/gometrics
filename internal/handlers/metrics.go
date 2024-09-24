package handlers

import (
	"context"

	"go.uber.org/zap"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/storage"
)

// Pingable определяет интерфейс с методом Ping для проверки доступности сервиса.
type Pingable interface {
	Ping() error
}

// Initable определяет интерфейс с методом Init для инициализации ресурсов.
type Initable interface {
	Init(ctx context.Context) error
}

// Storage определяет интерфейс для операций с хранилищем метрик.
type Storage interface {
	SetGauge(ctx context.Context, key string, value float64) error
	AddCounter(ctx context.Context, key string, value int64) error
	GetCounter(ctx context.Context, key string) (value int64, err error)
	GetGauge(ctx context.Context, key string) (value float64, err error)
	GetMetrics(ctx context.Context) (counter map[string]int64, gauge map[string]float64, err error)
	SetMetrics(ctx context.Context, counter map[string]int64, gauge map[string]float64) error
}

// MetricsHandler содержит логику обработки метрик и взаимодействия с хранилищем.
type MetricsHandler struct {
	storage Storage
	pinger  Pingable
	logger  *zap.SugaredLogger
}

// NewMetricsHandler создаёт и возвращает новый экземпляр MetricsHandler.
func NewMetricsHandler(storage Storage, pinger Pingable, logger *zap.SugaredLogger) *MetricsHandler {
	return &MetricsHandler{storage, pinger, logger}
}

// NewStorage создаёт и возвращает новое хранилище метрик в зависимости от конфигурации.
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

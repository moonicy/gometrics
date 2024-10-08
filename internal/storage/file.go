package storage

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
)

// Consumer определяет интерфейс для чтения событий из файла.
type Consumer interface {
	Open() error
	ReadEvent() (*file.Event, error)
	Close() error
}

// Producer определяет интерфейс для записи событий в файл.
type Producer interface {
	Open() error
	WriteEvent(event *file.Event) error
	Close() error
}

// FileStorage представляет хранилище метрик с использованием файловой системы.
type FileStorage struct {
	mem      *MemStorage
	consumer Consumer
	producer Producer
	cfg      config.ServerConfig
	mx       sync.Mutex
}

// NewFileStorage создаёт и возвращает новое файловое хранилище метрик.
func NewFileStorage(cfg config.ServerConfig, consumer Consumer, producer Producer) *FileStorage {
	mem := NewMemStorage()
	fs := &FileStorage{
		mem:      mem,
		consumer: consumer,
		producer: producer,
		cfg:      cfg,
	}
	return fs
}

// Init инициализирует файловое хранилище, выполняя восстановление и настройку синхронизации.
func (fs *FileStorage) Init(ctx context.Context) error {
	if fs.cfg.Restore {
		fs.Restore()
	}
	fs.RunSync()
	fs.WaitShutDown(ctx)

	return nil
}

// SetGauge устанавливает значение метрики типа gauge и сохраняет изменения в файл при необходимости.
func (fs *FileStorage) SetGauge(ctx context.Context, key string, value float64) error {
	err := fs.mem.SetGauge(ctx, key, value)
	if err != nil {
		return err
	}
	if fs.cfg.StoreInternal == 0 {
		err := fs.uploadToFile(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddCounter увеличивает значение метрики типа counter и сохраняет изменения в файл при необходимости.
func (fs *FileStorage) AddCounter(ctx context.Context, key string, value int64) error {
	err := fs.mem.AddCounter(ctx, key, value)
	if err != nil {
		return err
	}
	if fs.cfg.StoreInternal == 0 {
		err := fs.uploadToFile(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCounter возвращает текущее значение метрики типа counter.
func (fs *FileStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	return fs.mem.GetCounter(ctx, key)
}

// GetGauge возвращает текущее значение метрики типа gauge.
func (fs *FileStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	return fs.mem.GetGauge(ctx, key)
}

// GetMetrics возвращает все сохранённые метрики типа counter и gauge.
func (fs *FileStorage) GetMetrics(ctx context.Context) (map[string]int64, map[string]float64, error) {
	return fs.mem.GetMetrics(ctx)
}

// SetMetrics сохраняет переданные метрики в памяти.
func (fs *FileStorage) SetMetrics(ctx context.Context, counter map[string]int64, gauge map[string]float64) error {
	return fs.mem.SetMetrics(ctx, counter, gauge)
}

func (fs *FileStorage) uploadToFile(ctx context.Context) error {
	fs.mx.Lock()
	defer fs.mx.Unlock()
	counter, gauge, err := fs.GetMetrics(ctx)
	if err != nil {
		return err
	}

	event := &file.Event{
		Gauge:     gauge,
		Counter:   counter,
		Timestamp: time.Now().Unix(),
	}

	err = fs.producer.Open()
	if err != nil {
		return err
	}
	defer func(producer Producer) {
		err = producer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fs.producer)

	err = fs.producer.WriteEvent(event)
	if err != nil {
		fmt.Println("Error writing event:", err)
	}
	return nil
}

// RunSync запускает периодическую синхронизацию метрик с файлом.
func (fs *FileStorage) RunSync() {
	if fs.cfg.StoreInternal == 0 {
		return
	}
	go func() {
		time.Sleep(time.Duration(fs.cfg.StoreInternal) * time.Second)
		err := fs.uploadToFile(context.Background())
		if err != nil {
			log.Println("Error uploading file:", err)
		}
	}()
}

// Restore восстанавливает метрики из файла при запуске сервера.
func (fs *FileStorage) Restore() {
	err := fs.consumer.Open()
	if err != nil {
		panic(err)
	}
	data, err := fs.consumer.ReadEvent()
	if err != nil {
		panic(err)
	}
	err = fs.consumer.Close()
	if err != nil {
		panic(err)
	}
	if data != nil {
		fs.mem.gauge = data.Gauge
		fs.mem.counter = data.Counter
	}
}

// WaitShutDown ожидает завершения работы сервера и сохраняет метрики в файл.
func (fs *FileStorage) WaitShutDown(ctx context.Context) {
	go func() {
		<-ctx.Done()
		err := fs.uploadToFile(context.Background())
		if err != nil {
			panic(err)
		}
	}()
}

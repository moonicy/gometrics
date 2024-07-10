package storage

import (
	"context"
	"fmt"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
	"log"
	"sync"
	"time"
)

type Consumer interface {
	Open() error
	ReadEvent() (*file.Event, error)
	Close() error
}

type Producer interface {
	Open() error
	WriteEvent(event *file.Event) error
	Close() error
}

type FileStorage struct {
	mem      *MemStorage
	consumer Consumer
	producer Producer
	cfg      config.ServerConfig
	mx       sync.Mutex
}

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

func (fs *FileStorage) Init(ctx context.Context) error {
	if fs.cfg.Restore {
		fs.Restore()
	}
	fs.RunSync()
	fs.WaitShutDown(ctx)

	return nil
}

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

func (fs *FileStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	return fs.mem.GetCounter(ctx, key)
}

func (fs *FileStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	return fs.mem.GetGauge(ctx, key)
}

func (fs *FileStorage) GetMetrics(ctx context.Context) (map[string]int64, map[string]float64, error) {
	return fs.mem.GetMetrics(ctx)
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
	defer fs.producer.Close()

	err = fs.producer.WriteEvent(event)
	if err != nil {
		fmt.Println("Error writing event:", err)
	}
	return nil
}

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

func (fs *FileStorage) WaitShutDown(ctx context.Context) {
	go func() {
		<-ctx.Done()
		err := fs.uploadToFile(context.Background())
		if err != nil {
			panic(err)
		}
	}()
}

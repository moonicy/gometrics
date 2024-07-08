package storage

import (
	"context"
	"fmt"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/internal/file"
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

func (fs *FileStorage) SetGauge(key string, value float64) {
	fs.mem.SetGauge(key, value)
	if fs.cfg.StoreInternal == 0 {
		fs.uploadToFile()
	}
}

func (fs *FileStorage) AddCounter(key string, value int64) {
	fs.mem.AddCounter(key, value)
	if fs.cfg.StoreInternal == 0 {
		fs.uploadToFile()
	}
}

func (fs *FileStorage) GetCounter(key string) (int64, bool) {
	return fs.mem.GetCounter(key)
}

func (fs *FileStorage) GetGauge(key string) (float64, bool) {
	return fs.mem.GetGauge(key)
}

func (fs *FileStorage) GetMetrics() (map[string]int64, map[string]float64) {
	return fs.mem.GetMetrics()
}

func (fs *FileStorage) uploadToFile() {
	fs.mx.Lock()
	defer fs.mx.Unlock()
	counter, gauge := fs.GetMetrics()

	event := &file.Event{
		Gauge:     gauge,
		Counter:   counter,
		Timestamp: time.Now().Unix(),
	}

	err := fs.producer.Open()
	if err != nil {
		panic(err)
	}
	defer fs.producer.Close()

	err = fs.producer.WriteEvent(event)
	if err != nil {
		fmt.Println("Error writing event:", err)
	}
}

func (fs *FileStorage) RunSync() {
	if fs.cfg.StoreInternal == 0 {
		return
	}
	go func() {
		time.Sleep(time.Duration(fs.cfg.StoreInternal) * time.Second)
		fs.uploadToFile()
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
		fs.uploadToFile()
	}()
}

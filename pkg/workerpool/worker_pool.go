// Package workerpool предоставляет функциональность для управления пулом воркеров.
package workerpool

import (
	"log"
	"sync/atomic"
	"time"
)

// WorkerPool представляет пул воркеров для выполнения заданий.
type WorkerPool struct {
	workerCount int
	busyCount   atomic.Int64
	rateLimit   int
	chJob       chan Job
}

// Job представляет функцию задания, которое возвращает ошибку.
type Job func() error

// NewWorkerPool создаёт и возвращает новый пул воркеров с заданным количеством воркеров и ограничением скорости.
func NewWorkerPool(workerCount int, rateLimit int) *WorkerPool {
	chJobs := make(chan Job)
	return &WorkerPool{workerCount: workerCount, chJob: chJobs, rateLimit: rateLimit}
}

// AddJob добавляет новое задание в пул для выполнения.
func (wp *WorkerPool) AddJob(job Job) {
	wp.chJob <- job
}

// Run запускает воркеры и начинает обработку заданий из пула.
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workerCount; i++ {
		go func() {
			for {
				if wp.rateLimit != 0 && int(wp.busyCount.Load()) == wp.rateLimit {
					time.Sleep(1 * time.Millisecond)
					continue
				}
				wp.busyCount.Add(1)
				job, ok := <-wp.chJob
				if !ok {
					break
				}
				err := job()
				if err != nil {
					log.Println("job err: ", err)
					return
				}
				wp.busyCount.Add(-1)
			}
		}()
	}
}

// Close завершает работу пула воркеров и закрывает канал заданий.
func (wp *WorkerPool) Close() {
	close(wp.chJob)
}

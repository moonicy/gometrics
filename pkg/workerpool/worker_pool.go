package workerpool

import (
	"log"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	workerCount int
	busyCount   atomic.Int64
	rateLimit   int
	chJob       chan Job
}

type Job func() error

func NewWorkerPool(workerCount int, rateLimit int) *WorkerPool {
	chJobs := make(chan Job)
	return &WorkerPool{workerCount: workerCount, chJob: chJobs, rateLimit: rateLimit}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.chJob <- job
}

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

func (wp *WorkerPool) Close() {
	close(wp.chJob)
}

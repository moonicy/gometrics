package workerpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerPool(t *testing.T) {
	wp := NewWorkerPool(5, 10)
	assert.Equal(t, 5, wp.workerCount)
	assert.Equal(t, 10, wp.rateLimit)
	assert.NotNil(t, wp.chJob)
}

func TestWorkerPool_AddJob(t *testing.T) {
	wp := NewWorkerPool(1, 0)
	defer wp.Close()
	wp.Run()

	var executed int32
	job := func() error {
		atomic.AddInt32(&executed, 1)
		return nil
	}

	wp.AddJob(job)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(1), executed)
}

func TestWorkerPool_Close(t *testing.T) {
	wp := NewWorkerPool(1, 0)
	wp.Close()

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r.(error).Error(), "send on closed channel")
		} else {
			t.Errorf("Expected panic when adding job to closed channel")
		}
	}()

	wp.AddJob(func() error {
		return nil
	})
}

func TestWorkerPool_Run(t *testing.T) {
	wp := NewWorkerPool(3, 0)
	defer wp.Close()
	wp.Run()

	var wg sync.WaitGroup
	var counter int32
	totalJobs := 10

	for i := 0; i < totalJobs; i++ {
		wg.Add(1)
		wp.AddJob(func() error {
			defer wg.Done()
			atomic.AddInt32(&counter, 1)
			return nil
		})
	}

	wg.Wait()
	assert.Equal(t, int32(totalJobs), counter)
}

func TestWorkerPool_RateLimit(t *testing.T) {
	rateLimit := 2
	wp := NewWorkerPool(5, rateLimit)
	defer wp.Close()
	wp.Run()

	var maxBusyCount int64
	var wg sync.WaitGroup
	totalJobs := 10

	for i := 0; i < totalJobs; i++ {
		wg.Add(1)
		wp.AddJob(func() error {
			defer wg.Done()
			currentBusy := wp.busyCount.Load()
			if currentBusy > atomic.LoadInt64(&maxBusyCount) {
				atomic.StoreInt64(&maxBusyCount, currentBusy)
			}
			time.Sleep(50 * time.Millisecond)
			return nil
		})
	}

	wg.Wait()
	assert.Equal(t, int64(rateLimit), maxBusyCount)
}

func TestWorkerPool_JobError(t *testing.T) {
	wp := NewWorkerPool(1, 0)
	defer wp.Close()
	wp.Run()

	var executed int32
	job := func() error {
		atomic.AddInt32(&executed, 1)
		return errors.New("job error")
	}

	wp.AddJob(job)
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, int32(1), executed)
}

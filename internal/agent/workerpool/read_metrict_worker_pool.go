package workerpool

import (
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/workerpool"
	"time"
)

// RunReadMetrics запускает горутину для периодического чтения метрик и их сохранения в Report.
// Она использует пул воркеров для выполнения задач чтения метрик с заданным интервалом.
// При завершении возвращает функцию для корректного закрытия пула воркеров.
func RunReadMetrics(cfg config.AgentConfig, reader *agent.MetricsReader, mem *agent.Report, callback func()) func() {
	rwp := workerpool.NewWorkerPool(1, 1)
	rwp.Run()

	stop := make(chan struct{})

	go func() {
		defer callback()
		for {
			select {
			case <-stop:
				return
			default:
				rwp.AddJob(func() error {
					reader.Read(mem)
					return nil
				})
				time.Sleep(cfg.PollInterval)
			}
		}
	}()
	return func() {
		close(stop)
		ch := make(chan struct{})
		rwp.AddJob(func() error {
			reader.Read(mem)
			close(ch)
			return nil
		})
		<-ch
		rwp.Close()
	}
}

package workerpool

import (
	"github.com/moonicy/gometrics/internal/agent"
	"github.com/moonicy/gometrics/internal/client"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/workerpool"
	"time"
)

// RunSendReport запускает горутину для периодической отправки отчета с метриками на сервер.
// Она использует пул воркеров для управления количеством одновременных задач и ограничивает скорость отправки.
// При завершении возвращает функцию, которую можно вызвать для корректного закрытия пула воркеров.
func RunSendReport(cfg config.AgentConfig, client *client.Client, mem *agent.Report, callback func()) func() {
	cwp := workerpool.NewWorkerPool(5, cfg.RateLimit)
	cwp.Run()

	stop := make(chan struct{})

	go func() {

		defer callback()
		for {
			select {
			case <-stop:
				return
			default:
				cwp.AddJob(func() error {
					client.SendReport(mem)
					return nil
				})
				time.Sleep(cfg.ReportInterval)
			}
		}
	}()

	return func() {
		close(stop)
		ch := make(chan struct{})
		cwp.AddJob(func() error {
			client.SendReport(mem)
			ch <- struct{}{}
			return nil
		})
		<-ch
		cwp.Close()
	}
}

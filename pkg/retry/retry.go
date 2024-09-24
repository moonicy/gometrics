package retry

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// RetryHandle выполняет функцию do и повторяет попытку в случае возникновения RetryableError.
// Он повторяет попытки с увеличивающимися интервалами ожидания.
// Если ошибка не является RetryableError, она возвращается немедленно.
// Если после всех попыток ошибка не устранена, возвращает ошибку с сообщением о тайм-ауте.
func RetryHandle(do func() error) error {
	var err error
	for _, waitTime := range []time.Duration{0, 1, 3, 5} {
		time.Sleep(waitTime * time.Second)
		log.Println("wt: ", waitTime)
		err = do()
		if err == nil {
			return nil
		}
		var re *RetryableError
		if !errors.As(err, &re) {
			return err
		}
	}
	return fmt.Errorf("retry timed out, %w", err)
}

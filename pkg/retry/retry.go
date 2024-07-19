package retry

import (
	"errors"
	"fmt"
	"log"
	"time"
)

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

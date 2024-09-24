package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/moonicy/gometrics/pkg/retry"
)

// Consumer читает события из файла.
type Consumer struct {
	file     *os.File
	filename string
	scanner  *bufio.Scanner
}

// NewConsumer создаёт и возвращает новый Consumer для указанного файла.
func NewConsumer(filename string) *Consumer {
	return &Consumer{filename: filename}
}

// Open открывает файл для чтения и инициализирует сканер.
// В случае ошибок доступа выполняет повторные попытки.
func (c *Consumer) Open() error {
	var file *os.File
	var err error
	err = retry.RetryHandle(func() error {
		file, err = os.OpenFile(c.filename, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			if os.IsPermission(err) {
				return retry.NewRetryableError(err.Error())
			}
		}
		return err
	})
	if err != nil {
		return err
	}
	c.file = file
	c.scanner = bufio.NewScanner(file)
	return nil
}

// ReadEvent читает следующее событие из файла.
// Возвращает Event или ошибку при неудачном чтении.
func (c *Consumer) ReadEvent() (*Event, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()

	event := Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Close закрывает файл и освобождает связанные ресурсы.
func (c *Consumer) Close() error {
	defer func() {
		c.file = nil
		c.scanner = nil
	}()
	return c.file.Close()
}

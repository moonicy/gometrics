package file

import (
	"bufio"
	"encoding/json"
	"github.com/moonicy/gometrics/pkg/retry"
	"os"
)

type Consumer struct {
	file     *os.File
	filename string
	scanner  *bufio.Scanner
}

func NewConsumer(filename string) *Consumer {
	return &Consumer{filename: filename}
}

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

func (c *Consumer) Close() error {
	defer func() {
		c.file = nil
		c.scanner = nil
	}()
	return c.file.Close()
}

package file

import (
	"bufio"
	"encoding/json"
	"github.com/moonicy/gometrics/pkg/retry"
	"os"
)

type Producer struct {
	file     *os.File
	filename string
	writer   *bufio.Writer
}

func NewProducer(filename string) *Producer {
	return &Producer{filename: filename}
}

func (p *Producer) Open() error {
	var file *os.File
	var err error
	err = retry.RetryHandle(func() error {
		file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE, 0666)
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
	p.file = file
	p.writer = bufio.NewWriter(file)
	return nil
}

func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = p.file.Seek(0, 0)
	if err != nil {
		return err
	}
	if _, err = p.writer.Write(data); err != nil {
		return err
	}
	return p.writer.Flush()
}

func (p *Producer) Close() error {
	defer func() {
		p.file = nil
		p.writer = nil
	}()
	return p.file.Close()
}

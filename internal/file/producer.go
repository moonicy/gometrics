package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/moonicy/gometrics/pkg/retry"
)

// Producer представляет структуру для записи событий в файл.
type Producer struct {
	file     *os.File
	filename string
	writer   *bufio.Writer
}

// NewProducer создаёт и возвращает новый экземпляр Producer для указанного файла.
func NewProducer(filename string) *Producer {
	return &Producer{filename: filename}
}

// Open открывает файл для записи и инициализирует буферизированный writer.
// В случае ошибок доступа выполняет повторные попытки.
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

// WriteEvent записывает событие Event в файл.
// При записи перезаписывает содержимое файла.
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

// Close закрывает файл и освобождает связанные ресурсы.
func (p *Producer) Close() error {
	defer func() {
		p.file = nil
		p.writer = nil
	}()
	return p.file.Close()
}

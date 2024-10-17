// Package logger предоставляет функцию для создания нового логгера с использованием библиотеки zap.
package logger

import (
	"go.uber.org/zap"
	"log"
)

// NewLogger создаёт и возвращает новый zap.SugaredLogger для логирования.
func NewLogger() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	return *logger.Sugar()
}

// Package logger предоставляет функцию для создания нового логгера с использованием библиотеки zap.
package logger

import "go.uber.org/zap"

// NewLogger создаёт и возвращает новый zap.SugaredLogger для логирования.
func NewLogger() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return *logger.Sugar()
}

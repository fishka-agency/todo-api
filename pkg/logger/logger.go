package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger создает новый экземпляр логгера
func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()                          // конфигурация логгера
	config.EncoderConfig.TimeKey = "timestamp"                   // ключ времени
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // формат времени
	config.EncoderConfig.StacktraceKey = ""                      // Отключаем вывод стектрейса

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

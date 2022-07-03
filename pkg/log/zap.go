package log

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const logLevel zapcore.Level = zapcore.DebugLevel

type zapLogger struct {
	logger zap.SugaredLogger
}

func NewZapLogger() Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.TimeKey = "time"
	config.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05")
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(colorable.NewColorableStdout()),
		logLevel,
	))
	return logger.Sugar()
}

func (l *zapLogger) Close() {
	l.logger.Sync()
}

func (l *zapLogger) Fatalf(message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)
}

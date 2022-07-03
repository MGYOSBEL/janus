package log

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(logLevel string) Logger {
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.DebugLevel
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.TimeKey = "time"
	config.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout("02/01/2006 15:04:05")
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(colorable.NewColorableStdout()),
		level,
	))
	return &zapLogger{
		logger: logger.Sugar(),
	}
}

func (l *zapLogger) Close() {
	l.logger.Sync()
}

func (l *zapLogger) Panicf(message string, args ...interface{}) {
	l.logger.Panicf(message, args...)
}

func (l *zapLogger) Fatalf(message string, args ...interface{}) {
	l.logger.Fatalf(message, args...)
}

func (l *zapLogger) Errorf(message string, args ...interface{}) {
	l.logger.Errorf(message, args...)
}

func (l *zapLogger) Warnf(message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}

func (l *zapLogger) Infof(message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *zapLogger) Debugf(message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

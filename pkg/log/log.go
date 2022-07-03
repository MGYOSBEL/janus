package log

type Logger interface {
	Close()
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
	Panicf(message string, args ...interface{})
}

type LoggerType int32

const (
	ZAPLOGGER LoggerType = iota
)

func New(loggerType string, level string) Logger {
	lt := ParseType(loggerType)
	switch lt {
	case ZAPLOGGER:
		return NewZapLogger(level)
	default:
		return NewZapLogger(level)
	}
}

func ParseType(loggerType string) LoggerType {
	switch loggerType {
	case "ZapLogger":
		return ZAPLOGGER
	default:
		return ZAPLOGGER
	}

}

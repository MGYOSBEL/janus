package log

type Logger interface {
	Close()
	Infof(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
	Panicf(message string, args ...interface{})
}

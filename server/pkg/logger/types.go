package logger

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(template string, args ...interface{})
	Printf(format string, args ...any)
}

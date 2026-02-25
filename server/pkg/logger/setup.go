package logger

import "github.com/sirupsen/logrus"

var _logger *logrus.Logger

func Setup(level string) {
	_logger = NewLogger(level)
}

func GetLogger() *logrus.Logger {
	return _logger
}

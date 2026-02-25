package queue

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Entry
}

func NewLogger(log *logrus.Entry) *Logger {
	return &Logger{
		log: log,
	}
}

func (l *Logger) Print(v ...interface{}) {
	l.log.Print(v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.log.Println(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log.Fatalf(format, v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.log.Fatalln(v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.log.Panic(v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.log.Panicf(format, v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.log.Panicln(v...)
}

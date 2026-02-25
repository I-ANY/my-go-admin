package logger

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	//_ "github.com/zput/zxcTool/ztLog/zt_formatter"
	"os"
	"path"
	"runtime"
	"time"
)

func NewLogger(level string) *logrus.Logger {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalf("parse logger level error: %s, level: %s", err, level)
	}
	var l = &logrus.Logger{
		Out:          os.Stdout,
		Level:        logrusLevel,
		ReportCaller: true,
	}
	l.SetFormatter(&nested.Formatter{
		TimestampFormat: time.DateTime,
		NoColors:        false,
		NoFieldsColors:  true,
		ShowFullLevel:   false,
		CallerFirst:     false,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			fileName := path.Base(frame.File)
			return fmt.Sprintf(" [%v:%v]", fileName, frame.Line)
		},
	})
	return l
}

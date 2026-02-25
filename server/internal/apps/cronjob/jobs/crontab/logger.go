package crontab

import (
	"biz-auto-api/internal/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type PersistLogger struct {
	Logger       *logrus.Entry
	ExecRecordId int64
	DB           *gorm.DB
	keys         []string
}

func NewPersistLogger(logger *logrus.Entry, execRecordId int64, db *gorm.DB, keys ...string) *PersistLogger {
	//if !tools.InSlice("schedule_node", keys) {
	//	keys = append(keys, "schedule_node")
	//}
	return &PersistLogger{
		Logger:       logger,
		ExecRecordId: execRecordId,
		DB:           db,
		keys:         keys,
	}
}

func (l *PersistLogger) getEntryMessage(keys ...string) string {
	msg := ""
	logEntries := l.Logger.Data
	for _, key := range keys {
		if value, ok := logEntries[key]; ok {
			v, _ := value.(string)
			msg += fmt.Sprintf(" [%v:%v]", key, v)
		}
	}

	//keys := make([]string, 0)
	//logEntries := l.Logger.Data
	//for key, _ := range logEntries {
	//	keys = append(keys, key)
	//}
	//// 把key做排序
	//sort.Slice(keys, func(i, j int) bool {
	//	return keys[i] < keys[j]
	//})
	//for _, key := range keys {
	//	value := logEntries[key]
	//	v, _ := value.(string)
	//	msg += fmt.Sprintf(" [%v:%v]", key, v)
	//}

	return strings.TrimSpace(msg)
}

// 持久化日志
func (l *PersistLogger) persistence(msg string, level string) error {
	entryMessage := l.getEntryMessage(l.keys...) // "invoke_target"
	mergeMessage := fmt.Sprintf("%v %v", entryMessage, msg)
	return l.DB.Session(
		&gorm.Session{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	).Create(&models.CjJobExecLog{
		LogTime:      time.Now(),
		Level:        level,
		Message:      mergeMessage,
		ExecRecordId: l.ExecRecordId,
	}).Error
}

func (l *PersistLogger) Debugf(template string, args ...interface{}) {
	l.Logger.Debugf(template, args...)
	_ = l.persistence(fmt.Sprintf(template, args...), models.LogLevelDebug)
}

func (l *PersistLogger) Infof(template string, args ...interface{}) {
	l.Logger.Infof(template, args...)
	_ = l.persistence(fmt.Sprintf(template, args...), models.LogLevelInfo)
}

func (l *PersistLogger) Warnf(template string, args ...interface{}) {
	l.Logger.Warnf(template, args...)
	_ = l.persistence(fmt.Sprintf(template, args...), models.LogLevelWarn)
}

func (l *PersistLogger) Printf(template string, args ...any) {
	l.Logger.Infof(template, args...)
	_ = l.persistence(fmt.Sprintf(template, args...), models.LogLevelInfo)
}
func (l *PersistLogger) Errorf(template string, args ...interface{}) {
	l.Logger.Errorf(template, args...)
	_ = l.persistence(fmt.Sprintf(template, args...), models.LogLevelError)
}

type CronLogger struct {
	Logger *logrus.Entry
}

func NewCronLogger(logger *logrus.Entry) *CronLogger {
	return &CronLogger{
		Logger: logger,
	}
}

func (l *CronLogger) Printf(template string, args ...interface{}) {
	l.Logger.Infof(template, args...)
}
func (l *CronLogger) Info(template string, keysAndValues ...interface{}) {
	fields := logrus.Fields{}
	key := ""
	for i, e := range keysAndValues {
		if i%2 == 0 { // 索引为偶数则为key，索引为基数则为value
			key, _ = e.(string)
		} else {
			fields[key] = e
			key = ""
		}
	}
	l.Logger.WithFields(fields).Infof(template)
}
func (l *CronLogger) Error(err error, info string, keysAndValues ...interface{}) {
	fields := logrus.Fields{}
	key := ""
	for i, e := range keysAndValues {
		if i%2 == 0 { // 索引为偶数则为key，索引为基数则为value
			key, _ = e.(string)
		} else {
			fields[key] = e
			key = ""
		}
	}
	l.Logger.WithFields(fields).Errorf(fmt.Sprintf("%v: %s", info, err), keysAndValues...)
}

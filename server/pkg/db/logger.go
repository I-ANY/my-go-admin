package db

import (
	"biz-auto-api/pkg/consts"
	plogger "biz-auto-api/pkg/logger"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type GormLogger struct {
	Logger        *logrus.Logger
	SlowThreshold time.Duration
}

func NewGormLogger(slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		Logger:        plogger.GetLogger(),
		SlowThreshold: slowThreshold,
	}
}
func NewGormDefaultLogger() *GormLogger {
	return NewGormLogger(1 * time.Second)
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	requestId := ctx.Value(consts.RequestIdKey)
	log := g.Logger.WithField("engine", "mysql")
	if requestId != nil {
		log = log.WithField(strings.ToLower(consts.RequestIdKey), requestId.(string))
	}
	log.Infof(s, i...)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	requestId := ctx.Value(consts.RequestIdKey)
	log := g.Logger.WithField("engine", "mysql")
	if requestId != nil {
		log = log.WithField(strings.ToLower(consts.RequestIdKey), requestId.(string))
	}
	log.Warnf(s, i...)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	requestId := ctx.Value(consts.RequestIdKey)
	log := g.Logger.WithField("engine", "mysql")
	if requestId != nil {
		log = log.WithField(strings.ToLower(consts.RequestIdKey), requestId.(string))
	}
	log.Errorf(s, i...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	requestId := ctx.Value(consts.RequestIdKey)
	sql, affected := fc()
	elapsed := time.Since(begin)
	tms := float64(elapsed.Nanoseconds()) / 1e6
	log := g.Logger.WithField("engine", "mysql")
	if requestId != nil {
		log = log.WithField(strings.ToLower(consts.RequestIdKey), requestId.(string))
	}
	if err != nil {
		log.Errorf("[File] %v [%.3fms] [Affected=%v] %v[Error] %v [SQL] %v%v", utils.FileWithLineNum(), tms, affected, logger.Red, err, sql, logger.Reset)
	} else {
		if elapsed > g.SlowThreshold && g.SlowThreshold != 0 {
			log.Warnf("[File] %s %v[%.3fms,SLOW SQL >= %v]%v %v[Affected=%d]%v [SQL] %s", utils.FileWithLineNum(), logger.Red, tms, g.SlowThreshold, logger.Reset, logger.Blue, affected, logger.Reset, sql)
		} else {
			log.Debugf("[File] %s %v[%.3fms]%v %v[Affected=%d]%v [SQL] %s", utils.FileWithLineNum(), logger.Green, tms, logger.Reset, logger.Blue, affected, logger.Reset, sql)
		}
	}
}

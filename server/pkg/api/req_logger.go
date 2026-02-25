package api

import (
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/logger"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/gin-gonic/gin"
)

type loggerKey struct{}

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) *logrus.Entry {
	var log *logrus.Entry
	l, ok := c.Get(consts.LoggerKey)
	if ok {
		ok = false
		log, ok = l.(*logrus.Entry)
		if ok {
			return log
		}
	}
	//如果没有在上下文中放入logger
	requestId := GenerateMsgIDFromContext(c)
	log = logger.GetLogger().WithField(strings.ToLower(consts.RequestIdKey), requestId)
	return log
}

// SetRequestLogger 设置logger中间件
func SetRequestLogger(c *gin.Context) {
	requestId := GenerateMsgIDFromContext(c)
	log := logger.GetLogger().WithField(strings.ToLower(consts.RequestIdKey), requestId)
	c.Set(consts.LoggerKey, log)
}

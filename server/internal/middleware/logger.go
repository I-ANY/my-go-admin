package middleware

import (
	"biz-auto-api/pkg/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func SetContextLogger() func(c *gin.Context) {
	return func(c *gin.Context) {
		api.SetRequestLogger(c)
	}
}

// PrintRequestInfo 打印请求信息
func PrintRequestInfo() func(c *gin.Context) {
	return func(c *gin.Context) {
		l := api.GetRequestLogger(c)
		ip := c.ClientIP()
		begin := time.Now()
		c.Next()
		elapsed := time.Since(begin)
		msg := "uri: " + strings.ReplaceAll(c.Request.RequestURI, "%", "%%")
		tms := float64(elapsed.Nanoseconds()) / 1e6
		l = l.WithField("status", c.Writer.Status()).
			WithField("response-time", fmt.Sprintf("%vms", tms)).
			WithField("method", c.Request.Method).WithField("client-ip", ip)

		if c.Writer.Status() >= 500 {
			l.Errorf(msg)
		} else if c.Writer.Status() >= 400 {
			l.Warn(msg)
		} else {
			l.Infof(msg)
		}
	}
}

package middleware

import (
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetContextRequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		requestId := c.GetHeader(consts.RequestIdKey)
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set(consts.RequestIdKey, requestId)
	}
}

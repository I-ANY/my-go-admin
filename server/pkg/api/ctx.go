package api

import (
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
)

// GenerateMsgIDFromContext 生成msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	if requestId, exists := c.Get(consts.RequestIdKey); exists {
		return requestId.(string)
	}
	return ""
}

func GetUserIdFromContext(c *gin.Context) int64 {
	if userKey, exists := c.Get(consts.UserIdKey); exists {
		return userKey.(int64)
	}
	return 0
}

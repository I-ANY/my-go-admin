package service

import (
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type Service struct {
	DB     *gorm.DB
	EcdnDB *gorm.DB
	Log    *logrus.Entry
	Ctx    *gin.Context
	Redis  *redis.Client
	CK     *clickhouse.CK
}

func (c *Service) GetCurrentUserId() int64 {
	if userId, exists := c.Ctx.Get(consts.UserIdKey); exists {
		return userId.(int64)
	}
	return 0
}

func (c *Service) GetRequestId() string {
	if requestId, exists := c.Ctx.Get(consts.RequestIdKey); exists {
		return requestId.(string)
	} else {
		return ""
	}
}

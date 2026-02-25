package common

import (
	"biz-auto-api/pkg/clickhouse"
	pkgdb "biz-auto-api/pkg/db"
	pkgredis "biz-auto-api/pkg/redis"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	healthzStatusOK       = "OK"
	healthzStatusErr      = "ERROR"
	healthzTimeout        = "TIMEOUT"
	healthzTooManyRequest = "TOO_MANY_REQUEST"
	healthzSkipped        = "SKIPPED"
)

const (
	HealthzApi = "/healthz"
	//并发限制
	ConcurrencyLimit = 15
)

type healthzRes struct {
	Status     string `json:"status"`
	Database   string `json:"database,omitempty"`
	Redis      string `json:"redis,omitempty"`
	Clickhouse string `json:"clickhouse,omitempty"`
}

func GetHealthzFunc() gin.HandlerFunc {
	ch := make(chan struct{}, ConcurrencyLimit)
	return func(c *gin.Context) {
		// 并发控制
		select {
		case ch <- struct{}{}:
			defer func() { <-ch }()
		case <-time.After(3 * time.Second):
			res := healthzRes{
				Status:     healthzTooManyRequest,
				Database:   healthzSkipped,
				Redis:      healthzSkipped,
				Clickhouse: healthzSkipped,
			}
			c.JSON(http.StatusOK, res)
			return
		}
		res := healthzRes{
			Status: healthzStatusOK,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db := pkgdb.GetDB()
		if db != nil {
			err := db.WithContext(ctx).Raw("select 1").Error
			if err != nil {
				res.Database = healthzStatusErr
			} else {
				res.Database = healthzStatusOK
			}
		}
		redis := pkgredis.GetClient()
		if redis != nil {
			_, err := redis.Ping(ctx).Result()
			if err != nil {
				res.Redis = healthzStatusErr
			} else {
				res.Redis = healthzStatusOK
			}
		}
		ck := clickhouse.GetCK()
		if ck != nil {
			err := ck.WithContext(ctx).Raw("select 1").Error
			if err != nil {
				res.Clickhouse = healthzStatusErr
			} else {
				res.Clickhouse = healthzStatusOK
			}
		}
		c.JSON(http.StatusOK, res)
	}
}

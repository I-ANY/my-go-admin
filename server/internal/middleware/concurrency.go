package middleware

import (
	"biz-auto-api/pkg/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ConcurrencyLimit 并发限制中间件
func ConcurrencyLimit(concurrency int) gin.HandlerFunc {
	ch := make(chan struct{}, concurrency)
	return func(c *gin.Context) {
		select {
		case ch <- struct{}{}:
			defer func() {
				<-ch
			}()
			c.Next()
		default:
			// 429 Too Many Requests
			c.AbortWithStatusJSON(http.StatusTooManyRequests, &api.Response{
				Code:      http.StatusTooManyRequests,
				Msg:       "Too many requests",
				RequestId: api.GenerateMsgIDFromContext(c),
			})
		}
	}
}

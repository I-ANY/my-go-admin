package middleware

import (
	"biz-auto-api/pkg/api"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func Recover() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				var err error
				if e, ok := r.(error); ok {
					err = errors.WithStack(e)
				} else {
					err = errors.Errorf("%v", r)
				}
				l := api.GetRequestLogger(c)
				l.Errorf("%+v", err)
				api.Error(c, http.StatusInternalServerError, fmt.Sprintf("服务器异常：%s", err))
			}
		}()
		c.Next()
	}
}

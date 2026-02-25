package router

import (
	"biz-auto-api/internal/common"
	"biz-auto-api/internal/middleware"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func InitRouter() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		logger.GetLogger().Debugf("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	setRouter(engine)
	return engine
}

func setRouter(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		api.Error(c, http.StatusNotFound, "Resource not found")
		c.Abort()
	})
	engine.Use(
		middleware.SetContextRequestId(), // 生成RequestId
		middleware.SetContextLogger(),    // 生成请求的logger
		middleware.PrintRequestInfo(),    // 打印请求日志
		middleware.Recover(),             // 全局异常处理
		middleware.NoCache(),             // 缓存设置
		middleware.Options(),             // 跨域设置
		middleware.Secure(),              // 附加安全性
		middleware.Auth(),                // 认证
		middleware.AddOperaLog(),         // 添加操作日志
		middleware.PermissionCheck(),     // 鉴权
	)

	v1 := engine.Group(consts.ApiV1Prefix)
	SetProxyRouter(v1)
	SetBusinessRouter(v1)

	err := common.AddOrUpdateApis(engine)
	if err != nil {
		logger.GetLogger().Fatalf("%+v", errors.WithMessage(err, "add or update apis failed"))
	}
}

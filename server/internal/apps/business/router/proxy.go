package router

import (
	"biz-auto-api/internal/apps/business/apis"
	"github.com/gin-gonic/gin"
)

func SetProxyRouter(v1 *gin.RouterGroup) {
	flaskApiGroup := v1.Group("/flask-api")
	flaskApiForward := apis.FlaskApiForward{}
	{
		flaskApiGroup.Any("/proxy/*targetUrl", flaskApiForward.Forward)
	}

	proxyGroup := v1.Group("/proxy")
	ecdnProxyGroup := proxyGroup.Group("/ecdn")
	{
		ecdnProxyApi := apis.EcdnProxy{}
		ecdnProxyGroup.POST("/server/dispatchParams/report", ecdnProxyApi.ReportDispatchParamsReport)
	}

}

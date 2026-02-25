package router

import (
	"biz-auto-api/internal/apps/business/apis"
	"biz-auto-api/internal/common"

	"github.com/gin-gonic/gin"
)

func SetBusinessRouter(v1 *gin.RouterGroup) {
	businessApiGroup := v1.Group("/business")
	{
		businessApiGroup.GET(common.HealthzApi, common.GetHealthzFunc())
	}
	{
		demoApiGroup := businessApiGroup.Group("/demo")
		businessDemoApi := apis.DemoApi{}
		{
			demoApiGroup.POST("/list", businessDemoApi.GetDemo)
		}

	}
}

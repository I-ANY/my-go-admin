package router

import (
	"biz-auto-api/internal/apps/cronjob/apis"
	"biz-auto-api/internal/common"
	"github.com/gin-gonic/gin"
)

func SetCronjobRouter(v1 *gin.RouterGroup) {
	systemGroup := v1.Group("/cronjob")
	{
		systemGroup.GET(common.HealthzApi, common.GetHealthzFunc())
	}
	jobApi := apis.JobApi{}
	{
		systemGroup.GET("/job/list", jobApi.GetJobList)
		systemGroup.POST("/job", jobApi.AddJob)
		systemGroup.PUT("/job/:id", jobApi.UpdateJob)
		systemGroup.GET("/job/:id", jobApi.GetJob)
		systemGroup.DELETE("/job/:id", jobApi.DeleteJob)
		systemGroup.POST("/job/:id/exec", jobApi.ExecJob)
	}

	jobExecRecordApi := apis.JobExecRecord{}
	{
		systemGroup.GET("/job/exec-record/list", jobExecRecordApi.GetJobExecRecordList)
		systemGroup.GET("/job/exec/:id/logs", jobExecRecordApi.GetJobExecLog)
	}

}

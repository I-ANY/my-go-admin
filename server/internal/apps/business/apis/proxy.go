package apis

import (
	"biz-auto-api/internal/apps/business/service"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
)

type FlaskApiForward struct {
	api.Api
}

func (a FlaskApiForward) Forward(c *gin.Context) {
	s := service.FlaskApiForward{}
	err := s.Forward(c)
	if err != nil {
		a.OKWithBizCode(consts.BizCode500, err.Error())
	}
}

type EcdnProxy struct {
	api.Api
}

func (a EcdnProxy) ReportDispatchParamsReport(c *gin.Context) {
	s := service.EcdnProxy{}
	err := s.ReportDispatchParamsReport(c)
	if err != nil {
		a.OKWithBizCode(consts.BizCode500, err.Error())
	}
}

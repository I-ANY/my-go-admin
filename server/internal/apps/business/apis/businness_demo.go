package apis

import (
	"biz-auto-api/internal/apps/business/service"
	"biz-auto-api/internal/apps/business/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DemoApi struct {
	api.Api
}

// 获取zap客户管控平台的Cookie
func (z DemoApi) GetDemo(c *gin.Context) {
	s := service.BusinessDemoService{}
	var req = dto.DemoReq{}
	err := z.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.Form, binding.Query).
		Validate(req).
		Errors
	if err != nil {
		z.Logger.Error(err)
		z.OKWithBizCode(consts.BizCode400, err.Error())
	}
	res, err := s.Demo(&req)
	if err != nil {
		z.Logger.Error(err)
	}
	z.OKWithCodeAndData(consts.BizCode200, res, "success")
}

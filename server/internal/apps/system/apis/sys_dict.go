package apis

import (
	"biz-auto-api/internal/apps/system/service"
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SysDict struct {
	api.Api
}

func (a SysDict) GetAllDictData(c *gin.Context) {
	s := service.SysDict{}
	var req = dto.GetAllDictDataReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetAllDictData(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

package apis

import (
	"biz-auto-api/internal/apps/system/service"
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SysMenu struct {
	api.Api
}

func (a SysMenu) GetMenuTree(c *gin.Context) {
	s := service.SysMenu{}
	var req = dto.GetMenuTreeReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, binding.JSON).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetMenuTree(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, 1, 999999, "查询成功")
}

func (a SysMenu) UpdateMenu(c *gin.Context) {
	s := service.SysMenu{}
	var req = dto.UpdateMenuReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, binding.JSON, nil, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.UpdateMenu(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "编辑成功")
}

func (a SysMenu) AddMenu(c *gin.Context) {
	s := service.SysMenu{}
	var req = dto.AddMenuReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, binding.JSON, nil, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.AddMenu(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "添加成功")
}

func (a SysMenu) DeleteMenu(c *gin.Context) {
	s := service.SysMenu{}
	var req = dto.DeleteMenuReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, binding.JSON, nil, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.DeleteMenu(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	if !res.Success {
		a.OKWithBizCode(consts.BizCode300, res.Message)
		return
	}
	a.OK(res, "删除成功")
}

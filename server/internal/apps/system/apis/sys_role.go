package apis

import (
	"biz-auto-api/internal/apps/system/service"
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"regexp"
)

type SysRole struct {
	api.Api
}

func (a SysRole) GetRoleList(c *gin.Context) {
	s := service.SysRole{}
	var req = dto.GetRoleListReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.Form, binding.Query).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetRoleList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a SysRole) UpdateRole(c *gin.Context) {
	s := service.SysRole{}
	var req = dto.UpdateRoleReq{}
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
	res, err := s.UpdateRole(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "编辑成功")
}

func (a SysRole) AddRole(c *gin.Context) {
	s := service.SysRole{}
	var req = dto.AddRoleReq{}
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
	reg1 := regexp.MustCompile("[a-zA-Z_]{5,}")
	reg2 := regexp.MustCompile("^_+$")
	reg3 := regexp.MustCompile("^\\d+$")
	if !reg1.Match([]byte(req.Identify)) || reg2.Match([]byte(req.Identify)) || reg3.Match([]byte(req.Identify)) {
		a.OKWithBizCode(consts.BizCode400, "角色标识只能有字母下划线数字组成，且不能为纯数字和下划线")
		return
	}
	res, err := s.AddRole(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "添加成功")

}

func (a SysRole) DeleteRole(c *gin.Context) {
	s := service.SysRole{}
	var req = dto.DeleteRoleReq{}
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
	res, err := s.DeleteRole(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "删除成功")
}

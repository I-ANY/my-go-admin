package apis

import (
	"biz-auto-api/internal/apps/system/service"
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SysResource struct {
	api.Api
}

func (a SysResource) GetResourceList(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceListReq{}
	req.SetMaxPageSize(99999)
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
	res, err := s.GetResourceList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a SysResource) GetResourceTableList(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceTableListReq{}
	req.SetMaxPageSize(9999)
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
	res, err := s.GetResourceTableList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a SysResource) GetResourceTableField(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceTableFieldReq{}
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
	res, err := s.GetResourceTableField(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

func (a SysResource) AddResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.AddResourceReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.AddResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) UpdateResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.UpdateResourceReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.UpdateResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) GetResourceViewFormSchemas(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceViewFormSchemasReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetResourceViewFormSchemas(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) GetResourceViewTableColumns(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceViewTableColumnsReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetResourceViewTableColumns(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) GetResourceDetailList(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetResourceDetailListReq{}
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
	res, err := s.GetResourceDetail(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a SysResource) DeleteResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.DeleteResourceReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.DeleteResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) GetRoleResourceInfo(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetRoleResourceInfoReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetRoleResourceInfo(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

func (a SysResource) GetRoleResourceDetailList(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetRoleResourceDetailListReq{}
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
	res, err := s.GetRoleResourceDetailList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a SysResource) UpdateRoleResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.UpdateRoleResourceReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.UpdateRoleResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

func (a SysResource) GetRoleAuthedResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetRoleAuthedResourceReq{}
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
	res, err := s.GetRoleAuthedResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

func (a SysResource) GetBusinessResource(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.GetBusinessResourceReq{}
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
	res, err := s.GetBusinessResource(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

func (a SysResource) RoleResourceAuth(c *gin.Context) {
	s := service.SysResource{}
	var req = dto.RoleResourceAuthReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Bind(&req, nil, binding.JSON, binding.Query, binding.Form).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.RoleResourceAuth(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "操作成功")
}

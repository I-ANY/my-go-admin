package apis

import (
	"biz-auto-api/internal/apps/system/service"
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type SysBase struct {
	api.Api
}

func (a SysBase) UserLogin(c *gin.Context) {
	s := service.SysBase{}
	req := dto.SysBaseUserLoginReq{}
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
	res, err := s.UserLogin(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	if !res.Success {
		a.OKWithBizCode(http.StatusBadRequest, res.Message)
	} else {
		a.OK(res, "登录成功")
	}
}

func (a SysBase) UserLogout(c *gin.Context) {
	s := service.SysBase{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	err = s.UserLogout()
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(nil, "登出成功")
}

func (a SysBase) GetUserInfo(c *gin.Context) {
	var req = dto.SysBaseGetUserInfoReq{}
	s := service.SysBase{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetUserInfo(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "获取用户信息成功")
}

func (a SysBase) GetUserPermCode(c *gin.Context) {
	var req = dto.SysBaseGetUserPermCodeReq{}
	s := service.SysBase{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetUserPermCode(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "获取用户权限码成功")
}

func (a SysBase) GetUserMenu(c *gin.Context) {
	var req = dto.SysBaseGetUserMenuReq{}
	s := service.SysBase{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.GetUserMenu(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "获取用户菜单成功")
}

func (a SysBase) RefreshToken(c *gin.Context) {
	s := service.SysBase{}
	req := dto.RefreshTokenReq{}
	err := a.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Validate(req).
		Errors
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode400, err.Error())
		return
	}
	res, err := s.RefreshToken(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	if !res.Success {
		a.OKWithBizCode(http.StatusBadRequest, res.Message)
	} else {
		a.OK(res, "刷新成功")
	}
}

func (a SysBase) GetStarPortalLoginUrl(c *gin.Context) {
	s := service.SysBase{}
	req := dto.SysBaseGetStarPortalLoginUrlReq{}
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
	res, err := s.GetStarPortalLoginUrl(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "获取成功")

}

func (a SysBase) StarPortalLogin(c *gin.Context) {
	s := service.SysBase{}
	req := dto.SysBaseStarPortalLoginReq{}
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
	res, err := s.StarPortalLogin(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	if !res.Success {
		a.OKWithBizCode(http.StatusBadRequest, res.Message)
	} else {
		a.OK(res, "登录成功")
	}
}

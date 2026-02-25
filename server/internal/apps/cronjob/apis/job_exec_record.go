package apis

import (
	"biz-auto-api/internal/apps/cronjob/service"
	"biz-auto-api/internal/apps/cronjob/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type JobExecRecord struct {
	api.Api
}

func (a JobExecRecord) GetJobExecRecordList(c *gin.Context) {
	s := service.JobExecRecord{}
	var req = dto.GetJobExecRecordListReq{}
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
	res, err := s.GetJobExecRecordList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a JobExecRecord) GetJobExecLog(c *gin.Context) {
	s := service.JobExecRecord{}
	var req = dto.GetJobExecLogReq{}
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
	res, err := s.GetJobExecLogReq(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

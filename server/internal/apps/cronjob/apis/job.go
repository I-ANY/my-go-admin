package apis

import (
	"biz-auto-api/internal/apps/cronjob/service"
	"biz-auto-api/internal/apps/cronjob/service/dto"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/robfig/cron/v3"
)

type JobApi struct {
	api.Api
}

func (a JobApi) GetJobList(c *gin.Context) {
	s := service.Job{}
	var req = dto.GetJobListReq{}
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
	res, err := s.GetJobList(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.PageOK(res.Items, res.Total, req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (a JobApi) AddJob(c *gin.Context) {
	s := service.Job{}
	var req = dto.AddJobReq{}
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
	if !a.CheckCronExpr(req.CronExpression) {
		a.OKWithBizCode(consts.BizCode400, "非法的cron表达式: "+req.CronExpression)
		return
	}
	res, err := s.AddJob(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "添加成功")
}
func (a JobApi) CheckCronExpr(expr string) bool {
	p := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow, //| cron.Descriptor,
	)
	_, err := p.Parse(expr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func (a JobApi) UpdateJob(c *gin.Context) {
	s := service.Job{}
	var req = dto.UpdateJobReq{}
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
	if req.CronExpression != nil && !a.CheckCronExpr(*req.CronExpression) {
		a.OKWithBizCode(consts.BizCode400, "非法的cron表达式: "+*req.CronExpression)
		return
	}
	res, err := s.UpdateJob(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "编辑成功")
}

func (a JobApi) ExecJob(c *gin.Context) {
	s := service.Job{}
	var req = dto.ExecuteJobReq{}
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
	res, err := s.ExecJob(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "触发执行成功")
}

func (a JobApi) DeleteJob(c *gin.Context) {
	s := service.Job{}
	var req = dto.DeleteJobReq{}
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
	res, err := s.DeleteJob(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "删除成功")
}

func (a JobApi) GetJob(c *gin.Context) {
	s := service.Job{}
	var req = dto.GetJobReq{}
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
	res, err := s.GetJob(&req)
	if err != nil {
		a.Logger.Error(err)
		a.OKWithBizCode(consts.BizCode500, err.Error())
		return
	}
	a.OK(res, "查询成功")
}

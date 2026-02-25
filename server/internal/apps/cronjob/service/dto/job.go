package dto

import "biz-auto-api/pkg/dto"

type GetJobListReq struct {
	dto.PaginationReq
	JobName      string `json:"jobName" form:"jobName"`
	JobType      int64  `json:"jobType" form:"jobType" validate:"oneof=0 1 2"`
	InvokeTarget string `json:"invokeTarget" form:"invokeTarget"`
	ScheduleNode string `json:"scheduleNode" form:"scheduleNode"`
	RunStatus    int64  `json:"runStatus" form:"runStatus" validate:"oneof=0 1 2 3"`
	Status       int64  `json:"status" form:"status" validate:"oneof=0 1 2"`
}

type Job struct {
	Id             int64  `json:"id"`
	JobName        string `json:"jobName"` // 名称
	JobType        int64  `json:"jobType"`
	CronExpression string `json:"cronExpression"`
	InvokeTarget   string `json:"invokeTarget"`
	Args           string `json:"args"`
	Status         int64  `json:"status"`
	RunStatus      int64  `json:"runStatus"`
	EntryId        int64  `json:"entryId"`
	ScheduleNode   string `json:"scheduleNode" `
}

type GetJobListRes struct {
	Items []*Job `json:"items"`
	Total int64  `json:"total"`
}

type AddJobReq struct {
	JobName        string `json:"jobName" validate:"required"` // 名称
	JobType        int64  `json:"jobType" validate:"oneof=1 2"`
	CronExpression string `json:"cronExpression" validate:"required,cron"`
	InvokeTarget   string `json:"invokeTarget" validate:"required_if=JobType 2"` // exec必传
	Args           string `json:"args" validate:"required_if=JobType 1"`         // falsk-api 必传
	Status         int64  `json:"status" validate:"oneof=1 2"`
	ScheduleNode   string `json:"scheduleNode" validate:"required"`
}
type AddJobRes struct{}

type DeleteJobReq struct {
	Id int64 `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
}
type DeleteJobRes struct{}

type UpdateJobReq struct {
	Id             int64   `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
	JobName        *string `json:"jobName" validate:"omitnil,required"` // 名称
	JobType        *int64  `json:"jobType" validate:"omitnil,oneof=1 2"`
	CronExpression *string `json:"cronExpression" validate:"omitnil,required,cron"`
	InvokeTarget   *string `json:"invokeTarget" validate:"required_if=JobType 2"`
	Args           *string `json:"args" validate:"required_if=JobType 1"`
	Status         *int64  `json:"status" validate:"omitnil,oneof=1 2"`
	//ScheduleNode   *string `json:"scheduleNode" validate:"omitnil,required"`
}
type UpdateJobRes struct{}

type ExecuteJobReq struct {
	Id int64 `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
}
type ExecuteJobRes struct {
}

type GetJobReq struct {
	Id int64 `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
}
type GetJobRes struct {
	Job `json:",inline"`
}

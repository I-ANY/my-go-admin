package dto

import (
	"biz-auto-api/pkg/dto"
)

type GetJobExecRecordListReq struct {
	dto.PaginationReq
	JobId          int64  `json:"jobId" form:"jobId" uri:"jobId" validate:"gte=0"`
	RunStatus      int64  `json:"runStatus" form:"runStatus" validate:"oneof=0 1 2 3"`
	TriggerType    int64  `json:"triggerType" form:"triggerType" validate:"oneof=0 1 2"`
	StartTimeBegin string `json:"startTimeBegin" form:"startTimeBegin"`
	StartTimeEnd   string `json:"startTimeEnd" form:"startTimeEnd"`
}
type JobExecRecord struct {
	Id          int64  `json:"id"`
	RunStatus   int64  `json:"runStatus"`
	TriggerType int64  `json:"triggerType"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	JobId       int64  `json:"jobId"`
	LatencyTime string `json:"latencyTime"`
	Job         Job    `json:"job"`
}
type GetJobExecRecordListRes struct {
	Items []*JobExecRecord `json:"items"`
	Total int64            `json:"total"`
}

type GetJobExecLogReq struct {
	JobExecRecordId int64 `json:"id" form:"id" uri:"id" validate:"required,gt=0"`
	LastId          int64 `json:"lastId" form:"lastId" validate:"gte=0"`
}

type JobExecLog struct {
	Id      int64  `json:"id"`
	LogTime string `json:"logTime"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type GetJobExecLogRes struct {
	RunStatus int64         `json:"runStatus"`
	LastId    int64         `json:"lastId"`
	Logs      []*JobExecLog `json:"logs"`
}

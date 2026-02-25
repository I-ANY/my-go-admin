package models

import (
	"biz-auto-api/pkg/models"
	"time"
)

const (
	JobStatusEnable  = 1
	JobStatusDisable = 2
	RunStatusRunning = 1
	RunStatusSuccess = 2
	RunStatusFailed  = 3

	JobTypeFlaskApi = 1
	JobTypeExec     = 2

	JobNeedRemoveYes = 1
	JobNeedRemoveNo  = 2
)

type CjJob struct {
	models.BseId
	JobName        string             `json:"jobName" gorm:"size:255;not null"`                               // 名称
	JobType        int64              `json:"jobType" gorm:"size:1;default:2;comment:1 http，2 exec;not null"` // 任务类型
	CronExpression string             `json:"cronExpression" gorm:"size:255;not null"`                        // cron表达式
	InvokeTarget   string             `json:"invokeTarget" gorm:"size:255;not null"`                          // 调用目标
	Args           string             `json:"args" gorm:"type:text;"`                                         // 目标参数
	Status         int64              `json:"status" gorm:"size:1;comment:1启用，2禁用;default:1"`                 // 状态
	RunStatus      int64              `json:"RunStatus" gorm:"size:1;comment:1运行中，2执行成功，3,执行失败"`              // 运行状态
	EntryId        int64              `json:"entryId" gorm:"size:11;"`                                        // job启动时返回的id
	ScheduleNode   string             `json:"scheduleNode" gorm:"comment:调度节点;size:255;not null"`
	ExecRecord     []*CjJobExecRecord `json:"ExecRecord" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:JobId"`
	NeedRemove     int64              `json:"needRemove" gorm:"comment:标记是否需要在下一次加载定时任务时从未来执行列表中将其删除，删除后未来将不会再执行。1是，2否"`
	models.ControlBy
	models.ModelTime
}

func (CjJob) TableName() string {
	return "cj_job"
}

const (
	LogLevelInfo  = "INFO"
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
	LogLevelDebug = "DEBUG"
)

type CjJobExecLog struct {
	models.BseId
	LogTime      time.Time        `json:"logTime" gorm:"comment:日志打印时间"`
	Level        string           `json:"level" gorm:"comment:日志级别"`
	Message      string           `json:"message" gorm:"type:longtext;comment:日志信息"`
	ExecRecordId int64            `json:"execRecordId" gorm:"comment:执行记录Id;not null;index:,type:btree"`
	ExecRecord   *CjJobExecRecord `json:"execRecord" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:ExecRecordId"`
	models.ControlBy
	models.ModelTime
}

func (CjJobExecLog) TableName() string {
	return "cj_job_exec_log"
}

const (
	TriggerTypeAuto = iota + 1
	TriggerTypeManual
)

type CjJobExecRecord struct {
	models.BseId
	RunStatus   int64     `json:"runStatus" gorm:"comment:执行状态;comment:1运行中，2执行成功，3,执行失败"`
	TriggerType int64     `json:"triggerType" gorm:"comment:触发方式,1定时触发，2手动触发"`
	StartTime   time.Time `json:"startTime" gorm:"comment:开始时间"`
	EndTime     time.Time `json:"endTime" gorm:"comment:结束时间"`
	JobId       int64     `json:"jobId" gorm:"comment:关联的jobId;not null;index:,type:btree"`
	Job         *CjJob    `json:"job" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:JobId"`
	LatencyTime string    `json:"latencyTime" gorm:"size:128;comment:耗时"`

	ExecLog []*CjJobExecLog `json:"execLog" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:ExecRecordId"`
	models.ControlBy
	models.ModelTime
}

func (CjJobExecRecord) TableName() string {
	return "cj_job_exec_record"
}

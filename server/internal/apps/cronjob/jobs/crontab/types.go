package crontab

import (
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/config/types"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type JobMeta struct {
	InvokeTarget   string
	Name           string
	JobId          int64
	EntryId        int64
	CronExpression string
	Args           string
	ScheduleNode   string
	TriggerType    int64
}

type JobExecArgs struct {
	Args    string
	DB      *gorm.DB
	Log     *PersistLogger
	Config  types.CronjobConfig
	Redis   *redis.Client
	JobId   int64
	JobName string
	CK      *clickhouse.CK
}

type JobFlaskApiArgs struct {
	Uri     string                 `json:"uri"`
	Method  string                 `json:"method"`
	Query   map[string]interface{} `json:"query"`
	Header  map[string]interface{} `json:"header"`
	Body    interface{}            `json:"body"`
	Timeout int64                  `json:"timeout"`
}

type JobExec interface {
	Exec(args *JobExecArgs) error
}

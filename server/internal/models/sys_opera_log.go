package models

import (
	"biz-auto-api/pkg/models"
	"time"
)

const (
	SysLogHandleSourceUser    = 1
	SysLogHandleSourceProgram = 2
)

type SysOperaLog struct {
	models.BseId
	Uri           string    `json:"uri" gorm:"size:255;comment:请求uri;index:idx_req_time_source_uri,priority:10;"`
	Api           string    `json:"api" gorm:"size:255;comment:请求api"`
	RequestId     string    `json:"requestId" gorm:"size:255;comment:请求Id;index:idx_req_time_source_reqid,priority:10;"`
	RequestMethod string    `json:"requestMethod" gorm:"size:128;comment:请求方式 GET POST PUT DELETE;index:idx_req_time_source_method,priority:10;"`
	ClientIp      string    `json:"clientIp" gorm:"size:128;comment:客户端ip"`
	LatencyTime   string    `json:"latencyTime" gorm:"size:128;comment:耗时"`
	UserAgent     string    `json:"userAgent" gorm:"size:255;comment:ua"`
	ReqBody       string    `json:"reqBody" gorm:"type:longtext;comment:请求体数据"`
	JsonRes       string    `json:"JsonRes" gorm:"type:longtext;comment:返回数据"`
	HttpCode      string    `json:"httpCode" gorm:"size:10;comment:http状态码"`
	BizCode       string    `json:"bizCode" gorm:"size:10;comment:业务状态码"`
	RequestTime   time.Time `json:"requestTime" gorm:"comment:请求时间;index:idx_req_time_source_uri,priority:3;index:idx_req_time_source_method,priority:3;index:idx_req_time_source_user,priority:3;index:idx_req_time_source_reqid,priority:3;"`
	Handler       string    `json:"handler" gorm:"size:255;comment:处理函数"`
	HandleSource  int64     `json:"handleSource" gorm:"type:tinyint(1);comment:操作来源 1-用户 2-程序;index:idx_req_time_source_uri,priority:5;index:idx_req_time_source_method,priority:5;index:idx_req_time_source_user,priority:5;index:idx_req_time_source_reqid,priority:5;"`
	CreateBy      int64     `json:"createBy" gorm:"index;comment:创建者;index:idx_req_time_source_user,priority:10;"`
	UpdateBy      int64     `json:"updateBy" gorm:"index;comment:更新者"`
	models.ModelTime
}

func (SysOperaLog) TableName() string {
	return "sys_opera_log"
}

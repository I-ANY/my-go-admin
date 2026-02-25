package models

import "biz-auto-api/pkg/models"

type AlarmData struct {
	models.BseId
	Level  string `gorm:"column:level;comment:告警级别"`
	Topic  string `gorm:"column:topic;comment:告警类型"`
	Detail string `gorm:"column:detail;type:longtext;comment:告警详情"`
	Status int    `gorm:"column:status;comment:告警状态  0:失败, 1:成功"`
	models.ModelTime
}

func (AlarmData) TableName() string { return "business_alarm_data" }

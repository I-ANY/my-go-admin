package service

import (
	"biz-auto-api/internal/apps/cronjob/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type JobExecRecord struct {
	service.Service
}

func (s JobExecRecord) GetJobExecRecordList(req *dto.GetJobExecRecordListReq) (*dto.GetJobExecRecordListRes, error) {
	var result = &dto.GetJobExecRecordListRes{}
	db := s.DB
	log := s.Log
	var total int64
	var execRecords []*models.CjJobExecRecord

	err := db.Model(&models.CjJobExecRecord{}).Preload("Job", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "job_name")
	}).Joins("LEFT JOIN cj_job ON cj_job.id = cj_job_exec_record.job_id").
		Scopes(func(tx *gorm.DB) *gorm.DB {
			if req.JobId > 0 {
				tx = tx.Where("cj_job_exec_record.job_id = ?", req.JobId)
			}
			if req.TriggerType > 0 {
				tx = tx.Where("cj_job_exec_record.trigger_type = ?", req.TriggerType)
			}
			if req.RunStatus > 0 {
				tx = tx.Where("cj_job_exec_record.run_status = ?", req.RunStatus)
			}
			if len(req.StartTimeBegin) > 0 && len(req.StartTimeEnd) > 0 {
				tx = tx.Where("cj_job_exec_record.start_time >= ? and cj_job_exec_record.start_time <= ?", req.StartTimeBegin, req.StartTimeEnd)
			}
			return tx
		}).Where("cj_job.id IS NOT NULL and cj_job.deleted_at IS NULL"). // 只连接没有被删除的记录
		Count(&total).
		Order("cj_job_exec_record.id desc").
		Scopes(req.MakePagination()).
		Find(&execRecords).Error
	if err != nil {
		err = errors.Wrap(err, "list job exec records failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items = make([]*dto.JobExecRecord, 0, len(execRecords))
	for _, r := range execRecords {
		startTime := ""
		if r.StartTime.Unix() > 0 {
			startTime = r.StartTime.Format(time.DateTime)
		}
		endTime := ""
		if r.EndTime.Unix() > 0 {
			endTime = r.EndTime.Format(time.DateTime)
		}
		item := &dto.JobExecRecord{
			Id:          r.Id,
			JobId:       r.JobId,
			RunStatus:   r.RunStatus,
			TriggerType: r.TriggerType,
			StartTime:   startTime,
			EndTime:     endTime,
			LatencyTime: r.LatencyTime,
		}
		if r.Job.Id > 0 {
			item.Job = dto.Job{
				Id:      r.Job.Id,
				JobName: r.Job.JobName,
			}
		}
		items = append(items, item)
	}
	result.Total = total
	result.Items = items
	return result, nil
}

func (s JobExecRecord) GetJobExecLogReq(req *dto.GetJobExecLogReq) (*dto.GetJobExecLogRes, error) {
	var result = &dto.GetJobExecLogRes{}
	db := s.DB
	log := s.Log
	var execRecord = models.CjJobExecRecord{}
	var execLogs []*models.CjJobExecLog
	err := db.Select("run_status").First(&execRecord, req.JobExecRecordId).Error
	if err != nil {
		err = errors.Wrap(err, "query exec record failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = db.Where("exec_record_id = ? and id > ?", req.JobExecRecordId, req.LastId).Order("id asc").Find(&execLogs).Error
	if err != nil {
		err = errors.Wrap(err, "query exec record failed")
		log.Errorf("%+v", err)
		return nil, err
	}

	if len(execLogs) > 0 {
		result.LastId = execLogs[len(execLogs)-1].Id
	} else { // 未查询到日志信息则重新将LastId返回
		result.LastId = req.LastId
	}
	result.RunStatus = execRecord.RunStatus
	for _, l := range execLogs {
		logTime := ""
		if l.LogTime.Unix() > 0 {
			logTime = l.LogTime.Format(time.DateTime)
		}
		dl := &dto.JobExecLog{
			Id:      l.Id,
			LogTime: logTime,
			Level:   l.Level,
			Message: l.Message,
		}
		result.Logs = append(result.Logs, dl)
	}
	return result, nil
}

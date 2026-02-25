package service

import (
	"biz-auto-api/internal/apps/cronjob/jobs/crontab"
	"biz-auto-api/internal/apps/cronjob/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/config"
	pkgdb "biz-auto-api/pkg/db"
	pkgmodels "biz-auto-api/pkg/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Job struct {
	service.Service
}

func (s Job) GetJobList(req *dto.GetJobListReq) (*dto.GetJobListRes, error) {
	var result = &dto.GetJobListRes{}
	db := s.DB
	log := s.Log
	var total int64
	var jobs []*models.CjJob

	err := db.Model(&models.CjJob{}).Scopes(func(tx *gorm.DB) *gorm.DB {
		if len(req.JobName) > 0 {
			tx.Where("job_name LIKE ?", tools.FuzzyQuery(req.JobName))
		}
		if req.JobType > 0 {
			tx.Where("job_type = ?", req.JobType)
		}
		if len(req.InvokeTarget) > 0 {
			tx.Where("invoke_target LIKE ?", tools.FuzzyQuery(req.InvokeTarget))
		}
		if len(req.ScheduleNode) > 0 {
			tx.Where("schedule_node LIKE ?", tools.FuzzyQuery(req.ScheduleNode))
		}
		if req.Status > 0 {
			tx.Where("status = ?", req.Status)
		}
		if req.RunStatus > 0 {
			tx.Where("run_status = ?", req.RunStatus)
		}
		return tx
	}).Count(&total).Order("id desc ").Scopes(req.MakePagination()).Find(&jobs).Error
	if err != nil {
		err = errors.Wrap(err, "list jobs failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items = make([]*dto.Job, 0, len(jobs))
	err = copier.Copy(&items, jobs)
	if err != nil {
		err = errors.Wrap(err, "convert data failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.Total = total
	result.Items = items
	return result, nil
}

func (s Job) AddJob(req *dto.AddJobReq) (*dto.AddJobRes, error) {
	var result = &dto.AddJobRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var job = &models.CjJob{
		ControlBy: pkgmodels.ControlBy{
			UpdateBy: userId, CreateBy: userId,
		},
	}
	err := copier.Copy(job, req)
	if err != nil {
		err = errors.Wrap(err, "convert data failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	// JobTypeFlaskApi类型不需要传入InvokeTarget
	if req.JobType == models.JobTypeFlaskApi {
		job.InvokeTarget = ""
	}
	err = db.Create(job).Error
	if err != nil {
		err = errors.Wrap(err, "add job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}

func (s Job) DeleteJob(req *dto.DeleteJobReq) (*dto.DeleteJobRes, error) {
	var result = &dto.DeleteJobRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var job = models.CjJob{}
	err := db.First(&job, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		data := map[string]interface{}{
			"update_by":   userId,
			"need_remove": models.JobNeedRemoveYes, // 将其充未来执行列表中删除
		}
		err = tx.Model(&models.CjJob{}).Where("id = ?", job.Id).Updates(data).Delete(&models.CjJob{}).Error
		if err != nil {
			err = errors.Wrap(err, "delete job failed")
			log.Errorf("%+v", err)
			return err
		}
		//if job.EntryId > 0 {
		//	crontab.GetCrontab().Remove(cron.EntryID(job.EntryId))
		//}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s Job) UpdateJob(req *dto.UpdateJobReq) (*dto.UpdateJobRes, error) {
	if req.JobType != nil && *req.JobType == models.JobTypeFlaskApi && req.Args != nil && len(*req.Args) == 0 {
		return nil, errors.New("非法的参数：Args(required)")
	}
	var result = &dto.UpdateJobRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var job = models.CjJob{}
	err := db.First(&job, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	data := tools.StructToMap(req, "", false, "Id")
	data["update_by"] = userId
	data["need_remove"] = models.JobNeedRemoveYes // 将其充未来执行列表中删除
	err = db.Model(&models.CjJob{}).Where("id = ?", req.Id).Updates(data).Error
	if err != nil {
		err = errors.Wrap(err, "update job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}

func (s Job) ExecJob(req *dto.ExecuteJobReq) (*dto.ExecuteJobRes, error) {
	var result = &dto.ExecuteJobRes{}
	db := s.DB
	log := s.Log
	var job = models.CjJob{}
	err := db.First(&job, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var j cron.Job
	if job.JobType == models.JobTypeFlaskApi {
		fj := &crontab.FlaskApiJob{}
		fj.InvokeTarget = job.InvokeTarget
		fj.CronExpression = job.CronExpression
		fj.JobId = job.Id
		fj.Name = job.JobName
		fj.ScheduleNode = config.CronjobConfig.Cronjob.NodeName
		fj.DB = pkgdb.GetDB()
		fj.Args = job.Args
		fj.TriggerType = models.TriggerTypeManual
		j = fj
	} else if job.JobType == models.JobTypeExec {
		ej := &crontab.ExecJob{}
		ej.InvokeTarget = job.InvokeTarget
		ej.CronExpression = job.CronExpression
		ej.JobId = job.Id
		ej.Name = job.JobName
		ej.ScheduleNode = config.CronjobConfig.Cronjob.NodeName
		ej.Args = job.Args
		ej.DB = pkgdb.GetDB()
		ej.TriggerType = models.TriggerTypeManual
		j = ej
	} else {
		err = errors.Errorf("unknown job type %v", job.JobType)
		log.Warnf("%s", err)
		return nil, err
	}
	// 执行任务
	go j.Run()
	return result, nil
}

func (s Job) GetJob(req *dto.GetJobReq) (*dto.GetJobRes, error) {
	var result = &dto.GetJobRes{}
	db := s.DB
	log := s.Log
	var job = models.CjJob{}
	err := db.First(&job, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query job failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = copier.Copy(&result.Job, job)
	if err != nil {
		err = errors.Wrap(err, "convert job data failed")
		log.Errorf("%+v", err)
		return nil, err
	}

	return result, nil
}

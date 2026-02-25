package crontab

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/logger"
	pkgredis "biz-auto-api/pkg/redis"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type ExecJob struct {
	JobMeta
	DB *gorm.DB
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	l := logger.GetLogger().
		WithField("invoke_target", e.InvokeTarget).
		WithField("schedule_node", e.ScheduleNode).
		WithField("job_name", e.Name).
		WithField("engine", "crontab")
	// 创建一条执行记录
	var execRecord = models.CjJobExecRecord{
		JobId:       e.JobId,
		RunStatus:   models.RunStatusRunning,
		TriggerType: e.TriggerType,
		StartTime:   startTime,
	}
	err := e.DB.Omit("end_time").Create(&execRecord).Error
	if err != nil {
		panic(err)
	}
	// 创建持久化的logger
	log := NewPersistLogger(l, execRecord.Id, e.DB)
	log.Infof("job start with: %v", startTime.Format(time.DateTime))
	defer func() {
		var derr = err
		rcv := recover()
		if rcv != nil {
			if e1, ok := rcv.(error); ok {
				derr = errors.WithStack(e1)
			} else {
				derr = errors.WithStack(errors.Errorf("%s", e1))
			}
		}
		// 结束时间
		endTime := time.Now()
		// 执行耗时
		latencyTime := endTime.Sub(startTime)
		if derr == nil { // 执行成功
			log.Infof("job %s exec success , spend :%v", e.Name, latencyTime)
			err = UpdateJobExecInfo(e.DB, e.JobId, execRecord.Id, models.RunStatusSuccess, endTime, latencyTime)
		} else { // 执行失败
			log.Errorf("%+v", derr)
			log.Errorf("job %s exec failed, spend :%v", e.Name, latencyTime)
			// 更新job状态
			err = UpdateJobExecInfo(e.DB, e.JobId, execRecord.Id, models.RunStatusFailed, endTime, latencyTime)
		}
		if err != nil {
			log.Errorf("%s", errors.WithMessage(err, "update job exec info failed"))
		}
	}()

	// 任务启动更新job状态
	err = e.DB.Model(&models.CjJob{}).Where("id = ? ", e.JobId).Update("run_status", models.RunStatusRunning).Error
	if err != nil {
		err = errors.Wrap(err, "update job status failed")
		return
	}
	var execJob = JobList[e.InvokeTarget]
	if execJob == nil { // 没有找到对应的InvokeTarget
		err = errors.Errorf("invoke target not found: %s", e.InvokeTarget)
		return
	}
	args := &JobExecArgs{
		Log:     log,
		DB:      e.DB,
		Args:    e.Args,
		Config:  *config.CronjobConfig,
		Redis:   pkgredis.GetClient(),
		JobId:   e.JobId,
		JobName: e.Name,
		CK:      clickhouse.GetCK(),
	}
	err = execJob.Exec(args)
	if err != nil {
		err = errors.WithMessage(err, "exec job failed")
		return
	}
	return
}

func UpdateJobExecInfo(db *gorm.DB, jobId, execRecordId int64, runStatus int64, endTime time.Time, latencyTime time.Duration) error {
	err := db.Model(&models.CjJob{}).Where("id = ? ", jobId).Updates(map[string]interface{}{
		"run_status": runStatus,
	}).Error
	if err != nil {
		return errors.Wrap(err, "update job run status failed")
	}
	err = db.Model(&models.CjJobExecRecord{}).Where("id = ? ", execRecordId).Updates(map[string]interface{}{
		"run_status":   runStatus,
		"end_time":     endTime,
		"latency_time": latencyTime.String(),
	}).Error
	if err != nil {
		return errors.Wrap(err, "update job exec record info failed")
	}
	return nil
}

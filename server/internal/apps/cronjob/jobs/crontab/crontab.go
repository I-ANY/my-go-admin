package crontab

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/tools"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"
	"sync"
	"time"
)

var JobList map[string]JobExec
var (
	crontab *cron.Cron
	once    sync.Once
)

func GetCrontab() *cron.Cron {
	return crontab
}
func Setup(scheduleNode string, db *gorm.DB, jobLoadIntervalSeconds int) {
	log := logger.GetLogger().WithField("schedule_node", scheduleNode).WithField("engine", "crontab")
	cronLogger := NewCronLogger(log)
	once.Do(func() {
		crontab = cron.New(
			cron.WithSeconds(),
			cron.WithLogger(cron.VerbosePrintfLogger(cronLogger)),
			cron.WithChain(cron.Recover(cronLogger), cron.SkipIfStillRunning(cronLogger)),
		)
		log.Info("create cron job executor success")
	})
	setup(scheduleNode, db, log, jobLoadIntervalSeconds)
}

func setup(scheduleNode string, db *gorm.DB, log *logrus.Entry, jobLoadIntervalSeconds int) {

	// 启动清空所有的entry_id,need_remove字段
	err := db.Model(models.CjJob{}).
		Where("schedule_node =?", scheduleNode).
		Updates(map[string]interface{}{
			"entry_id":    nil,
			"need_remove": nil,
		}).Error
	if err != nil {
		err = errors.Wrap(err, "clean job entry_id failed")
		log.Fatalf("%+v", err)
	}

	// 修改运行中job的状态，并且查询出job的ID
	err = db.Model(models.CjJob{}).
		Where("run_status =? and schedule_node =?", models.RunStatusRunning, scheduleNode).Updates(map[string]interface{}{
		"run_status": models.RunStatusFailed,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update job status failed")
		log.Fatalf("%+v", err)
	}
	log.Info("update job run status success")

	// 获取当前节点所有的job，修改它们的执行任务状态：运行中 --> 失败
	var runningJobIds = make([]int64, 0)
	err = db.Model(models.CjJob{}).Where("schedule_node =?", scheduleNode).Pluck("id", &runningJobIds).Error
	if err != nil {
		err = errors.Wrap(err, "query job failed")
		log.Fatalf("%+v", err)
	}
	// 修改job执行记录的运行状态
	err = db.Model(models.CjJobExecRecord{}).
		Where("run_status =? and job_id in ?", models.RunStatusRunning, runningJobIds).
		Updates(map[string]interface{}{
			"run_status": models.RunStatusFailed,
			"end_time":   time.Now().Format(time.DateTime),
		}).Error
	if err != nil {
		err = errors.Wrap(err, "update job exec record run status failed")
		log.Fatalf("%+v", err)
	}
	log.Info("update job exec record run status success")

	// 只查询当节点的job
	var jobs []*models.CjJob
	err = db.Model(models.CjJob{}).Where("status =? and schedule_node =?", models.JobStatusEnable, scheduleNode).Find(&jobs).Error
	if err != nil {
		err = errors.Wrap(err, "query jobs failed")
		log.Fatalf("%+v", err)
	}
	log.Infof("jobCore scheduleNode: %v, total: %v", scheduleNode, len(jobs))

	err = AddJobs(db, crontab, jobs, scheduleNode)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	// 其中任务
	crontab.Start()
	log.Infof("crontab start success")
	// 定时读取数据库内的定时任务列表并且将定时任务加入到执行任务中
	go AddJobCyclic(db, crontab, scheduleNode, jobLoadIntervalSeconds)
}

// AddJob 添加任务 AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(c *cron.Cron, cronExpression string, job cron.Job) (int, error) {
	if job == nil {
		return 0, errors.New("jon can not be empty")
	}
	id, err := c.AddJob(cronExpression, job)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	EntryId := int(id)
	return EntryId, nil
}

// Remove 移除任务
func Remove(c *cron.Cron, entryID int) {
	c.Remove(cron.EntryID(entryID))
}

// 任务停止
//func Stop() chan bool {
//	ch := make(chan bool)
//	go func() {
//		global.GADMCron.Stop()
//		ch <- true
//	}()
//	return ch
//}

func AddJobs(db *gorm.DB, crontab *cron.Cron, jobs []*models.CjJob, scheduleNode string) error {
	log := logger.GetLogger().WithField("schedule_node", scheduleNode)
	var err error
	for _, jb := range jobs {
		var entryId = 0
		if jb.JobType == models.JobTypeFlaskApi {
			fj := &FlaskApiJob{}
			fj.InvokeTarget = jb.InvokeTarget
			fj.CronExpression = jb.CronExpression
			fj.JobId = jb.Id
			fj.Name = jb.JobName
			fj.ScheduleNode = scheduleNode
			fj.DB = db
			fj.Args = jb.Args
			fj.TriggerType = models.TriggerTypeAuto
			entryId, err = AddJob(crontab, jb.CronExpression, fj)
		} else if jb.JobType == models.JobTypeExec {
			ej := &ExecJob{}
			ej.InvokeTarget = jb.InvokeTarget
			ej.CronExpression = jb.CronExpression
			ej.JobId = jb.Id
			ej.Name = jb.JobName
			ej.ScheduleNode = scheduleNode
			ej.Args = jb.Args
			ej.DB = db
			ej.TriggerType = models.TriggerTypeAuto
			entryId, err = AddJob(crontab, jb.CronExpression, ej)
		} else {
			log.Warnf("Unknown jobType=%v", jb.JobType)
			continue
		}
		if err != nil {
			err = errors.WithMessage(err, "add job failed")
			return err
		}
		err = db.Model(models.CjJob{}).Where("id =?", jb.Id).Update("entry_id", entryId).Error
		if err != nil {
			err = errors.Wrap(err, "update job info failed")
			return err
		}
		log.Infof("Add job %v success", jb.JobName)
	}
	return nil
}

func AddJobCyclic(db *gorm.DB, crontab *cron.Cron, scheduleNode string, jobLoadIntervalSeconds int) {
	log := logger.GetLogger().WithField("schedule_node", scheduleNode)
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("%+v", e)
		}
	}()
	ticker := time.NewTicker(time.Second * time.Duration(jobLoadIntervalSeconds))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			silentDB := db.Session(&gorm.Session{Logger: ormlogger.Default.LogMode(ormlogger.Silent)})
			err := RemoveJobFromCrontab(db, log, scheduleNode)
			if err != nil {
				log.Errorf("%+v", err)
				continue
			}
			// 只查询当节点的job，添加job
			var jobs []*models.CjJob
			err = silentDB.Model(models.CjJob{}).Where("status =? and schedule_node =? and ( entry_id is null or entry_id=0 or entry_id='')", models.JobStatusEnable, scheduleNode).Find(&jobs).Error
			if err != nil {
				err = errors.Wrap(err, "query jobs failed")
				log.Errorf("%+v", err)
				continue
			}
			log.Infof("load job count=%v...", len(jobs))
			err = AddJobs(db, crontab, jobs, scheduleNode)
			if err != nil {
				err = errors.WithMessage(err, "add jobs failed")
				log.Errorf("%+v", err)
			}
		}
	}
}

func RemoveJobFromCrontab(db *gorm.DB, log *logrus.Entry, scheduleNode string) error {
	// 把需要删除的job查询出来删了
	var needRemoveJobs []*models.CjJob
	err := db.Model(models.CjJob{}).Unscoped().Where("need_remove=? and schedule_node =?", models.JobNeedRemoveYes, scheduleNode).Find(&needRemoveJobs).Error
	if err != nil {
		return errors.Wrap(err, "query need deleted jobs failed")
	}
	if len(needRemoveJobs) > 0 {
		log.Infof("need deleted jobs from crontab count=%d", len(needRemoveJobs))
		// 从执行列表中删除
		for _, job := range needRemoveJobs {
			if job.EntryId > 0 {
				GetCrontab().Remove(cron.EntryID(job.EntryId))
			}
		}
		log.Infof("jobs delete from crontab success")
		needRemoveJobIds := tools.GetSlice(needRemoveJobs, func(e *models.CjJob) int64 {
			return e.Id
		})
		// 修改这些job的 need_remove  entry_id
		err = db.Unscoped().Model(models.CjJob{}).Where("id in ?", needRemoveJobIds).Updates(map[string]interface{}{
			"need_remove": nil,
			"entry_id":    nil,
		}).Error
		if err != nil {
			return errors.Wrap(err, "update need deleted jobs info failed")
		}
	}
	return nil
}

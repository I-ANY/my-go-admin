package jobs

import "biz-auto-api/internal/apps/cronjob/jobs/crontab"

// ExampleJob
// 新添加的job 必须按照以下格式定义，并实现Exec函数
type ExampleJob struct {
}

func NewExampleJob() *ExampleJob {
	return &ExampleJob{}
}

func (j *ExampleJob) Exec(args *crontab.JobExecArgs) error {
	// TODO 待完成
	args.Log.Infof("ExampleJob")
	return nil
}

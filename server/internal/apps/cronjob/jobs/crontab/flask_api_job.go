package crontab

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/clients/flask_api"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/tools"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FlaskApiJob 任务类型 http
type FlaskApiJob struct {
	JobMeta
	DB *gorm.DB
}

// Run http 任务接口
func (f *FlaskApiJob) Run() {
	startTime := time.Now()
	l := logger.GetLogger().
		WithField("schedule_node", f.ScheduleNode).
		WithField("job_name", f.Name).
		WithField("engine", "crontab")
	// 创建一条执行记录
	var execRecord = models.CjJobExecRecord{
		JobId:       f.JobId,
		RunStatus:   models.RunStatusRunning,
		TriggerType: f.TriggerType,
		StartTime:   startTime,
	}
	err := f.DB.Omit("end_time").Create(&execRecord).Error
	if err != nil {
		panic(err)
	}
	// 创建持久化的logger
	log := NewPersistLogger(l, execRecord.Id, f.DB)
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
			log.Infof("job %s exec success , spend :%v", f.Name, latencyTime)
			err = UpdateJobExecInfo(f.DB, f.JobId, execRecord.Id, models.RunStatusSuccess, endTime, latencyTime)
		} else { // 执行失败
			log.Errorf("%+v", derr)
			log.Errorf("job %s exec failed, spend :%v", f.Name, latencyTime)
			// 更新job状态
			err = UpdateJobExecInfo(f.DB, f.JobId, execRecord.Id, models.RunStatusFailed, endTime, latencyTime)
		}
		if err != nil {
			log.Errorf("%s", errors.WithMessage(err, "update job exec info failed"))
		}
	}()
	// 任务启动更新job状态
	err = f.DB.Model(&models.CjJob{}).Where("id = ? ", f.JobId).Update("run_status", models.RunStatusRunning).Error
	if err != nil {
		err = errors.Wrap(err, "update job status failed")
		return
	}
	c := config.CronjobConfig.FlaskApi
	log.Infof("request args = %v", f.Args)
	if len(f.Args) == 0 {
		err = errors.New("args can not be empty")
		return
	}
	flaskApiArgs := &JobFlaskApiArgs{}
	err = json.Unmarshal([]byte(f.Args), flaskApiArgs)
	if err != nil {
		err = errors.Wrapf(err, "unmarshal failed args= %v", f.Args)
		return
	}
	flaskApiArgs.Method = strings.ToUpper(flaskApiArgs.Method)
	// 解析参数
	header, query, body, err := f.ParseArgs(flaskApiArgs)
	if err != nil {
		err = errors.WithMessage(err, "parse args failed")
		return
	}
	flaskApiClient, err := flask_api.NewFlaskApiClient(c.Url, c.ApiToken)
	if err != nil {
		err = errors.WithMessage(err, "create flask api client failed")
		return
	}
	// 设置超时时间
	if flaskApiArgs.Timeout == 0 {
		flaskApiArgs.Timeout = 15
	}
	flaskApiClient.Timeout = time.Second * time.Duration(flaskApiArgs.Timeout)
	log.Infof("request info : host = %v, uri = %v, method = %v, header = %v, query = %v, body = %s",
		c.Url, flaskApiArgs.Uri, flaskApiArgs.Method, header, query, body)
	bs, err := flaskApiClient.HttpRequest(flaskApiArgs.Method, flaskApiArgs.Uri, header, query, nil, body)
	if err != nil {
		err = errors.WithMessage(err, "request flask api failed")
		return
	}
	bodyStr := string(bs)
	if !gjson.Valid(bodyStr) {
		err = errors.Errorf("response not json: %s", string(bs))
		return
	}
	if gjson.Get(bodyStr, "code").Int() != 200 {
		err = errors.Errorf("request flask api failed: %s", string(bs))
		return
	}
	log.Infof("request flask api response: %v", string(bs))
}

func (f *FlaskApiJob) ParseArgs(flaskApiArgs *JobFlaskApiArgs) (http.Header, url.Values, []byte, error) {
	var header http.Header = map[string][]string{}
	var query url.Values = map[string][]string{}
	var body []byte
	var err error
	var allowMethods = []string{"POST", "GET", "PUT", "DELETE"}
	// 校验URI
	if len(flaskApiArgs.Uri) == 0 {
		err = errors.New("uri can not be empty")
		return nil, nil, nil, err
	}
	// 校验请求方法
	if !tools.InSlice(flaskApiArgs.Method, allowMethods) {
		err = errors.Errorf("request method [%s] is not allowed, must be in %v", flaskApiArgs.Method, allowMethods)
		return nil, nil, nil, err
	}
	// 封装header
	for hk, hv := range flaskApiArgs.Header {
		if v, ok := hv.(string); ok {
			header[hk] = []string{v}
		} else {
			err = errors.Errorf("header param key=%v value=%v, value is not string", hk, hv)
			return nil, nil, nil, err
		}
	}
	// 封装query
	for qk, qv := range flaskApiArgs.Query {
		if v, ok := qv.(string); ok {
			query[qk] = []string{v}
		} else {
			err = errors.Errorf("query param key=%v value=%v, value is not string", qk, qv)
			return nil, nil, nil, err
		}
	}
	// 序列化body
	if flaskApiArgs.Body != nil {
		body, err = json.Marshal(flaskApiArgs.Body)
		if err != nil {
			err = errors.Wrapf(err, "marshal failed, body=%v", flaskApiArgs.Body)
			return nil, nil, nil, err
		}
	}
	return header, query, body, nil
}

package inspect

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/clients/operator"
	"biz-auto-api/pkg/logger"
	pkgmodels "biz-auto-api/pkg/models"
	"biz-auto-api/pkg/tools"
	"encoding/base64"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Handler struct {
	log          logger.Logger
	db           *gorm.DB
	userId       int64
	taskId       int64
	servers      []*models.BusinessServer
	operatorType int64
	sync         bool
	pool         *ants.Pool
	TaskName     string
}

func (h *Handler) SetTaskId(taskId int64) {
	h.taskId = taskId
}
func NewHandler(db *gorm.DB, log logger.Logger, userId, operatorType int64, sync bool, taskName string, servers []*models.BusinessServer) (*Handler, error) {
	if db == nil || log == nil {
		return nil, errors.New("db or log is nil")
	}
	getPoolSize := func(serverNumber int) int {
		minimumSize := 30
		maximumSize := 300
		if len(servers) < minimumSize {
			return len(servers)
		}
		poolSize := serverNumber / 3
		if poolSize < minimumSize {
			poolSize = minimumSize
		}
		if poolSize > maximumSize {
			poolSize = maximumSize
		}
		return poolSize
	}
	log.Infof("pool size: %v", getPoolSize(len(servers)))

	// 创建创建 线程池
	pool, err := ants.NewPool(getPoolSize(len(servers)),
		ants.WithPanicHandler(func(i interface{}) { log.Errorf("%+v", i) }),
		ants.WithNonblocking(false),
		ants.WithPreAlloc(true),
		ants.WithLogger(log),
	)
	if err != nil {
		err = errors.Wrap(err, "create ants pool failed")
		return nil, err
	}
	return &Handler{
		log:          log,
		db:           db,
		userId:       userId,
		servers:      servers,
		operatorType: operatorType,
		sync:         sync,
		pool:         pool,
		TaskName:     taskName,
	}, nil
}

func (h *Handler) Start() error {
	// 新建task
	task := models.BusinessInspectTask{
		OperatorType: tools.ToPointer(h.operatorType),
		Status:       tools.ToPointer(int64(models.InspectExecing)),
		StartTime:    tools.ToPointer(time.Now()),
		ControlBy:    pkgmodels.ControlBy{CreateBy: h.userId, UpdateBy: h.userId},
		ServerCount:  tools.ToPointer(int64(len(h.servers))),
		TaskName:     tools.ToPointer(h.TaskName),
	}
	err := h.db.Create(&task).Error
	if err != nil {
		return errors.Wrap(err, "add task failed")
	}
	h.SetTaskId(task.Id)
	if h.sync {
		return h.InspectServers()
	} else {
		go h.InspectServers()
	}
	return nil
}
func (h *Handler) InspectServers() error {
	defer h.pool.ReleaseTimeout(time.Minute * 3)
	h.log.Infof("start, trigger inspect server(count=%v)....", len(h.servers))
	for _, server := range h.servers {
		if len(tools.ToValue(server.Hostname)) == 0 {
			continue
		}
		h.log.Infof("trigger inspect server: %v", *server.Hostname)
		err := h.pool.Submit(h.GetInspectFunc(server))
		if err != nil {
			h.log.Errorf("%+v", errors.Wrap(err, "trigger inspect server failed"))
		}
	}
	h.log.Infof("end, trigger inspect server")
	return nil
}
func (h *Handler) GetInspectFunc(server *models.BusinessServer) func() {
	return func() {
		var err error
		inspectResult := &models.BusinessInspectResult{
			TaskId:            tools.ToPointer(h.taskId),
			FrankID:           server.FrankID,
			Hostname:          server.Hostname,
			Business:          server.Business,
			Owner:             server.Owner,
			Status:            tools.ToPointer(int64(models.InspectExecing)),
			StartTime:         tools.ToPointer(time.Now()),
			PeakUtilization95: server.EveningPeak95Utilization,
			NetworkSpeed:      server.SpeedNow,
			PlannedBandwidth:  server.BwPlan,
			Bandwidth:         server.BwTotal,
			Remark:            nil,
		}
		err = h.db.Create(inspectResult).Error
		if err != nil {
			err = errors.Wrapf(err, "create inspect result failed, hostname=%v", *inspectResult.Hostname)
			h.log.Errorf("%+v", err)
			return
		}

		defer func() {
			// 执行报错
			data := make(map[string]interface{})
			if err != nil {
				data["status"] = models.InspectExecFailed
				data["remark"] = err.Error()
				data["finish_time"] = time.Now()
			}
			if e1 := recover(); e1 != nil {
				data["status"] = models.InspectExecFailed
				data["remark"] = fmt.Sprintf("%+v", e1)
				data["finish_time"] = time.Now()
			}
			if len(data) > 0 {
				e2 := h.db.Model(models.BusinessInspectResult{}).Where("id =?", inspectResult.Id).Updates(data).Error
				if e2 != nil {
					e2 = errors.Wrapf(e2, "update inspect result failed, hostname=%v", *inspectResult.Hostname)
					h.log.Errorf("%+v", e2)
				}
			}
		}()
		argsStr := fmt.Sprintf("%s %s", strconv.FormatInt(inspectResult.Id, 10), strconv.FormatInt(h.taskId, 10))
		// 将ID进行编码
		argsEncode := base64.StdEncoding.EncodeToString([]byte(argsStr))
		_, err = operator.ExecScript("business_inspect.sh", *server.Hostname, "general", argsEncode, time.Second*30)
		if err != nil {
			err = errors.WithMessagef(err, "inspect server failed, hostname=%v", *server.Hostname)
			h.log.Errorf("%+v", err)
		}
	}
}

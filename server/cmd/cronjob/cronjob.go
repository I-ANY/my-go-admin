package cronjob

import (
	"biz-auto-api/internal/apps/cronjob/jobs"
	"biz-auto-api/internal/apps/cronjob/jobs/crontab"
	"biz-auto-api/internal/apps/cronjob/router"
	"biz-auto-api/pkg/casbin"
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/config/types"
	"biz-auto-api/pkg/consts"
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/grpc/clients"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/queue"
	"biz-auto-api/pkg/redis"
	"biz-auto-api/pkg/snowid"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configYml string
	nodeName  string
	Cmd       = &cobra.Command{
		Use:     "cronjob",
		Aliases: []string{"c"},
		Short:   "Start cronjob app",
		Example: "biz-auto-api cronjob -c config-cronjob.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

// func init() {
// 	Cmd.Flags().StringVarP(&configYml, "config", "c", "config-cronjob.yaml", "Start cronjob with provided configuration file")
// 	Cmd.Flags().StringVarP(&nodeName, "node-name", "n", "", "specifies the current node name")
// }

func setup() {
	// 1.加载配置
	config.SetupCronjobConfig(configYml)
	var c = config.CronjobConfig
	// 2.加载logger
	logger.Setup(c.Log.Level)
	// 3.初始化数据库连接
	db.Setup(
		c.Mysql.Host,
		c.Mysql.Database,
		c.Mysql.Username,
		c.Mysql.Password,
		c.Mysql.Port,
		c.Mysql.MaxIdleCon,
		c.Mysql.MaxOpenCon,
	)
	// 初始化casbin
	casbin.Setup(db.GetDB())
	// 未传flage参数，尝试获取环境变量
	if len(nodeName) == 0 {
		if nn, ok := os.LookupEnv("NODE_NAME"); ok {
			nodeName = nn
		}
	}
	// 如果有值则覆盖配置文件中的值
	if len(nodeName) > 0 {
		c.Cronjob.NodeName = nodeName
	}
	redis.Setup(
		c.Redis.Username,
		c.Redis.Password,
		c.Redis.Host,
		c.Redis.DB,
		c.Redis.Port,
		c.Redis.Timeout,
		c.Redis.PoolSize,
	)
	// connect to clickhouse
	clickhouse.Setup(
		c.Clickhouse.Host,
		c.Clickhouse.Database,
		c.Clickhouse.Username,
		c.Clickhouse.Password,
		c.Clickhouse.Port,
		c.Clickhouse.MaxIdleCon,
		c.Clickhouse.MaxOpenCon,
		false,
	)
	err := clients.SetupVpnGrpcClient(c.VpnClient.Address, time.Duration(c.VpnClient.ConnTimeoutSecond)*time.Second)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	log := logger.GetLogger()
	c := config.CronjobConfig
	engine := router.InitRouter()
	//创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", c.Cronjob.Host, c.Cronjob.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(c.Cronjob.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Cronjob.WriteTimeout) * time.Second,
	}
	//启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Cronjob server start failed: %+v", errors.WithStack(err))
		}
	}()
	log.Infof("cronjob server listen at: %v", fmt.Sprintf("%v:%v", c.Cronjob.Host, c.Cronjob.Port))
	crontab.Setup(c.Cronjob.NodeName, db.GetDB(), c.Cronjob.JobLoadIntervalSeconds)
	jobs.RegisterJobs()
	InitSpeedLimitQueue(c, log)
	snowid.Setup("cronjob", redis.GetClient())
	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Start Shutdown crontab...")
	stopC := crontab.GetCrontab().Stop()

	log.Info("Start shutdown server...")
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %+v", errors.WithStack(err))
	}
	log.Info("Server Shutdown success")

	select {
	case <-stopC.Done():
		log.Info("Crontab Shutdown  success")
	case <-ctx.Done():
		log.Warn("Crontab Shutdown timeout，Forced stop")
	}
	return nil
}

func InitSpeedLimitQueue(c *types.CronjobConfig, log *logrus.Logger) {
	queueName := consts.QueneName_SpeedLimitTask
	queueLogger := queue.NewLogger(logger.GetLogger().WithField("type", "queue"))
	q, err := queue.NewQueue(c.Redis.Username,
		c.Redis.Password,
		c.Redis.Host,
		queueName,
		int64(c.Redis.Port),
		int64(c.Redis.DB),
		60*60, //结果保留1小时
		nil,
		queueLogger,
	)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = queue.Queue.AddQueue(queueName, q)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

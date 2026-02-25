package business

import (
	"biz-auto-api/internal/apps/business/router"
	"biz-auto-api/pkg/casbin"
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/config"
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/grpc/clients"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/redis"
	"biz-auto-api/pkg/snowid"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configYml string
	Cmd       = &cobra.Command{
		Use:     "business",
		Aliases: []string{"b"},
		Short:   "Start business app",
		Example: "biz-auto-api business -c config-business.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&configYml, "config", "c", "config-business.yaml", "Start business with provided configuration file")
}

func setup() {
	// 1.加载配置
	config.SetupBusinessConfig(configYml)
	var c = config.BusinessConfig
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
	redis.Setup(
		c.Redis.Username,
		c.Redis.Password,
		c.Redis.Host,
		c.Redis.DB,
		c.Redis.Port,
		c.Redis.Timeout,
		c.Redis.PoolSize,
	)
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
	err := clients.SetupAuthGrpcClient(c.AuthClient.Address, time.Duration(c.AuthClient.ConnTimeoutSecond)*time.Second)
	if err != nil {
		logger.GetLogger().Fatalf("%+v", errors.WithMessage(err, "new auth grpc client failed"))
	}
	logger.GetLogger().Infof("new auth grpc client success")
	snowid.Setup("business", redis.GetClient())
}

func run() error {
	log := logger.GetLogger()
	c := config.BusinessConfig
	engine := router.InitRouter()
	//创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", c.Business.Host, c.Business.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(c.Business.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Business.WriteTimeout) * time.Second,
	}
	//启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("business server start failed: %+v", errors.WithStack(err))
		}
	}()
	log.Infof("business server listen at: %v", fmt.Sprintf("%v:%v", c.Business.Host, c.Business.Port))
	snowid.Setup("business", redis.GetClient())
	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Info("Start shutdown server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %+v", errors.WithStack(err))
	}
	log.Info("Server Shutdown success")
	return nil
}

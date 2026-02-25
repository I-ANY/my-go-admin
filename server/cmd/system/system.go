package system

import (
	"biz-auto-api/internal/apps/system/router"
	"biz-auto-api/pkg/casbin"
	"biz-auto-api/pkg/config"
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/logger"
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
		Use:     "system",
		Aliases: []string{"s"},
		Short:   "Start system app",
		Example: "biz-auto-api system -c config-system.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&configYml, "config", "c", "config-system.yaml", "Start system with provided configuration file")
}

func setup() {
	// 1.加载配置
	config.SetupSystemConfig(configYml)
	var c = config.SystemConfig
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
	//snowid.Setup()
	//redis.Setup(
	//	c.Redis.Username,
	//	c.Redis.Password,
	//	c.Redis.Host,
	//	c.Redis.DB,
	//	c.Redis.Port,
	//	c.Redis.Timeout,
	//	c.Redis.PoolSize,
	//)
}

func run() error {
	log := logger.GetLogger()
	c := config.SystemConfig
	engine := router.InitRouter()
	//创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", c.System.Host, c.System.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(c.System.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.System.WriteTimeout) * time.Second,
	}
	//启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("System server start failed: %+v", errors.WithStack(err))
		}
	}()
	log.Infof("System server listen at: %v", fmt.Sprintf("%v:%v", c.System.Host, c.System.Port))
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

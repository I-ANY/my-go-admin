package migrate

import (
	"biz-auto-api/cmd/migrate/migration"
	"biz-auto-api/pkg/config"
	db "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/logger"
	"biz-auto-api/pkg/models"
	"biz-auto-api/pkg/tools"
	"bytes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	_ "biz-auto-api/cmd/migrate/migration/version"
	"strconv"
	"text/template"
	"time"
)

var (
	configYml string
	generate  bool
	Cmd       = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"m"},
		Short:   "Initialize the database",
		Example: "biz-auto-api migrate -c config-migrate.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	Cmd.Flags().StringVarP(&configYml, "config", "c", "config-migrate.yaml", "Migrate database with provided configuration file")
	Cmd.Flags().BoolVarP(&generate, "generate", "g", false, "Generate migration file")
}

func setup() {
	// 1.加载配置
	config.SetupMigrateConfig(configYml)
	var c = config.MigrateConfig
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
}

func run() {
	if generate {
		if err := genFile(); err != nil {
			logger.GetLogger().Fatalf("%+v", err)
		}
		logger.GetLogger().Info("generate migrate file success")
	} else {
		logger.GetLogger().Info("migrate database start...")
		if err := migrateModel(); err != nil {
			logger.GetLogger().Fatalf("%+v", errors.WithMessage(err, "migrate database failed"))
		}
		logger.GetLogger().Info("migrate database success")
	}
}

func migrateModel() error {
	d := db.GetDB()
	//初始化数据库时候用
	d.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	err := d.Debug().AutoMigrate(&models.Migration{})
	if err != nil {
		return errors.WithStack(err)
	}
	migration.Migrate.SetDb(d.Debug())
	migration.Migrate.Migrate()
	return nil
}

func genFile() error {
	t1, err := template.ParseFiles("template/migrate.template")
	if err != nil {
		return errors.WithStack(err)
	}
	m := map[string]string{}
	m["GenerateTime"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	m["Package"] = "version"
	m["Module"] = "admin"
	var b1 bytes.Buffer
	err = t1.Execute(&b1, m)
	if err != nil {
		return errors.WithStack(err)
	}
	err = tools.Create(b1, "./cmd/migrate/migration/version/"+m["GenerateTime"]+"_migrate.go")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

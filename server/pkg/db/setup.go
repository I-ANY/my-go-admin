package db

import (
	"biz-auto-api/pkg/logger"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB
var _ecdnDB *gorm.DB

func Setup(
	host,
	database,
	username,
	password string,
	port,
	maxIdleCon,
	maxOpenCon int,
) {
	db, err := NewDB(host, database, username, password, port, maxIdleCon, maxOpenCon)
	if err != nil {
		logger.GetLogger().WithField("engine", "mysql").Fatalf("mysql database connect error: %+v", err)
	}
	_db = db
	logger.GetLogger().WithField("engine", "mysql").Info("init database connection success")
}

func initDB(connString string, maxIdleCon, maxOpenCon int) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connString, // DSN data source name
		DefaultStringSize:         256,        // string 类型字段的默认长度
		DisableDatetimePrecision:  true,       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		SkipInitializeWithVersion: false,      // 根据版本自动配置
	}), &gorm.Config{
		Logger: NewGormLogger(time.Second / 2),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		// 启用更新时间戳功能
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleCon) // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(maxOpenCon) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 60)
	return db.Debug(), nil // 开启sql语句打印
}

func GetDB() *gorm.DB {
	return _db
}

func NewDB(host,
	database,
	username,
	password string,
	port,
	maxIdleCon,
	maxOpenCon int) (*gorm.DB, error) {
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", strconv.Itoa(port), ")/", database, "?charset=utf8mb4&timeout=30s&parseTime=true&&loc=Asia%2FShanghai"}, "")
	return initDB(dsn, maxIdleCon, maxOpenCon)
}

func GetEcdnDB() *gorm.DB {
	return _ecdnDB
}

func SetupEcdnDB(
	host,
	database,
	username,
	password string,
	port,
	maxIdleCon,
	maxOpenCon int,
) {
	db, err := NewDB(host, database, username, password, port, maxIdleCon, maxOpenCon)
	if err != nil {
		logger.GetLogger().WithField("engine", "mysql").Fatalf("mysql database connect error: %+v", err)
	}
	_ecdnDB = db
	logger.GetLogger().WithField("engine", "mysql").Info("init ecdn database connection success")
}

package clickhouse

import (
	"biz-auto-api/pkg/logger"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/driver/clickhouse"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CK struct {
	*gorm.DB
}

var _ck *CK

func Setup(
	host,
	database,
	username,
	password string,
	port,
	maxIdleCon,
	maxOpenCon int,
	readOnly bool,
) {
	ck, err := NewClickhouse(host, database, username, password, port, maxIdleCon, maxOpenCon, readOnly)
	if err != nil {
		logger.GetLogger().WithField("engine", "clickhouse").Fatalf("click database connect error: %+v", err)
	}
	_ck = ck
	logger.GetLogger().WithField("engine", "clickhouse").Info("init click connection success")
}

func initDB(connString string, maxIdleCon, maxOpenCon int) (*CK, error) {
	ck, err := gorm.Open(clickhouse.Open(connString), &gorm.Config{
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
	sqlDB, _ := ck.DB()
	sqlDB.SetMaxIdleConns(maxIdleCon) // 设置连接池，空闲
	sqlDB.SetMaxOpenConns(maxOpenCon) // 打开
	sqlDB.SetConnMaxLifetime(time.Second * 60)
	return &CK{
		DB: ck.Debug(),
	}, nil // 开启sql语句打印
}

func GetCK() *CK {
	return _ck
}
func (c *CK) WithContext(ctx context.Context) *CK {
	return &CK{
		DB: c.DB.WithContext(ctx),
	}
}

func NewClickhouse(host,
	database,
	username,
	password string,
	port,
	maxIdleCon,
	maxOpenCon int, readOnly bool) (*CK, error) {
	u := &url.URL{
		Scheme: "clickhouse",
		User:   url.UserPassword(username, password),
		Host:   host + ":" + strconv.Itoa(port),
		Path:   "/" + database,
	}
	// 添加查询参数
	params := url.Values{}
	params.Add("dial_timeout", "30s")
	params.Add("read_timeout", "120s")
	params.Add("max_open_conns", strconv.Itoa(maxIdleCon))
	params.Add("max_idle_conns", strconv.Itoa(maxOpenCon))
	u.RawQuery = params.Encode()
	if !readOnly {
		params.Add("max_query_size", "1024000000")
	}

	return initDB(u.String(), maxIdleCon, maxOpenCon)
}

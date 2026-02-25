package api

import (
	"biz-auto-api/pkg/clickhouse"
	"biz-auto-api/pkg/consts"
	db "biz-auto-api/pkg/db"
	pkgredis "biz-auto-api/pkg/redis"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type Api struct {
	Context *gin.Context
	Logger  *logrus.Entry
	DB      *gorm.DB
	EcdnDB  *gorm.DB
	Redis   *redis.Client
	CK      *clickhouse.CK
	Errors  error
}

func (e *Api) AddError(err error) {
	if e.Errors == nil {
		e.Errors = err
	} else if err != nil {
		e.Logger.Error(err)
		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	}
}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	e.Logger = GetRequestLogger(c)
	return e
}

// GetLogger 获取上下文提供的日志
func (e *Api) GetLogger() *logrus.Entry {
	return GetRequestLogger(e.Context)
}

// Bind 参数校验
func (e *Api) Bind(d interface{}, bindings ...binding.Binding) *Api {
	var err error
	for i := range bindings {
		if bindings[i] == nil {
			err = e.Context.ShouldBindUri(d)
		} else {
			err = e.Context.ShouldBindWith(d, bindings[i])
		}
		if err != nil && err.Error() == "EOF" {
			e.Logger.Warn("request body is not present anymore. ")
			err = nil
			continue
		}
		if err != nil {
			e.AddError(err)
			break
		}
	}
	return e
}

// Validate 校验数据
func (e *Api) Validate(struc interface{}) *Api {
	v := validator.New()
	err := v.Struct(struc)
	if err != nil {
		msg := "参数不合法: "
		es := err.(validator.ValidationErrors)
		errFields := []string{}
		for _, e := range es {
			m := ""
			if len(e.Param()) > 0 && !strings.Contains(e.ActualTag(), e.Param()) {
				m = fmt.Sprintf("%v(%v=%v)", e.StructField(), e.ActualTag(), e.Param())
			} else {
				m = fmt.Sprintf("%v(%v)", e.StructField(), e.ActualTag())
			}
			if !tools.InSlice(m, errFields) {
				errFields = append(errFields, m)
			}
		}
		// 拼接不符合要求的字段
		msg += strings.Join(errFields, ",")
		e.AddError(errors.New(msg))
	}
	return e
}

// GetOrm 获取Orm DB
func (e *Api) GetOrm() *gorm.DB {
	return e.DB
}

// MakeOrm 设置Orm DB
func (e *Api) MakeOrm() *Api {
	//var err error
	//if e.Logger == nil {
	//	err = errors.New("at MakeOrm logger is nil")
	//	e.AddError(err)
	//	return e
	//}
	e.DB = db.GetDB().WithContext(e.Context)
	if db.GetEcdnDB() != nil {
		e.EcdnDB = db.GetEcdnDB().WithContext(e.Context)
	}
	e.Redis = pkgredis.GetClient()
	e.CK = clickhouse.GetCK()
	return e
}

func (e *Api) MakeService(c *service.Service) *Api {
	c.Log = e.Logger
	c.DB = e.DB
	c.EcdnDB = e.EcdnDB
	c.Ctx = e.Context
	c.Redis = e.Redis
	c.CK = e.CK
	return e
}

func (e *Api) OKWithBizCode(code consts.BizCode, msg string) {
	OKWithBizCode(e.Context, code, msg)
}

func (e *Api) Error(code int, msg string) {
	Error(e.Context, code, msg)
}

// OK 通常成功数据处理
func (e *Api) OK(data interface{}, msg string) {
	OK(e.Context, data, msg)
}

// PageOK 分页数据处理
func (e *Api) PageOK(result interface{}, count, pageIndex, pageSize int64, msg string) {
	PageOK(e.Context, result, count, pageIndex, pageSize, msg)
}

func (e *Api) OKWithCodeAndData(code consts.BizCode, data interface{}, msg string) {
	OKWithCodeAndData(e.Context, code, data, msg)
}

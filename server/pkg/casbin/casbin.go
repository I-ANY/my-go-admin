package casbin

import (
	"biz-auto-api/pkg/logger"
	pkgmodels "biz-auto-api/pkg/models"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	ormlogger "gorm.io/gorm/logger"
)

var (
	_enforcer *casbin.SyncedEnforcer
	once      = sync.Once{}
)

func Setup(db *gorm.DB) {
	db = db.Session(&gorm.Session{
		Logger: ormlogger.Default.LogMode(ormlogger.Silent),
	})
	log := logger.GetLogger().WithField("engine", "casbin")
	once.Do(func() {
		md, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && ignoreCase(r.act, p.act)
`)
		if err != nil {
			log.Fatalf("New casbin model failed: %+v", errors.WithStack(err))
		}
		csbRuleModel := pkgmodels.CasbinRule{}
		adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, csbRuleModel, csbRuleModel.TableName())
		if err != nil {
			log.Fatalf("New casbin adapter failed: %+v", errors.WithStack(err))
		}
		enforcer, err := casbin.NewSyncedEnforcer(md, adapter)
		if err != nil {
			log.Fatalf("New casbin enforcer failed: %+v", errors.WithStack(err))
		}
		enforcer.AddFunction("ignoreCase", IgnoreCaseFunc)
		err = enforcer.LoadPolicy()
		if err != nil {
			log.Fatalf("Load casbin policy failed: %+v", errors.WithStack(err))
		}
		_enforcer = enforcer
		go loadPolicyCyclic()
	})
	log.Info("Init casbin enforcer success")
}

func GetEnforcer() *casbin.SyncedEnforcer {
	return _enforcer
}

func loadPolicyCyclic() {
	for {
		select {
		case <-time.After(time.Second * 15):
			err := _enforcer.LoadPolicy()
			if err != nil {
				logger.GetLogger().WithField("engine", "casbin").Errorf("Load casbin policy failed: %+v", errors.WithStack(err))
			}
		}
	}
}

// IgnoreCase key不区分大小写
func IgnoreCase(key1 string, key2 string) bool {
	return strings.ToLower(key1) == strings.ToLower(key2)
}
func IgnoreCaseFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	return (bool)(IgnoreCase(name1, name2)), nil
}

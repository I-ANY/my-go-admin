package migration

import (
	"biz-auto-api/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"path/filepath"
	"sort"
	"sync"
)

var Migrate = &Migration{
	version: make(map[string]func(db *gorm.DB, version string) error),
}

type Migration struct {
	db      *gorm.DB
	version map[string]func(db *gorm.DB, version string) error
	mutex   sync.Mutex
}

func (e *Migration) GetDb() *gorm.DB {
	return e.db
}

func (e *Migration) SetDb(db *gorm.DB) {
	e.db = db
}

func (e *Migration) SetVersion(k string, f func(db *gorm.DB, version string) error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.version[k] = f
}

func (e *Migration) Migrate() {
	versions := make([]string, 0)
	for k := range e.version {
		versions = append(versions, k)
	}
	if !sort.StringsAreSorted(versions) {
		sort.Strings(versions)
	}
	var err error
	for _, v := range versions {
		var count int64
		err = e.db.Table("sys_migration").Where("version = ?", v).Count(&count).Error
		if err != nil {
			logger.GetLogger().Fatalf("%+v", errors.WithStack(err))
		}
		if count > 0 { // 已经执行过迁移，跳过
			continue
		}
		err = (e.version[v])(e.db.Debug(), v)
		if err != nil {
			logger.GetLogger().Fatalf("%+v", errors.WithStack(err))
		}
	}
}

func GetFilename(s string) string {
	s = filepath.Base(s)
	return s[:13]
}

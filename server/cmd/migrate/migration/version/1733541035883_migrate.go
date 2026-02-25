package version

import (
	"biz-auto-api/cmd/migrate/migration"
	"biz-auto-api/internal/models"
	pkgmodels "biz-auto-api/pkg/models"
	"runtime"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1733541035883Test)
}

func _1733541035883Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Migrator().AutoMigrate(
			new(pkgmodels.CasbinRule),
			new(models.SysApi),
			new(models.SysMenu),
			new(models.SysDept),
			new(models.SysRole),
			new(models.SysUser),
			new(models.SysDictType),
			new(models.SysDictData),
			new(models.SysOperaLog),
			new(models.CjJob),
			new(models.CjJobExecRecord),
			new(models.CjJobExecLog),
			new(models.SysOperaLog),
		)
		if err != nil {
			return err
		}
		return tx.Create(&pkgmodels.Migration{
			Version: version,
		}).Error
	})
}

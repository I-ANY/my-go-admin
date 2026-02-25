package version

import (
	"biz-auto-api/cmd/migrate/migration"
	"biz-auto-api/internal/models"
	pkgmodels "biz-auto-api/pkg/models"
	"gorm.io/gorm"
	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1756785431429Test)
}

func _1756785431429Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var err error
		err = tx.Debug().Migrator().AutoMigrate(
			new(models.SysRoleResource),
			new(models.SysResourceType),
			new(models.SysResourceTypeField),
		)
		if err != nil {
			return err
		}

		return tx.Create(&pkgmodels.Migration{
			Version: version,
		}).Error
	})
}

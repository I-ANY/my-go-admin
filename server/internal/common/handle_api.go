package common

import (
	"biz-auto-api/internal/models"
	pkgdb "biz-auto-api/pkg/db"
	"biz-auto-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

func AddOrUpdateApis(engin *gin.Engine) error {
	db := pkgdb.GetDB()
	routeInfos := engin.Routes()
	return db.Transaction(func(tx *gorm.DB) error {
		for _, routeInfo := range routeInfos {
			s := strings.Split(routeInfo.Handler, ".")
			handler := s[1] + "." + s[2]
			if len(handler) >= 3 && strings.HasSuffix(handler, "-fm") {
				handler = handler[0 : len(handler)-3]
			}
			modelApi := &models.SysApi{
				Path:    routeInfo.Path,
				Method:  routeInfo.Method,
				Handler: handler,
			}
			var total int64
			var a *models.SysApi
			err := tx.Model(&models.SysApi{}).Where("path = ? and  method = ?", modelApi.Path, modelApi.Method).Count(&total).Find(&a).Error
			if err != nil {
				return errors.Wrap(err, "查询Api信息失败")
			}
			if total == 0 {
				err = tx.Model(&models.SysApi{}).Create(&modelApi).Error
			} else {
				err = tx.Model(&models.SysApi{}).Where("id=?", a.Id).Updates(&modelApi).Error
			}
			if err != nil {
				return errors.Wrap(err, "添加或者更改Api信息失败")
			}
		}
		logger.GetLogger().Info("add or update apis success")
		return nil
	})
}

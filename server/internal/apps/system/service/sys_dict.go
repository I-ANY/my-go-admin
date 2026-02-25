package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/service"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysDict struct {
	service.Service
}

func (s SysDict) GetAllDictData(req *dto.GetAllDictDataReq) ([]*dto.DictType, error) {
	var result = make([]*dto.DictType, 0)
	db := s.DB
	log := s.Log
	var dictTypes = make([]*models.SysDictType, 0)

	err := db.Model(models.SysDictType{}).Preload("DictData", func(tx *gorm.DB) *gorm.DB {
		// 只查询启用的数据
		return tx.Where("status = ?", models.DictDataStatusEnable).Order("sort asc ,id asc")
	}).Where("status = ?", models.DictTypeStatusEnable).Order("sort asc ,id asc").Find(&dictTypes).Error
	if err != nil {
		err = errors.Wrap(err, "query all dict data fail")
		log.Errorf("%+v", err)
		return nil, err
	}
	err = copier.Copy(&result, dictTypes)
	if err != nil {
		err = errors.Wrap(err, "convert dict data failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	return result, nil
}

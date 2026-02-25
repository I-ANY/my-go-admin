package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysApi struct {
	service.Service
}

func (s SysApi) GetApiList(req *dto.GetApiListReq) (*dto.GetApiListRes, error) {
	var result = &dto.GetApiListRes{}
	db := s.DB
	log := s.Log
	var total int64
	var apis []*models.SysApi
	err := db.Model(models.SysApi{}).Scopes(func(tx *gorm.DB) *gorm.DB {
		if len(req.Search) > 0 {
			tx.Where("path like ? or handler like ?", tools.FuzzyQuery(req.Search), tools.FuzzyQuery(req.Search))
		}
		return tx
	}).Count(&total).Scopes(req.MakePagination()).Find(&apis).Error
	if err != nil {
		err = errors.Wrap(err, "list apis failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items []*dto.SysApi
	err = copier.Copy(&items, apis)
	if err != nil {
		err = errors.Wrap(err, "convert apis failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.Items = items
	result.Total = total
	return result, nil
}

package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/service"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type SysDept struct {
	service.Service
}

func (s SysDept) GetDeptTree(req *dto.GetDeptTreeReq) (*dto.GetDeptTreeRes, error) {

	var result = &dto.GetDeptTreeRes{}
	db := s.DB
	log := s.Log
	var total int64
	var depts []*models.SysDept
	err := db.Model(&models.SysDept{}).Count(&total).Order("order_no desc").Find(&depts).Error
	if err != nil {
		err = errors.Wrap(err, "query dept failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items = make([]*dto.Dept, 0, len(depts))
	err = copier.Copy(&items, depts)
	if err != nil {
		err = errors.Wrap(err, "convert dept failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	result.Items = s.BuildDeptTree(items, 0)
	result.Total = total
	return result, nil
}

func (s SysDept) BuildDeptTree(depts []*dto.Dept, parentDept int64) []*dto.Dept {
	var result = make([]*dto.Dept, 0)
	for _, dept := range depts {
		if dept.ParentDept == parentDept {
			result = append(result, dept)
			dept.Children = s.BuildDeptTree(depts, dept.Id)
		}
	}
	return result
}

func GetDeptChildrenIds(depts []*models.SysDept, deptId int64) []int64 {
	childrenIds := make([]int64, 0)
	for _, dept := range depts {
		if dept.ParentDept == deptId {
			childrenIds = append(childrenIds, dept.Id)
			childrenIds = append(childrenIds, GetDeptChildrenIds(depts, dept.Id)...)
		}
	}
	return childrenIds
}

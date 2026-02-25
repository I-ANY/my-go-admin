package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/common"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/casbin"
	"biz-auto-api/pkg/consts"
	pkgmodels "biz-auto-api/pkg/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysRole struct {
	service.Service
}

func (s SysRole) GetRoleList(req *dto.GetRoleListReq) (*dto.GetRoleListRes, error) {
	var result = &dto.GetRoleListRes{}
	db := s.DB
	log := s.Log
	var total int64
	var roles []*models.SysRole
	err := db.Model(models.SysRole{}).Preload("Menus").Scopes(func(tx *gorm.DB) *gorm.DB {
		if len(req.Search) > 0 {
			tx.Where("name like ? or identify like ?", tools.FuzzyQuery(req.Search), tools.FuzzyQuery(req.Search))
		}
		if req.Status > 0 {
			tx.Where("status = ?", req.Status)
		}
		return tx
	}).Count(&total).Scopes(req.MakePagination()).Find(&roles).Error
	if err != nil {
		err = errors.Wrap(err, "list roles failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, role := range roles {
		item := &dto.Role{
			Id:       role.Id,
			Name:     role.Name,
			Status:   role.Status,
			Identify: role.Identify,
			Remark:   role.Remark,
		}
		for _, menu := range role.Menus {
			item.MenuIds = append(item.MenuIds, menu.Id)
		}
		result.Items = append(result.Items, item)
	}
	result.Total = total
	return result, nil
}

func (s SysRole) UpdateRole(req *dto.UpdateRoleReq) (*dto.UpdateRoleRes, error) {
	var result = &dto.UpdateRoleRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var role = models.SysRole{}
	err := db.Model(role).First(&role, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	if role.Identify == consts.AdminRoleIdentify {
		return nil, errors.New("管理员和基础角色不允许被修改")
	}

	menus := make([]*models.SysMenu, 0, len(req.MenuIds))
	for _, menuId := range req.MenuIds {
		menu := models.SysMenu{}
		err = db.First(&menu, menuId).Error
		if err != nil {
			err = errors.Wrapf(err, "query menu(id=%v) failed", menuId)
			log.Errorf("%+v", err)
			return nil, err
		}
		menus = append(menus, &models.SysMenu{
			BseId: pkgmodels.BseId{
				Id: menuId,
			}})
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		data := tools.StructToMap(req, "", false, "MenuIds")
		data["update_by"] = userId
		err = tx.Model(&role).Where("id = ?", req.Id).Updates(data).Error
		if err != nil {
			err = errors.Wrap(err, "update role failed")
			log.Errorf("%+v", err)
			return err
		}
		if req.MenuIds != nil {
			err = tx.Model(&role).Association("Menus").Clear()
			if err != nil {
				err = errors.Wrap(err, "update role's menu failed")
				log.Errorf("%+v", err)
				return err
			}
			err = tx.Model(&role).Association("Menus").Append(menus)
			if err != nil {
				err = errors.Wrap(err, "update role's menu failed")
				log.Errorf("%+v", err)
				return err
			}
		}
		// 修改角色可能导致角色的权限变更，需要刷新角色的权限
		err = common.RefreshRolesPermission([]string{role.Identify}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh role permission failed")
			log.Errorf("%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	go casbin.GetEnforcer().LoadPolicy()
	return result, nil
}

func (s SysRole) AddRole(req *dto.AddRoleReq) (*dto.AddRoleRes, error) {
	var result = &dto.AddRoleRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var role = &models.SysRole{
		ControlBy: pkgmodels.ControlBy{
			UpdateBy: userId, CreateBy: userId,
		},
	}
	err := copier.Copy(role, req)
	if err != nil {
		err = errors.Wrap(err, "convert role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	menus := make([]*models.SysMenu, 0, len(req.MenuIds))
	for _, menuId := range req.MenuIds {
		menu := models.SysMenu{}
		err = db.First(&menu, menuId).Error
		if err != nil {
			err = errors.Wrapf(err, "query menu(id=%v) failed", menuId)
			log.Errorf("%+v", err)
			return nil, err
		}
		menus = append(menus, &models.SysMenu{
			BseId: pkgmodels.BseId{
				Id: menuId,
			}})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.SysRole{}).Create(role).Error
		if err != nil {
			err = errors.Wrap(err, "Create role failed")
			log.Errorf("%+v", err)
			return err
		}

		err = tx.Model(role).Association("Menus").Append(menus)
		if err != nil {
			err = errors.Wrap(err, "add role's menu failed")
			log.Errorf("%+v", err)
			return err
		}
		// 修改角色可能导致角色的权限变更，需要刷新角色的权限
		err = common.RefreshRolesPermission([]string{req.Identify}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh role permission failed")
			log.Errorf("%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	go casbin.GetEnforcer().LoadPolicy()
	return result, nil
}

func (s SysRole) DeleteRole(req *dto.DeleteRoleReq) (*dto.DeleteRoleRes, error) {
	db := s.DB
	log := s.Log
	var result = &dto.DeleteRoleRes{}
	var role = models.SysRole{}
	err := db.Model(&role).First(&role, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query role failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	if role.Identify == consts.AdminRoleIdentify || role.Identify == consts.BaseRoleIdentify {
		return nil, errors.New("管理员和基础角色不允许被删除")
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(models.SysRole{}).Unscoped().Delete(&role, req.Id).Error
		if err != nil {
			err = errors.Wrap(err, "delete role failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&role).Association("Menus").Clear()
		if err != nil {
			err = errors.Wrap(err, "delete role's menu failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&role).Association("Users").Clear()
		if err != nil {
			err = errors.Wrap(err, "delete role's user failed")
			log.Errorf("%+v", err)
			return err
		}
		// 删除角色可能导致角色的权限变更，需要刷新角色的权限
		err = common.RefreshRolesPermission([]string{role.Identify}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh role permission failed")
			log.Errorf("%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	go casbin.GetEnforcer().LoadPolicy()
	return result, nil
}

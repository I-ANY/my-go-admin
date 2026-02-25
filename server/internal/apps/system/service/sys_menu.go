package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/common"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/casbin"
	pkgmodels "biz-auto-api/pkg/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SysMenu struct {
	service.Service
}

func (s SysMenu) GetMenuTree(req *dto.GetMenuTreeReq) (*dto.GetMenuTreeRes, error) {
	var result = &dto.GetMenuTreeRes{}
	db := s.DB
	log := s.Log
	var dbMenus []*models.SysMenu
	var total int64
	err := db.Model(&models.SysMenu{}).Preload("Apis").Count(&total).Order("order_no desc,id asc").Find(&dbMenus).Error
	if err != nil {
		log.Errorf("Query all menu failed: %+v", errors.WithStack(err))
		return nil, errors.WithMessage(errors.WithStack(err), "Query all menu failed")
	}
	menus := BuildMenu(dbMenus)
	result.Items = menus
	result.Total = total
	return result, nil
}

func (s SysMenu) UpdateMenu(req *dto.UpdateMenuReq) (*dto.UpdateMenuRes, error) {
	var result = &dto.UpdateMenuRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var menu = models.SysMenu{}

	err := db.Preload("Apis").Preload("Roles").Model(&models.SysMenu{}).First(&menu, req.Id).Error
	if err != nil {
		err = errors.Wrapf(err, "query menu failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	// 获取这个菜单的修改会影响到那些角色的权限
	roleIdentifies := make([]string, 0)
	for _, role := range menu.Roles {
		roleIdentifies = append(roleIdentifies, role.Identify)
	}
	roleIdentifies = tools.RemoveDuplication(roleIdentifies)
	// 更新菜单关联api
	apis := make([]*models.SysApi, 0, len(req.ApiIds))
	for _, apiId := range req.ApiIds {
		api := models.SysApi{}
		err = db.First(&api, apiId).Error
		if err != nil {
			err = errors.Wrapf(err, "query api(id=%v) failed", apiId)
			log.Errorf("%+v", err)
			return nil, err
		}
		apis = append(apis, &models.SysApi{
			BseId: pkgmodels.BseId{
				Id: apiId,
			},
		})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		data := tools.StructToMap(req, "", false, "ApiIds")
		data["update_by"] = userId
		for key, value := range data {
			if isTrue, ok := value.(*bool); ok {
				if *isTrue {
					data[key] = 1
				} else {
					data[key] = 0
				}
			}
		}
		err = tx.Model(&models.SysMenu{}).Where("id=?", req.Id).Updates(data).Error
		if err != nil {
			err = errors.Wrapf(err, "update menu failed")
			log.Errorf("%+v", err)
			return err
		}
		if req.ApiIds != nil {
			err = tx.Model(&menu).Association("Apis").Clear()
			if err != nil {
				err = errors.Wrapf(err, "clear menu's api failed")
				log.Errorf("%+v", err)
				return err
			}
			err = tx.Model(&menu).Association("Apis").Append(apis)
			if err != nil {
				err = errors.Wrapf(err, "add menu's api failed")
				log.Errorf("%+v", err)
				return err
			}
		}
		// 修改菜单可能导致角色的权限变更，需要刷新角色的权限
		err = common.RefreshRolesPermission(roleIdentifies, tx)
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
	// 重新加载权限
	go casbin.GetEnforcer().LoadPolicy()
	return result, nil
}

func (s SysMenu) AddMenu(req *dto.AddMenuReq) (*dto.AddMenuRes, error) {
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var result = &dto.AddMenuRes{}

	var menu = models.SysMenu{}
	err := copier.Copy(&menu, req)
	if err != nil {
		err = errors.Wrap(err, "convert menu failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	menu.CreateBy = userId
	menu.UpdateBy = userId
	apis := make([]*models.SysApi, 0, len(req.ApiIds))
	for _, apiId := range req.ApiIds {
		api := models.SysApi{}
		err = db.First(&api, apiId).Error
		if err != nil {
			err = errors.Wrapf(err, "query api(id=%v) failed", apiId)
			log.Errorf("%+v", err)
			return nil, err
		}
		apis = append(apis, &models.SysApi{
			BseId: pkgmodels.BseId{
				Id: apiId,
			},
		})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&menu).Error
		if err != nil {
			err = errors.Wrapf(err, "add menu failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&menu).Association("Apis").Append(apis)
		if err != nil {
			err = errors.Wrapf(err, "add menu's api failed")
			log.Errorf("%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s SysMenu) DeleteMenu(req *dto.DeleteMenuReq) (*dto.DeleteMenuRes, error) {
	db := s.DB
	log := s.Log
	//userId := s.GetCurrentUserId()
	var result = &dto.DeleteMenuRes{
		Success: true,
	}
	var menu = models.SysMenu{}

	var total int64
	// 查询是否有子菜单
	err := db.Model(&menu).Where("parent_menu=?", req.Id).Count(&total).Error
	if err != nil {
		err = errors.Wrapf(err, "Query children menu failed")
		log.Errorf("： %+v", err)
		return nil, err
	}
	if total > 0 {
		result.Success = false
		result.Message = "请先删除子菜单"
		return result, nil
	}

	err = db.Preload("Apis").Preload("Roles").Model(&models.SysMenu{}).First(&menu, req.Id).Error
	if err != nil {
		err = errors.Wrapf(err, "query menu failed")
		log.Errorf("%+v", err)
		return nil, err
	}

	// 获取这个菜单的修改会影响到那些角色的权限
	roleIdentifies := make([]string, 0)
	for _, role := range menu.Roles {
		roleIdentifies = append(roleIdentifies, role.Identify)
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		// 物理删除
		err = tx.Model(&menu).Where("id=?", req.Id).Unscoped().Delete(&menu).Error
		if err != nil {
			err = errors.Wrapf(err, "delete menu failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&menu).Association("Apis").Clear()
		if err != nil {
			err = errors.Wrapf(err, "clear menu's apis failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&menu).Association("Roles").Clear()
		if err != nil {
			err = errors.Wrapf(err, "clear menu's roles failed")
			log.Errorf("%+v", err)
			return err
		}
		// 修改菜单可能导致角色的权限变更，需要刷新角色的权限
		err = common.RefreshRolesPermission(roleIdentifies, tx)
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

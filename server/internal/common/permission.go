package common

import (
	"biz-auto-api/internal/models"
	pkgmodels "biz-auto-api/pkg/models"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

func RefreshUsersRole(userIds []int64, db *gorm.DB) error {
	// 删除用户的所有角色绑定
	err := db.Model(&pkgmodels.CasbinRule{}).Where("ptype = 'g' and v0 in ?", userIds).Delete(&pkgmodels.CasbinRule{}).Error
	if err != nil {
		return errors.WithStack(err)
	}
	// 查询用户信息和关联的角色信息，被禁用的用户和角色不加进去
	var users = make([]*models.SysUser, 0, len(userIds))
	err = db.Model(models.SysUser{}).Preload("Roles", func(tx *gorm.DB) *gorm.DB {
		return tx.Where(" status = ? ", models.RoleStatusEnable)
	}).Where(" status = ?  and id  in ?", models.UserStatusEnable, userIds).Find(&users).Error
	if err != nil {
		return errors.WithStack(err)
	}
	ruleMap := make(map[string]*pkgmodels.CasbinRule)
	for _, user := range users {
		for _, role := range user.Roles {
			r := &pkgmodels.CasbinRule{
				Ptype: "g",
				V0:    strconv.Itoa(int(user.Id)),
				V1:    role.Identify,
			}
			// 用map去重，防止后面批量添加唯一键报错
			ruleMap[fmt.Sprintf("%v__%v", user.Id, role.Identify)] = r
		}
	}
	rules := make([]*pkgmodels.CasbinRule, 0)
	for _, r := range ruleMap {
		rules = append(rules, r)
	}
	err = db.Model(pkgmodels.CasbinRule{}).CreateInBatches(rules, 50).Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func RefreshRolesPermission(roleIdentifies []string, db *gorm.DB) error {
	// 删除角色API的权限
	err := db.Model(&pkgmodels.CasbinRule{}).Where("ptype = 'p' and v0  in ?", roleIdentifies).Delete(&pkgmodels.CasbinRule{}).Error
	if err != nil {
		return errors.WithStack(err)
	}
	// 查询这些角色最新API权限信息，被禁用的角色和菜单不添加进权限表
	var roles = make([]*models.SysRole, 0)
	err = db.Model(&models.SysRole{}).Preload("Menus", func(tx *gorm.DB) *gorm.DB {
		//被禁用的菜单权限不添加
		return tx.Where("status = ? ", models.MenuStatusEnable)
	}).Preload("Menus.Apis").Where("status = ? and identify in ? ", models.RoleStatusEnable, roleIdentifies).Find(&roles).Error
	if err != nil {
		return errors.WithStack(err)
	}
	rolePermissions := make(map[string]map[int64]*models.SysApi)
	for _, role := range roles {
		roleApis := make(map[int64]*models.SysApi)
		for _, menu := range role.Menus {
			for _, api := range menu.Apis {
				roleApis[api.Id] = api // 多个菜单可能会需要一个菜单的权限，所以这里API重复则覆盖
			}
		}
		rolePermissions[role.Identify] = roleApis
	}
	ruleMap := make(map[string]*pkgmodels.CasbinRule)
	for roleIdentify, apis := range rolePermissions {
		for _, api := range apis {
			r := &pkgmodels.CasbinRule{
				Ptype: "p",
				V0:    roleIdentify,
				V1:    api.Path,
				V2:    api.Method,
			}
			// 用map去重，防止后面批量添加报错
			ruleMap[fmt.Sprintf("%v__%v__%v", roleIdentify, api.Path, api.Method)] = r
		}
	}
	rules := make([]*pkgmodels.CasbinRule, 0)
	for _, r := range ruleMap {
		rules = append(rules, r)
	}
	err = db.Model(pkgmodels.CasbinRule{}).CreateInBatches(rules, 50).Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

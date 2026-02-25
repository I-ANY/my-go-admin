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

type SysUser struct {
	service.Service
}

func (s SysUser) GetUserList(req *dto.GetListUserReq) (*dto.GetListUserRes, error) {
	var result = &dto.GetListUserRes{}
	db := s.DB
	log := s.Log
	var total int64
	var users []*models.SysUser
	var deptIds []int64
	var err error

	// 查询部门下的 子部门ID + 当前请求的部门ID
	if req.DeptId > 0 {
		var depts []*models.SysDept
		err := db.Model(&models.SysDept{}).Find(&depts).Error
		if err != nil {
			err = errors.Wrap(err, "query dept failed")
			log.Errorf("%+v", err)
			return nil, err
		}
		deptIds = GetDeptChildrenIds(depts, req.DeptId)
		deptIds = append(deptIds, req.DeptId)
	}
	if req.RoleId == 0 {
		// 单表查询
		users, total, err = s.commonQueryUser(req, deptIds)
	} else {
		// 连表查询
		users, total, err = s.roleQueryUser(req, deptIds)
	}
	if err != nil {
		err = errors.WithMessage(err, "query user failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items = make([]*dto.User, 0, len(users))
	for _, user := range users {
		item := &dto.User{
			Id:       user.Id,
			Username: user.Username,
			NickName: user.NickName,
			Email:    tools.ToValue(user.Email),
			Tel:      tools.ToValue(user.Tel),
			Status:   user.Status,
			Source:   user.Source,
		}
		if user.Dept != nil {
			dept := &dto.Dept{}
			err = copier.Copy(dept, user.Dept)
			if err != nil {
				err = errors.Wrap(err, "convert dept failed")
				log.Errorf("%+v", err)
				return nil, err
			}
			item.DeptId = user.DeptId
			item.Dept = dept
		}
		if len(user.Roles) > 0 {
			roles := make([]*dto.Role, 0, len(user.Roles))
			err = copier.Copy(&roles, user.Roles)
			if err != nil {
				err = errors.Wrap(err, "convert roles failed")
				log.Errorf("%+v", err)
				return nil, err
			}
			item.Roles = roles
			for _, role := range user.Roles {
				item.RoleIds = append(item.RoleIds, role.Id)
			}
		}
		items = append(items, item)
	}
	result.Items = items
	result.Total = total
	return result, nil
}

func (s SysUser) UpdateUser(req *dto.UpdateUserReq) (*dto.UpdateUserRes, error) {
	var result = &dto.UpdateUserRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var user = models.SysUser{}

	if req.Id == s.GetCurrentUserId() && req.Status != nil && *req.Status == models.UserStatusDisable {
		return nil, errors.New("用户不允许禁用自己")
	}

	err := db.Model(&user).Where("id =? ", req.Id).First(&user).Error
	if err != nil {
		err = errors.Wrap(err, "query user failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	data := tools.StructToMap(req, "", false, "RoleIds", "Password")
	data["update_by"] = userId
	if req.Password != nil && len(*req.Password) > 0 {
		data["password"] = tools.GetEncryptedPassword(*req.Password)
	}
	roles := make([]*models.SysRole, 0, len(req.RoleIds))
	for _, roleId := range req.RoleIds {
		role := models.SysRole{}
		err = db.First(&role, roleId).Error
		if err != nil {
			err = errors.Wrapf(err, "query role(id=%v) failed", roleId)
			log.Errorf("%+v", err)
			return nil, err
		}
		roles = append(roles, &models.SysRole{
			BseId: pkgmodels.BseId{
				Id: roleId,
			}})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		// 编辑用户信息
		err = tx.Model(&user).Where("id = ? ", req.Id).Updates(data).Error
		if err != nil {
			err = errors.Wrap(err, "update user failed")
			log.Errorf("%+v", err)
			return err
		}
		// 未传RoleIds字段则不更新
		if req.RoleIds != nil {
			// 编辑用户关联的角色信息
			err = tx.Model(&user).Association("Roles").Clear()
			if err != nil {
				err = errors.Wrap(err, "update user's role failed")
				log.Errorf("%+v", err)
				return err
			}
			err = tx.Model(&user).Association("Roles").Append(roles)
			if err != nil {
				err = errors.Wrap(err, "update user's role failed")
				log.Errorf("%+v", err)
				return err
			}
		}
		// 需要刷新权限表
		err = common.RefreshUsersRole([]int64{req.Id}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh user's role relation failed")
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

func (s SysUser) AddUser(req *dto.AddUserReq) (*dto.AddUserRes, error) {
	var result = &dto.AddUserRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var user = &models.SysUser{
		ControlBy: pkgmodels.ControlBy{
			UpdateBy: userId, CreateBy: userId,
		},
	}
	err := copier.Copy(user, req)
	if err != nil {
		err = errors.Wrap(err, "convert user failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	user.Password = tools.GetEncryptedPassword(req.Password)
	roles := make([]*models.SysRole, 0, len(req.RoleIds))
	for _, roleId := range req.RoleIds {
		role := models.SysRole{}
		err = db.First(&role, roleId).Error
		if err != nil {
			err = errors.Wrapf(err, "query role(id=%v) failed", roleId)
			log.Errorf("%+v", err)
			return nil, err
		}
		roles = append(roles, &models.SysRole{
			BseId: pkgmodels.BseId{
				Id: roleId,
			}})
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.SysUser{}).Omit("open_id").Create(user).Error
		if err != nil {
			err = errors.Wrap(err, "Create user failed")
			log.Errorf("%+v", err)
			return err
		}

		err = tx.Model(user).Association("Roles").Append(roles)
		if err != nil {
			err = errors.Wrap(err, "add user's role failed")
			log.Errorf("%+v", err)
			return err
		}
		// 需要刷新权限表
		err = common.RefreshUsersRole([]int64{user.Id}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh user permission failed")
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

func (s SysUser) DeleteUser(req *dto.DeleteUserReq) (*dto.DeleteUserRes, error) {
	db := s.DB
	log := s.Log
	var result = &dto.DeleteUserRes{}
	var user = models.SysUser{}
	err := db.Model(&user).First(&user, req.Id).Error
	if err != nil {
		err = errors.Wrap(err, "query user failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	if req.Id == s.GetCurrentUserId() {
		return nil, errors.New("用户不允许删除自己")
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(models.SysUser{}).Unscoped().Delete(&user, req.Id).Error
		if err != nil {
			err = errors.Wrap(err, "delete user failed")
			log.Errorf("%+v", err)
			return err
		}
		err = tx.Model(&user).Association("Roles").Clear()
		if err != nil {
			err = errors.Wrap(err, "delete user's role failed")
			log.Errorf("%+v", err)
			return err
		}
		// 需要刷新权限表
		err = common.RefreshUsersRole([]int64{req.Id}, tx)
		if err != nil {
			err = errors.WithMessage(err, "refresh user permission failed")
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

func (s SysUser) roleQueryUser(req *dto.GetListUserReq, deptIds []int64) ([]*models.SysUser, int64, error) {
	var (
		total int64
		users []*models.SysUser
		db    = s.DB
	)
	// 子查询，获取角色关联的用户id
	subQuery := db.Table("sys_user_role").Select("sys_user_id").
		Where("sys_role_id = ?", req.RoleId).Distinct("sys_user_id")

	err := db.Model(models.SysUser{}).
		Preload("Roles").
		Preload("Dept").
		Joins("join (?) as r on sys_user.id = r.sys_user_id", subQuery).
		Scopes(func(tx *gorm.DB) *gorm.DB {
			if req.Status != 0 {
				tx.Where("sys_user.status=?", req.Status)
			}
			if len(req.Username) > 0 {
				tx.Where("sys_user.username like ?", tools.FuzzyQuery(req.Username))
			}
			if len(req.NickName) > 0 {
				tx.Where("sys_user.nick_name like ?", tools.FuzzyQuery(req.NickName))
			}
			if len(req.Search) > 0 {
				tx.Where("sys_user.username like ? or sys_user.nick_name like ?", tools.FuzzyQuery(req.Search), tools.FuzzyQuery(req.Search))
			}
			if len(req.Email) > 0 {
				tx.Where("sys_user.email like ?", tools.FuzzyQuery(req.Email))
			}
			if len(req.Tel) > 0 {
				tx.Where("sys_user.tel like ?", tools.FuzzyQuery(req.Tel))
			}
			if len(deptIds) > 0 {
				tx.Where("sys_user.dept_id in ?", deptIds)
			}
			if req.Source > 0 {
				tx.Where("sys_user.source = ?", req.Source)
			}
			return tx
		}).Count(&total).Scopes(req.MakePagination()).Find(&users).Error
	return users, total, errors.WithStack(err)

}

func (s SysUser) commonQueryUser(req *dto.GetListUserReq, deptIds []int64) ([]*models.SysUser, int64, error) {
	var (
		total int64
		users []*models.SysUser
		db    = s.DB
	)
	err := db.Model(&models.SysUser{}).Preload("Roles").Preload("Dept").Scopes(func(tx *gorm.DB) *gorm.DB {
		if req.Status != 0 {
			tx.Where("status=?", req.Status)
		}
		if len(req.Username) > 0 {
			tx.Where("username like ?", tools.FuzzyQuery(req.Username))
		}
		if len(req.NickName) > 0 {
			tx.Where("nick_name like ?", tools.FuzzyQuery(req.NickName))
		}
		if len(req.Search) > 0 {
			tx.Where("username like ? or nick_name like ?", tools.FuzzyQuery(req.Search), tools.FuzzyQuery(req.Search))
		}

		if len(req.Email) > 0 {
			tx.Where("email like ?", tools.FuzzyQuery(req.Email))
		}
		if len(req.Tel) > 0 {
			tx.Where("tel like ?", tools.FuzzyQuery(req.Tel))
		}
		if len(deptIds) > 0 {
			tx.Where("dept_id in ?", deptIds)
		}
		if req.Source > 0 {
			tx.Where("source = ?", req.Source)
		}
		return tx
	}).Count(&total).Scopes(req.MakePagination()).Find(&users).Error
	return users, total, errors.WithStack(err)
}

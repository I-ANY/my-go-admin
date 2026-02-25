package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/models"
	starPortal "biz-auto-api/pkg/clients/star_portal"
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

type SysBase struct {
	service.Service
}

func (s SysBase) UserLogin(req *dto.SysBaseUserLoginReq) (*dto.SysBaseUserLoginRes, error) {
	var res = &dto.SysBaseUserLoginRes{}
	db := s.DB
	log := s.Log
	var user *models.SysUser
	err := db.Model(&models.SysUser{}).Where("username=? and password=?", req.Username, tools.GetEncryptedPassword(req.Password)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Success = false
			res.Message = "用户名或密码不正确"
			log.Warn("用户名或密码不正确")
			return res, nil
		} else {
			log.Errorf("Login failed: %+v", errors.WithStack(err))
			return nil, errors.Errorf("login failed：%s", err)
		}
	}

	if user.Status == models.UserStatusDisable {
		res.Success = false
		res.Message = "用户已被禁用，请联系管理员"
		log.Warnf("用户 %v 已被禁用，请联系管理员", req.Username)
		return res, nil
	}
	tokenExpireAt := time.Now().Unix() + int64(config.SystemConfig.System.TokenExpireSec)
	token, err := tools.CreateToken(strconv.Itoa(int(user.Id)), user.Username, config.SystemConfig.System.TokenExpireSec)
	if err != nil {
		log.Errorf("Create token failed: %+v", errors.WithStack(err))
		return nil, errors.Errorf("Token创建失败：%s", err)
	}
	log.Infof("user id = %v login success", user.Id)
	s.Ctx.Set(consts.JwtTokenKey, token)
	return &dto.SysBaseUserLoginRes{
		Success:  true,
		Token:    token,
		Message:  "登录成功",
		ExpireAt: tokenExpireAt,
	}, nil
}
func (s SysBase) RefreshToken(req *dto.RefreshTokenReq) (*dto.SysBaseUserLoginRes, error) {
	log := s.Log
	userId := s.GetCurrentUserId()
	tokenExpireSec := config.SystemConfig.System.TokenExpireSec
	tokenExpireAt := time.Now().Add(time.Second * time.Duration(tokenExpireSec)).Unix()
	var user = models.SysUser{}
	if err := s.DB.Where("status = ? ", models.UserStatusEnable).Find(&user, userId).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	token, err := tools.CreateToken(strconv.Itoa(int(userId)), user.Username, tokenExpireSec)
	if err != nil {
		log.Errorf("Create token failed: %+v", errors.WithStack(err))
		return nil, errors.Errorf("Token创建失败：%s", err)
	}
	log.Infof("user id = %v refresh token success", userId)
	return &dto.SysBaseUserLoginRes{
		Success:  true,
		Token:    token,
		Message:  "刷新成功",
		ExpireAt: tokenExpireAt,
	}, nil
}
func (s SysBase) UserLogout() error {
	log := s.Log
	log.Logger.Infof("user id = %v logout sueccess", s.GetCurrentUserId())
	return nil
}

func (s SysBase) GetUserInfo(req *dto.SysBaseGetUserInfoReq) (*dto.SysBaseGetUserInfoRes, error) {
	var result = &dto.SysBaseGetUserInfoRes{}
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var user *models.SysUser
	err := db.Model(&models.SysUser{}).Preload("Roles").Where("id=? and status = ? ", userId, models.UserStatusEnable).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("获取用户信息失败，用户未找到")
			return nil, errors.New("获取用户信息失败，用户未找到")
		}
		log.Errorf("查询用户信息失败：%v", errors.WithStack(err))
		return nil, errors.Errorf("查询用户信息失败：%s", err.Error())
	}
	result.Id = user.Id
	result.Username = user.Username
	result.NickName = user.NickName
	result.Email = tools.ToValue(user.Email)
	result.Tel = tools.ToValue(user.Tel)
	result.Avatar = user.Avatar
	result.Source = user.Source
	for _, r := range user.Roles {
		result.Roles = append(result.Roles, r.Identify)
	}
	log.Info("查询用户信息成功")
	return result, nil
}

func (s SysBase) GetUserPermCode(req *dto.SysBaseGetUserPermCodeReq) ([]string, error) {
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var permissionCodes = make([]string, 0)
	var user *models.SysUser
	log.Infof("查询用户（%v）权限编码", userId)
	err := db.Model(&models.SysUser{}).Preload("Roles", "status=?", models.RoleStatusEnable).
		Preload("Roles.Menus", "status=?", models.MenuStatusEnable).
		Where("id=?", userId).Find(&user).Error
	if err != nil {
		err = errors.Wrapf(err, "list user menu failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	for _, role := range user.Roles {
		// 如果当前用户未管理员的话就直接返回*，代表所有的权限
		if role.Identify == consts.AdminRoleIdentify {
			permissionCodes = []string{"*"}
			return permissionCodes, nil
		}
		for _, menu := range role.Menus {
			if menu.Permission != "" {
				permissionCodes = append(permissionCodes, menu.Permission)
			}
		}
	}
	permissionCodes = tools.RemoveDuplication(permissionCodes)
	return permissionCodes, nil
}

func (s SysBase) GetUserMenu(req *dto.SysBaseGetUserMenuReq) ([]*dto.Menu, error) {
	db := s.DB
	log := s.Log
	userId := s.GetCurrentUserId()
	var user *models.SysUser
	var dbMenus []*models.SysMenu
	var result = make([]*dto.Menu, 0)

	log.Infof("List userId=%v menus", userId)
	err := db.Model(&models.SysUser{}).Preload("Roles", "identify = ?", consts.AdminRoleIdentify).Where("status =? and id = ?", models.UserStatusEnable, userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或者被禁用")
		}
		err = errors.Wrapf(err, "查询用户信息失败")
		log.Errorf("%+v", err)
		return nil, err
	}
	// 用户为管理员 直接查询所有可用的菜单
	if len(user.Roles) > 0 {
		err = db.Model(&models.SysMenu{}).Where("status=? and type in ?", models.MenuStatusEnable, []int{models.MenuTypeDir, models.MenuTypeMenu}).Order("order_no desc").Find(&dbMenus).Error
		if err != nil {
			err = errors.Wrapf(err, "查询用户菜单失败")
			log.Errorf("%+v", err)
			return nil, err
		}
	} else { // 不是管理员 查询用户---> 角色 ---> 菜单
		err = db.Model(&models.SysUser{}).Preload("Roles", "status=?", models.RoleStatusEnable).
			Preload("Roles.Menus", "status=? and type in ?", models.MenuStatusEnable, []int{models.MenuTypeDir, models.MenuTypeMenu}, func(tx *gorm.DB) *gorm.DB {
				return tx.Order("order_no desc")
			}).
			Where("id=?", userId).Find(&user).Error
		if err != nil {
			err = errors.Wrapf(err, "查询用户菜单失败")
			log.Errorf("%+v", err)
			return nil, err
		}
		if user == nil || len(user.Roles) == 0 {
			return result, nil
		}
		dbMenus = GetMenus(user)
	}
	sort.Slice(dbMenus, func(i, j int) bool {
		if dbMenus[i].OrderNo == dbMenus[j].OrderNo {
			return dbMenus[i].Id < dbMenus[j].Id
		}
		return dbMenus[i].OrderNo > dbMenus[j].OrderNo
	})
	result = BuildMenu(dbMenus)
	return result, nil
}

func (s SysBase) GetStarPortalLoginUrl(req *dto.SysBaseGetStarPortalLoginUrlReq) (*dto.SysBaseGetStarPortalLoginUrlRes, error) {
	//db := s.DB
	//log := s.Log
	result := &dto.SysBaseGetStarPortalLoginUrlRes{}
	//redirectUrl := url.QueryEscape(req.RedirectUri)

	// TODO: state 这个字段暂时没有用上，目前只生成一个时间戳
	state := time.Now().Unix()
	starPortalConf := config.SystemConfig.StarPortal
	loginUrl := fmt.Sprintf("%v/application/o/authorize/?redirect_uri=%v&client_id=%v&response_type=code&scope=openid+profile+email&state=%v&prompt=consent",
		starPortalConf.Url,
		req.RedirectUri,
		starPortalConf.ClientId,
		state,
	)
	result.LoginUrl = loginUrl
	return result, nil
}

func (s SysBase) StarPortalLogin(req *dto.SysBaseStarPortalLoginReq) (*dto.SysBaseUserLoginRes, error) {
	db := s.DB
	log := s.Log
	res := &dto.SysBaseUserLoginRes{}

	starPortalConf := config.SystemConfig.StarPortal
	// 通过code从星云平台获取用户token
	starPortalClient, err := starPortal.NewStarPortalClient(starPortalConf.Url, "")
	if err != nil {
		err = errors.WithMessage(err, "new star portal client failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	userToken, err := starPortalClient.GetUserToken(req.Code, starPortalConf.ClientId, starPortalConf.ClientSecret, req.RedirectUri)
	if err != nil {
		err = errors.WithMessage(err, "登录失败，从星云平台获取token失败")
		log.Errorf("%+v", err)
		res.Success = false
		res.Message = "登录失败，从星云平台获取token失败"
		return res, nil
	}
	// 通过用户token从星云平台获取用户信息
	starPortalClient2, _ := starPortal.NewStarPortalClient(starPortalConf.Url, userToken.AccessToken)
	userInfo, err := starPortalClient2.GetUserInfo()
	if err != nil {
		err = errors.WithMessage(err, "登录失败，从星云平台获取用户信息失败")
		log.Errorf("%+v", err)
		res.Success = false
		res.Message = "登录失败，从星云平台获取用户信息失败"
		return res, nil
	}

	// 通过星云平台中的userid在数据库该用户信息，创建token
	email := userInfo.Email
	if len(email) == 0 {
		res.Success = false
		res.Message = "登录失败，星云平台用户信息异常"
		return res, nil
	}
	var user = models.SysUser{}
	err = db.Where("source = ? and email = ?", models.UserSourceStarPortal, email).First(&user).Error
	if err != nil {
		// 数据库中不存在该用户
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Success = false
			res.Message = "用户不存在，请联系管理员"
			log.Warn("用户不存在，请联系管理员")
			return res, nil
		} else {
			log.Errorf("Login failed: %+v", errors.WithStack(err))
			return nil, errors.Errorf("login failed：%s", err)
		}
	}

	if user.Status == models.UserStatusDisable {
		res.Success = false
		res.Message = "用户已被禁用，请联系管理员"
		log.Warnf("用户 %v 已被禁用，请联系管理员", user.NickName)
		return res, nil
	}
	// 生成token
	tokenExpireAt := time.Now().Unix() + int64(config.SystemConfig.System.TokenExpireSec)
	token, err := tools.CreateToken(strconv.Itoa(int(user.Id)), user.Username, config.SystemConfig.System.TokenExpireSec)
	if err != nil {
		log.Errorf("Create token failed: %+v", errors.WithStack(err))
		return nil, errors.Errorf("Token创建失败：%s", err)
	}
	log.Infof("user id = %v star portal login success", user.Id)

	// 将用户的token设置到context中
	s.Ctx.Set(consts.JwtTokenKey, token)
	return &dto.SysBaseUserLoginRes{
		Success:  true,
		Token:    token,
		Message:  "登录成功",
		ExpireAt: tokenExpireAt,
	}, nil
}

func GetMenus(user *models.SysUser) (menus []*models.SysMenu) {
	var hasAppendMenu = make(map[int64]struct{})
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			// 当这个菜单没有添加进menus时需要添加进去，否则无需添加
			if _, exist := hasAppendMenu[menu.Id]; !exist {
				menus = append(menus, menu)
				// 记录menus中已经添加过这个菜单
				hasAppendMenu[menu.Id] = struct{}{}
			}
		}
	}
	return menus
}

func BuildMenu(dbMenus []*models.SysMenu) []*dto.Menu {
	pbMenu := ConvertDbMenu2DtoMenu(dbMenus)
	return FindChildrenMenu(pbMenu, 0)
}

// FindChildrenMenu 查找菜单menuId的子菜单
func FindChildrenMenu(menus []*dto.Menu, menuId int64) []*dto.Menu {
	var result []*dto.Menu
	for _, menu := range menus {
		if menu.ParentMenu == menuId {
			childrenMenu := FindChildrenMenu(menus, menu.Id)
			menu.Children = childrenMenu
			result = append(result, menu)
		}
	}
	return result
}

// ConvertDbMenu2DtoMenu
//
//	@Description: 将数据模型的menu转换成dto中的menu
//	@receiver l
//	@param dbMenus
//	@return []*dto.Menu
func ConvertDbMenu2DtoMenu(dbMenus []*models.SysMenu) []*dto.Menu {
	var result []*dto.Menu
	for _, dbMenu := range dbMenus {
		meta := &dto.MenuMeta{
			Title:               dbMenu.Title,
			IgnoreKeepAlive:     dbMenu.IgnoreKeepAlive,
			Icon:                dbMenu.Icon,
			HideChildrenInMenu:  dbMenu.HideChildrenInMenu,
			HideMenu:            dbMenu.HideMenu,
			OrderNo:             dbMenu.OrderNo,
			IgnoreRoute:         dbMenu.IgnoreRoute,
			HidePathForChildren: dbMenu.HidePathForChildren,
		}
		pbMenu := &dto.Menu{
			Id:         dbMenu.Id,
			Path:       dbMenu.Path,
			Name:       dbMenu.Name,
			Component:  dbMenu.Component,
			Redirect:   dbMenu.Redirect,
			ParentMenu: int64(dbMenu.ParentMenu),
			Permission: dbMenu.Permission,
			Type:       dbMenu.Type,
			Meta:       meta,
			MenuTitle:  meta.Title,
			Status:     int64(dbMenu.Status),
		}
		// 增加菜单关联的API
		apiIds := make([]int64, 0, len(dbMenu.Apis))
		for _, api := range dbMenu.Apis {
			apiIds = append(apiIds, api.Id)
		}
		pbMenu.ApiIds = apiIds
		result = append(result, pbMenu)
	}
	return result
}

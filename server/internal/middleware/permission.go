package middleware

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/casbin"
	"biz-auto-api/pkg/consts"
	pkgdb "biz-auto-api/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type PermissionCheckReq struct {
	UserId int64  `json:"userId" validate:"required,gt=0"`
	Url    string `json:"url" validate:"required"`
	Method string `json:"method" validate:"required"`
}
type PermissionCheckRes struct {
	Success bool `json:"success"`
}

func PermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		// 判断当前请求路径是否为不需要认证的URL，如果是的话就不需要后面的权限校验，直接跳过
		if value, exists := c.Get(consts.UrlIsNoAuthKey); exists && value.(string) == consts.UrlIsNoAuthValue {
			c.Next()
			return
		}
		userId := api.GetUserIdFromContext(c)
		success, err := check(&PermissionCheckReq{UserId: userId, Url: c.Request.URL.Path, Method: c.Request.Method}, c)
		if err != nil {
			log.Errorf("权限校验异常: %v", err)
			api.Error(c, http.StatusInternalServerError, "权限校验异常："+err.Error())
			return
		}
		if !success {
			log.WithField("user-id", userId).WithField("path", c.Request.URL.Path).WithField("method", c.Request.Method).Warnf("权限拒绝")
			api.Error(c, http.StatusForbidden, "权限拒绝")
			return
		}
		c.Next()
	}
}

func check(req *PermissionCheckReq, c *gin.Context) (bool, error) {
	db := pkgdb.GetDB().WithContext(c)
	log := api.GetRequestLogger(c)

	// 校验是否为管理员
	isAdmin, err := casbin.GetEnforcer().HasGroupingPolicy(strconv.Itoa(int(req.UserId)), consts.AdminRoleIdentify)
	if err != nil {
		return false, errors.Wrapf(err, "check user role type failed")
	}
	// 管理员不需要校验权限
	if isAdmin {
		return true, nil
	}
	// 非管理员使用casbin校验
	hasPermission, err := casbin.GetEnforcer().Enforce(strconv.Itoa(int(req.UserId)), req.Url, req.Method)
	if err != nil {
		return false, errors.WithStack(err)
	}
	// 失败时判断用户状态，不存在？被禁用？
	if !hasPermission {
		var user = models.SysUser{}
		err = db.Model(&models.SysUser{}).Where(" id = ?", req.UserId).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, errors.New("用户不存在")
			}
			log.Errorf("%v", errors.WithStack(err))
		} else {
			if user.Status != models.UserStatusEnable {
				return false, errors.New("用户被禁用")
			}
		}
	}
	return hasPermission, nil
}

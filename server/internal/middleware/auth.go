package middleware

import (
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Url struct {
	UrlRegx string
	Method  string
	Desc    string
}

type ParseTokenReq struct {
	Token string `json:"token" validate:"required"`
}
type ParseTokenRes struct {
	Message  string `json:"message"`
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	Success  bool   `json:"success"`
}

var (
	// NoAuthUrl 无需认证的接口
	NoAuthUrl = []Url{
		{
			UrlRegx: consts.ApiV1Prefix + "/system/login$",
			Method:  http.MethodPost,
			Desc:    "登入接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/system/star-portal/login-url$",
			Method:  http.MethodGet,
			Desc:    "获取星云平台登录地址",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/system/star-portal/login$",
			Method:  http.MethodPost,
			Desc:    "星云平台登录",
		},
		//{
		//	UrlRegx:    consts.ApiV1Prefix + "/system/logout",
		//	Method: http.MethodPost,
		//	Desc:   "登出接口",
		//},
		{
			UrlRegx: consts.ApiV1Prefix + "/price/ping",
			Method:  http.MethodGet,
			Desc:    "连通性测试接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/billing/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/business/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/cronjob/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/network/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/price/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/system/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/ops/healthz",
			Method:  http.MethodGet,
			Desc:    "健康状态检查接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/proxy/ecdn/server/dispatchParams/report",
			Method:  http.MethodPost,
			Desc:    "ECDN-修改服务器带宽限速",
		},
	}
)

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		// 请求接口无需认证
		url := strings.TrimSuffix(c.Request.URL.Path, "/")

		if MatchOneOfUrl(NoAuthUrl, url, c.Request.Method) {
			c.Set(consts.UrlIsNoAuthKey, consts.UrlIsNoAuthValue)
			c.Next()
			return
		}
		token := c.Request.Header.Get(consts.TokenKeyInHeader)
		// 未找到token
		if token == "" {
			api.Error(c, http.StatusUnauthorized, "身份未认证，请先登录")
			return
		}

		checkTokenRes, err := parseToken(strings.TrimPrefix(token, "Bearer "), c)
		if err != nil {
			log.Errorf("%+v", err)
			api.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		// 解析不成功
		if !checkTokenRes.Success {
			log.Warnf(checkTokenRes.Message)
			api.Error(c, http.StatusUnauthorized, checkTokenRes.Message)
			return
		}

		c.Set(consts.UserIdKey, checkTokenRes.UserId)
		c.Set(consts.UsernameKey, checkTokenRes.Username)
		c.Set(consts.JwtTokenKey, token)
		c.Next()
	}
}
func parseToken(token string, c *gin.Context) (*ParseTokenRes, error) {
	log := api.GetRequestLogger(c)
	claims, err := tools.ParseToken(token)
	if err != nil {
		// token 过期
		if strings.Contains(err.Error(), "token has invalid claims") {
			log.Infof("登录已过期，请重新登陆")
			return &ParseTokenRes{Success: false, Message: "登录已过期，请重新登陆"}, nil
		}
		err = errors.WithMessage(err, "token解析失败")
		return nil, err
	}

	// TODO 认证不需要查询数据库中用户的信息，因为后续在做鉴权时会鉴权失败，再去查询失败的原因，无需每次请求都查询数据库
	//db := pkgdb.GetDB().WithContext(c)
	//var user *models.SysUser
	//err = db.Model(&models.SysUser{}).Where("id=?", claims.ID).First(&user).Error
	//if err != nil {
	//	// 用户不存在
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		log.Warnf("用户 id = %v 不存在", claims.ID)
	//		return &ParseTokenRes{Success: false, Message: "用户不存在"}, nil
	//	}
	//	// 查询失败
	//	err = errors.Wrapf(err, "Query user id=%v info failed", claims.ID)
	//	return nil, err
	//}
	//// 用户被禁用
	//if user.Status == models.UserStatusDisable {
	//	log.Warnf("用户id=%v已被禁用，请联系管理员", claims.ID)
	//	return &ParseTokenRes{Success: false, Message: "用户已被禁用，请联系管理员"}, nil
	//}
	userId, err := strconv.Atoi(claims.ID)
	if err != nil {
		err = errors.Wrapf(err, "Convert userId=%v failed", userId)
		return nil, err
	}
	return &ParseTokenRes{
		Success:  true,
		UserId:   int64(userId),
		Username: claims.Issuer,
	}, nil
}

func MatchOneOfUrl(urls []Url, url string, method string) bool {
	for _, u := range urls {
		if regexp.MustCompile(u.UrlRegx).Match([]byte(url)) && u.Method == method {
			return true
		}
	}
	return false
}

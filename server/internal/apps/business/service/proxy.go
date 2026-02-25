package service

import (
	"biz-auto-api/pkg/config"
	"biz-auto-api/pkg/service"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ApiForward struct {
	service.Service
}

func (s ApiForward) Forward(c *gin.Context) error {
	flaskApiConf := config.BusinessConfig.FlaskApi
	var target = flaskApiConf.Url
	proxyUrl, err := url.Parse(target)
	if err != nil {
		return errors.WithStack(err)
	}
	c.Request.URL.Path = c.Param("targetUrl")
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	originalDirector := proxy.Director         // 先将原本的处理函数缓存
	proxy.Director = func(req *http.Request) { // 重新赋值新的处理函数
		originalDirector(req)                                            // 执行原本的处理函数
		req.Header.Set("Authorization", "Bearer "+flaskApiConf.ApiToken) // 修改鉴权token
	}
	proxy.ServeHTTP(c.Writer, c.Request)
	return nil
}

// ForwardToPath 转发请求到指定路径
func (s ApiForward) ForwardToPath(c *gin.Context, targetPath string) error {
	flaskApiConf := config.BusinessConfig.FlaskApi
	var target = flaskApiConf.Url
	proxyUrl, err := url.Parse(target)
	if err != nil {
		return errors.WithStack(err)
	}
	// 保存原始查询参数
	c.Request.URL.Path = targetPath
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	originalDirector := proxy.Director         // 先将原本的处理函数缓存
	proxy.Director = func(req *http.Request) { // 重新赋值新的处理函数
		originalDirector(req)                                            // 执行原本的处理函数
		req.Header.Set("Authorization", "Bearer "+flaskApiConf.ApiToken) // 修改鉴权token
	}
	proxy.ServeHTTP(c.Writer, c.Request)
	return nil
}

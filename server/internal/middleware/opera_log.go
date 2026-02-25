package middleware

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/api"
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/db"
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	// NoLogUrl 不需要记录日志的url
	NoLogUrl = []Url{
		{
			UrlRegx: consts.ApiV1Prefix + "/business/k/hdd/capacity/list$",
			Method:  http.MethodPost,
			Desc:    "K业务hdd容量信息查询",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/ops/script-configs/tree$",
			Method:  http.MethodPost,
			Desc:    "OPS-查询操作菜单",
		},
	}
	// NoLogResponseUrl 不需要记录响应体的url
	NoLogResponseUrl = []Url{
		//{
		//	UrlRegx: consts.ApiV1Prefix + "/business/k/hdd/capacity/export$",
		//	Method:  http.MethodPost,
		//	Desc:    "导出hdd设备信息",
		//},
		//{
		//	UrlRegx: consts.ApiV1Prefix + "/business/b/traffic/export$",
		//	Method:  http.MethodPost,
		//	Desc:    "导出B业务流量数据",
		//},
	}
	// LogOnErrorUrl 只在错误时记录的url
	LogOnErrorUrl = []Url{
		{
			UrlRegx: consts.ApiV1Prefix + "/billing/b/traffic/report$",
			Method:  http.MethodPost,
			Desc:    "BX业务流量数据推送接口",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/flask-api/proxy/tencent/mac/upload$",
			Method:  http.MethodPost,
			Desc:    "腾讯MAC数据上传",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/flask-api/proxy/tencent/delivery/upload$",
			Method:  http.MethodPost,
			Desc:    "推送腾讯交付数据",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/business/k/hdd/capacity$",
			Method:  http.MethodPost,
			Desc:    "腾讯HDD数据上传",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/billing/k/traffic/report$",
			Method:  http.MethodPost,
			Desc:    "腾讯流量数据上报",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/billing/la/traffic/report$",
			Method:  http.MethodPost,
			Desc:    "LA流量数据上报",
		},
		{
			UrlRegx: consts.ApiV1Prefix + "/business/inspect/result/report$",
			Method:  http.MethodPost,
			Desc:    "巡检数据上报",
		},
	}
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	// 限制单次写入的最大长度和总长度
	if w.body != nil && w.body.Len() < 1024*4 && len(s) < 1024*4 {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func AddOperaLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 过滤不需要记录日志的url
		url := strings.TrimSuffix(c.Request.URL.Path, "/")
		if MatchOneOfUrl(NoLogUrl, url, c.Request.Method) {
			c.Next()
			return
		}

		reqBody := ""
		// 自定义响应体，用于获取响应体内容
		var crw *CustomResponseWriter
		// 如果不需要记录返回响应体，则使用nil
		if MatchOneOfUrl(NoLogResponseUrl, url, c.Request.Method) {
			crw = &CustomResponseWriter{body: nil, ResponseWriter: c.Writer}
		} else {
			crw = &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		}

		log := api.GetRequestLogger(c)
		startTime := time.Now()
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			bf := bytes.NewBuffer(nil)
			wt := bufio.NewWriter(bf)
			_, err := io.Copy(wt, c.Request.Body)
			if err != nil {
				log.Warnf("copy body error, %s", err.Error())
				err = nil
			}
			rb, _ := io.ReadAll(bf)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rb))
			reqBody = string(rb)
			c.Writer = crw
		}

		c.Next()
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			requestId := ""
			if r, exist := c.Get(consts.RequestIdKey); exist {
				requestId, _ = r.(string)
			}
			var userId int64 = 0
			if r, exist := c.Get(consts.UserIdKey); exist {
				userId, _ = r.(int64)
			}

			// 如果没有获取从ctx中直接获取到用户ID，说明可能是登录接口，或者不需要认证的接口，
			// 这边再次尝试从ctx中拿一下token，解析出里面的用户，因为用户登录完成后会重新设置token到ctx中
			if userId == 0 {
				if token, ok := c.Get(consts.JwtTokenKey); ok {
					res, _ := parseToken(token.(string), c)
					if res != nil {
						userId = res.UserId
					}
				}
			}
			elapsed := time.Since(startTime)
			latencyTime := float64(elapsed.Nanoseconds()) / 1e6
			clientIp := c.ClientIP()
			uri := strings.ReplaceAll(c.Request.RequestURI, "%", "%%")
			requestMethod := c.Request.Method
			userAgent := c.Request.UserAgent()
			httpCode := strconv.Itoa(c.Writer.Status())

			handleSource := models.SysLogHandleSourceProgram
			if isBrowserUserAgent(userAgent) {
				handleSource = models.SysLogHandleSourceUser
			}
			// body不为空，则获取响应体内容
			jsonRes := "{}"
			if crw.body != nil {
				if crw.body.Len() <= 1024*8 {
					jsonRes = crw.body.String()
				} else {
					jsonRes = `{"msg":"响应体长度超出限制"}`
				}
			}
			// 获取业务状态码
			bizCode := gjson.Get(jsonRes, "code").String()
			// 获取handler
			handler := ""
			if httpCode == "404" {
				handler = "NoHandler"
			} else {
				s := strings.Split(c.HandlerName(), ".")
				if len(s) >= 3 {
					handler = s[1] + "." + s[2]
					handler = handler[0 : len(handler)-3]
				}
			}
			// LogOnErrorUrl 中的请求只记录非200的情况
			if MatchOneOfUrl(LogOnErrorUrl, url, requestMethod) && bizCode == "200" && httpCode == "200" {
				return
			}
			operaLog := &models.SysOperaLog{
				Uri:           uri,
				Api:           c.FullPath(),
				RequestId:     requestId,
				RequestMethod: requestMethod,
				ClientIp:      clientIp,
				LatencyTime:   fmt.Sprintf("%v", latencyTime),
				UserAgent:     userAgent,
				ReqBody:       reqBody,
				JsonRes:       jsonRes,
				HttpCode:      httpCode,
				BizCode:       bizCode,
				CreateBy:      userId,
				UpdateBy:      userId,
				RequestTime:   startTime,
				Handler:       handler,
				HandleSource:  int64(handleSource),
			}
			if int(handleSource) == models.SysLogHandleSourceUser {
				err := db.GetDB().Session(
					&gorm.Session{
						Logger: logger.Default.LogMode(logger.Silent),
					},
				).Create(&operaLog).Error
				if err != nil {
					log.Warnf("add opera log error, %s", err.Error())
				}
			}
		}
	}
}

func isBrowserUserAgent(userAgent string) bool {
	if userAgent == "" {
		return false
	}
	// 常见浏览器关键词（不区分大小写）
	browserKeywords := []string{
		"Mozilla", "Chrome", "Safari", "Firefox", "Edge", "Opera",
		"WebKit", "Gecko", "Mobile", "Android", "iPhone",
	}
	userAgent = strings.ToLower(userAgent)
	for _, keyword := range browserKeywords {
		if strings.Contains(userAgent, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

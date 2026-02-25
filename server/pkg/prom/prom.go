package prom

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Prometheus struct {
	PromUrl string
	Auth    string
}

func NewProm(promUrl, auth string) *Prometheus {
	api, _ := url.JoinPath(promUrl, "/api/v1/query")
	return &Prometheus{
		api,
		auth,
	}
}

func (prom *Prometheus) Fetch(sql string, timeout time.Duration) ([]byte, error) {
	// 创建一个带有超时的http.Client
	client := &http.Client{
		Timeout: timeout,
	}

	// 构建查询URL
	queryUrl := fmt.Sprintf("%s?query=%s", prom.PromUrl, url.QueryEscape(sql))

	// 创建新的请求
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加Basic认证头
	if prom.Auth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(prom.Auth))
		req.Header.Add("Authorization", basicAuth)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Prometheus返回错误状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	return body, nil
}

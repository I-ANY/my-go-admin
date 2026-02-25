package le_register

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	http.Client
	url         string
	username    string
	password    string
	tokenExpiry time.Time
	token       string
	log         *logrus.Entry
}

var LeSnSomeRegisterFailed = errors.New("sn some register failed")

func NewClient(url, username, password string, log *logrus.Entry) *Client {
	url = strings.TrimSuffix(url, "/")
	return &Client{
		url:      url,
		username: username,
		password: password,
		log:      log,
	}
}

func (c *Client) Auth() error {
	req, err := http.NewRequest(http.MethodGet, c.url+"/api/app_token", nil)
	if err != nil {
		return fmt.Errorf("build auth request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Token     string `json:"token"`
			ExpiredAt string `json:"expired_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}
	if out.Code != 200 || out.Data.Token == "" {
		return fmt.Errorf("auth failed: code=%d msg=%s", out.Code, out.Message)
	}

	c.token = out.Data.Token

	if out.Data.ExpiredAt != "" {
		if t, err := time.Parse(time.RFC3339Nano, out.Data.ExpiredAt); err == nil {
			c.tokenExpiry = t
		}
	}

	return nil
}

func (c *Client) Register(sns []string) error {
	// 如果 token 为空或已过期，尝试刷新
	if c.token == "" || (!c.tokenExpiry.IsZero() && time.Now().After(c.tokenExpiry)) {
		if err := c.Auth(); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
	}

	// 定义并直接赋值
	requestBody := struct {
		Token string   `json:"token"`
		Sns   []string `json:"sns"`
	}{
		Token: c.token,
		Sns:   sns,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"/api/app_authorization", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("build auth request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	// {"code":200,"message":"成功","data":{"requested_count":4,"authorized_count":4}}
	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			RequestedCount  int `json:"requested_count"`
			AuthorizedCount int `json:"authorized_count"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}

	c.log.Infof("code: %d, message: %s,requested count: %d, authorized count: %d", out.Code, out.Message, out.Data.RequestedCount, out.Data.AuthorizedCount)

	if out.Code != 200 {
		return fmt.Errorf("auth failed: code=%d msg=%s", out.Code, out.Message)
	}

	if out.Data.RequestedCount != out.Data.AuthorizedCount || len(sns) != out.Data.AuthorizedCount || len(sns) != out.Data.RequestedCount {
		return fmt.Errorf("auth failed: requested=%d authorized=%d: %w", out.Data.RequestedCount, out.Data.AuthorizedCount, LeSnSomeRegisterFailed)
	}

	return nil

}

func (c *Client) Query(sns []string) (map[string]bool, string, error) {
	// 如果 token 为空或已过期，尝试刷新
	if c.token == "" || (!c.tokenExpiry.IsZero() && time.Now().After(c.tokenExpiry)) {
		if err := c.Auth(); err != nil {
			return nil, "", fmt.Errorf("refresh token: %w", err)
		}
	}

	// 构造请求体，只需要 token
	body := struct {
		Token string `json:"token"`
	}{Token: c.token}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"/api/app_authorized", bytes.NewBuffer(data))
	if err != nil {
		return nil, "", fmt.Errorf("build query request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("query request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("query failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			DownloadURL string `json:"download_url"`
			ExpiredAt   string `json:"expired_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, "", fmt.Errorf("decode query response: %w", err)
	}
	if out.Code != 200 || out.Data.DownloadURL == "" {
		return nil, "", fmt.Errorf("query failed: code=%d msg=%s", out.Code, out.Message)
	}

	// 下载 download_url 指向的清单
	listResp, err := c.Get(out.Data.DownloadURL)
	if err != nil {
		return nil, out.Data.DownloadURL, fmt.Errorf("download list failed: %w", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(listResp.Body)
		return nil, out.Data.DownloadURL, fmt.Errorf("download list failed: status=%d body=%s", listResp.StatusCode, string(b))
	}

	// 目标 SNS 集合
	want := make(map[string]struct{}, len(sns))
	for _, sn := range sns {
		sn = strings.TrimSpace(sn)
		if sn != "" {
			want[sn] = struct{}{}
		}
	}

	// 结果：sn -> 是否存在
	result := make(map[string]bool, len(want))
	for sn := range want {
		result[sn] = false
	}

	// 只关心“当前时间 +24 小时”的日期（忽略时分秒），格式：YYYYMMDD
	targetDate := time.Now().In(time.Local).Add(24 * time.Hour).Format("20060102")

	scanner := bufio.NewScanner(listResp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		// 每行形如：<sn> <tab/space> <yyyymmddHHMMSS>
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		sn := fields[0]
		if _, ok := want[sn]; ok {
			// 仅当该 SN 的日期部分等于 targetDate 时才视为存在
			if !result[sn] && len(fields) >= 2 {
				ts := fields[1]
				if len(ts) >= 8 && ts[:8] == targetDate {
					result[sn] = true
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, out.Data.DownloadURL, fmt.Errorf("scan list: %w", err)
	}

	return result, out.Data.DownloadURL, nil
}

func (c *Client) DownloadUrl() (string, error) {

	// 如果 token 为空或已过期，尝试刷新
	if c.token == "" || (!c.tokenExpiry.IsZero() && time.Now().After(c.tokenExpiry)) {
		if err := c.Auth(); err != nil {
			return "", fmt.Errorf("refresh token: %w", err)
		}
	}

	// 构造请求体，只需要 token
	body := struct {
		Token string `json:"token"`
	}{Token: c.token}

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.url+"/api/app_authorized", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("build query request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("query request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("query failed: status=%d body=%s", resp.StatusCode, string(b))
	}

	var out struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			DownloadURL string `json:"download_url"`
			ExpiredAt   string `json:"expired_at"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode query response: %w", err)
	}
	if out.Code != 200 || out.Data.DownloadURL == "" {
		return "", fmt.Errorf("query failed: code=%d msg=%s", out.Code, out.Message)
	}
	return out.Data.DownloadURL, nil
}

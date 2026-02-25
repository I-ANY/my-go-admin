package overprovisioning

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Prometheus struct {
	PromURL  string
	PromAuth string
}

type OwnerResp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Owner string `json:"owner"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type HostsResp struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				ID              string `json:"id"`
				Business        string `json:"business"`
				BwPlan          string `json:"bwPlan"`
				Hostname        string `json:"hostname"`
				Sn              string `json:"sn"`
				Isp             string `json:"isp"`
				Location        string `json:"location"`
				Owner           string `json:"owner"`
				Status          string `json:"status"`
				BwSingle        string `json:"bwsingle"`
				BwCount         string `json:"bwcount"`
				Origin          string `json:"origin"`
				CactiNotes      string `json:"cactiNotes"`
				Day95           string `json:"day95"`
				Evening95       string `json:"evening95"`
				Kvm             string `json:"kvm"`             //是否KVM设备（1-是，2-否, 3-否）
				Parent          string `json:"parent"`          //宿主机名
				Interprovincial string `json:"interprovincial"` // 是否跨省
				IP              string `json:"ip"`              // IP地址
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func NewProm(promURL, promAuth string) *Prometheus {
	return &Prometheus{
		PromURL:  promURL,
		PromAuth: promAuth,
	}
}

func (prom *Prometheus) Fetch(sql string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	queryUrl := fmt.Sprintf("%sapi/v1/query?query=%s", prom.PromURL, url.QueryEscape(sql))

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加Basic认证头
	if prom.PromAuth != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(prom.PromAuth))
		req.Header.Add("Authorization", basicAuth)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	return body, nil
}

// GetHostsByBiz query hosts by biz
func (prom *Prometheus) GetHostsByBiz(biz string) (*HostsResp, error) {
	hostSQL := fmt.Sprintf(`mfy_hosts_ecdn_info{origin="自建",business="%s"}`, biz)
	bytes, err := prom.Fetch(hostSQL, time.Second*30)
	if err != nil {
		return nil, err
	}
	hosts := HostsResp{}
	_ = json.Unmarshal(bytes, &hosts)
	return &hosts, nil
}

// GetOwners query owner
func (prom *Prometheus) GetOwners() (*OwnerResp, error) {
	sql := `count(mfy_hosts_ecdn_info{origin="自建"}) by (owner)`
	data, err := prom.Fetch(sql, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("fetch owners from prometheus failed: %w", err)
	}
	var owners OwnerResp
	if err := json.Unmarshal(data, &owners); err != nil {
		return nil, fmt.Errorf("unmarshal owner response failed: %w", err)
	}
	return &owners, nil
}

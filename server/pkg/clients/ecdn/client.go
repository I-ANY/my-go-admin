package ecdn

import (
	"biz-auto-api/pkg/clients"
	"biz-auto-api/pkg/tools"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type EcdnClient struct {
	http.Client
	url string
	key string
}

func NewEcdnClient(url, key string) (*EcdnClient, error) {
	if len(url) == 0 {
		return nil, errors.Errorf("url can not be empty")
	}
	return &EcdnClient{
		url: strings.TrimSuffix(url, "/"),
		key: key,
	}, nil
}
func (c *EcdnClient) GetServerByBusiness(business string) ([]*Server, error) {
	return c.GetServers([]string{business}, "", "", "")
}
func (c *EcdnClient) GetServerByBusinesses(businesses []string) ([]*Server, error) {
	return c.GetServers(businesses, "", "", "")
}

func (c *EcdnClient) GetServers(businesses []string, hostname, sn, frankID string) ([]*Server, error) {
	api := "/p2pcdn-go-admin/api/server/serversPub"
	formParams := url.Values{}

	if len(businesses) > 0 {
		for _, business := range businesses {
			formParams.Add("businesses[]", business)
		}
	}
	if len(hostname) > 0 {
		formParams.Add("hostname", hostname)
	}
	if len(sn) > 0 {
		formParams.Add("sn", sn)
	}
	if len(frankID) > 0 {
		formParams.Add("frankID", frankID)
	}
	res, err := c.HttpRequest(http.MethodGet, api, nil, nil, formParams, nil)
	if err != nil {
		return nil, errors.Wrap(err, "query ecdn servers failed")
	}
	if !gjson.Valid(string(res)) {
		return nil, errors.Errorf("response not json: %s", string(res))
	}
	if gjson.Get(string(res), "code").Int() != 0 {
		return nil, errors.Errorf("get ecdn servers failed: %s", string(res))
	}
	var servers = make([]*Server, 0)
	serverData := gjson.Get(string(res), "data.data")
	if !serverData.IsArray() {
		return nil, nil
	}
	err = json.Unmarshal([]byte(serverData.String()), &servers)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal ecdn servers failed")
	}
	return servers, nil
}

func (c *EcdnClient) HttpRequest(httpMethod, path string, header http.Header, query url.Values, formParams url.Values, body []byte) ([]byte, error) {
	if header == nil {
		header = make(map[string][]string)
	}
	if _, ok := header["X-App-Name"]; !ok {
		hostname, _ := os.Hostname()
		header.Set("X-App-Name", hostname)
	}
	if formParams == nil {
		formParams = url.Values{}
	}
	now := time.Now().Unix()
	if len(formParams.Get("timestamp")) == 0 {
		formParams.Set("timestamp", fmt.Sprintf("%v", now))
	}
	frankID := formParams.Get("frankID")
	if len(c.key) > 0 && len(header.Get("Authorization")) == 0 {
		token := tools.MD5Str(fmt.Sprintf("%v%v%v", frankID, c.key, now))
		header.Set("Authorization", "Bearer "+token)
	}
	return clients.HttpRequest(c.Client, c.url, path, httpMethod, header, query, formParams, body)
}

// DifIsp 异网下发
func (c *EcdnClient) DifIsp(req *DifIspReq) (*DifIspData, error) {
	api := "/p2pcdn-go-admin/api/autoSystem/difIspTask"
	ts := time.Now().Unix()
	query := url.Values{}
	query["frankIDS"] = []string{req.FrankIDS}
	query["carrier"] = []string{fmt.Sprintf("%v", req.Carrier)}
	query["provincial"] = []string{req.Provincial}
	query["phone"] = []string{req.Phone}
	query["remind"] = []string{fmt.Sprintf("%v", req.Remind)}
	query["note"] = []string{req.Note}
	query["timestamp"] = []string{fmt.Sprintf("%v", ts)}

	header := make(http.Header)
	token := tools.MD5Str(fmt.Sprintf("%v%v", c.key, ts))
	header.Set("Authorization", "Bearer "+token)
	c.Client.Timeout = time.Second * 30
	res, err := clients.HttpRequest(c.Client, c.url, api, http.MethodPost, header, query, nil, nil)
	if err != nil {
		return nil, err
	}
	var resData = &difIspRes{}
	err = json.Unmarshal(res, resData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resData.Code != 0 {
		return nil, errors.Errorf("request failed: status_code=%v res=%s", resData.Code, resData.Message)
	}
	return resData.Data, nil
}

// GetServerIndicators 查询指定主机在某小时内的流量指标
// hostname: 主机名
// queryTime: 查询时间，格式 YYYYMMDDHH
func (c *EcdnClient) GetServerIndicators(hostname, queryTime string) (*ServerIndicatorsData, error) {
	api := "/p2pcdn-go-admin/api/server/serversPub/indicators"
	formParams := url.Values{}
	if len(hostname) == 0 || len(queryTime) == 0 {
		return nil, errors.New("hostname and queryTime can not be empty")
	}
	formParams.Set("hostname", hostname)
	formParams.Set("queryTime", queryTime)

	res, err := c.HttpRequest(http.MethodGet, api, nil, nil, formParams, nil)
	if err != nil {
		return nil, errors.Wrap(err, "query ecdn indicators failed")
	}
	if !gjson.ValidBytes(res) {
		return nil, errors.Errorf("response not json: %s", string(res))
	}

	var resp serverIndicatorsResp
	if err := json.Unmarshal(res, &resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal ecdn indicators failed")
	}
	if resp.Code != 0 {
		if resp.Code == 404 {
			return nil, nil
		} else {
			return nil, nil
		}
	}
	return resp.Data, nil
}

// GetNodeUtilization 查询指定节点机房某天的利用率（新统计接口，带客户端鉴权）
func (c *EcdnClient) GetNodeUtilization(req *NodeUtilizationQueryReq) (*NodeUtilizationQueryResp, error) {
	if req == nil {
		return nil, errors.New("request can not be nil")
	}
	api := "/p2pcdn-go-admin/api/statistic/node/utilization/query"
	formParams := url.Values{}
	if req.Date != "" {
		formParams.Set("date", req.Date)
	} else {
		return nil, errors.New("date can not be empty")
	}
	if req.Owner != "" {
		formParams.Set("owner", req.Owner)
	} else {
		return nil, errors.New("owner can not be empty")
	}
	if req.Isp != "" {
		formParams.Set("isp", req.Isp)
	} else {
		return nil, errors.New("isp can not be empty")
	}
	if req.Location != "" {
		formParams.Set("location", req.Location)
	} else {
		return nil, errors.New("location can not be empty")
	}

	res, err := c.HttpRequest(http.MethodGet, api, nil, nil, formParams, nil)
	if err != nil {
		return nil, errors.Wrap(err, "query ecdn node utilization failed")
	}
	if !gjson.ValidBytes(res) {
		return nil, errors.Errorf("response not json: %s", string(res))
	}
	var resp NodeUtilizationQueryResp
	if err := json.Unmarshal(res, &resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal ecdn node utilization failed")
	}

	if resp.Code != 0 {
		if resp.Code == 404 {
			return nil, nil
		} else if resp.Code == 429 {
			// 等待2秒后重试
			time.Sleep(time.Second * 2)
			fmt.Printf("get node utilization rate limit %s|%s|%s, waiting 2 seconds and retry\n", req.Owner, req.Isp, req.Location)
			return c.GetNodeUtilization(req)
		} else {
			return nil, errors.Errorf("request failed: status_code=%v res=%s", resp.Code, resp.Message)
		}
	}
	return &resp, nil
}

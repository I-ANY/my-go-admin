package la_traffic

import (
	"biz-auto-api/pkg/clients"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	http.Client
	url      string
	username string
	password string
}

func NewClient(url, username, password string) *Client {
	url = strings.TrimSuffix(url, "/")
	return &Client{
		url:      url,
		username: username,
		password: password,
	}
}

// GetDeviceTraffic 时间范围为左闭右开
func (c *Client) GetDeviceTraffic(deviceId string, startTime, endTime int64) ([]*DeviceTrafficItem, error) {
	path := "/api/flow"
	query := url.Values{}
	query.Add("deviceid", deviceId)
	query.Add("starttime", strconv.Itoa(int(startTime)))
	query.Add("endtime", strconv.Itoa(int(endTime)))
	res, err := c.HttpRequest(http.MethodGet, path, nil, query, nil, nil)
	if err != nil {
		return nil, err
	}
	trafficData := CommonResponse[[]*DeviceTrafficItem]{}
	err = json.Unmarshal(res, &trafficData)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal data failed")
	}
	if trafficData.Code != 0 {
		return nil, errors.Errorf("get la traffic data failed: code: %v, message: %v", trafficData.Code, trafficData.Msg)
	}
	return trafficData.Data, nil
}

func (c *Client) GetDeviceTrafficAtMoment(deviceId string, moment int64) (*DeviceTrafficItem, error) {
	// 左闭右开
	traffics, err := c.GetDeviceTraffic(deviceId, moment, moment+300)
	if err != nil {
		return nil, errors.WithMessagef(err, "deviceId=%v, ts=%v, time=%v", deviceId, moment, time.Unix(moment, 0).Format(time.DateTime))
	}
	for _, traffic := range traffics {
		if traffic.Timestamp == time.Unix(moment, 0).Format(time.DateTime) {
			return traffic, nil
		}
	}
	return nil, nil
}

func (c *Client) HttpRequest(httpMethod, path string, header http.Header, query url.Values, formParams url.Values, body []byte) ([]byte, error) {
	if header == nil {
		header = make(map[string][]string)
	}
	auth := c.username + ":" + c.password
	auth = base64.StdEncoding.EncodeToString([]byte(auth))
	header.Set("Authorization", "Basic "+auth)
	return clients.HttpRequest(c.Client, c.url, path, httpMethod, header, query, formParams, body)
}

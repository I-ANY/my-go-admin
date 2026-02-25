package flask_api

import (
	"biz-auto-api/pkg/clients"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"strings"
)

type FlaskApiClient struct {
	http.Client
	url      string
	apiToken string
}

func NewFlaskApiClient(url, apiToken string) (*FlaskApiClient, error) {
	if len(url) == 0 {
		return nil, errors.Errorf("url can not be empty")
	}
	return &FlaskApiClient{
		url:      strings.TrimSuffix(url, "/"),
		apiToken: apiToken,
	}, nil
}

func (c *FlaskApiClient) HttpRequest(httpMethod, path string, header http.Header, query url.Values, formParams url.Values, body []byte) ([]byte, error) {
	if header == nil {
		header = make(map[string][]string)
	}
	if len(c.apiToken) > 0 && len(header.Get("Authorization")) == 0 {
		header.Set("Authorization", "Bearer "+c.apiToken)
	}
	return clients.HttpRequest(c.Client, c.url, path, httpMethod, header, query, formParams, body)
}

func (c *FlaskApiClient) DeliveryDevices(payload *DeliveryPayload) (*DeliveryResp, error) {
	if payload == nil {
		return nil, errors.New("payload can not be nil")
	}
	bs, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal failed, payload=%v", payload)
	}
	response, err := c.HttpRequest(http.MethodPost, "/tencent/delivery/submit", nil, nil, nil, bs)
	if err != nil {
		return nil, err
	}
	if !gjson.Valid(string(response)) {
		return nil, errors.Errorf("response not json: %s", string(response))
	}
	var res = DeliveryResp{}
	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response failed")
	}
	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}
	return &res, nil
}

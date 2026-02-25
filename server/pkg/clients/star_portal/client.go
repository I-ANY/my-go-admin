package starPortal

import (
	"biz-auto-api/pkg/clients"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type StarPortalClient struct {
	http.Client
	url      string
	apiToken string
}

func NewStarPortalClient(url, apiToken string) (*StarPortalClient, error) {
	if len(url) == 0 {
		return nil, errors.Errorf("url can not be empty")
	}
	return &StarPortalClient{
		Client: http.Client{
			Timeout: time.Second * 30,
		},
		url:      strings.TrimSuffix(url, "/"),
		apiToken: apiToken,
	}, nil
}

func (c *StarPortalClient) GetUsers(force bool) ([]*User, error) {
	var headers = http.Header{}
	if force {
		headers = http.Header{
			"X-Force-Fresh": []string{strconv.Itoa(int(time.Now().UnixMilli()))},
		}
	}
	bs, err := c.HttpRequest(http.MethodGet, GetUserApi, headers, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	us := gjson.Get(string(bs), "users")
	if !us.IsArray() {
		return nil, errors.Errorf("get users failed: %s", string(bs))
	}
	var users []*User
	err = json.Unmarshal([]byte(us.String()), &users)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal users failed")
	}
	return users, nil
}

func (c *StarPortalClient) GetUserToken(code, clientId, clientSecret, redirectUri string) (*UserTokenInfo, error) {
	formData := url.Values{
		"client_id":     []string{clientId},
		"client_secret": []string{clientSecret},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"redirect_uri":  []string{redirectUri},
	}
	bs, err := c.HttpRequest(http.MethodPost, GetUserToken, nil, nil, formData, nil)
	if err != nil {
		return nil, err
	}
	var userTokenInfo *UserTokenInfo
	err = json.Unmarshal(bs, &userTokenInfo)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal user token failed")
	}
	return userTokenInfo, nil
}

func (c *StarPortalClient) GetUserInfo() (*UserInfo, error) {
	bs, err := c.HttpRequest(http.MethodGet, GetUserInfo, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var userToken *UserInfo
	err = json.Unmarshal(bs, &userToken)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal user info failed")
	}
	return userToken, nil
}

func (c *StarPortalClient) HttpRequest(httpMethod, path string, header http.Header, query url.Values, formParams url.Values, body []byte) ([]byte, error) {
	if header == nil {
		header = make(map[string][]string)
	}
	if len(c.apiToken) > 0 && len(header.Get("Authorization")) == 0 {
		header.Set("Authorization", "Bearer "+c.apiToken)
	}
	return clients.HttpRequest(c.Client, c.url, path, httpMethod, header, query, formParams, body)
}

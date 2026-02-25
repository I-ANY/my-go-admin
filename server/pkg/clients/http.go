package clients

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func HttpRequest(client http.Client,
	baseUrl, apiPath, httpMethod string,
	header http.Header,
	query url.Values,
	formParams url.Values,
	body []byte) ([]byte, error) {
	if header == nil {
		header = make(map[string][]string)
	}
	if query == nil {
		query = make(url.Values)
	}
	var buffer io.Reader
	if body != nil {
		buffer = bytes.NewBuffer(body)
	}
	if formParams != nil {
		buffer = strings.NewReader(formParams.Encode())
		header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	} else {
		header["Content-Type"] = []string{"application/json"}
	}
	// 解析路径
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	u.Path += apiPath
	u.RawQuery = query.Encode()
	request, err := http.NewRequest(httpMethod, u.String(), buffer)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	request.Header = header
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()
	resData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if response.StatusCode != 200 {
		return nil, errors.Errorf("request failed: status_code=%v res=%s", response.StatusCode, resData)
	}

	return resData, nil
}

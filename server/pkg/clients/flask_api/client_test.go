package flask_api

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

const (
	u        = "http://10.1.10.173:5000"
	apitoken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIwNTA1OTkxMzQsImlhdCI6MTczNTIzOTEzNCwiaWRlbnRpdHkiOjF9.2m7VsxxOxsE50BNwZH2-iPGgnHCAs5FCL-mIoRRMmWQ"
)

func TestFlaskApiClient(t *testing.T) {
	flaskApiClient, err := NewFlaskApiClient(u, apitoken)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := flaskApiClient.HttpRequest(http.MethodGet, "/tencent/k/device", nil, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}

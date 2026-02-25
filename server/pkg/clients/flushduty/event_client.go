package flushduty

import (
	"biz-auto-api/pkg/clients"
	"biz-auto-api/pkg/tools"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type EventClient struct {
	url            string
	integrationKey string
	client         http.Client
}

func NewEventClient(url, integrationKey string) *EventClient {
	return &EventClient{
		url:            url,
		integrationKey: integrationKey,
		client:         http.Client{},
	}
}
func (c *EventClient) SendAlert(event *EventReqPayload) (*EventResp, error) {

	bs, err := json.Marshal(event)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	query := url.Values{
		"integration_key": []string{c.integrationKey},
	}
	resp, err := clients.HttpRequest(c.client, c.url, "/event/push/alert/standard", http.MethodPost, nil, query, nil, bs)
	if err != nil {
		return nil, err
	}
	response := &EventResp{}
	err = json.Unmarshal(resp, response)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if response.Error != nil {
		return nil, errors.Errorf("recover failed, resp: %v", string(bs))
	}
	return response, nil
}
func (c *EventClient) RecoverAlert(titleRule, alertKey string) error {
	payload := &EventReqPayload{
		TitleRule:   tools.ToPointer(titleRule),
		AlertKey:    tools.ToPointer(alertKey),
		EventStatus: tools.ToPointer(EventStatus_Ok),
	}
	_, err := c.SendAlert(payload)
	if err != nil {
		return err
	}
	return nil
}

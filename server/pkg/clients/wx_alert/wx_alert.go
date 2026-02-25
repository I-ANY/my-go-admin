package wx_alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type textPayload struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content       string   `json:"content"`
		MentionedList []string `json:"mentioned_list,omitempty"`
	} `json:"text"`
}

func SendText(webhookURL string, content string, mentionAll bool) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook url is empty")
	}
	if content == "" {
		return fmt.Errorf("content is empty")
	}

	var payload textPayload
	payload.MsgType = "text"
	payload.Text.Content = content
	if mentionAll {
		payload.Text.MentionedList = []string{"@all"}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook http status=%d", resp.StatusCode)
	}

	// 企业微信机器人返回 {"errcode":0,"errmsg":"ok"}
	var out struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("decode webhook response: %w", err)
	}
	if out.ErrCode != 0 {
		return fmt.Errorf("webhook error: errcode=%d errmsg=%s", out.ErrCode, out.ErrMsg)
	}
	return nil
}

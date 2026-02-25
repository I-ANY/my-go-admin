package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Sender interface {
	SendMsg(msg string, webhookUrl string, users ...string) error
}

type Tencent struct {
	MsgType   string
	ImgBase64 string
	ImgMd5    string
}

// Tencent文件上传返回结果结构体
type UploadResp struct {
	ErrCode   int    `json:"errCode"`
	ErrMsg    string `json:"errMsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func (t *Tencent) SetType(msgType string) {
	t.MsgType = msgType
}

func (t *Tencent) SendFile(msg string, webhookUrl string, fileName string) error {
	key := strings.Split(webhookUrl, "=")[1]
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=file", key)
	fileContent := bytes.NewReader([]byte(msg))
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// 添加文件字段
	part, err := writer.CreateFormFile("media", fileName)
	if err != nil {
		return fmt.Errorf("create form file failed: %w", err)
	}
	_, err = io.Copy(part, fileContent)
	if err != nil {
		return fmt.Errorf("copy file content failed: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("close writer failed: %w", err)
	}
	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 设置客户端超时
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http post failed: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}
	var uploadResp UploadResp
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}
	if uploadResp.ErrCode != 0 {
		return fmt.Errorf("upload failed: %s", uploadResp.ErrMsg)
	}
	err = t.SendMsg(uploadResp.MediaID, webhookUrl)
	if err != nil {
		return fmt.Errorf("send msg failed: %w", err)
	}
	return nil
}

func (t *Tencent) SendMsg(msg string, webhookUrl string, users ...string) error {
	var content map[string]any
	switch t.MsgType {
	case "text":
		content = map[string]any{
			"msgtype": "text",
			"text": map[string]any{
				"content":               msg,
				"mentioned_mobile_list": users,
			},
		}
	case "markdown":
		content = map[string]any{
			"msgtype": "markdown",
			"markdown": map[string]any{
				"content": msg,
			},
		}
	case "markdown_v2":
		content = map[string]any{
			"msgtype": "markdown_v2",
			"markdown_v2": map[string]any{
				"content": msg,
			},
		}
	case "file":
		content = map[string]any{
			"msgtype": "file",
			"file": map[string]any{
				"media_id": msg,
			},
		}
	case "image":
		content = map[string]any{
			"msgtype": "image",
			"image": map[string]any{
				"base64": t.ImgBase64,
				"md5":    t.ImgMd5,
			},
		}
	}

	// 将内容转换为JSON格式
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 发送HTTP POST请求
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonContent))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

type DingTalk struct {
	MsgType string
}

func (t *DingTalk) SetType(msgType string) {
	t.MsgType = msgType
}

func (t *DingTalk) SendMsg(msg string, webhookUrl string, users ...string) error {
	var content map[string]any
	switch t.MsgType {
	case "text":
		content = map[string]any{
			"msgtype": "text",
			"text": map[string]any{
				"content": msg,
			},
			"at": map[string]any{
				"atMobiles": users,
			},
		}
	case "markdown":
		content = map[string]any{
			"msgtype": "markdown",
			"markdown": map[string]any{
				"title": "业务变更",
				"text":  msg,
			},
		}
	}

	// 将内容转换为JSON格式
	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	// 发送HTTP POST请求
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonContent))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

package operator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	FrpAddress_Prod = "frp.xaidc.com:6789"
	FrpAddress_Dev  = "frp-dev.xaidc.com:6789"
)

func ExecScript(scriptName, hostname, busiType, scriptArgsEncode string, timeout time.Duration) (*RemoteExecResponse, error) {
	// hostname 包含 “_” 则使用 FrpAddress_Dev
	if len(hostname) >= 8 && (hostname)[7:8] == "_" {
		return ExecScriptWithDomain(FrpAddress_Dev, scriptName, hostname, busiType, scriptArgsEncode, timeout)
	}
	return ExecScriptWithDomain(FrpAddress_Prod, scriptName, hostname, busiType, scriptArgsEncode, timeout)
}

func ExecScriptWithDomain(domain, scriptName, hostname, busiType, scriptArgsEncode string, timeout time.Duration) (*RemoteExecResponse, error) {
	if timeout <= 0 {
		timeout = time.Minute * 6
	}
	scriptPath := "agent"
	baseURL := fmt.Sprintf("http://%s-agent.%s/%s/%s/%s",
		strings.ToLower(hostname), domain, scriptPath, busiType, scriptName)
	// 判断是否要加参数
	var fullURL string
	if strings.TrimSpace(scriptArgsEncode) != "" {
		fullURL = fmt.Sprintf("%s?conf=%s", baseURL, scriptArgsEncode)
	} else {
		fullURL = baseURL
	}

	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(fullURL)
	if err != nil {
		// 识别超时错误，返回明确的超时信息
		if isTimeoutErr(err) {
			return nil, errors.Wrap(err, "request timeout")
		}
		return nil, errors.Errorf("remote script execution request address is abnormal: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("remote script execution response body is abnormal: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.Errorf("ops agent 连接失败, 请检查设备是否在线或ops agent是否正常运行")
	}

	var result *RemoteExecResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.Errorf("remote script execution response body is abnormal: %v", err)
	}

	if result.Code != 0 {
		return nil, errors.Errorf("code: %d, result: %v", result.Code, result.Result)
	}

	return result, nil
}

// isTimeoutErr 判断是否为超时错误
func isTimeoutErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
		return true
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline") {
		return true
	}
	return false
}

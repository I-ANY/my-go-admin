package switch_collector

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type SwitchSSHSession struct {
	session *ssh.Session
	stdin   io.WriteCloser
	out     *bytes.Buffer
	lock    sync.Mutex
}

func (s *SwitchSSH) NewSession() (*SwitchSSHSession, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // disable echoing
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	//建立伪终端
	err = session.RequestPty("xterm", 80, 100, modes)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 只重定向输出，输入通过管道控制
	var o bytes.Buffer
	session.Stdout = &o
	session.Stderr = &o
	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "get stdin pipe failed")
	}
	err = session.Shell()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	se := &SwitchSSHSession{session: session, stdin: stdin, out: &o}
	_, err = se.Exec(" ", time.Second*10, true, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return se, nil
}
func (s *SwitchSSHSession) Close() error {
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}

func (s *SwitchSSHSession) Exec(cmd string, timeout time.Duration, removeFirstLine, removeLastLine bool, prompts ...string) ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 记录执行命令前的缓冲区大小
	initialSize := s.out.Len()
	_, err := fmt.Fprintf(s.stdin, "%s\n", cmd)
	if err != nil {
		return nil, errors.Wrap(err, "send command failed")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	// 创建一个缓冲区来存储本次命令的输出
	var commandOutput bytes.Buffer
	for {
		select {
		case <-ctx.Done():
			return commandOutput.Bytes(), ctx.Err()
		default:
			time.Sleep(time.Millisecond * 500)
			// 获取当前缓冲区的所有数据
			currentData := s.out.Bytes()
			if len(currentData) <= initialSize {
				continue
			}
			newData := currentData[initialSize:]
			initialSize += len(newData)
			commandOutput.Write(newData)
			if !s.EndWithPrompt(newData, prompts...) {
				continue
			}
			outputBytes := commandOutput.Bytes()
			// 移除提示符行
			lines := bytes.Split(outputBytes, []byte("\n"))
			// 如果第一行以 cmd 结尾则移除
			if len(lines) > 0 {
				// 检查第一行是否包含命令（更宽松的匹配）
				firstLine := bytes.TrimSpace(lines[0])
				cmdBytes := []byte(cmd)
				// 多种方式检查是否包含命令
				if (bytes.Equal(firstLine, cmdBytes) ||
					bytes.HasSuffix(firstLine, cmdBytes)) && removeFirstLine {
					// 移除第一行
					lines = lines[1:]
				}
			}
			if len(lines) > 0 {
				var resultLines [][]byte
				if removeLastLine {
					resultLines = lines[:len(lines)-1] // 移除最后一行（提示符）
				} else {
					resultLines = lines
				}
				result := bytes.Join(resultLines, []byte("\n"))
				result = bytes.TrimSpace(result)
				if s.HasError(string(result)) {
					return cleanTerminalOutput(result), errors.Errorf("command exec failed")
				} else {
					return cleanTerminalOutput(result), nil
				}
			}
			return cleanTerminalOutput(outputBytes), nil
		}
	}
}

// HasError 检查输出中是否包含错误信息
func (s *SwitchSSHSession) HasError(output string) bool {
	// 常见的错误关键词模式
	var commonErrorPatterns = []string{
		`^%.*$`,                         // 以%开头的错误（华为/华三等设备）
		`^Error:.*$`,                    // Error: 开头的错误
		`\bError\b`,                     // 单词边界Error
		`\bInvalid\b`,                   // Invalid命令或参数
		`\bUnknown\b.*\bcommand\b`,      // Unknown command
		`\bUnrecognized\b.*\bcommand\b`, // Unrecognized command
		`\bIncomplete\b.*\bcommand\b`,   // Incomplete command
		`\bAmbiguous\b.*\bcommand\b`,    // Ambiguous command
		`\bFailed\b`,                    // Failed
		`\bDenied\b`,                    // Denied
		`\bSyntax error\b`,              // 语法错误
		`\bBad parameter\b`,             // 参数错误
		`\bCommand rejected\b`,          // 命令被拒绝
		`\bConfiguration error\b`,       // 配置错误
		`\bResource unavailable\b`,      // 资源不可用
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 首先检查是否为误报
		if s.isFalsePositive(line) {
			continue
		}
		// 检查错误模式
		for _, pattern := range commonErrorPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(line) {
				return true
			}
		}
	}
	return false
}

// isFalsePositive 判断是否为误报
func (s *SwitchSSHSession) isFalsePositive(line string) bool {
	// 直接的误报模式（完全匹配这些模式的行不是真正的错误）
	directFalsePositives := []string{
		`Error-free`,
		`No error`,
		`%age`,
		`percent`,
		`Error Count`,
		`Error Statistics`,
		`Error Information`,
		`No Error Statistics`,
		`Zero Error Count`,
	}

	for _, pattern := range directFalsePositives {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			return true
		}
	}
	// 上下文误报（错误关键词在特定上下文中不是真正的错误）
	contextFalsePositives := map[string]string{
		"Error": "Error-free|No error|Error Statistics|Error Count",
		"%":     "%age|percent",
	}
	for errorPattern, contextPattern := range contextFalsePositives {
		matchedError, _ := regexp.MatchString(errorPattern, line)
		matchedContext, _ := regexp.MatchString(contextPattern, line)
		if matchedError && matchedContext {
			return true
		}
	}
	return false
}

// EndWithPrompt bs是否结束了
func (s *SwitchSSHSession) EndWithPrompt(bs []byte, prompts ...string) bool {
	// 更全面的提示符模式
	commonPrompts := []string{
		"<[^>\n\r]{1,100}>$",
		"\\[[^\\]\n\r]{1,100}\\]$",
		"\\[[Yy]/[Nn]\\]:$",
		"\\[yes/no\\]:$",
	}
	commonPrompts = append(commonPrompts, prompts...)
	for _, prompt := range commonPrompts {
		re := regexp.MustCompile(prompt)
		if re.Match(bs) {
			return true
		}
	}
	return false
}

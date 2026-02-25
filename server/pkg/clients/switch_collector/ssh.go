package switch_collector

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/ssh"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewSwitchSSHClientWithPrivateKey(username, ip string, port int, privateKey string) (*SwitchSSH, error) {
	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, errors.WithMessage(errors.WithStack(err), "parse private key failed")
	}
	authMethod := []ssh.AuthMethod{ssh.PublicKeys(key)}
	return NewSwitchSSHClientWithAuthMethods(username, ip, port, authMethod)
}

func NewSwitchSSHClientWithPassword(username, ip string, port int, password string) (*SwitchSSH, error) {
	authMethod := []ssh.AuthMethod{ssh.Password(password)}
	return NewSwitchSSHClientWithAuthMethods(username, ip, port, authMethod)
}

func NewSwitchSSHClientWithAuthMethods(username, ip string, port int, authMethod []ssh.AuthMethod) (*SwitchSSH, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 10,
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            authMethod,
		Config: ssh.Config{
			// 添加更多密钥交换算法
			KeyExchanges: []string{
				"diffie-hellman-group1-sha1",           // 添加传统算法
				"diffie-hellman-group14-sha1",          // 添加传统算法
				"diffie-hellman-group-exchange-sha1",   // 添加传统算法
				"diffie-hellman-group-exchange-sha256", // 添加传统算法
				"curve25519-sha256@libssh.org",
				"ecdh-sha2-nistp256",
				"ecdh-sha2-nistp384",
				"ecdh-sha2-nistp521",
				"diffie-hellman-group14-sha256",
			},
			// 也可以添加更多密码算法
			Ciphers: []string{
				"aes128-cbc", // 传统算法
				"aes192-cbc", // 传统算法
				"aes256-cbc", // 传统算法
				"aes128-ctr",
				"aes192-ctr",
				"aes256-ctr",
				"3des-cbc", // 传统算法
			},
			// 添加更多MAC算法
			MACs: []string{
				"hmac-sha1", // 传统算法
				"hmac-sha2-256",
				"hmac-sha2-512",
				"hmac-md5", // 传统算法（不安全）
			},
		},
	}
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", ip, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, errors.Wrap(err, "connect to server failed")
	}
	return &SwitchSSH{
		client: sshClient,
	}, nil
}

type SwitchSSH struct {
	client  *ssh.Client
	session *SwitchSSHSession
}

func (s *SwitchSSH) InitSession() error {
	session, err := s.NewSession()
	if err != nil {
		return nil
	}
	s.session = session
	return nil
}
func (s *SwitchSSH) Session() *SwitchSSHSession {
	return s.session
}
func (s *SwitchSSH) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}
func (s *SwitchSSH) GetQosTemplate() ([]*QosTemplate, error) {
	var result = make([]*QosTemplate, 0)
	cmd := "dis current-configuration  | include qos | include bps | include car"
	bs, err := s.Session().Exec(cmd, time.Second*10, true, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	re := regexp.MustCompile(`^qos car (.*) cir (\d+) ([kmg]?bps)$`)
	templates := strings.Split(string(bs), "\n")
	for _, template := range templates {
		template = strings.TrimSpace(template)
		if len(template) == 0 {
			continue
		}
		matches := re.FindStringSubmatch(template)
		if len(matches) < 4 {
			continue
		}
		qos := &QosTemplate{
			Name:    matches[1], // 名称
			Ratebps: 0,          // 速率值
		}
		// 解析速率值
		rateValue, err := strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse qos template speed failed")
		}
		// 转换为bps
		rateValue, err = s.Tobps(rateValue, matches[3])
		if err != nil {
			return nil, errors.Wrap(err, "parse qos template speed failed")
		}
		qos.Ratebps = rateValue
		result = append(result, qos)
	}
	return result, nil
}
func (s *SwitchSSH) Tobps(value int64, unit string) (int64, error) {
	var bpsValue int64 = 0
	switch unit {
	case "bps":
		bpsValue = value * 1
	case "kbps":
		bpsValue = value * 1000
	case "mbps":
		bpsValue = value * 1000 * 1000
	case "gbps":
		bpsValue = value * 1000 * 1000 * 1000
	default:
		return 0, errors.Errorf("invalid qos template speed unit: %v", unit)
	}
	return bpsValue, nil
}

func (s *SwitchSSH) GetIfConfigs() ([]*IfConfig, error) {
	var result = make([]*IfConfig, 0)
	cmd := "dis cu int"
	bs, err := s.Session().Exec(cmd, time.Second*10, true, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ifConfigs := strings.Split(string(bs), "#")
	interfaceRe := regexp.MustCompile("^interface (.+)$")
	qosInRe := regexp.MustCompile("^qos car inbound (.+)$")
	for _, ifConfigStr := range ifConfigs {
		ifConfigStr = strings.TrimSpace(ifConfigStr)
		if len(ifConfigStr) == 0 {
			continue
		}
		lines := strings.Split(ifConfigStr, "\n")
		ifConfig := &IfConfig{}
		if len(lines) >= 1 {
			line := strings.TrimSpace(lines[0])
			matches := interfaceRe.FindStringSubmatch(line)
			if len(matches) >= 2 {
				ifConfig.Name = strings.TrimSpace(matches[1])
			} else {
				// 未获取到接口名字跳过
				continue
			}
		}
		for _, line := range lines {
			line = strings.TrimSpace(line)
			matches := qosInRe.FindStringSubmatch(line)
			if len(matches) >= 2 {
				ifConfig.QosInTemplateName = strings.TrimSpace(matches[1])
			}
		}
		result = append(result, ifConfig)
	}
	qosTemplate, err := s.GetQosTemplate()
	if err != nil {
		return nil, err
	}
	qosTemplateMap := make(map[string]*QosTemplate)
	for _, t := range qosTemplate {
		if len(t.Name) == 0 {
			continue
		}
		qosTemplateMap[t.Name] = t
	}
	for _, ifConfig := range result {
		ifConfig.QosInTemplate = qosTemplateMap[ifConfig.QosInTemplateName]
	}
	return result, nil
}

// UnlimitedInfInSpeed 取消网卡In方向的限速
func (s *SwitchSSH) UnlimitedInfInSpeed(ifName string) ([]byte, error) {
	return s.UnlimitedIfRangeInSpeed(ifName, false)
}

func (s *SwitchSSH) LimitIfInSpeed(ifName string, mbps int64) ([]byte, error) {
	return s.LimitIfRangeInSpeed(ifName, false, mbps)
}

// UnlimitedIfRangeInSpeed 端口范围限速
func (s *SwitchSSH) UnlimitedIfRangeInSpeed(ifNameRange string, isRange bool) ([]byte, error) {
	if isRange {
		ifNameRange = "range " + ifNameRange
	}
	cmds := []string{
		"system-view",
		fmt.Sprintf("int %s", ifNameRange),
		"undo qos car inbound",
		"commit",
		//"quit",
		//"quit",
		//"save",
		//"y",
	}
	out := make([]byte, 0)
	for _, c := range cmds {
		result, err := s.Session().Exec(c, commandExecuteTimeout, false, false)
		if len(result) > 0 {
			out = append(out, result...)
		}
		if err != nil {
			return out, errors.WithMessagef(err, "exec command: %v failed", c)
		}
	}
	return out, nil
}

// LimitIfRangeInSpeed 端口范围限速
func (s *SwitchSSH) LimitIfRangeInSpeed(ifNameRange string, isRange bool, mbps int64) ([]byte, error) {
	qosTemplateName := s.GenerateQosTemplateName(mbps)
	if isRange {
		ifNameRange = "range " + ifNameRange
	}
	cmds := []string{
		"system-view",
		fmt.Sprintf("qos car %s cir %s kbps", qosTemplateName, strconv.FormatInt(mbps*Mbps/Kbps, 10)),
		fmt.Sprintf("int %s", ifNameRange),
		fmt.Sprintf("qos car inbound %v", qosTemplateName),
		"commit",
		//"quit",
		//"quit",
		//"save",
		//"y",
	}
	out := make([]byte, 0)
	for _, c := range cmds {
		result, err := s.Session().Exec(c, commandExecuteTimeout, false, false)
		if len(result) > 0 {
			out = append(out, result...)
		}
		if err != nil {
			return out, errors.WithMessagef(err, "exec command: %v failed", c)
		}
	}
	return out, nil
}

// 可以添加一个清理函数来移除ANSI转义序列
func cleanTerminalOutput(output []byte) []byte {
	// 移除ANSI转义序列的正则表达式
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][0-9;]*\x07`)
	cleaned := ansiRegex.ReplaceAll(output, []byte{})
	return cleaned
}
func (s *SwitchSSH) GenerateQosTemplateName(mbps int64) string {
	str := strconv.FormatInt(mbps, 10)
	return fmt.Sprintf("%sM", str)
}

func (s *SwitchSSH) MustLimitIfInSpeed(ifName string, speedbps int64) (out []byte, err error) {
	mbps := FloorbpsToNearest100MbpsDecimal(speedbps)
	if mbps < MinimumMbps {
		mbps = MinimumMbps
	}
	out, err = s.LimitIfInSpeed(ifName, mbps)
	if err != nil {
		return
	}
	ifConfigs, err := s.GetIfConfigs()
	if err != nil {
		return
	}
	for _, ifConfig := range ifConfigs {
		if ifConfig.Name == ifName {
			if ifConfig.QosInTemplate != nil && ifConfig.QosInTemplate.Ratebps == mbps*Mbps {
				return
			}
		}
	}
	return out, errors.New("limit if in speed failed")
}
func (s *SwitchSSH) MustUnlimitedIfInSpeed(ifName string) (out []byte, err error) {
	out, err = s.UnlimitedInfInSpeed(ifName)
	if err != nil {
		return
	}
	ifConfigs, err := s.GetIfConfigs()
	if err != nil {
		return
	}
	for _, ifConfig := range ifConfigs {
		if ifConfig.Name == ifName {
			if ifConfig.QosInTemplate == nil && ifConfig.QosInTemplateName == "" {
				return
			}
		}
	}
	return out, errors.New("unlimited if in speed failed")
}

// CeilbpsToNearest100MbpsDecimal 按照100M向上取整
func CeilbpsToNearest100MbpsDecimal(speedbps int64) int64 {
	if speedbps <= 0 {
		return 0
	}
	var speed = decimal.NewFromInt(speedbps)
	m := speed.Div(decimal.NewFromInt(100 * Mbps)).Ceil().Mul(decimal.NewFromInt(100)).IntPart()
	return m
}

// FloorbpsToNearest100MbpsDecimal 按照100M向下取整
func FloorbpsToNearest100MbpsDecimal(speedbps int64) int64 {
	if speedbps <= 0 {
		return 0
	}
	var speed = decimal.NewFromInt(speedbps)
	m := speed.Div(decimal.NewFromInt(100 * Mbps)).Floor().Mul(decimal.NewFromInt(100)).IntPart()
	return m
}

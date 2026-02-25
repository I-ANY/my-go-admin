package switch_collector

import "time"

const (
	oidIfIndex      = ".1.3.6.1.2.1.2.2.1.1"     // ifIndex
	oidIfName       = ".1.3.6.1.2.1.31.1.1.1.1"  // ifName
	oidIfOperStatus = ".1.3.6.1.2.1.2.2.1.8"     // ifOperStatus
	oidIfAlias      = ".1.3.6.1.2.1.31.1.1.1.18" // ifAlias
	oidIfDescr      = ".1.3.6.1.2.1.2.2.1.2"     // ifDescr (备用描述)
	oidIfOut        = ".1.3.6.1.2.1.31.1.1.1.10" // ifOutOctets
	oidIfIn         = ".1.3.6.1.2.1.31.1.1.1.6"  // ifInOctets
)
const (
	IfStatusUp             = 1
	IfStatusDown           = 2
	IfStatusTesting        = 3
	IfStatusUnknown        = 4
	IfStatusDormant        = 5
	IfStatusNotPresent     = 6
	IfStatusLowerLayerDown = 7
)
const (
	IfStatusNameUp             string = "up"
	IfStatusNameDown           string = "down"
	IfStatusNameTesting        string = "testing"
	IfStatusNameUnknown        string = "unknown"
	IfStatusNameDormant        string = "dormant"
	IfStatusNameNotPresent     string = "notPresent"
	IfStatusNameLowerLayerDown string = "lowerLayerDown"
)

// 端口状态映射表
var IfStatusMapping = map[int64]string{
	IfStatusUp:             IfStatusNameUp,
	IfStatusDown:           IfStatusNameDown,
	IfStatusTesting:        IfStatusNameTesting,
	IfStatusUnknown:        IfStatusNameUnknown,
	IfStatusDormant:        IfStatusNameDormant,
	IfStatusNotPresent:     IfStatusNameNotPresent,
	IfStatusLowerLayerDown: IfStatusNameLowerLayerDown,
}

const (
	uplinkTag           = "uplink"
	downlinkTag         = "downlink"
	speedLimitTagPrefix = "mf_xs"
)

const (
	// ReasonableByteDiffLimit 合理的流量差值上限（bytes）
	ReasonableByteDiffLimit = 1e12 // 1TB
	// 如果采样间隔是5秒，这相当于 1.6 Tbps 的接口全速运行（转换后）
)
const (
	// BulkWalkRedundancyFactor BulkWalk操作的冗余系数
	//- 10% 的冗余足以应对：
	//- 隐藏的管理接口
	//- 临时的虚拟接口
	//- SNMP实现的差异
	//- 网络波动导致的数据不完整
	BulkWalkRedundancyFactor = 1.0
)

const (
	bps  = 1
	Kbps = 1000 * bps
	Mbps = 1000 * Kbps
	Gbps = 1000 * Mbps
	Tbps = 1000 * Gbps
)
const (
	MinimumMbps = 15
	// MinimumLimitSpeed 最小限速至15M
	MinimumLimitSpeed = MinimumMbps * Mbps
)

// ImportantBusiness 重要业务列表，不能在这些端口上执行任何操作
var ImportantBusiness = []string{
	"MF80",
	"M80",
	"ZIDC",
	"KIDC",
	"ZCDN",
	"BD302",
	"SA302",
	"123_CDN",
	"ZA",
	"MCF",
	"LAF",
	"LA_DX",
	"ZP_DX",
}

const (
	commandExecuteTimeout = 15 * time.Second
)

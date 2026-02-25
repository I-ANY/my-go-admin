package switch_collector

import (
	"github.com/gosnmp/gosnmp"
	"time"
)

// IfInfo 交换机端口信息结构体
type IfInfo struct {
	Index        int64    // 端口索引
	Name         string   // 端口名称
	Status       string   // 端口状态
	Tags         []string // 端口标签
	IsUplink     bool     // 是否为上联口
	IsDownlink   bool     // 是否为下联口
	IsSpeedLimit bool     // 是否限速端口
	PeakBusiness string   // 削峰业务
	Alias        string
	Speed        *IfTrafficSpeed // 端口速度
	//RoomId      int64  // 机房ID
	//SwitchId    int64  // 交换机ID
	//Description string   // 端口描述
}

type IfTrafficCounter struct {
	InTime  time.Time
	OutTime time.Time
	Index   int64  // 端口索引
	InByte  uint64 //byte
	OutByte uint64 //byte
}
type IfTrafficSpeed struct {
	T      time.Time
	Index  int64
	Inbps  float64 // bps
	Outbps float64 // bps
}

type SnmpResult struct {
	Results    []gosnmp.SnmpPDU
	FinishTime time.Time
	StartTime  time.Time
}

type QosTemplate struct {
	Name    string
	Ratebps int64
}

type IfConfig struct {
	Name              string
	QosInTemplateName string
	QosInTemplate     *QosTemplate
}

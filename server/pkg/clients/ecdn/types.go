package ecdn

type Server struct {
	Hostname                 *string  `json:"hostname"`
	Sn                       *string  `json:"sn"`
	Carrier                  *string  `json:"carrier"`  // 运营商
	BwTotal                  *int64   `json:"bwTotal"`  // 总带宽
	BwPlan                   *int64   `json:"bwPlan"`   // 规划带宽
	BwCount                  *int64   `json:"bwCount"`  // 线路数
	BwSingle                 *int64   `json:"bwSingle"` // 线路带宽
	Owner                    *string  `json:"owner"`    // 供应商
	BusinessID               *int64   `json:"businessID" db:"business_id"`
	Business                 *string  `json:"business"`
	Status                   *int64   `json:"status"` // status 0-待审核 1-已上线 4-已审核 5-未通过 6-休眠
	Location                 *string  `json:"location"`
	FrankID                  *string  `json:"frankID" db:"frank_id"`
	Notes                    *string  `json:"notes"`
	ServerID                 *string  `json:"serverID" db:"server_id"`
	DeployStatus             *string  `json:"deployStatus"`             // 业务流状态
	RoomType                 *int64   `json:"roomType"`                 // 机房类型 1-IDC 2-ACDN 3-PCDN
	Province                 *string  `json:"province"`                 // 省份
	IsInterprovincial        *int64   `json:"isInterprovincial"`        // 可跨省 0-未知 1-是 2-否
	IP                       *string  `json:"ip" db:"ip"`               // ip地址
	ReteMirabileStatus       *string  `json:"reteMirabileStatus"`       // 异网标签
	Online                   *bool    `json:"online"`                   // 是否在线 true 在线 false 离线
	AllDay95Utilization      *float64 `json:"allDay95Utilization"`      // 日95利用率
	EveningPeak95Utilization *float64 `json:"eveningPeak95Utilization"` // 晚高峰95利用率
	ProvStatus               *string  `json:"provStatus"`               // 跨省标签
	CactiNotes               *string  `json:"cactiNotes"`               // 监控标签
	SpeedNow                 *int64   `json:"speedNow"`                 // 当前速率
	TimePassed               *string  `json:"timePassed"`               // 审核时间
	Parent                   *string  `json:"parent"`                   // 宿主机
	Origin                   *int64   `json:"origin"`                   // 机房归属 1-自建 2-招募
	ChargeMode               *int64   `json:"chargeMode"`               // 计费类型 1-买断 2-95 3-单机95
}

// 指标数据：某个时间点的上行流量
type Indicator struct {
	Timestamp int64 `json:"timestamp"`
	Up        int64 `json:"up"`
}

// GetServerIndicators 接口返回的数据结构
type ServerIndicatorsData struct {
	Hostname               string      `json:"hostname"`
	Business               string      `json:"business"`
	Carrier                string      `json:"carrier"`
	Owner                  string      `json:"owner"`
	Province               string      `json:"province"`
	Region                 string      `json:"region"`
	ReteMirabileCarrier    string      `json:"reteMirabileCarrier"`    // 异网状态，本网/其他运营商
	ReteMirabileProvincial string      `json:"reteMirabileProvincial"` // 跨省省份
	ReteMirabileRegion     string      `json:"ReteMirabileRegion"`     // 跨省大区
	Provincial             bool        `json:"provincial"`
	Indicators             []Indicator `json:"indicators"`
}

type serverIndicatorsResp struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    *ServerIndicatorsData `json:"data"`
}

type DifIspReq struct {
	FrankIDS   string
	Carrier    int64
	Provincial string
	Phone      string
	Remind     bool
	Note       string
}

type DifIspData struct {
	TaskInfo map[string]struct {
		FrankID   string `json:"frankID"`
		StartTime string `json:"startTime"`
		Carrier   string `json:"carrier"`
		Province  string `json:"province"`
		FalseFlag bool   `json:"falseFlag"`
	} `json:"taskInfo"`
}
type difIspRes struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    *DifIspData `json:"data"`
}

// 节点利用率相关
type NodeUtilizationQueryReq struct {
	Date     string `json:"date"`
	Owner    string `json:"owner,omitempty"`
	Isp      string `json:"isp,omitempty"`
	Location string `json:"location,omitempty"`
}

type NodeUtilizationQueryResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Node     []NodeUtilizationItem `json:"node"`
		Business []NodeUtilizationItem `json:"business"`
	} `json:"data"`
}

type NodeUtilizationItem struct {
	StatisticType int              `json:"statisticType"` // 1=节点 2=保底 3=削峰
	Data          interface{}      `json:"data"`          // 兼容字段
	Graph         UtilizationGraph `json:"graph"`
	Max           float64          `json:"max"`
	Mode          int              `json:"mode"`
	Stat          UtilizationStat  `json:"stat"`
}

type UtilizationGraph struct {
	X []string `json:"x"`
	Y []YData  `json:"y"`
}

type UtilizationStat struct {
	Date             string                 `json:"date"`
	Owner            string                 `json:"owner"`
	OwnerName        string                 `json:"ownerName"`
	Location         string                 `json:"location"`
	Isp              string                 `json:"isp"`
	Business         int                    `json:"business"`
	BwUsageRateDay   float64                `json:"bwUsageRateDay"`
	BwUsageRateNight float64                `json:"bwUsageRateNight"`
	BwNight95Just    string                 `json:"bwNight95Just"`
	BwPlan           string                 `json:"bwPlan"`
	BwFree           string                 `json:"bwFree"`
	Businesses       []BusinessUtilization  `json:"businesses"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

type BusinessUtilization struct {
	Date             string  `json:"date"`
	Owner            string  `json:"owner"`
	OwnerName        string  `json:"ownerName"`
	Location         string  `json:"location"`
	Isp              string  `json:"isp"`
	Business         int     `json:"business"`
	BusinessName     string  `json:"businessName"`
	BwUsageRateDay   float64 `json:"bwUsageRateDay"`
	BwUsageRateNight float64 `json:"bwUsageRateNight"`
	BwNight95Just    string  `json:"bwNight95Just"`
	BwPlan           string  `json:"bwPlan"`
	BwFree           string  `json:"bwFree"`
	IsMain           bool    `json:"isMain"`
	ServerCount      int     `json:"serverCount"`
}

type YData struct {
	DataType      string      `json:"dataType"`
	DataSource    string      `json:"dataSource"`
	Data          []float64   `json:"data"`
	SendBandwidth interface{} `json:"sendBandwidth"`
	PlanBandwidth interface{} `json:"planBandwidth"`
}

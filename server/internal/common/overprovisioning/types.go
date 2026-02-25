package overprovisioning

// Config 超配检测配置
type Config struct {
	PromURL  string
	PromAuth string
}

// PrometheusHostInfo Prometheus 主机信息
type PrometheusHostInfo struct {
	ID              string `json:"id"`
	Hostname        string `json:"hostname"`
	Sn              string `json:"sn"`
	BwPlan          string `json:"bwPlan"`
	BwCount         string `json:"bwcount"`
	BwSingle        string `json:"bwsingle"`
	Day95           string `json:"day95"`
	Evening95       string `json:"evening95"`
	Interprovincial string `json:"interprovincial"`
	Kvm             string `json:"kvm"`
	Location        string `json:"location"`
	Owner           string `json:"owner"`
	IP              string `json:"ip"`
	CactiNotes      string `json:"cactiNotes"`
}

// DetectionResult 检测结果
type DetectionResult struct {
	IsOverProvisioned    bool               `json:"is_over_provisioned"`
	OverProvisionedItems []string           `json:"over_provisioned_items"`
	Details              string             `json:"details"`
	StandardValues       map[string]float64 `json:"standard_values"`
	CurrentValues        map[string]float64 `json:"current_values"`
}

// PromHost Prometheus 主机信息（别名，用于兼容）
type PromHost = PrometheusHostInfo

// BizWithRules 表示存在超配规则的业务小类
type BizWithRules struct {
	BusinessId       int64   `gorm:"column:business_id"`
	BusinessName     *string `gorm:"column:business_name"`
	BusinessCategory *string `gorm:"column:business_category"`
}

// ExecArg 执行参数
type ExecArg struct {
	ExecType string // host=物理机, kvm=KM宿主机, all=全部
}

const (
	SysDisk        = "系统盘"
	Memory         = "内存"
	SSD            = "SSD"
	HDD            = "HDD"
	TotalDataDisk  = "业务盘"
	PlanBw         = "规划带宽"
	SubmitBw       = "提交带宽"
	Day95          = "95值"
	Evening95      = "晚高峰95值"
	DiskType       = "磁盘类型"
	StorageBwRatio = "存储带宽比"
)

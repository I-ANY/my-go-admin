package base

type SpeedLimit struct {
	WorkerConcurrency int64  `mapstructure:"workerConcurrency" default:"15"`
	BatchPortSize     int64  `mapstructure:"batchPortSize" default:"5"` // 限速 - 一次性操作的端口数量
	AlertHook         string `mapstructure:"alertHook"`
	AlertInterval     int64  `mapstructure:"alertInterval" default:"10"` // 告警间隔
}

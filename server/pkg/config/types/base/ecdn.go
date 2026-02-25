package base

type Ecdn struct {
	Url       string `mapstructure:"url"`
	Key       string `mapstructure:"key"`
	DifIspKey string `mapstructure:"difIspKey"`
	// IndicatorsConcurrency ECDN 指标接口并发度，为 0 或未配置时按主机数量自适应
	IndicatorsConcurrency int    `mapstructure:"indicatorsConcurrency"`
	EditServerBwPlanKey   string `mapstructure:"editServerBwPlanKey"`
}

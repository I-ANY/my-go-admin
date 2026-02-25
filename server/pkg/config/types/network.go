package types

import "biz-auto-api/pkg/config/types/base"

type NetworkConfig struct {
	Network               Network         `mapstructure:"network"`
	Mysql                 base.Mysql      `mapstructure:"mysql"`
	EcdnDB                base.Mysql      `mapstructure:"ecdnDB"`
	Log                   base.Log        `mapstructure:"log"`
	TmpDir                string          `mapstructure:"tmpDir" default:"/app/tmp"`
	Redis                 base.Redis      `mapstructure:"redis"`
	Clickhouse            base.Clickhouse `mapstructure:"clickhouse"`
	VpnClient             base.GrpcClient `mapstructure:"vpnClient"`
	DSCPWorkerConcurrency int64           `mapstructure:"dscpWorkerConcurrency" default:"3"`
	SpeedLimitJob         base.SpeedLimit `mapstructure:"speedLimitJob"`
}
type Network struct {
	Host         string `mapstructure:"host" default:"127.0.0.1"`
	Port         int    `mapstructure:"port" default:"8093"`
	ReadTimeout  int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout int    `mapstructure:"writeTimeout" default:"600"`
}

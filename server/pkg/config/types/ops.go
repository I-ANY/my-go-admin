package types

import "biz-auto-api/pkg/config/types/base"

type OpsConfig struct {
	Ops        Ops             `mapstructure:"ops"`
	Mysql      base.Mysql      `mapstructure:"mysql"`
	Log        base.Log        `mapstructure:"log"`
	TmpDir     string          `mapstructure:"tmpDir" default:"/app/tmp"`
	Redis      base.Redis      `mapstructure:"redis"`
	AuthClient base.GrpcClient `mapstructure:"authClient"`
	//Clickhouse base.Clickhouse `mapstructure:"clickhouse"`
}
type Ops struct {
	Host         string `mapstructure:"host" default:"127.0.0.1"`
	Port         int    `mapstructure:"port" default:"8096"`
	ReadTimeout  int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout int    `mapstructure:"writeTimeout" default:"600"`
}

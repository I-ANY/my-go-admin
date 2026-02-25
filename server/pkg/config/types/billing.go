package types

import "biz-auto-api/pkg/config/types/base"

type BillingConfig struct {
	Billing    Billing         `mapstructure:"billing"`
	Mysql      base.Mysql      `mapstructure:"mysql"`
	Log        base.Log        `mapstructure:"log"`
	TmpDir     string          `mapstructure:"tmpDir" default:"/app/tmp"`
	Redis      base.Redis      `mapstructure:"redis"`
	Clickhouse base.Clickhouse `mapstructure:"clickhouse"`
}
type Billing struct {
	Host         string `mapstructure:"host" default:"127.0.0.1"`
	Port         int    `mapstructure:"port" default:"8088"`
	ReadTimeout  int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout int    `mapstructure:"writeTimeout" default:"600"`
}

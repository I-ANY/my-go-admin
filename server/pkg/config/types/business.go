package types

import "biz-auto-api/pkg/config/types/base"

type BusinessConfig struct {
	Business       Business            `mapstructure:"business"`
	Mysql          base.Mysql          `mapstructure:"mysql"`
	Log            base.Log            `mapstructure:"log"`
	FlaskApi       base.FlaskApi       `mapstructure:"flaskApi"`
	TmpDir         string              `mapstructure:"tmpDir" default:"/app/tmp"`
	Redis          base.Redis          `mapstructure:"redis"`
	Ecdn           base.Ecdn           `mapstructure:"ecdn"`
	EcdnDB         base.Mysql          `mapstructure:"ecdnDB"`
	Clickhouse     base.Clickhouse     `mapstructure:"clickhouse"`
	Le             base.Le             `mapstructure:"le"`
	SaAlertWebhook base.SaAlertWebhook `mapstructure:"saAlertWebhook"`
	AliOSS         base.AliOSS         `mapstructure:"aliOSS"`
	ZapBusiness    base.ZapBusiness    `mapstructure:"zapBusiness"`
	AuthClient     base.GrpcClient     `mapstructure:"authClient"`
	Prom           Prom                `mapstructure:"prom"`
}
type Business struct {
	Host         string `mapstructure:"host" default:"127.0.0.1"`
	Port         int    `mapstructure:"port" default:"8088"`
	ReadTimeout  int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout int    `mapstructure:"writeTimeout" default:"600"`
}

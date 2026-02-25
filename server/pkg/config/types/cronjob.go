package types

import "biz-auto-api/pkg/config/types/base"

type CronjobConfig struct {
	Cronjob               Cronjob             `mapstructure:"cronjob"`
	Mysql                 base.Mysql          `mapstructure:"mysql"`
	Log                   base.Log            `mapstructure:"log"`
	StarPortal            base.StarPortal     `mapstructure:"starPortal"`
	FlaskApi              base.FlaskApi       `mapstructure:"flaskApi"`
	Ecdn                  base.Ecdn           `mapstructure:"ecdn"`
	PrometheusUrl         string              `mapstructure:"prometheusUrl"`
	P2p                   string              `mapstructure:"p2p"`
	AutoOpsDomain         string              `mapstructure:"autoOpsDomain"`
	Redis                 base.Redis          `mapstructure:"redis"`
	EcdnDB                base.Mysql          `mapstructure:"ecdnDB"`
	EcdnCk                base.Clickhouse     `mapstructure:"ecdnCk"`
	SaAlertWebhook        base.SaAlertWebhook `mapstructure:"saWebhook"`
	Clickhouse            base.Clickhouse     `mapstructure:"clickhouse"`
	La                    base.LA             `mapstructure:"la"`
	Le                    base.Le             `mapstructure:"le"`
	EcdnInfluxDB          base.InfluxdbV1     `mapstructure:"ecdnInfluxDB"`
	Webhook               WebhookConfig       `mapstructure:"Webhook"`
	Prom                  Prom                `mapstructure:"prom"`
	PromCustom            Prom                `mapstructure:"prom-custom"`
	VpnClient             base.GrpcClient     `mapstructure:"vpnClient"`
	RoomSpeedLimitWebHook string              `mapstructure:"roomSpeedLimitWebHook"`
	LaSlaAlertWebhook     string              `mapstructure:"laSlaAlertWebhook"`
}
type Cronjob struct {
	Host                   string `mapstructure:"host" default:"127.0.0.1"`
	Port                   int    `mapstructure:"port" default:"8089"`
	ReadTimeout            int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout           int    `mapstructure:"writeTimeout" default:"600"`
	NodeName               string `mapstructure:"nodeName"`
	JobLoadIntervalSeconds int    `mapstructure:"jobLoadIntervalSeconds" default:"30"`
}

type WebhookConfig struct {
	KTrafficKey string `mapstructure:"KTrafficKey"`
}

type Prom struct {
	Url  string `mapstructure:"url"`
	Auth string `mapstructure:"auth"`
}

package types

import "biz-auto-api/pkg/config/types/base"

type PriceConfig struct {
	Price      Price           `mapstructure:"price"`
	Mysql      base.Mysql      `mapstructure:"mysql"`
	EcdnMysql  base.Mysql      `mapstructure:"ecdnDB"`
	EcdnUrl    string          `mapstructure:"ecdnUrl"`
	Log        base.Log        `mapstructure:"log"`
	Prom       Prom            `mapstructure:"prom"`
	TmpDir     string          `mapstructure:"tmpDir" default:"/app/tmp"`
	StarPortal base.StarPortal `mapstructure:"starPortal"`
}
type Price struct {
	Host         string `mapstructure:"host" default:"127.0.0.1"`
	Port         int    `mapstructure:"port" default:"8090"`
	ReadTimeout  int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout int    `mapstructure:"writeTimeout" default:"600"`
	NodeName     string `mapstructure:"nodeName"`
}

package types

import "biz-auto-api/pkg/config/types/base"

type SystemConfig struct {
	System     System          `mapstructure:"system"`
	Mysql      base.Mysql      `mapstructure:"mysql"`
	Log        base.Log        `mapstructure:"log"`
	StarPortal base.StarPortal `mapstructure:"starPortal"`
	FlaskApi   base.FlaskApi   `mapstructure:"flaskApi"`
	Redis      base.Redis      `mapstructure:"redis"`
}
type System struct {
	Host           string `mapstructure:"host" default:"127.0.0.1"`
	Port           int    `mapstructure:"port" default:"8088"`
	ReadTimeout    int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout   int    `mapstructure:"writeTimeout" default:"600"`
	TokenExpireSec int    `mapstructure:"tokenExpireSec" default:"86400"`
}

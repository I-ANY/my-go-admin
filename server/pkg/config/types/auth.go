package types

import "biz-auto-api/pkg/config/types/base"

type AuthConfig struct {
	Auth   Auth       `mapstructure:"auth" `
	Log    base.Log   `mapstructure:"log"`
	TmpDir string     `mapstructure:"tmpDir" default:"/app/tmp"`
	Mysql  base.Mysql `mapstructure:"mysql"`
	Redis  base.Redis `mapstructure:"redis"`
}
type Auth struct {
	Host                     string `mapstructure:"host" default:"127.0.0.1"`
	Port                     int    `mapstructure:"port" default:"8095"`
	ReadTimeout              int    `mapstructure:"readTimeout" default:"600"`
	WriteTimeout             int    `mapstructure:"writeTimeout" default:"600"`
	PolicySyncIntervalSecond int64  `mapstructure:"policySyncIntervalSecond" default:"60"`
	PolicyCacheTTLSecond     int64  `mapstructure:"policyCacheTTLSecond" default:"60"`
}

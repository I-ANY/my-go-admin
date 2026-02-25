package base

type Clickhouse struct {
	Host       string `mapstructure:"host" default:"127.0.0.1"`
	Port       int    `mapstructure:"port" default:"9000"`
	Database   string `mapstructure:"database" default:"biz_auto"`
	Username   string `mapstructure:"username" default:"default"`
	Password   string `mapstructure:"password" default:"123456"`
	MaxIdleCon int    `mapstructure:"maxIdleCon" default:"10"`
	MaxOpenCon int    `mapstructure:"maxOpenCon" default:"100"`
}

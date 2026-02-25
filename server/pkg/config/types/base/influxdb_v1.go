package base

type InfluxdbV1 struct {
	Addr     string `mapstructure:"addr" default:"127.0.0.1"`
	Database string `mapstructure:"database"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
}

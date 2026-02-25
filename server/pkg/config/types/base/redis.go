package base

type Redis struct {
	Host     string `mapstructure:"host" default:"127.0.0.1"`
	Port     int    `mapstructure:"port" default:"6379"`
	Username string `mapstructure:"username" default:""`
	Password string `mapstructure:"password" default:"123456"`
	PoolSize int    `mapstructure:"poolSize" default:"20"`
	Timeout  int    `mapstructure:"timeout" default:"10"`
	DB       int    `mapstructure:"db" default:"0"`
}

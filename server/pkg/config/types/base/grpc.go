package base

type GrpcClient struct {
	Address           string `mapstructure:"address" default:"127.0.0.1:9999"`
	ConnTimeoutSecond int    `mapstructure:"connTimeoutSecond" default:"30"`
}

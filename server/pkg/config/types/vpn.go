package types

import "biz-auto-api/pkg/config/types/base"

type VpnConfig struct {
	Vpn    Vpn      `mapstructure:"vpn" `
	Log    base.Log `mapstructure:"log"`
	Switch Switch   `mapstructure:"switch"`
	TmpDir string   `mapstructure:"tmpDir" default:"/app/tmp"`
}
type Vpn struct {
	Host string `mapstructure:"host" default:"127.0.0.1"`
	Port int    `mapstructure:"port" default:"8094"`
}

type Switch struct {
	PrivateKeyFile string `mapstructure:"privateKeyFile" default:"/app/iaas-private.key"`
	PrivateKey     string `mapstructure:"-"`
	Port           int    `mapstructure:"port" default:"22"`
	Username       string `mapstructure:"username" default:"root"`
}

package types

import "biz-auto-api/pkg/config/types/base"

type MigrateConfig struct {
	Mysql base.Mysql `mapstructure:"mysql"`
	Log   base.Log   `mapstructure:"log"`
}

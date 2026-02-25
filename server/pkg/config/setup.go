package config

import (
	"biz-auto-api/pkg/config/types"
	"biz-auto-api/pkg/logger"
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	SystemConfig   *types.SystemConfig
	CronjobConfig  *types.CronjobConfig
	PriceConfig    *types.PriceConfig
	MigrateConfig  *types.MigrateConfig
	BusinessConfig *types.BusinessConfig
	BillingConfig  *types.BillingConfig
	OpsConfig      *types.OpsConfig
	NetworkConfig  *types.NetworkConfig
	VpnConfig      *types.VpnConfig
	AuthConfig     *types.AuthConfig
)

func SetupMigrateConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.MigrateConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	MigrateConfig = &config
	log.Infof("load config from %v success", configYaml)
}

func SetupSystemConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.SystemConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	SystemConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupBusinessConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.BusinessConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	BusinessConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupBillingConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.BillingConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	BillingConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupOpsConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.OpsConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	OpsConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupNetworkConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.NetworkConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	NetworkConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupVpnConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.VpnConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	bs, err := os.ReadFile(config.Switch.PrivateKeyFile)
	if err != nil {
		log.Fatalf("%+v", errors.Wrap(err, "load switch private key failed"))
	}
	config.Switch.PrivateKey = string(bs)
	VpnConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupAuthConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.AuthConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	AuthConfig = &config
	log.Infof("load config from %v success", configYaml)
}
func SetupCronjobConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.CronjobConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	CronjobConfig = &config
	log.Infof("load config from %v success", configYaml)
}

func SetupPriceConfig(configYaml string) {
	log := logger.NewLogger("debug").WithField("engine", "config")
	config := types.PriceConfig{}
	if err := setup(&config, configYaml); err != nil {
		log.Fatalf("%+v", errors.WithMessage(err, "setup failed"))
	}
	PriceConfig = &config
	log.Infof("load config from %v success", configYaml)
}

func setup(config interface{}, configYaml string) error {
	viper.SetConfigFile(configYaml)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "read config %s failed", configYaml)
	}
	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrapf(err, "unmarshal config %s failed", configYaml)
	}
	setDefaults(config)
	return nil
}

func setDefaults(p interface{}) {
	log := logger.NewLogger("debug")
	val := reflect.ValueOf(p).Elem()
	typ := reflect.TypeOf(p).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("default")
		if field.Kind() == reflect.Struct {
			setDefaults(field.Addr().Interface())
		} else if field.Kind() == reflect.Slice {
			// 处理切片赋默认值
			if tag != "" && (field.IsZero() || field.Len() == 0) {
				tempSlice := make([]interface{}, 0)
				err := json.Unmarshal([]byte(tag), &tempSlice)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "Unmarshal to []string failed "))
				}
				// 根据目标类型创建具体类型的切片
				targetType := field.Type().Elem() // 获取切片元素类型
				resultSlice := reflect.MakeSlice(field.Type(), len(tempSlice), len(tempSlice))
				// 类型转换赋值
				for i, item := range tempSlice {
					val := reflect.ValueOf(item)
					// 处理类型不匹配的情况（如JSON数字转float64，但目标是int）
					if !val.Type().ConvertibleTo(targetType) {
						// 尝试通过json重新序列化/反序列化转换类型
						jsonData, _ := json.Marshal(item)
						newVal := reflect.New(targetType)
						if err := json.Unmarshal(jsonData, newVal.Interface()); err == nil {
							val = newVal.Elem()
						}
					}
					if val.Type().ConvertibleTo(targetType) {
						resultSlice.Index(i).Set(val.Convert(targetType))
					} else {
						log.Errorf("cannot convert %v to %s", val.Type(), targetType)
					}
				}
				field.Set(resultSlice)
			}

		} else if field.Interface() == reflect.Zero(field.Type()).Interface() && tag != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(tag)
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint32, reflect.Uint64:
				defaultVal, err := strconv.Atoi(tag)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "convert %v to int/int32/int64/uin8/uint32/uin64 failed", tag))
				}
				field.SetInt(int64(defaultVal))
			case reflect.Float32, reflect.Float64:
				defaultVal, err := strconv.ParseFloat(tag, 64)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "convert %v to float32/float64 failed", tag))
				}
				field.SetFloat(defaultVal)
			case reflect.Bool:
				defaultVal, err := strconv.ParseBool(tag)
				if err != nil {
					log.Fatalf("%+v", errors.Wrapf(err, "convert %v to bool failed", tag))
				}
				field.SetBool(defaultVal)
			}
		}
	}
}

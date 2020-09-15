package config

import (
	"aws-observability.io/collector/pkg/consts"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/service"
	"os"
)

// GetCfgFactory returns AOC/Otel config
func GetCfgFactory() func(otelViper *viper.Viper, f component.Factories) (*configmodels.Config, error) {
	return func(otelViper *viper.Viper, f component.Factories) (*configmodels.Config, error) {
		// AOC supports loading yaml config from SSM parameter store
		if ssmConfigContent, ok := os.LookupEnv(consts.AOC_CONFIG_CONTENT); ok &&
			os.Getenv(consts.RUN_IN_CONTAINER) == consts.RUN_IN_CONTAINER_TRUE {
			fmt.Printf("Reading json consts from from environment: %v = %v\n",
				consts.AOC_CONFIG_CONTENT, ssmConfigContent)
			return sSMConfigLoader(otelViper, f, ssmConfigContent)
		}

		// use OTel yaml consts from input
		otelCfg, err := service.FileLoaderConfigFactory(otelViper, f)
		if err != nil {
			return nil, err
		}
		return otelCfg, nil
	}
}

// sSMConfigLoader set AOC/Otel config from SSM parameter store
func sSMConfigLoader(v *viper.Viper,
	factories component.Factories,
	configContent string) (*configmodels.Config, error) {
	v.SetConfigType(consts.YAML)
	var configBytes = []byte(configContent)
	err := v.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		return nil, fmt.Errorf("error loading SSM consts file %v", err)
	}
	return config.Load(v, factories)
}

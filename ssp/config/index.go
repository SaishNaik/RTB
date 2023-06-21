package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadMainConfig(filePath string) *MainConfig {
	var mainConfig MainConfig
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&mainConfig)
	if err != nil {
		panic(fmt.Errorf("Fatal error marshal config file: %s \n", err))
	}

	return &mainConfig
}

package cmd

import (
	"dynamicflare/service"
	"fmt"

	"github.com/spf13/viper"
)

const (
	AuthKeyVariableName   = "DYNAMICFLARE_AUTH_KEY"
	AuthEmailVariableName = "DYNAMICFLARE_AUTH_EMAIL"
)

func loadConfig(configFile string) *service.Config {
	viper.SetConfigFile(configFile)
	viper.BindEnv(AuthEmailVariableName)
	viper.BindEnv(AuthKeyVariableName)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	t := service.Config{}
	viper.Unmarshal(&t)
	return &t
}

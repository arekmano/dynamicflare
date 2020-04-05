package cmd

import (
	"fmt"

	"github.com/arekmano/dynamicflare/service"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const (
	// AuthKeyVariableName the environment variable name of the auth key for Cloudflare
	AuthKeyVariableName = "DYNAMICFLARE_AUTH_KEY"
	// AuthEmailVariableName the environment variable name of the auth email for Cloudflare
	AuthEmailVariableName = "DYNAMICFLARE_AUTH_EMAIL"
)

func setup(verbose bool, configFile string) *service.Config {
	if err := validate(); err != nil {
		panic(err)
	}
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	config := loadConfig(configFile)
	logrus.
		Debugf("loaded configuration file: %+v", config)
	return config
}

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

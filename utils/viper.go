package utils

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

func SetupViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("UrlShortener")

	viper.SetDefault("port", "80")
	viper.SetDefault("redis.network", "tcp")
	viper.SetDefault("redis.ip", "127.0.0.1:6379")
	viper.SetDefault("redis.poolSize", 10)
	viper.SetDefault("configFile", "")

	configFile := viper.GetString("configFile")
	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Panicf("Unable to read config file: %s", err)
		}
	}
}

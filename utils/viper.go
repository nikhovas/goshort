package utils

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

func SetupViper(configFile string) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("goshort")

	viper.SetDefault("port", "80")
	viper.SetDefault("token", "")
	viper.SetDefault("redis.ip", "127.0.0.1:6379")
	viper.SetDefault("redis.poolSize", 10)

	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Panicf("Unable to read config file: %s", err)
		}
	}
}

package main

import (
	"github.com/spf13/viper"
	"goshort/src"
	"goshort/utils"
)

func main() {
	utils.SetupViper()
	src.AppObject = src.App{}
	//app := src.App{}

	src.AppObject.Initialize(
		viper.GetString("redis.network"),
		viper.GetString("redis.ip"),
		viper.GetInt("redis.poolSize"))

	src.AppObject.Run(":" + viper.GetString("port"))
}

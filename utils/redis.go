package utils

import (
	"github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
)

func ConnectToRedis() *radix.Pool {
	pool, err := radix.NewPool(
		viper.GetString("redis.network"),
		viper.GetString("redis.ip"),
		viper.GetInt("redis.poolSize"))
	if err != nil {
		// handle error
	}

	return pool
}

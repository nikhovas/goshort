package main

import (
	"bytes"
	"github.com/spf13/viper"
	"goshort/kernel"
	"goshort/modules/dbModules"
	"goshort/modules/inputModules"
	"goshort/modules/logModules"
	"goshort/types"
)

//var yamlExample = []byte(`
//inputs:
// server:
//   name: serverInput
//   ip: ''
//   port: 80
//   mode: tcp
//database:
// redis:
//   name: redisDatabase
//   ip: 127.0.0.1:6379
//   port: 6379
//   mode: tcp
//   pool_size: 10
//loggers:
// kafka:
//   name: kafkaLogger
//   ip: localhost
//   port: 9999
//   topic: goshort_logs
// console:
//middlewares:
// - url_normalizer
//limits:
// max_connections: 2000
//`)

var yamlExample = []byte(`
inputs:
 server:
   name: serverInput
   ip: ''
   port: 80
   mode: tcp
database:
 in_memory:
   name: inMemory
   max_elements: 100
loggers:
 kafka:
   name: kafkaLogger
   ip: localhost
   port: 9999
   topic: goshort_logs
 console:
middlewares:
 - url_normalizer
limits:
 max_connections: 2000
`)

func main() {
	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(bytes.NewBuffer(yamlExample))

	var kernelInstance kernel.Kernel

	var inputs []types.InputControllerInterface
	inputsConfig := viper.GetStringMap("inputs")
	for k, v := range inputsConfig {
		var input types.InputControllerInterface
		switch k {
		case "server":
			input = &inputModules.Server{Kernel: &kernelInstance}
		default:
			break
		}
		if input != nil {
			_ = input.Init(v.(map[string]interface{}))
			inputs = append(inputs, input)
		}
	}

	var database types.UrlControllerInterface
	databaseConfig := viper.GetStringMap("database")
	for k, v := range databaseConfig {
		switch k {
		case "redis":
			database = &dbModules.Redis{Kernel: &kernelInstance}
		case "in_memory":
			database = &dbModules.InMemory{Kernel: &kernelInstance}
		default:
			break
		}
		if database != nil {
			_ = database.Init(v.(map[string]interface{}))
			break
		}
	}

	consoleLog := logModules.Console{Kernel: &kernelInstance}
	loggers := []types.LoggerInterface{&consoleLog}

	kernelInstance = kernel.Kernel{
		InputControllers:    inputs,
		UrlController:       database,
		Loggers:             loggers,
		Middlewares:         []types.MiddlewareInterface{},
		ExtraordinaryLogger: &consoleLog,
	}

	err := kernelInstance.Run()
	if err != nil {
		print(err.Error())
	}
}

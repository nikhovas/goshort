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
   port: 9092
   topic: test
   max_retry_attempts: 50
 console:
   name: consoleLogger
   extra_logger: true
   common_logger: true
middlewares:
 - url_normalizer
limits:
 max_connections: 2000
`)

func main() {
	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(bytes.NewBuffer(yamlExample))

	var kernelInstance kernel.Kernel

	inputControllerCreators := map[string]func(kernel *kernel.Kernel) types.InputControllerInterface{
		"server": inputModules.CreateServer,
	}

	urlControllerCreators := map[string]func(kernel *kernel.Kernel) types.UrlControllerInterface{
		"redis":     dbModules.CreateRedis,
		"in_memory": dbModules.CreateInMemory,
	}

	loggerCreators := map[string]func(kernel *kernel.Kernel) types.LoggerInterface{
		"console": logModules.CreateConsole,
		"kafka":   logModules.CreateKafka,
	}

	var inputs []types.InputControllerInterface
	inputsConfig := viper.GetStringMap("inputs")
	for k, v := range inputsConfig {
		creator, ok := inputControllerCreators[k]
		if !ok {
			continue
		} else {
			input := creator(&kernelInstance)
			_ = input.Init(v.(map[string]interface{}))
			inputs = append(inputs, input)
		}
	}

	var database types.UrlControllerInterface
	databaseConfig := viper.GetStringMap("database")
	for k, v := range databaseConfig {
		creator, ok := urlControllerCreators[k]
		if !ok {
			continue
		} else {
			database = creator(&kernelInstance)
			_ = database.Init(v.(map[string]interface{}))
			break
		}
	}

	var loggers []types.LoggerInterface
	loggersConfig := viper.GetStringMap("loggers")
	for k, v := range loggersConfig {
		creator, ok := loggerCreators[k]
		if !ok {
			continue
		} else {
			logger := creator(&kernelInstance)
			_ = logger.Init(v.(map[string]interface{}))
			loggers = append(loggers, logger)
		}
	}

	consoleLog := logModules.Console{Kernel: &kernelInstance}
	kernelInstance = kernel.Kernel{
		InputControllers:    inputs,
		UrlController:       database,
		Loggers:             loggers,
		Middlewares:         []types.MiddlewareInterface{},
		ExtraordinaryLogger: &consoleLog,
	}

	err := kernelInstance.Run(true)
	if err != nil {
		print(err.Error())
	}

	//sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//go func() {
	//	sig := <-sigs
	//	fmt.Println()
	//	fmt.Println(sig)
	//	done <- true
	//}()
	//
	//fmt.Println("awaiting signal")
	//<-done
	//fmt.Println("exiting")
}

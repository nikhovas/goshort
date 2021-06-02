package logModules

import (
	"encoding/json"
	"goshort/kernel"
	"goshort/types"
)

type Console struct {
	name   string
	Kernel *kernel.Kernel
}

func (controller *Console) Init(config map[string]interface{}) error {
	controller.name = config["name"].(string)
	return nil
}

func (controller *Console) Run() error {
	return nil
}

func (controller *Console) Send(le types.Log) {
	b, _ := json.Marshal(le.ToMap())
	println(string(b))
}

func (controller *Console) GetName() string {
	return controller.name
}

func (controller *Console) GetType() string {
	return "ConsoleLogger"
}

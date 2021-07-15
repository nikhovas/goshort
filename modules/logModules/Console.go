package logModules

import (
	"encoding/json"
	"goshort/kernel"
	"goshort/types"
)

type Console struct {
	types.LoggerBase
	name   string
	Kernel *kernel.Kernel
}

func CreateConsole(kernel *kernel.Kernel) types.LoggerInterface {
	return &Console{Kernel: kernel}
}

func (controller *Console) Init(config map[string]interface{}) error {
	_ = controller.LoggerBase.Init(config)
	controller.name = config["name"].(string)
	return nil
}

func (controller *Console) Run() error {
	defer controller.Kernel.OperationDone()
	return nil
}

func (controller *Console) Stop() error {
	return nil
}

func (controller *Console) Send(le types.Log) error {
	b, _ := json.Marshal(le.ToMap())
	println(string(b))
	return nil
}

func (controller *Console) GetName() string {
	return controller.name
}

func (controller *Console) GetType() string {
	return "ConsoleLogger"
}

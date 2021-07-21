package logModules

import (
	"encoding/json"
	"goshort/kernel/utils/other"
	"goshort/types"
	"sync"
)

type Console struct {
	types.LoggerBase
	name   string
	Kernel types.KernelInterface
}

func CreateConsole(kernel types.KernelInterface) types.LoggerInterface {
	return &Console{Kernel: kernel}
}

func (controller *Console) Init(config map[string]interface{}) error {
	_ = controller.LoggerBase.Init(config)
	controller.name = config["name"].(string)
	return nil
}

func (controller *Console) Run(wg *sync.WaitGroup) error {
	wg.Done()
	controller.IsAvailableVal = 1
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

func (controller *Console) SendError(err error) error {
	return controller.Send(other.InterfaceToLogWrapper(err))
}

func (controller *Console) SendBatch(batch *types.LoggingQueueNode) error {
	for batch != nil {
		err := controller.Send(batch.Log)
		if err != nil {
			return err
		}
		batch = batch.Next
	}
	return nil
}

func (controller *Console) GetName() string {
	return controller.name
}

func (controller *Console) GetType() string {
	return "ConsoleLogger"
}

package kernel

import (
	"goshort/types"
	"sync"
)

type InputKernel struct {
	types.ModuleBase
	Inputs []types.InputInterface
	Kernel *Kernel
}

func (inputKernel *InputKernel) Init(config map[string]interface{}) error {
	for k, v := range config {
		creator, ok := inputKernel.Kernel.InputCreators[k]
		if !ok {
			continue
		} else {
			input := creator(inputKernel.Kernel)
			_ = input.Init(v.(map[string]interface{}))
			inputKernel.Inputs = append(inputKernel.Inputs, input)
		}
	}
	return nil
}

func (inputKernel *InputKernel) Run(wg *sync.WaitGroup) error {
	wg.Add(len(inputKernel.Inputs))

	for _, input := range inputKernel.Inputs {
		if err := input.Run(wg); err != nil {
			_ = inputKernel.Kernel.Logger.SendError(err)
			inputKernel.Kernel.ErrorManage(err, input)
		}
	}

	wg.Done()
	inputKernel.Kernel.SetModuleRunState(inputKernel)
	return nil
}

func (inputKernel *InputKernel) Stop() error {
	for _, input := range inputKernel.Inputs {
		if err := input.Stop(); err != nil {
			_ = inputKernel.Kernel.Logger.SendError(err)
			inputKernel.Kernel.ErrorManage(err, input)
		}
	}
	inputKernel.Kernel.SetModuleStopState(inputKernel)
	return nil
}

func (inputKernel *InputKernel) GetName() string {
	return "Kernel.Input"
}

func (inputKernel *InputKernel) GetType() string {
	return "Kernel.Input"
}

package kernel

import (
	"goshort/types"
	"sync"
)

type MiddlewareKernel struct {
	types.ModuleBase

	Middlewares []types.MiddlewareInterface
	Kernel      *Kernel
}

func (middlewareKernel *MiddlewareKernel) Init(config map[string]interface{}) error {
	for k, v := range config {
		creator, ok := middlewareKernel.Kernel.MiddlewareCreators[k]
		if !ok {
			continue
		} else {
			input := creator(middlewareKernel.Kernel)
			_ = input.Init(v.(map[string]interface{}))
			middlewareKernel.Middlewares = append(middlewareKernel.Middlewares, input)
		}
	}
	return nil
}

func (middlewareKernel *MiddlewareKernel) Run(wg *sync.WaitGroup) error {
	wg.Add(len(middlewareKernel.Middlewares))

	for _, middleware := range middlewareKernel.Middlewares {
		if err := middleware.Run(wg); err != nil {
			_ = middlewareKernel.Kernel.Logger.SendError(err)
			middlewareKernel.Kernel.ErrorManage(err, middleware)
		}
	}

	wg.Done()
	middlewareKernel.Kernel.SetModuleRunState(middlewareKernel)
	return nil
}

func (middlewareKernel *MiddlewareKernel) Stop() error {
	for _, middleware := range middlewareKernel.Middlewares {
		if err := middleware.Stop(); err != nil {
			_ = middlewareKernel.Kernel.Logger.SendError(err)
			middlewareKernel.Kernel.ErrorManage(err, middleware)
		}
	}
	middlewareKernel.Kernel.SetModuleStopState(middlewareKernel)
	return nil
}

func (middlewareKernel *MiddlewareKernel) Exec(url *types.Url) error {
	for _, middleware := range middlewareKernel.Middlewares {
		err := middleware.Exec(url)
		if err != nil {
			if middleware.BreakOnError() {
				return err
			} else {
				_ = middlewareKernel.Kernel.Logger.SendError(err)
			}
		}
	}

	return nil
}

func (middlewareKernel *MiddlewareKernel) BreakOnError() bool {
	return true
}

func (middlewareKernel *MiddlewareKernel) GetName() string {
	return "Kernel.Middleware"
}

func (middlewareKernel *MiddlewareKernel) GetType() string {
	return "Kernel.Middleware"
}

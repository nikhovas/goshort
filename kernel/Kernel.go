package kernel

import (
	"goshort/types"
	errors2 "goshort/types/errors"
	"sync"
	"sync/atomic"
)

type Kernel struct {
	types.ModuleBase
	Input               types.InputInterface
	Logger              types.LoggerInterface
	Database            types.DatabaseInterface
	Middleware          types.MiddlewareInterface
	Reconnection        ReconnectionKernel
	Signal              SignalKernel
	DefaultRedirectCode int
	nextOperationNumber uint64
	Stopped             bool
	wg                  sync.WaitGroup

	InputCreators      map[string]func(kernel *Kernel) types.InputInterface
	LoggerCreators     map[string]func(kernel types.KernelInterface) types.LoggerInterface
	DatabaseCreators   map[string]func(kernel *Kernel) types.DatabaseInterface
	MiddlewareCreators map[string]func(kernel *Kernel) types.MiddlewareInterface
}

func (kernel *Kernel) Init(config map[string]interface{}) error {
	kernel.Logger = &LoggingKernel{Kernel: kernel}
	kernel.Database = &DatabaseKernel{Kernel: kernel}
	kernel.Input = &InputKernel{Kernel: kernel}
	kernel.Middleware = &MiddlewareKernel{Kernel: kernel}
	kernel.Reconnection = ReconnectionKernel{Kernel: kernel}
	kernel.Signal = SignalKernel{Kernel: kernel}

	_ = kernel.Logger.Init(config["loggers"].(map[string]interface{}))
	_ = kernel.Middleware.Init(config["middlewares"].(map[string]interface{}))
	_ = kernel.Database.Init(config["database"].(map[string]interface{}))
	_ = kernel.Input.Init(config["inputs"].(map[string]interface{}))

	return nil
}

func (kernel *Kernel) Run(_ *sync.WaitGroup) error {
	_ = kernel.GetNextOperationNumber()

	kernel.wg.Add(4)

	_ = kernel.Logger.Run(&kernel.wg)
	_ = kernel.Middleware.Run(&kernel.wg)
	_ = kernel.Database.Run(&kernel.wg)
	_ = kernel.Input.Run(&kernel.wg)
	_ = kernel.Signal.Run()

	kernel.SetModuleRunState(kernel.Logger)
	_ = kernel.Stop()
	return nil
}

func (kernel *Kernel) Stop() error {
	_ = kernel.Input.Stop()
	_ = kernel.Database.Stop()
	_ = kernel.Middleware.Stop()
	_ = kernel.Logger.Stop()

	kernel.wg.Wait()
	kernel.SetModuleStopState(kernel)

	return nil
}

func (kernel *Kernel) GetNextOperationNumber() uint64 {
	return atomic.AddUint64(&kernel.nextOperationNumber, 1)
}

func (kernel *Kernel) ErrorManage(err error, module types.ModuleInterface) {
	err2, ok := err.(*errors2.BadConnectionError)
	if ok && err2.Retryable {
		kernel.Reconnection.Start(module)
	}
}

func (kernel *Kernel) GetName() string {
	return "Kernel"
}

func (kernel *Kernel) GetType() string {
	return "Kernel"
}

func (kernel *Kernel) ManageModuleStateChange(module types.ModuleInterface, stateName string) {
	_ = kernel.Logger.Send(&types.ModuleStateChangeLog{
		ModuleName: module.GetName(),
		ModuleType: module.GetType(),
		State:      stateName,
	})
}

func (kernel *Kernel) SetModuleRunState(module types.ModuleInterface) {
	kernel.ManageModuleStateChange(module, "Run")
}

func (kernel *Kernel) SetModuleStopState(module types.ModuleInterface) {
	kernel.ManageModuleStateChange(module, "Stop")
}

func (kernel *Kernel) GetLoggerCreators() map[string]func(kernel types.KernelInterface) types.LoggerInterface {
	return kernel.LoggerCreators
}
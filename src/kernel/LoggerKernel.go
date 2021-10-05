package kernel

import (
	"goshort/src/kernel/utils/other"
	"goshort/src/types"
	errors2 "goshort/src/types/errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type LoggingKernel struct {
	types.ModuleBase

	Loggers []types.LoggerInterface
	Kernel  types.KernelInterface

	head      unsafe.Pointer
	isWorking int32

	waitingLogs int64
}

func (loggingKernel *LoggingKernel) SendBatch(_ *types.LoggingQueueNode) error {
	panic("implement me")
}

func (loggingKernel *LoggingKernel) Init(config map[string]interface{}) error {
	for k, v := range config {
		creator, ok := loggingKernel.Kernel.GetLoggerCreators()[k]
		if !ok {
			continue
		} else {
			logger := creator(loggingKernel.Kernel)
			_ = logger.Init(v.(map[string]interface{}))
			loggingKernel.Loggers = append(loggingKernel.Loggers, logger)
		}
	}
	return nil
}

func (loggingKernel *LoggingKernel) Run(wg *sync.WaitGroup) error {
	wg.Add(len(loggingKernel.Loggers))

	var erroredLoggers []struct {
		types.LoggerInterface
		error
	}
	for _, logger := range loggingKernel.Loggers {
		if err := logger.Run(wg); err != nil {

			erroredLoggers = append(erroredLoggers, struct {
				types.LoggerInterface
				error
			}{logger, err})
		}
	}

	wg.Done()

	for _, erroredLogger := range erroredLoggers {
		err := erroredLogger.error
		_ = loggingKernel.SendError(err)
		loggingKernel.Kernel.ErrorManage(err, erroredLogger.LoggerInterface)
	}

	loggingKernel.Kernel.SetModuleRunState(loggingKernel)
	loggingKernel.SetAvailable()

	return nil
}

func (loggingKernel *LoggingKernel) Stop() error {
	for atomic.LoadInt64(&loggingKernel.waitingLogs) != 0 {
		time.Sleep(time.Second)
	}

	loggingKernel.SetUnavailableAndTryGetReconnectionControl()
	loggingKernel.SetDeath()
	for _, logger := range loggingKernel.Loggers {
		if err := logger.Stop(); err != nil {
			_ = loggingKernel.SendError(err)
			loggingKernel.Kernel.ErrorManage(err, logger)
		}
	}
	loggingKernel.Kernel.SetModuleStopState(loggingKernel)
	return nil
}

func (loggingKernel *LoggingKernel) work() {
	batch := (*types.LoggingQueueNode)(atomic.SwapPointer(&loggingKernel.head, nil))

	var counter int64 = int64(-batch.Len())

	for _, logger := range loggingKernel.Loggers {
		if !logger.IsAvailable() {
			continue
		}

		err := logger.SendBatch(batch)
		if err != nil {
			_ = loggingKernel.SendError(err)
			loggingKernel.Kernel.ErrorManage(err, logger)
		}
	}

	atomic.AddInt64(&loggingKernel.waitingLogs, counter)

	atomic.StoreInt32(&loggingKernel.isWorking, 0)

	if loggingKernel.IsAvailable() {
		loggingKernel.tryWork()
	}
}

func (loggingKernel *LoggingKernel) tryWork() {
	if atomic.LoadPointer(&loggingKernel.head) != nil && atomic.SwapInt32(&loggingKernel.isWorking, 1) == 0 {
		go loggingKernel.work()
	}
}

func (loggingKernel *LoggingKernel) LogToLogger(logger types.LoggerInterface, log types.Log) {
	err := logger.Send(log)
	if err != nil {
		_ = loggingKernel.SendError(err)
		loggingKernel.Kernel.ErrorManage(err, logger)
	}
}

func (loggingKernel *LoggingKernel) _log(log types.Log) {
	for _, logger := range loggingKernel.Loggers {
		if !logger.IsAvailable() {
			continue
		}

		loggingKernel.LogToLogger(logger, log)
	}
}

func (loggingKernel *LoggingKernel) Send(element types.Log) error {
	log, ok := element.(types.Log)
	if !ok {
		err, ok := element.(error)
		if ok {
			log = &errors2.SimpleErrorWrapper{Err: err}
		} else {
			return nil
		}
	}

	if !loggingKernel.IsAvailable() {
		for atomic.LoadInt32(&loggingKernel.isWorking) == 1 {
			time.Sleep(time.Second)
		}
		loggingKernel._log(log)
		return nil
	}

	atomic.AddInt64(&loggingKernel.waitingLogs, 1)

	item := new(types.LoggingQueueNode)
	item.Log = log
	item.Next = (*types.LoggingQueueNode)(atomic.LoadPointer(&loggingKernel.head))
	for atomic.CompareAndSwapPointer(&loggingKernel.head, (unsafe.Pointer)(item.Next), (unsafe.Pointer)(item)) {
	}

	loggingKernel.tryWork()
	return nil
}

func (loggingKernel *LoggingKernel) SendError(err error) error {
	return loggingKernel.Send(other.InterfaceToLogWrapper(err))
}

func (loggingKernel *LoggingKernel) LogErrorForward(params ...interface{}) interface{} {
	paramsLength := len(params)
	err, ok := params[paramsLength-1].(error)
	if !ok {
		return params
	}

	_ = loggingKernel.SendError(err)
	return params
}

func (loggingKernel *LoggingKernel) GetName() string {
	return "Kernel.Logger"
}

func (loggingKernel *LoggingKernel) GetType() string {
	return "Kernel.Logger"
}

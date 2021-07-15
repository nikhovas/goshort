package kernel

import (
	"goshort/types"
	errors2 "goshort/types/errors"
)

func (kernel *Kernel) loggerErrorsLoggingFunction(log types.Log) {
	kernel.Log(false, true, log)
}

func (kernel *Kernel) standardControllerErrorsLoggingFunction(log types.Log) {
	kernel.Log(false, false, log)
}

func (kernel *Kernel) LogToLogger(logger types.LoggerInterface, element interface{}) {
	log, ok := element.(types.Log)
	if !ok {
		err, ok := element.(error)
		if ok {
			log = &errors2.SimpleErrorWrapper{Err: err}
		} else {
			return
		}
	}

	err := logger.Send(log)
	if err != nil {
		kernel.Log(false, true, err)
		kernel.BadConnectionErrorManage(err, logger, kernel.loggerErrorsLoggingFunction)
	}
}

func (kernel *Kernel) Log(clientLog bool, extraLog bool, element interface{}) {
	for _, logger := range kernel.Loggers {
		if (clientLog && !logger.ClientConnectionLogs()) || (!clientLog && !logger.SystemLogs()) {
			continue
		}
		if (extraLog && !logger.IsExtraLogger()) || (!extraLog && !logger.IsCommonLogger()) {
			continue
		}
		if !logger.IsAvailable() {
			continue
		}

		kernel.LogToLogger(logger, element)
	}
}

func (kernel *Kernel) SystemLog(element interface{}) {
	kernel.Log(false, false, element)
}

//func (kernel *Kernel) LogError(err error) error {
//	e, ok := err.(types.AdvancedError)
//	if ok {
//		return e
//	} else {
//		return err
//	}
//}

func (kernel *Kernel) LogErrorForward(params ...interface{}) interface{} {
	paramsLength := len(params)
	err, ok := params[paramsLength-1].(error)
	if !ok {
		return params
	}

	kernel.SystemLog(err)
	//_ = kernel.LogError(err)
	return params
}

func (kernel *Kernel) ClientConnectionLog(element interface{}) {
	kernel.Log(true, false, element)
}

func (kernel *Kernel) BadConnectionErrorManage(err error, module types.ModuleInterface, reconnectionLoggerFunction func(log types.Log)) {
	err2, ok := err.(*errors2.BadConnectionError)
	if ok && err2.Retryable {
		kernel.StartReconnection(module, reconnectionLoggerFunction)
	}
}

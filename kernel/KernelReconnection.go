package kernel

import (
	"goshort/types"
	errors2 "goshort/types/errors"
	"time"
)

type ReconnectionTaskStateChange struct {
	TaskId     uint64
	ModuleName string
	EventName  string
}

func (task ReconnectionTaskStateChange) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Kernel.Reconnection.StateChange",
		"event":      task.EventName,
		"type":       "log",
		"moduleName": task.ModuleName,
		"taskId":     task.TaskId,
	}
}

type ReconnectionLogModuleInfo struct {
	ModuleName string
	ModuleType string
}

type SuccessReconnection struct {
	ReconnectionLogModuleInfo
}

func (sr SuccessReconnection) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.Success",
		"type":       "log",
		"moduleName": sr.ModuleName,
		"moduleType": sr.ModuleType,
	}
}

type BadReconnectionAttempt struct {
	ReconnectionLogModuleInfo
	Attempt int
	Limit   int
	Error   error
}

func (bra *BadReconnectionAttempt) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.BadAttempt",
		"type":       "error",
		"moduleName": bra.ModuleName,
		"moduleType": bra.ModuleType,
		"attempt":    bra.Attempt,
		"limit":      bra.Limit,
		"error":      bra.Error,
	}
}

type ReconnectionAttemptsLimit struct {
	ReconnectionLogModuleInfo
	Limit int
}

func (ral *ReconnectionAttemptsLimit) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":       "Common.Reconnection.ReconnectionAttemptsLimit",
		"type":       "log",
		"moduleName": ral.ModuleName,
		"moduleType": ral.ModuleType,
	}
}

func getModuleInfo(module types.ModuleInterface) ReconnectionLogModuleInfo {
	return ReconnectionLogModuleInfo{
		ModuleName: module.GetName(),
		ModuleType: module.GetType(),
	}
}

// TODO: place to the code
var ReconnectionAlert = &errors2.GenericLog{Name: "Common.Reconnection.ReconnectionAlert", IsError: false}

func (kernel *Kernel) reconnectionCore(module types.ModuleInterface, loggerFunction func(log types.Log)) {
	defer kernel.wg.Done()

	operationId := kernel.GetNextOperationNumber()

	maxAttempts := module.GetMaxRetryAttempts()
	retryInterval := module.GetRetryAttemptInterval()

	loggerFunction(&ReconnectionTaskStateChange{
		TaskId:     operationId,
		ModuleName: module.GetName(),
		EventName:  "start",
	})
	defer loggerFunction(&ReconnectionTaskStateChange{
		TaskId:     operationId,
		ModuleName: module.GetName(),
		EventName:  "end",
	})

	for currentAttempt := 0; currentAttempt < maxAttempts; currentAttempt += 1 {
		if kernel.Stopped {
			// TODO: logging here
			return
		}
		err := module.TryReconnect()
		if err == nil {
			module.SetAvailable()
			loggerFunction(&SuccessReconnection{getModuleInfo(module)})
			return
		} else {
			loggerFunction(&BadReconnectionAttempt{
				ReconnectionLogModuleInfo: getModuleInfo(module),
				Error:                     err,
				Attempt:                   currentAttempt + 1,
				Limit:                     maxAttempts,
			})
		}
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}

	module.SetDeath()
	loggerFunction(&ReconnectionAttemptsLimit{
		ReconnectionLogModuleInfo: getModuleInfo(module),
		Limit:                     maxAttempts,
	})
}

func (kernel *Kernel) StartReconnection(module types.ModuleInterface, loggerFunction func(log types.Log)) {
	if module.SetUnavailableAndTryGetReconnectionControl() {
		kernel.wg.Add(1)
		go kernel.reconnectionCore(module, loggerFunction)
	}
}
